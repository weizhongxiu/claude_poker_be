package game

import (
	"context"
	"time"
)

const (
	defaultActionTimeoutSec = 30 // seconds per player action
	autoFoldDelay           = 100 * time.Millisecond
)

// FSMCallbacks is injected by the logic layer so the FSM can notify about state changes
// without importing database packages (keeps engine pure).
type FSMCallbacks struct {
	// OnAction is called after each valid player action (for DB logging and WS broadcast).
	OnAction func(state *GameState, action PlayerAction)
	// OnStageChange is called when the stage advances (for WS broadcast).
	OnStageChange func(state *GameState)
	// OnHandEnd is called when a hand is fully settled.
	OnHandEnd func(result HandEndResult)
}

// HandFSM manages the state machine for one hand.
type HandFSM struct {
	state     *GameState
	actionCh  chan PlayerAction
	timerCh   chan timerSignal
	callbacks FSMCallbacks
	cancel    context.CancelFunc
}

type timerSignal struct {
	seatNo int
	gameID int64
}

// NewHandFSM creates a new FSM for the given game state.
func NewHandFSM(state *GameState, cb FSMCallbacks) *HandFSM {
	return &HandFSM{
		state:     state,
		actionCh:  make(chan PlayerAction, 16),
		timerCh:   make(chan timerSignal, 4),
		callbacks: cb,
	}
}

// SubmitAction enqueues a player action (non-blocking).
func (f *HandFSM) SubmitAction(action PlayerAction) {
	select {
	case f.actionCh <- action:
	default:
	}
}

// Run starts the hand FSM. It blocks until the hand ends or ctx is cancelled.
func (f *HandFSM) Run(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	f.cancel = cancel
	defer cancel()

	// Stage 0: post blinds automatically
	f.postBlinds()

	// Start at PreFlop
	f.state.Stage = StagePreFlop
	f.dealHoleCards()
	f.setFirstToAct(StagePreFlop)
	f.notifyStageChange()
	f.startActionTimer(ctx)

	for {
		select {
		case <-ctx.Done():
			return

		case sig := <-f.timerCh:
			// Timeout: auto-fold (or auto-check if no bet)
			if sig.gameID != f.state.GameID || sig.seatNo != f.state.CurrentSeat {
				continue
			}
			f.handleTimeout()

		case action := <-f.actionCh:
			if !f.isValidTurn(action) {
				continue
			}
			f.processAction(action)
		}

		if f.isHandEnd() {
			f.finishHand()
			return
		}

		if f.isRoundEnd() {
			if !f.advanceStage() {
				// No more stages (all-in run-out), go to showdown
				f.goToShowdown()
				f.finishHand()
				return
			}
			if f.state.Stage == StageShowdown {
				f.finishHand()
				return
			}
			f.notifyStageChange()
			f.startActionTimer(ctx)
		}
	}
}

// --- Blind posting (Stage 0) ---

func (f *HandFSM) postBlinds() {
	s := f.state
	sbSeat := f.nextActiveSeat(s.DealerSeat)
	bbSeat := f.nextActiveSeat(sbSeat)

	// Set positions
	f.assignPositions()

	// Post ante for all players
	if s.Ante > 0 {
		for _, p := range s.Players {
			ante := min64(s.Ante, p.Chips)
			p.Chips -= ante
			p.ForcedBet += ante
			s.Pot += ante
			f.recordAction(p.SeatNo, ActionAnte, ante)
		}
	}

	// Post small blind
	if p, ok := s.Players[sbSeat]; ok {
		sb := min64(s.SmallBlind, p.Chips)
		p.Chips -= sb
		p.ForcedBet += sb
		p.Bet += sb
		s.Pot += sb
		if p.Chips == 0 {
			p.Status = PlayerAllIn
		}
		f.recordAction(sbSeat, ActionBlind, sb)
	}

	// Post big blind
	if p, ok := s.Players[bbSeat]; ok {
		bb := min64(s.BigBlind, p.Chips)
		p.Chips -= bb
		p.ForcedBet += bb
		p.Bet += bb
		s.Pot += bb
		if p.Chips == 0 {
			p.Status = PlayerAllIn
		}
		f.recordAction(bbSeat, ActionBlind, bb)
	}
}

// --- Hole card dealing ---

func (f *HandFSM) dealHoleCards() {
	s := f.state
	// Deal 2 cards to each active player in seat order starting from SB
	sbSeat := f.nextActiveSeat(s.DealerSeat)
	for i := 0; i < 2; i++ {
		seat := sbSeat
		for range s.SeatOrder {
			if p, ok := s.Players[seat]; ok && p.Status == PlayerActive {
				if len(s.Deck) > 0 {
					p.HoleCards = append(p.HoleCards, s.Deck[0])
					s.Deck = s.Deck[1:]
				}
			}
			seat = f.nextActiveSeat(seat)
			if seat == sbSeat {
				break
			}
		}
	}
}

// --- Action processing ---

func (f *HandFSM) processAction(action PlayerAction) {
	s := f.state
	p, ok := s.Players[action.SeatNo]
	if !ok || p.Status != PlayerActive {
		return
	}

	currentBet := f.maxBet()
	startTime := time.Now()

	switch action.Action {
	case ActionFold:
		p.Status = PlayerFolded
		p.FoldStage = s.Stage
		f.recordAction(action.SeatNo, ActionFold, 0)

	case ActionCheck:
		if p.Bet < currentBet {
			return // invalid: must call/fold/raise
		}
		f.recordAction(action.SeatNo, ActionCheck, 0)

	case ActionCall:
		toCall := currentBet - p.Bet
		toCall = min64(toCall, p.Chips)
		p.Chips -= toCall
		p.Bet += toCall
		p.TotalBet += toCall
		s.Pot += toCall
		if s.Stage == StagePreFlop {
			p.IsVPIP = true
		}
		if p.Chips == 0 {
			p.Status = PlayerAllIn
		}
		f.recordAction(action.SeatNo, ActionCall, toCall)

	case ActionBet:
		if currentBet > 0 {
			return // invalid, must raise
		}
		amount := min64(action.Amount, p.Chips)
		p.Chips -= amount
		p.Bet += amount
		p.TotalBet += amount
		s.Pot += amount
		if s.Stage == StagePreFlop {
			p.IsVPIP = true
			p.IsPFR = true
		}
		if p.Chips == 0 {
			p.Status = PlayerAllIn
		}
		f.recordAction(action.SeatNo, ActionBet, amount)

	case ActionRaise:
		toCall := currentBet - p.Bet
		if action.Amount <= currentBet {
			return // raise must be > current bet
		}
		raiseTotal := min64(action.Amount, p.Chips+p.Bet)
		pay := raiseTotal - p.Bet
		if pay > p.Chips {
			pay = p.Chips
		}
		p.Chips -= pay
		p.Bet = raiseTotal
		p.TotalBet += pay
		s.Pot += pay
		_ = toCall
		if s.Stage == StagePreFlop {
			p.IsVPIP = true
			p.IsPFR = true
		}
		if p.Chips == 0 {
			p.Status = PlayerAllIn
		}
		f.recordAction(action.SeatNo, ActionRaise, raiseTotal)

	case ActionAllIn:
		amount := p.Chips
		p.Bet += amount
		p.TotalBet += amount
		p.Chips = 0
		p.Status = PlayerAllIn
		s.Pot += amount
		if s.Stage == StagePreFlop {
			p.IsVPIP = true
			if p.Bet > currentBet {
				p.IsPFR = true
			}
		}
		f.recordAction(action.SeatNo, ActionAllIn, amount)
	}

	_ = startTime
	if f.callbacks.OnAction != nil {
		f.callbacks.OnAction(s, action)
	}

	// Advance current seat
	f.state.CurrentSeat = f.nextToAct(action.SeatNo)
}

func (f *HandFSM) handleTimeout() {
	s := f.state
	p, ok := s.Players[s.CurrentSeat]
	if !ok {
		return
	}
	// Auto-check if no bet to call, otherwise auto-fold
	currentBet := f.maxBet()
	var auto PlayerAction
	auto.SeatNo = s.CurrentSeat
	auto.UserID = p.UserID
	if p.Bet >= currentBet {
		auto.Action = ActionCheck
	} else {
		auto.Action = ActionFold
	}
	f.processAction(auto)
}

func (f *HandFSM) recordAction(seatNo, action int, amount int64) {
	f.state.ActionSeq++
}

// --- Round/Stage logic ---

func (f *HandFSM) isHandEnd() bool {
	active := 0
	for _, p := range f.state.Players {
		if p.Status == PlayerActive {
			active++
		}
	}
	return active <= 1 && f.allInCount() == 0 || f.onlyOneNotFolded()
}

func (f *HandFSM) onlyOneNotFolded() bool {
	notFolded := 0
	for _, p := range f.state.Players {
		if p.Status != PlayerFolded {
			notFolded++
		}
	}
	return notFolded <= 1
}

func (f *HandFSM) allInCount() int {
	count := 0
	for _, p := range f.state.Players {
		if p.Status == PlayerAllIn {
			count++
		}
	}
	return count
}

// isRoundEnd returns true when all active players have acted and bets are equal.
func (f *HandFSM) isRoundEnd() bool {
	maxBet := f.maxBet()
	for _, p := range f.state.Players {
		if p.Status == PlayerActive {
			if p.Bet < maxBet {
				return false
			}
		}
	}
	// Everyone has acted at least once
	return f.state.CurrentSeat == -1 || f.everyoneActed()
}

func (f *HandFSM) everyoneActed() bool {
	// Simple check: current seat has looped back to the aggressor or first player
	// In practice tracked via round trip flag; simplified here
	return true
}

func (f *HandFSM) advanceStage() bool {
	s := f.state
	switch s.Stage {
	case StagePreFlop:
		s.Stage = StageFlop
		f.dealCommunity(3)
	case StageFlop:
		s.Stage = StageTurn
		f.dealCommunity(1)
	case StageTurn:
		s.Stage = StageRiver
		f.dealCommunity(1)
	case StageRiver:
		s.Stage = StageShowdown
		return true
	default:
		return false
	}
	// Reset bets for new round
	f.resetRoundBets()
	f.setFirstToAct(s.Stage)
	return true
}

func (f *HandFSM) goToShowdown() {
	// Run out remaining community cards
	s := f.state
	needed := 5 - len(s.CommunityCards)
	if needed > 0 {
		f.dealCommunity(needed)
	}
	s.Stage = StageShowdown
}

func (f *HandFSM) dealCommunity(n int) {
	s := f.state
	// Burn one card
	if len(s.Deck) > 0 {
		s.Deck = s.Deck[1:]
	}
	for i := 0; i < n && len(s.Deck) > 0; i++ {
		s.CommunityCards = append(s.CommunityCards, s.Deck[0])
		s.Deck = s.Deck[1:]
	}
}

func (f *HandFSM) resetRoundBets() {
	for _, p := range f.state.Players {
		p.Bet = 0
	}
}

// --- Showdown & Settlement ---

func (f *HandFSM) finishHand() {
	s := f.state
	start := time.Now()

	// Mark showdown participants
	for _, p := range s.Players {
		if p.Status != PlayerFolded {
			p.WentToSD = true
		}
	}

	// Evaluate hands
	results := make(map[int]HandResult)
	for seatNo, p := range s.Players {
		if p.Status != PlayerFolded && len(p.HoleCards) > 0 {
			results[seatNo] = EvalBest5(p.HoleCards, s.CommunityCards)
		}
	}

	// Calculate pots
	bets := make(map[int]int64)
	var foldedSeats []int
	for seatNo, p := range s.Players {
		totalIn := p.ForcedBet + p.TotalBet
		if totalIn > 0 {
			bets[seatNo] = totalIn
		}
		if p.Status == PlayerFolded {
			foldedSeats = append(foldedSeats, seatNo)
		}
	}

	pots, _ := CalcPots(bets, foldedSeats)

	// Determine winners and distribute
	playerEndStates := make(map[int]*PlayerEndState)
	for seatNo, p := range s.Players {
		playerEndStates[seatNo] = &PlayerEndState{
			UserID:     p.UserID,
			SeatNo:     seatNo,
			Position:   p.Position,
			HoleCards:  p.HoleCards,
			ForcedBet:  p.ForcedBet,
			TotalBet:   p.TotalBet,
			ChipsStart: p.Chips + p.ForcedBet + p.TotalBet, // reconstruct start
			IsVPIP:     p.IsVPIP,
			IsPFR:      p.IsPFR,
			WentToSD:   p.WentToSD,
			FoldStage:  p.FoldStage,
			IsShowCard: p.Status != PlayerFolded,
		}
		if res, ok := results[seatNo]; ok {
			playerEndStates[seatNo].HandResult = res
		}
	}

	var potResults []PotResult
	isSplit := false
	dealerLeft := f.nextActiveSeat(s.DealerSeat)

	for idx, pot := range pots {
		// Find best hand among eligible players
		var bestResult HandResult
		var bestSeats []int
		for _, seat := range pot.EligibleSeats {
			res, ok := results[seat]
			if !ok {
				continue
			}
			cmp := CompareHands(res, bestResult)
			if bestResult.Rank == 0 || cmp > 0 {
				bestResult = res
				bestSeats = []int{seat}
			} else if cmp == 0 {
				bestSeats = append(bestSeats, seat)
			}
		}

		// Only one eligible player (others folded): they win without showdown
		if len(pot.EligibleSeats) == 1 {
			bestSeats = pot.EligibleSeats
		}

		if len(bestSeats) > 1 {
			isSplit = true
		}

		shares := SplitPot(pot.Amount, bestSeats, dealerLeft)
		for _, share := range shares {
			if pe, ok := playerEndStates[share.SeatNo]; ok {
				pe.ChipsEnd += share.Amount
				pe.IsWinner = true
				pe.Result += share.Amount
			}
		}

		potType := 1
		potIndex := 0
		if idx > 0 {
			potType = 2
			potIndex = idx
		}
		potResults = append(potResults, PotResult{
			PotType:  potType,
			PotIndex: potIndex,
			Amount:   pot.Amount,
			Winners:  shares,
			WinRank:  bestResult.Rank,
			WinDesc:  bestResult.Desc,
		})
	}

	// Compute final ChipsEnd (start chips - invested + winnings)
	for _, pe := range playerEndStates {
		invested := pe.ForcedBet + pe.TotalBet
		pe.ChipsEnd = s.Players[pe.SeatNo].Chips + pe.ChipsEnd // remaining + won
		pe.Result = pe.ChipsEnd - pe.ChipsStart + invested      // net P&L
	}

	// Build player slice
	var players []*PlayerEndState
	for _, pe := range playerEndStates {
		players = append(players, pe)
	}

	// Build snapshots (one per stage played)
	var snapshots []StageSnapshot
	for stage := StagePreFlop; stage <= s.Stage; stage++ {
		snapshots = append(snapshots, BuildStageSnapshot(s, stage, 0, s.ActionSeq))
	}

	result := HandEndResult{
		GameID:         s.GameID,
		HandNo:         s.HandNo,
		ShuffleSeed:    s.ShuffleSeed,
		CommunityCards: s.CommunityCards,
		RunTwiceUsed:   s.RunTwiceUsed,
		RunTwiceBoard2: s.RunTwiceBoard2,
		IsSplitPot:     isSplit,
		DurationMs:     int(time.Since(start).Milliseconds()),
		Players:        players,
		Pots:           potResults,
		Snapshots:      snapshots,
	}

	if f.callbacks.OnHandEnd != nil {
		f.callbacks.OnHandEnd(result)
	}
}

// --- Utility ---

func (f *HandFSM) maxBet() int64 {
	var max int64
	for _, p := range f.state.Players {
		if p.Bet > max {
			max = p.Bet
		}
	}
	return max
}

func (f *HandFSM) nextActiveSeat(fromSeat int) int {
	order := f.state.SeatOrder
	n := len(order)
	// Find fromSeat index
	start := 0
	for i, s := range order {
		if s == fromSeat {
			start = i
			break
		}
	}
	for i := 1; i <= n; i++ {
		seat := order[(start+i)%n]
		if p, ok := f.state.Players[seat]; ok && p.Status == PlayerActive {
			return seat
		}
	}
	return fromSeat
}

func (f *HandFSM) nextToAct(fromSeat int) int {
	order := f.state.SeatOrder
	n := len(order)
	start := 0
	for i, s := range order {
		if s == fromSeat {
			start = i
			break
		}
	}
	for i := 1; i <= n; i++ {
		seat := order[(start+i)%n]
		if p, ok := f.state.Players[seat]; ok && p.Status == PlayerActive {
			return seat
		}
	}
	return -1
}

func (f *HandFSM) setFirstToAct(stage int) {
	s := f.state
	if stage == StagePreFlop {
		// UTG acts first (left of BB = left of left of dealer's left)
		sbSeat := f.nextActiveSeat(s.DealerSeat)
		bbSeat := f.nextActiveSeat(sbSeat)
		s.CurrentSeat = f.nextActiveSeat(bbSeat)
		if s.IsStraddled {
			s.CurrentSeat = f.nextActiveSeat(s.StraddleSeat)
		}
	} else {
		// Post-flop: first active to left of dealer
		s.CurrentSeat = f.nextActiveSeat(s.DealerSeat)
	}
}

func (f *HandFSM) isValidTurn(action PlayerAction) bool {
	return action.SeatNo == f.state.CurrentSeat
}

func (f *HandFSM) notifyStageChange() {
	if f.callbacks.OnStageChange != nil {
		f.callbacks.OnStageChange(f.state)
	}
}

func (f *HandFSM) startActionTimer(ctx context.Context) {
	s := f.state
	if s.CurrentSeat < 0 {
		return
	}
	deadline := time.Duration(defaultActionTimeoutSec) * time.Second
	s.ActionDeadline = time.Now().Add(deadline).UnixMilli()
	gameID := s.GameID
	seatNo := s.CurrentSeat

	go func() {
		select {
		case <-ctx.Done():
			return
		case <-time.After(deadline):
			select {
			case f.timerCh <- timerSignal{seatNo: seatNo, gameID: gameID}:
			default:
			}
		}
	}()
}

func (f *HandFSM) assignPositions() {
	s := f.state
	seats := s.SeatOrder
	n := len(seats)
	if n == 0 {
		return
	}

	// Find dealer index
	dealerIdx := 0
	for i, seat := range seats {
		if seat == s.DealerSeat {
			dealerIdx = i
			break
		}
	}

	posMap := []int{PosSB, PosBB, PosUTG, PosUTG1, PosMP, PosHJ, PosCO}

	for offset := 0; offset < n; offset++ {
		seat := seats[(dealerIdx+1+offset)%n]
		if p, ok := s.Players[seat]; ok {
			if offset == 0 {
				p.Position = PosSB
			} else if offset == 1 {
				p.Position = PosBB
			} else if offset < len(posMap) {
				p.Position = posMap[offset]
			} else {
				p.Position = PosCO
			}
		}
	}
	if p, ok := s.Players[s.DealerSeat]; ok {
		p.Position = PosBTN
	}
}
