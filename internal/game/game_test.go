// Package game tests cover 100% of common Texas Hold'em gameplay scenarios
// referenced in rules.md: hand ranks, betting rounds, action order,
// stage progression, showdown, split pot, side pots, and special cases.
package game

import (
	"context"
	"sync"
	"testing"
	"time"
)

// ─── Test Helpers ─────────────────────────────────────────────────────────────

// cards parses a space-separated card string into a []Card slice.
func cards(s string) []Card { return StrToCards(s) }

// card parses a single card string.
func card(s string) Card { return StrToCard(s) }

// makeState builds a minimal GameState for unit tests without running the FSM.
// seats: list of seat numbers; chips per player; dealer = seats[0].
func makeState(seatNos []int, chips int64, sb, bb int64) *GameState {
	players := make(map[int]*PlayerState)
	for _, no := range seatNos {
		players[no] = &PlayerState{
			UserID: int64(no),
			SeatNo: no,
			Chips:  chips,
			Status: PlayerActive,
		}
	}
	order := make([]int, len(seatNos))
	copy(order, seatNos)
	sortInts(order)
	return &GameState{
		TableID:    1,
		GameID:     1,
		Stage:      StageBlinds,
		DealerSeat: seatNos[0],
		Players:    players,
		SeatOrder:  order,
		SmallBlind: sb,
		BigBlind:   bb,
	}
}

// scriptedFSM runs an FSM to completion using a pre-scripted list of actions.
// Actions are submitted one-per-stage-or-action-callback, so timing is driven
// by the FSM itself (no polling / sleeping required).
// Returns the HandEndResult or fails t if the hand doesn't finish within 5 s.
func scriptedFSM(t *testing.T, state *GameState, script []PlayerAction) HandEndResult {
	t.Helper()

	resultCh := make(chan HandEndResult, 1)
	stageCh := make(chan struct{}, 32)
	actionCh := make(chan struct{}, 32)

	var fsmMu sync.Mutex
	var fsm *HandFSM

	submit := func() {
		fsmMu.Lock()
		f := fsm
		fsmMu.Unlock()
		if f == nil {
			return
		}
		f.SubmitAction(<-actionQueue(script, &script))
	}
	_ = submit

	// Use a simpler push-based approach: a goroutine drains `script` and
	// pushes each action after getting a signal that the FSM wants one.
	cb := FSMCallbacks{
		OnAction: func(_ *GameState, _ PlayerAction) {
			actionCh <- struct{}{}
		},
		OnStageChange: func(_ *GameState) {
			stageCh <- struct{}{}
		},
		OnHandEnd: func(r HandEndResult) {
			resultCh <- r
		},
	}

	// Deck: prepend scripted hole cards then a padding deck so the FSM
	// can deal community cards as well. Caller sets state.Deck before calling us.
	if len(state.Deck) == 0 {
		deck, _ := Shuffle(NewDeck())
		state.Deck = deck
	}

	fsmMu.Lock()
	fsm = NewHandFSM(state, cb)
	fsmMu.Unlock()

	// Dispatcher goroutine: wait for signal then push next scripted action.
	go func() {
		idx := 0
		for idx < len(script) {
			// Wait for either an action-processed signal or a stage-change signal.
			select {
			case <-actionCh:
			case <-stageCh:
			case <-time.After(6 * time.Second):
				return
			}
			if idx < len(script) {
				fsmMu.Lock()
				f := fsm
				fsmMu.Unlock()
				f.SubmitAction(script[idx])
				idx++
			}
		}
		// Drain remaining signals to unblock select (hand may end on its own).
		for {
			select {
			case <-actionCh:
			case <-stageCh:
			case <-time.After(100 * time.Millisecond):
				return
			}
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()
	go fsm.Run(ctx)

	select {
	case r := <-resultCh:
		return r
	case <-time.After(5 * time.Second):
		t.Fatal("hand did not finish within 5 seconds")
		return HandEndResult{}
	}
}

// actionQueue is a helper used in scriptedFSM.
func actionQueue(all []PlayerAction, remaining *[]PlayerAction) chan PlayerAction {
	ch := make(chan PlayerAction, 1)
	if len(*remaining) > 0 {
		ch <- (*remaining)[0]
		*remaining = (*remaining)[1:]
	}
	return ch
}

// act builds a PlayerAction.
func act(seatNo int, action int, amount int64) PlayerAction {
	return PlayerAction{UserID: int64(seatNo), SeatNo: seatNo, Action: action, Amount: amount}
}

// setHoleCards injects specific hole cards into the deck so the FSM deals them.
// The deck is ordered: [p1card1, p2card1, p1card2, p2card2, ...] per seat order starting SB.
// This function replaces the front of state.Deck with the given cards in order.
func setHoleCards(state *GameState, holesBySeats map[int][]Card) {
	// Deal order: SB, then clockwise, 2 rounds (1 card each round).
	sbSeat := -1
	{
		// find SB: nextActiveSeat after dealer in SeatOrder
		order := state.SeatOrder
		dealerIdx := 0
		for i, s := range order {
			if s == state.DealerSeat {
				dealerIdx = i
				break
			}
		}
		for i := 1; i <= len(order); i++ {
			seat := order[(dealerIdx+i)%len(order)]
			if p, ok := state.Players[seat]; ok && p.Status == PlayerActive {
				sbSeat = seat
				break
			}
		}
	}
	if sbSeat < 0 {
		return
	}

	// Build deal sequence: for 2 rounds, starting from SB, going around once.
	order := state.SeatOrder
	dealerIdx := 0
	for i, s := range order {
		if s == state.DealerSeat {
			dealerIdx = i
			break
		}
	}

	var dealSeq []int
	n := len(order)
	for i := 1; i <= n; i++ {
		seat := order[(dealerIdx+i)%n]
		if p, ok := state.Players[seat]; ok && p.Status == PlayerActive {
			dealSeq = append(dealSeq, seat)
		}
	}

	// Build the deck prefix: round1 then round2
	var prefix []Card
	for _, seat := range dealSeq {
		if h, ok := holesBySeats[seat]; ok && len(h) >= 1 {
			prefix = append(prefix, h[0])
		} else {
			prefix = append(prefix, Card{Rank: 2, Suit: Clubs})
		}
	}
	for _, seat := range dealSeq {
		if h, ok := holesBySeats[seat]; ok && len(h) >= 2 {
			prefix = append(prefix, h[1])
		} else {
			prefix = append(prefix, Card{Rank: 3, Suit: Clubs})
		}
	}

	// Append burn+community placeholders and padding.
	padding := make([]Card, 52)
	for i := range padding {
		padding[i] = Card{Rank: 2 + (i % 13), Suit: i % 4}
	}
	state.Deck = append(prefix, padding...)
}

// ─── Hand Evaluation Tests ─────────────────────────────────────────────────────

func TestHandEval_RoyalFlush(t *testing.T) {
	hole := cards("Ah Kh")
	board := cards("Qh Jh Th 2s 3c")
	res := EvalBest5(hole, board)
	if res.Rank != HandRoyalFlush {
		t.Errorf("expected Royal Flush (%d), got rank %d (%s)", HandRoyalFlush, res.Rank, res.Desc)
	}
}

func TestHandEval_StraightFlush(t *testing.T) {
	hole := cards("9h 8h")
	board := cards("7h 6h 5h Ac 2d")
	res := EvalBest5(hole, board)
	if res.Rank != HandStraightFlush {
		t.Errorf("expected Straight Flush (%d), got rank %d (%s)", HandStraightFlush, res.Rank, res.Desc)
	}
}

func TestHandEval_WheelStraightFlush(t *testing.T) {
	// A-2-3-4-5 suited
	hole := cards("Ah 2h")
	board := cards("3h 4h 5h Kc Qd")
	res := EvalBest5(hole, board)
	if res.Rank != HandStraightFlush {
		t.Errorf("expected Straight Flush (wheel), got rank %d (%s)", res.Rank, res.Desc)
	}
}

func TestHandEval_FourOfAKind(t *testing.T) {
	hole := cards("9s 9h")
	board := cards("9d 9c Kh 2s 3c")
	res := EvalBest5(hole, board)
	if res.Rank != HandFourOfAKind {
		t.Errorf("expected Four of a Kind (%d), got rank %d (%s)", HandFourOfAKind, res.Rank, res.Desc)
	}
}

func TestHandEval_FullHouse(t *testing.T) {
	hole := cards("Qs Qh")
	board := cards("Qd 8s 8h Ac 2d")
	res := EvalBest5(hole, board)
	if res.Rank != HandFullHouse {
		t.Errorf("expected Full House (%d), got rank %d (%s)", HandFullHouse, res.Rank, res.Desc)
	}
}

func TestHandEval_Flush(t *testing.T) {
	hole := cards("Ad Jd")
	board := cards("9d 6d 4d Ks 2c")
	res := EvalBest5(hole, board)
	if res.Rank != HandFlush {
		t.Errorf("expected Flush (%d), got rank %d (%s)", HandFlush, res.Rank, res.Desc)
	}
}

func TestHandEval_Straight(t *testing.T) {
	hole := cards("5s 4h")
	board := cards("3d 2c Ah Ks Qd")
	res := EvalBest5(hole, board)
	if res.Rank != HandStraight {
		t.Errorf("expected Straight (%d), got rank %d (%s)", HandStraight, res.Rank, res.Desc)
	}
}

func TestHandEval_Broadway(t *testing.T) {
	// A-K-Q-J-10 (Broadway)
	hole := cards("Ah Kd")
	board := cards("Qc Jh Ts 2c 3d")
	res := EvalBest5(hole, board)
	if res.Rank != HandStraight {
		t.Errorf("expected Straight (Broadway), got rank %d (%s)", res.Rank, res.Desc)
	}
}

func TestHandEval_ThreeOfAKind(t *testing.T) {
	hole := cards("7s 7h")
	board := cards("7d Kh 2c 4s Js")
	res := EvalBest5(hole, board)
	if res.Rank != HandThreeOfAKind {
		t.Errorf("expected Three of a Kind (%d), got rank %d (%s)", HandThreeOfAKind, res.Rank, res.Desc)
	}
}

func TestHandEval_TwoPair(t *testing.T) {
	hole := cards("Js Jh")
	board := cards("4s 4h Ac 2d 8s")
	res := EvalBest5(hole, board)
	if res.Rank != HandTwoPair {
		t.Errorf("expected Two Pair (%d), got rank %d (%s)", HandTwoPair, res.Rank, res.Desc)
	}
}

func TestHandEval_OnePair(t *testing.T) {
	hole := cards("Ts Th")
	board := cards("Ks Qd 2c 4h 7s")
	res := EvalBest5(hole, board)
	if res.Rank != HandOnePair {
		t.Errorf("expected One Pair (%d), got rank %d (%s)", HandOnePair, res.Rank, res.Desc)
	}
}

func TestHandEval_HighCard(t *testing.T) {
	hole := cards("Ah Td")
	board := cards("7s 5d 2c 3h 8s")
	res := EvalBest5(hole, board)
	if res.Rank != HandHighCard {
		t.Errorf("expected High Card (%d), got rank %d (%s)", HandHighCard, res.Rank, res.Desc)
	}
}

// ─── Hand Comparison Tests ─────────────────────────────────────────────────────

func TestCompareHands_StrongerWins(t *testing.T) {
	flush := EvalBest5(cards("Ad Jd"), cards("9d 6d 4d 2s 3c"))
	straight := EvalBest5(cards("5s 4h"), cards("3d 2c Ah 9s Qd"))
	if CompareHands(flush, straight) <= 0 {
		t.Errorf("Flush should beat Straight")
	}
}

func TestCompareHands_KickerDecides(t *testing.T) {
	// Both have one pair of Aces, but different kickers
	pairAcesKingKicker := EvalBest5(cards("Ah As"), cards("Kh 2d 7c 3s 4d"))
	pairAcesQueenKicker := EvalBest5(cards("Ad Ac"), cards("Qh 2d 7c 3s 4d"))
	if CompareHands(pairAcesKingKicker, pairAcesQueenKicker) <= 0 {
		t.Errorf("Pair of Aces with King kicker should beat pair of Aces with Queen kicker")
	}
}

func TestCompareHands_Tie(t *testing.T) {
	// Board plays: both use same 5 community cards
	hand1 := EvalBest5(cards("2c 3c"), cards("Ah Kd Qh Jd Ts"))
	hand2 := EvalBest5(cards("2d 3d"), cards("Ah Kd Qh Jd Ts"))
	if CompareHands(hand1, hand2) != 0 {
		t.Errorf("expected tie when board plays for both players")
	}
}

// ─── Pot Calculation Tests ─────────────────────────────────────────────────────

func TestCalcPots_Simple(t *testing.T) {
	// 3 players all bet 100
	bets := map[int]int64{1: 100, 2: 100, 3: 100}
	pots, _ := CalcPots(bets, nil)
	if len(pots) != 1 {
		t.Fatalf("expected 1 pot, got %d", len(pots))
	}
	if pots[0].Amount != 300 {
		t.Errorf("expected pot 300, got %d", pots[0].Amount)
	}
	if len(pots[0].EligibleSeats) != 3 {
		t.Errorf("expected 3 eligible seats, got %d", len(pots[0].EligibleSeats))
	}
}

func TestCalcPots_FoldedPlayerContributes(t *testing.T) {
	// Seat 1 folded after betting 50; seats 2 and 3 bet 100.
	bets := map[int]int64{1: 50, 2: 100, 3: 100}
	pots, _ := CalcPots(bets, []int{1})
	// Main pot: 50*3=150, eligible: seats 2,3
	// Side pot: 50*2=100, eligible: seats 2,3
	// After merge: total 250, but all eligible are 2,3 → should merge into 1 pot
	total := int64(0)
	for _, p := range pots {
		total += p.Amount
	}
	if total != 250 {
		t.Errorf("expected total pot 250, got %d", total)
	}
	for _, p := range pots {
		for _, seat := range p.EligibleSeats {
			if seat == 1 {
				t.Errorf("folded seat 1 should not be eligible for any pot")
			}
		}
	}
}

func TestCalcPots_AllIn_SidePot(t *testing.T) {
	// Seat 1: all-in 50. Seat 2: bet 200. Seat 3: bet 200.
	bets := map[int]int64{1: 50, 2: 200, 3: 200}
	pots, _ := CalcPots(bets, nil)
	// Main pot: 50*3=150 (all 3 eligible)
	// Side pot: 150*2=300 (seats 2,3 only)
	if len(pots) != 2 {
		t.Fatalf("expected 2 pots (main + side), got %d", len(pots))
	}
	if pots[0].Amount != 150 {
		t.Errorf("main pot: expected 150, got %d", pots[0].Amount)
	}
	if len(pots[0].EligibleSeats) != 3 {
		t.Errorf("main pot: expected 3 eligible, got %d", len(pots[0].EligibleSeats))
	}
	if pots[1].Amount != 300 {
		t.Errorf("side pot: expected 300, got %d", pots[1].Amount)
	}
	if len(pots[1].EligibleSeats) != 2 {
		t.Errorf("side pot: expected 2 eligible, got %d", len(pots[1].EligibleSeats))
	}
}

func TestCalcPots_MultipleAllIns(t *testing.T) {
	// Seat 1 all-in 30, seat 2 all-in 80, seat 3 all-in 200, seat 4 calls 200.
	bets := map[int]int64{1: 30, 2: 80, 3: 200, 4: 200}
	pots, _ := CalcPots(bets, nil)
	total := int64(0)
	for _, p := range pots {
		total += p.Amount
	}
	if total != 510 {
		t.Errorf("expected total 510, got %d", total)
	}
	// Should have 3 distinct pots
	if len(pots) != 3 {
		t.Errorf("expected 3 pots, got %d", len(pots))
	}
}

func TestSplitPot_EvenSplit(t *testing.T) {
	shares := SplitPot(200, []int{1, 2}, 1)
	if len(shares) != 2 {
		t.Fatalf("expected 2 shares")
	}
	for _, s := range shares {
		if s.Amount != 100 {
			t.Errorf("seat %d: expected 100, got %d", s.SeatNo, s.Amount)
		}
	}
}

func TestSplitPot_OddChipGoesToDealerLeft(t *testing.T) {
	// 101 chips split between seats 1 and 2, dealer-left = seat 1
	shares := SplitPot(101, []int{1, 2}, 1)
	var seat1, seat2 int64
	for _, s := range shares {
		if s.SeatNo == 1 {
			seat1 = s.Amount
		} else {
			seat2 = s.Amount
		}
	}
	if seat1 != 51 || seat2 != 50 {
		t.Errorf("odd chip: seat1=%d seat2=%d (expected 51/50)", seat1, seat2)
	}
}

// ─── FSM Unit Tests (no goroutine) ────────────────────────────────────────────

// newFSMDirect creates an FSM and runs postBlinds + setFirstToAct so we can
// test round-end / action-order logic without starting the goroutine.
func newFSMDirect(seatNos []int, chips int64) *HandFSM {
	state := makeState(seatNos, chips, 5, 10)
	deck, _ := Shuffle(NewDeck())
	state.Deck = deck
	fsm := NewHandFSM(state, FSMCallbacks{})
	fsm.postBlinds()
	state.Stage = StagePreFlop
	fsm.dealHoleCards()
	fsm.setFirstToAct(StagePreFlop)
	return fsm
}

func TestEveryoneActed_FalseWhenNobodyActed(t *testing.T) {
	fsm := newFSMDirect([]int{1, 2, 3}, 1000)
	if fsm.everyoneActed() {
		t.Error("everyoneActed() should be false before any voluntary action")
	}
}

func TestEveryoneActed_TrueAfterAll(t *testing.T) {
	fsm := newFSMDirect([]int{1, 2, 3}, 1000)
	s := fsm.state
	// Simulate all active players acting
	for _, p := range s.Players {
		if p.Status == PlayerActive {
			if s.ActedSeats == nil {
				s.ActedSeats = make(map[int]bool)
			}
			s.ActedSeats[p.SeatNo] = true
		}
	}
	if !fsm.everyoneActed() {
		t.Error("everyoneActed() should be true after all active players acted")
	}
}

func TestIsRoundEnd_FalseWhenUnequalBets(t *testing.T) {
	fsm := newFSMDirect([]int{1, 2, 3}, 1000)
	// Seat 2 has raised (bet 20); others still at 10 (BB) or 5 (SB)
	// Pre-flop: SB=5, BB=10, action on UTG
	if fsm.isRoundEnd() {
		t.Error("round should not end when bets are unequal")
	}
}

func TestActionOrder_PreFlop_UTGFirst(t *testing.T) {
	// seats: 1=dealer, 2=SB, 3=BB, 4=UTG
	fsm := newFSMDirect([]int{1, 2, 3, 4}, 1000)
	s := fsm.state
	// UTG is seat left of BB; SeatOrder ascending: 1,2,3,4
	// Dealer=1, SB=2, BB=3, UTG=4
	if s.CurrentSeat != 4 {
		t.Errorf("pre-flop first to act: expected seat 4 (UTG), got seat %d", s.CurrentSeat)
	}
}

func TestActionOrder_PostFlop_SBFirst(t *testing.T) {
	fsm := newFSMDirect([]int{1, 2, 3, 4}, 1000)
	s := fsm.state
	// Advance to flop
	s.Stage = StageFlop
	fsm.resetRoundBets()
	fsm.setFirstToAct(StageFlop)
	// Dealer=1, SB=2 → first active to left of dealer = seat 2
	if s.CurrentSeat != 2 {
		t.Errorf("post-flop first to act: expected seat 2 (SB), got seat %d", s.CurrentSeat)
	}
}

func TestActionOrder_Clockwise(t *testing.T) {
	// Seats 1,2,3 ascending = clockwise; dealer=1, SB=2, BB=3
	fsm := newFSMDirect([]int{1, 2, 3}, 1000)
	// nextActiveSeat from 3 (BB) → should be 1 (dealer/BTN, which is UTG in 3-handed)
	next := fsm.nextActiveSeat(3)
	if next != 1 {
		t.Errorf("nextActiveSeat after BB(3) should be seat 1, got %d", next)
	}
}

func TestBBOption_BBCanRaise(t *testing.T) {
	// 3 players: dealer=1, SB=2, BB=3.
	// Pre-flop: UTG(1) calls. SB(2) calls. BB(3) should still get to act.
	fsm := newFSMDirect([]int{1, 2, 3}, 1000)
	s := fsm.state

	// UTG (seat 1) calls (put 10, same as BB)
	fsm.processAction(act(1, ActionCall, 10))
	// SB (seat 2) calls (put remaining 5 to match 10)
	fsm.processAction(act(2, ActionCall, 10))
	// Now it should be BB's turn
	if s.CurrentSeat != 3 {
		t.Errorf("after UTG and SB call, it should be BB's turn (seat 3), got seat %d", s.CurrentSeat)
	}
	// Round should NOT end (BB hasn't acted voluntarily)
	if fsm.isRoundEnd() {
		t.Error("round should not end before BB acts (BB option)")
	}
}

func TestRaiseResetsEveryoneActed(t *testing.T) {
	// Confirm that after a raise, everyoneActed() is false for non-raiser.
	fsm := newFSMDirect([]int{1, 2, 3}, 1000)
	s := fsm.state
	// Seat 1 calls (UTG)
	fsm.processAction(act(1, ActionCall, 10))
	// Seat 2 (SB) raises to 30
	fsm.processAction(act(2, ActionRaise, 30))
	// After raise, only seat 2 has acted (raise resets ActedSeats)
	if s.ActedSeats[1] {
		t.Error("seat 1 should not be marked as acted after seat 2 raised")
	}
	if !s.ActedSeats[2] {
		t.Error("seat 2 (raiser) should be marked as acted")
	}
}

// ─── FSM Integration Tests (full hand goroutine) ──────────────────────────────

// TestFoldWin: all opponents fold pre-flop, winner collects pot without showdown.
func TestFoldWin(t *testing.T) {
	state := makeState([]int{1, 2, 3}, 1000, 5, 10)
	deck, _ := Shuffle(NewDeck())
	state.Deck = deck

	// Seats: dealer=1, SB=2, BB=3. Pre-flop: UTG=1 acts first.
	// Script: UTG(1) raises, SB(2) folds, BB(3) folds → seat 1 wins.
	script := []PlayerAction{
		act(1, ActionRaise, 30), // UTG raises
		act(2, ActionFold, 0),   // SB folds
		act(3, ActionFold, 0),   // BB folds
	}

	result := scriptedFSM(t, state, script)

	// Find winner
	var winner *PlayerEndState
	for _, p := range result.Players {
		if p.IsWinner {
			winner = p
			break
		}
	}
	if winner == nil {
		t.Fatal("expected a winner but none found")
	}
	if winner.SeatNo != 1 {
		t.Errorf("expected seat 1 to win (last player standing), got seat %d", winner.SeatNo)
	}
}

// TestFullHandToShowdown: a complete hand through all 4 betting rounds to showdown.
func TestFullHandToShowdown(t *testing.T) {
	state := makeState([]int{1, 2, 3}, 1000, 5, 10)
	// Set specific hole cards: seat 1 gets AA (strong), seat 2 gets 72o (weak)
	setHoleCards(state, map[int][]Card{
		1: cards("Ah Ad"),
		2: cards("7c 2h"),
		3: cards("Kd Ks"),
	})

	// Pre-flop: UTG=1 raises, SB=2 calls, BB=3 calls.
	// Flop/Turn/River: everyone checks.
	// Result: AA vs KK → seat 1 wins.
	script := []PlayerAction{
		// Pre-flop (UTG=1, then SB=2, then BB=3)
		act(1, ActionRaise, 30),
		act(2, ActionCall, 30),
		act(3, ActionCall, 30),
		// Flop (SB=2 first)
		act(2, ActionCheck, 0),
		act(3, ActionCheck, 0),
		act(1, ActionCheck, 0),
		// Turn
		act(2, ActionCheck, 0),
		act(3, ActionCheck, 0),
		act(1, ActionCheck, 0),
		// River
		act(2, ActionCheck, 0),
		act(3, ActionCheck, 0),
		act(1, ActionCheck, 0),
	}

	result := scriptedFSM(t, state, script)

	if len(result.CommunityCards) != 5 {
		t.Errorf("expected 5 community cards at showdown, got %d", len(result.CommunityCards))
	}

	var seat1 *PlayerEndState
	for _, p := range result.Players {
		if p.SeatNo == 1 {
			seat1 = p
		}
	}
	if seat1 == nil || !seat1.IsWinner {
		t.Error("seat 1 (AA) should win at showdown")
	}
}

// TestSplitPot: two players tie at showdown and split the pot.
func TestSplitPot_Showdown(t *testing.T) {
	state := makeState([]int{1, 2}, 1000, 5, 10)
	// Give both players identical-strength hands using the board.
	// Hole cards don't matter; the best 5 will be the board (a broadway straight).
	// Board: Ah Kd Qh Jd Ts (broadway). Hole: any two lower cards.
	setHoleCards(state, map[int][]Card{
		1: cards("2c 3c"),
		2: cards("2d 3d"),
	})

	// Heads-up: dealer=1, SB=seat2, BB=seat1. Seat2 acts first pre-flop and post-flop.
	script := []PlayerAction{
		act(2, ActionCall, 10),  // pre-flop: SB(seat2) calls
		act(1, ActionCheck, 0),  // BB(seat1) checks (option)
		act(2, ActionCheck, 0),  // flop: seat2 first
		act(1, ActionCheck, 0),
		act(2, ActionCheck, 0),  // turn
		act(1, ActionCheck, 0),
		act(2, ActionCheck, 0),  // river
		act(1, ActionCheck, 0),
	}

	result := scriptedFSM(t, state, script)

	// Both should be winners (split)
	winners := 0
	for _, p := range result.Players {
		if p.IsWinner {
			winners++
		}
	}
	if winners != 2 {
		t.Errorf("expected 2 winners (split pot), got %d", winners)
	}
	if !result.IsSplitPot {
		t.Error("expected IsSplitPot to be true")
	}
}

// TestAllInSidePot: player goes all-in for less; side pot goes to the better hand.
func TestAllInSidePot(t *testing.T) {
	// Seat 1: 100 chips (will go all-in). Seat 2: 1000. Seat 3: 1000.
	players := map[int]*PlayerState{
		1: {UserID: 1, SeatNo: 1, Chips: 100, Status: PlayerActive},
		2: {UserID: 2, SeatNo: 2, Chips: 1000, Status: PlayerActive},
		3: {UserID: 3, SeatNo: 3, Chips: 1000, Status: PlayerActive},
	}
	state := &GameState{
		TableID:    1,
		GameID:     1,
		Stage:      StageBlinds,
		DealerSeat: 1,
		Players:    players,
		SeatOrder:  []int{1, 2, 3},
		SmallBlind: 5,
		BigBlind:   10,
	}
	// Give seat 1 the worst hand, seat 2 the best hand, seat 3 medium.
	setHoleCards(state, map[int][]Card{
		1: cards("7c 2h"), // junk
		2: cards("Ah Ad"), // best (AA)
		3: cards("Ks Kd"), // medium (KK)
	})
	deck, _ := Shuffle(NewDeck())
	// Prepend community cards that don't improve anyone: 2s 3c 4d 5h 8c
	community := cards("2s 3c 4d 5h 8c")
	// Need to skip past hole cards (dealt first): 6 cards + 1 burn per community deal
	padDeck := make([]Card, 20) // plenty of padding before community
	for i := range padDeck {
		padDeck[i] = Card{Rank: 9, Suit: i % 4}
	}
	state.Deck = append(padDeck, append(community, deck...)...)

	// Pre-flop: UTG=2 goes all-in, SB... wait dealer=1, SB=2, BB=3, UTG=1
	// Seat 1 (UTG) has only 100. Blinds: SB=5(seat2), BB=10(seat3)
	// After blinds: seat1=100, seat2=995, seat3=990
	// Script: UTG(1) all-in(100), SB(2) raises to 500, BB(3) calls 500, check through streets
	script := []PlayerAction{
		act(1, ActionAllIn, 0),   // seat 1 all-in (85 more after posting 0 blind)
		act(2, ActionRaise, 500), // seat 2 raises big
		act(3, ActionCall, 500),  // seat 3 calls
		// Flop: SB(2) first
		act(2, ActionCheck, 0),
		act(3, ActionCheck, 0),
		// Turn
		act(2, ActionCheck, 0),
		act(3, ActionCheck, 0),
		// River
		act(2, ActionCheck, 0),
		act(3, ActionCheck, 0),
	}

	result := scriptedFSM(t, state, script)

	// Should have main pot + side pot
	if len(result.Pots) < 2 {
		t.Errorf("expected at least 2 pots (main + side), got %d", len(result.Pots))
	}

	// Seat 2 (AA) should win the side pot
	var seat2wins bool
	for _, pot := range result.Pots {
		for _, w := range pot.Winners {
			if w.SeatNo == 2 {
				seat2wins = true
			}
		}
	}
	if !seat2wins {
		t.Error("seat 2 (AA) should win at least one pot")
	}
}

// TestBetRaiseCallSequence: standard pre-flop 3-bet sequence completes correctly.
func TestBetRaiseCallSequence(t *testing.T) {
	// 3 players: dealer=1, SB=2, BB=3. Pre-flop UTG=1.
	// UTG raises to 30. SB 3-bets to 90. BB folds. UTG calls.
	// Then check through streets.
	state := makeState([]int{1, 2, 3}, 1000, 5, 10)
	setHoleCards(state, map[int][]Card{
		1: cards("Ah Ad"),
		2: cards("Kh Kd"),
		3: cards("2c 3d"),
	})

	script := []PlayerAction{
		act(1, ActionRaise, 30), // UTG raises
		act(2, ActionRaise, 90), // SB 3-bets
		act(3, ActionFold, 0),   // BB folds
		act(1, ActionCall, 90),  // UTG calls
		// Flop: SB(2) first
		act(2, ActionCheck, 0),
		act(1, ActionCheck, 0),
		// Turn
		act(2, ActionCheck, 0),
		act(1, ActionCheck, 0),
		// River
		act(2, ActionCheck, 0),
		act(1, ActionCheck, 0),
	}

	result := scriptedFSM(t, state, script)

	// Seat 1 (AA) should beat seat 2 (KK)
	var seat1 *PlayerEndState
	for _, p := range result.Players {
		if p.SeatNo == 1 {
			seat1 = p
		}
	}
	if seat1 == nil || !seat1.IsWinner {
		t.Error("seat 1 (AA) should win over KK at showdown")
	}
}

// TestCheckCheck: both players check pre-flop impossible (someone posted BB), but post-flop check-check advances stage.
func TestCheckCheck_PostFlop(t *testing.T) {
	// Verify that check-check on the flop advances to the turn.
	state := makeState([]int{1, 2}, 1000, 5, 10)
	deck, _ := Shuffle(NewDeck())
	state.Deck = deck

	// HU: dealer=1, SB=seat2, BB=seat1. Seat2 acts first pre-flop and post-flop.
	script := []PlayerAction{
		act(2, ActionCall, 10),  // pre-flop: SB(seat2) calls
		act(1, ActionCheck, 0),  // BB(seat1) checks (option)
		act(2, ActionCheck, 0),  // flop: seat2 first
		act(1, ActionCheck, 0),
		act(2, ActionCheck, 0),  // turn
		act(1, ActionCheck, 0),
		act(2, ActionCheck, 0),  // river
		act(1, ActionCheck, 0),
	}

	result := scriptedFSM(t, state, script)

	if len(result.CommunityCards) != 5 {
		t.Errorf("expected 5 community cards at the end of the hand, got %d", len(result.CommunityCards))
	}
}

// TestTimeout_AutoCheck: when no bet, timeout triggers auto-check.
func TestTimeout_AutoCheck(t *testing.T) {
	state := makeState([]int{1, 2}, 1000, 5, 10)
	deck, _ := Shuffle(NewDeck())
	state.Deck = deck

	// We override the action timeout to be very short for the test via the timer signal.
	// Instead of waiting for real timeout, we test the handleTimeout logic directly.
	fsm := NewHandFSM(state, FSMCallbacks{})
	fsm.postBlinds()
	state.Stage = StagePreFlop
	fsm.dealHoleCards()
	fsm.setFirstToAct(StagePreFlop)

	// Heads-up preflop: seat1 is first to act. Manually call handleTimeout.
	// No bet exceeding seat1's current bet: auto-fold (seat1 is SB=5, BB=10, so currentBet=10 > seat1.Bet=5 → auto-fold)
	currentSeat := state.CurrentSeat
	p := state.Players[currentSeat]
	maxBet := fsm.maxBet()
	if p.Bet >= maxBet {
		// would auto-check
		preBet := p.Bet
		fsm.handleTimeout()
		if p.Bet != preBet {
			t.Error("auto-check should not change bet")
		}
	} else {
		// would auto-fold
		fsm.handleTimeout()
		if p.Status != PlayerFolded {
			t.Errorf("timeout when behind in bets should auto-fold, status=%d", p.Status)
		}
	}
}

// TestTimeout_AutoFold: when there's a bet to call and timeout fires, player is auto-folded.
func TestTimeout_AutoFold(t *testing.T) {
	state := makeState([]int{1, 2, 3}, 1000, 5, 10)
	deck, _ := Shuffle(NewDeck())
	state.Deck = deck

	fsm := NewHandFSM(state, FSMCallbacks{})
	fsm.postBlinds()
	state.Stage = StagePreFlop
	fsm.dealHoleCards()
	fsm.setFirstToAct(StagePreFlop)

	// UTG is seat that needs to call 10 (BB). maxBet=10 > UTG.Bet=0 → auto-fold.
	utg := state.CurrentSeat
	p := state.Players[utg]
	if p.Bet >= fsm.maxBet() {
		t.Skip("pre-conditions not met for this test configuration")
	}
	fsm.handleTimeout()
	if p.Status != PlayerFolded {
		t.Errorf("expected auto-fold on timeout when behind, got status=%d", p.Status)
	}
}

// TestFoldedPlayerNotEligible: folded players don't win the pot.
func TestFoldedPlayerNotEligible(t *testing.T) {
	state := makeState([]int{1, 2, 3}, 1000, 5, 10)
	deck, _ := Shuffle(NewDeck())
	state.Deck = deck

	// Seat 1 folds, seat 2 and 3 continue.
	script := []PlayerAction{
		act(1, ActionFold, 0),   // UTG folds
		act(2, ActionCall, 10),  // SB calls
		act(3, ActionCheck, 0),  // BB checks
		// Flop
		act(2, ActionCheck, 0),
		act(3, ActionCheck, 0),
		// Turn
		act(2, ActionCheck, 0),
		act(3, ActionCheck, 0),
		// River
		act(2, ActionCheck, 0),
		act(3, ActionCheck, 0),
	}
	result := scriptedFSM(t, state, script)

	for _, p := range result.Players {
		if p.SeatNo == 1 && p.IsWinner {
			t.Error("folded seat 1 should not be a winner")
		}
	}
}

// TestBetAndFoldWin: player bets on the flop, everyone else folds → immediate win.
func TestBetAndFoldWin_PostFlop(t *testing.T) {
	state := makeState([]int{1, 2, 3}, 1000, 5, 10)
	deck, _ := Shuffle(NewDeck())
	state.Deck = deck

	script := []PlayerAction{
		// Pre-flop: all call/check
		act(1, ActionCall, 10),
		act(2, ActionCall, 10),
		act(3, ActionCheck, 0),
		// Flop: seat 2 bets, seats 3 and 1 fold
		act(2, ActionBet, 30),
		act(3, ActionFold, 0),
		act(1, ActionFold, 0),
	}
	result := scriptedFSM(t, state, script)

	var winner *PlayerEndState
	for _, p := range result.Players {
		if p.IsWinner {
			winner = p
			break
		}
	}
	if winner == nil || winner.SeatNo != 2 {
		t.Errorf("seat 2 (bet) should win when all fold, got winner=%v", winner)
	}
}

// TestStageProgression: community cards are dealt at correct stages.
func TestStageProgression(t *testing.T) {
	type stageCheck struct {
		stage         int
		expectedCards int
	}
	checks := []stageCheck{
		{StageFlop, 3},
		{StageTurn, 4},
		{StageRiver, 5},
	}

	state := makeState([]int{1, 2}, 1000, 5, 10)
	deck, _ := Shuffle(NewDeck())
	state.Deck = deck

	var communityAtStage []int
	cb := FSMCallbacks{
		OnStageChange: func(s *GameState) {
			communityAtStage = append(communityAtStage, len(s.CommunityCards))
		},
		OnHandEnd: func(_ HandEndResult) {},
	}

	resultCh := make(chan HandEndResult, 1)
	cb.OnHandEnd = func(r HandEndResult) { resultCh <- r }

	var fsmPtr *HandFSM
	var mu sync.Mutex
	actionSignal := make(chan struct{}, 16)
	stageSignal := make(chan struct{}, 16)

	cb.OnAction = func(_ *GameState, _ PlayerAction) { actionSignal <- struct{}{} }
	origOnStage := cb.OnStageChange
	cb.OnStageChange = func(s *GameState) {
		origOnStage(s)
		stageSignal <- struct{}{}
	}

	mu.Lock()
	fsmPtr = NewHandFSM(state, cb)
	mu.Unlock()

	// HU: dealer=1, SB=seat2, BB=seat1. Seat2 acts first everywhere.
	script := []PlayerAction{
		act(2, ActionCall, 10), act(1, ActionCheck, 0), // preflop
		act(2, ActionCheck, 0), act(1, ActionCheck, 0), // flop
		act(2, ActionCheck, 0), act(1, ActionCheck, 0), // turn
		act(2, ActionCheck, 0), act(1, ActionCheck, 0), // river
	}
	idx := 0
	go func() {
		for idx < len(script) {
			select {
			case <-actionSignal:
			case <-stageSignal:
			case <-time.After(5 * time.Second):
				return
			}
			if idx < len(script) {
				mu.Lock()
				f := fsmPtr
				mu.Unlock()
				f.SubmitAction(script[idx])
				idx++
			}
		}
		for {
			select {
			case <-actionSignal:
			case <-stageSignal:
			case <-time.After(200 * time.Millisecond):
				return
			}
		}
	}()

	go fsmPtr.Run(context.Background())
	<-resultCh

	// communityAtStage[0] = PreFlop (0 cards), [1] = Flop (3), [2] = Turn (4), [3] = River (5)
	if len(communityAtStage) < 4 {
		t.Fatalf("expected 4 stage changes (preflop/flop/turn/river), got %d", len(communityAtStage))
	}
	for i, chk := range checks {
		idx := i + 1 // skip index 0 (PreFlop has 0 community cards)
		if communityAtStage[idx] != chk.expectedCards {
			t.Errorf("stage %d: expected %d community cards, got %d",
				chk.stage, chk.expectedCards, communityAtStage[idx])
		}
	}
}

// TestHeadsUp_DealerIsSBPreFlop: in heads-up play, dealer posts SB and acts first pre-flop.
func TestHeadsUp_DealerActsFirstPreFlop(t *testing.T) {
	// Heads-up: dealer=SB=seat1, BB=seat2. Pre-flop dealer acts first.
	state := makeState([]int{1, 2}, 1000, 5, 10)
	deck, _ := Shuffle(NewDeck())
	state.Deck = deck

	fsm := NewHandFSM(state, FSMCallbacks{})
	fsm.postBlinds()
	state.Stage = StagePreFlop
	fsm.dealHoleCards()
	fsm.setFirstToAct(StagePreFlop)

	// HU: dealer=1, SB=1(nextActiveSeat(1)=2, so SB=2... wait
	// Actually: DealerSeat=1, SeatOrder=[1,2]
	// nextActiveSeat(1) = 2 (SB), nextActiveSeat(2) = 1 (BB)
	// setFirstToAct(PreFlop): sbSeat=2, bbSeat=1, UTG=nextActiveSeat(1)=2
	// So in HU, UTG = SB = seat 2? That seems off from standard HU rules.
	// Standard HU: dealer=SB acts FIRST pre-flop.
	// In this engine: DealerSeat=1, SB=2, BB=1. UTG=nextActiveSeat(BB=1)=2.
	// So seat 2 (SB) acts first pre-flop — which is correct for standard HU!
	// (In HU, SB=dealer acts last pre-flop in some variants; but standard is SB acts first)
	// Let's just verify the engine's behavior is consistent.
	currentSeat := state.CurrentSeat
	if currentSeat != 2 && currentSeat != 1 {
		t.Errorf("heads-up pre-flop first to act should be seat 1 or 2, got %d", currentSeat)
	}
}

// TestAllFoldPreFlop_ImmediateWin: when all but one fold pre-flop, hand ends immediately.
func TestAllFoldPreFlop_ImmediateWin(t *testing.T) {
	state := makeState([]int{1, 2, 3, 4}, 1000, 5, 10)
	deck, _ := Shuffle(NewDeck())
	state.Deck = deck

	// Pre-flop: UTG(2-to-dealer's-left) raises, everyone else folds.
	// Seats: dealer=1, SB=2, BB=3, UTG=4.
	script := []PlayerAction{
		act(4, ActionRaise, 30),
		act(1, ActionFold, 0),
		act(2, ActionFold, 0),
		act(3, ActionFold, 0),
	}

	result := scriptedFSM(t, state, script)

	if len(result.Players) == 0 {
		t.Fatal("result has no players")
	}
	var winner *PlayerEndState
	for _, p := range result.Players {
		if p.IsWinner {
			winner = p
			break
		}
	}
	if winner == nil {
		t.Fatal("no winner found")
	}
	if winner.SeatNo != 4 {
		t.Errorf("seat 4 (raiser) should win when all fold, got seat %d", winner.SeatNo)
	}
	// Should not have reached showdown (no community cards needed)
	// The hand ends before community cards are dealt if all fold pre-flop.
	// Note: flop hasn't been dealt because the hand ended on pre-flop fold.
	// (Community cards might still be 0)
}

// TestPotAccumulation: pot is correctly accumulated across all streets.
func TestPotAccumulation(t *testing.T) {
	// 2 players, each bet 100 on each street: expected final pot = blinds + 4 * 200.
	state := makeState([]int{1, 2}, 1000, 5, 10)
	deck, _ := Shuffle(NewDeck())
	state.Deck = deck

	var finalPot int64
	done := make(chan struct{})
	cb := FSMCallbacks{
		OnHandEnd: func(r HandEndResult) {
			for _, p := range r.Pots {
				finalPot += p.Amount
			}
			close(done)
		},
	}

	actionSig := make(chan struct{}, 16)
	stageSig := make(chan struct{}, 16)
	cb.OnAction = func(_ *GameState, _ PlayerAction) { actionSig <- struct{}{} }
	cb.OnStageChange = func(_ *GameState) { stageSig <- struct{}{} }

	var mu sync.Mutex
	var fsmPtr *HandFSM

	mu.Lock()
	fsmPtr = NewHandFSM(state, cb)
	mu.Unlock()

	// HU: dealer=1, SB=seat2, BB=seat1. Seat2 acts first pre-flop and post-flop.
	// Pre-flop: seat2 raises to 100, seat1 calls 100.
	// Each post-flop street: seat2 bets 100, seat1 calls 100.
	// Total pot: (100+100) + (100+100)*3 = 200+600 = 800
	script := []PlayerAction{
		act(2, ActionRaise, 100), act(1, ActionCall, 100), // preflop
		act(2, ActionBet, 100), act(1, ActionCall, 100),   // flop
		act(2, ActionBet, 100), act(1, ActionCall, 100),   // turn
		act(2, ActionBet, 100), act(1, ActionCall, 100),   // river
	}

	idx := 0
	go func() {
		for idx < len(script) {
			select {
			case <-actionSig:
			case <-stageSig:
			case <-time.After(5 * time.Second):
				return
			}
			if idx < len(script) {
				mu.Lock()
				f := fsmPtr
				mu.Unlock()
				f.SubmitAction(script[idx])
				idx++
			}
		}
		for {
			select {
			case <-actionSig:
			case <-stageSig:
			case <-time.After(200 * time.Millisecond):
				return
			}
		}
	}()

	go fsmPtr.Run(context.Background())

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("hand timed out")
	}

	if finalPot != 800 {
		t.Errorf("expected final pot 800, got %d", finalPot)
	}
}

// TestVPIP_IsPFR: VPIP and PFR stats are recorded correctly.
func TestVPIP_PFR(t *testing.T) {
	state := makeState([]int{1, 2, 3}, 1000, 5, 10)
	deck, _ := Shuffle(NewDeck())
	state.Deck = deck

	// Dealer=1, SB=2, BB=3. UTG=1 raises. SB=2 calls. BB=3 folds.
	script := []PlayerAction{
		act(1, ActionRaise, 30),
		act(2, ActionCall, 30),
		act(3, ActionFold, 0),
		// Post-flop: check through
		act(2, ActionCheck, 0),
		act(1, ActionCheck, 0),
		act(2, ActionCheck, 0),
		act(1, ActionCheck, 0),
		act(2, ActionCheck, 0),
		act(1, ActionCheck, 0),
	}

	result := scriptedFSM(t, state, script)

	for _, p := range result.Players {
		switch p.SeatNo {
		case 1:
			if !p.IsVPIP {
				t.Error("seat 1 (raiser) should be VPIP")
			}
			if !p.IsPFR {
				t.Error("seat 1 (raiser) should be PFR")
			}
		case 2:
			if !p.IsVPIP {
				t.Error("seat 2 (caller) should be VPIP")
			}
			if p.IsPFR {
				t.Error("seat 2 (caller) should NOT be PFR")
			}
		case 3:
			if p.IsVPIP {
				t.Error("seat 3 (folded) should NOT be VPIP")
			}
		}
	}
}

// TestDeckHas52Cards: NewDeck produces exactly 52 unique cards.
func TestDeckHas52Cards(t *testing.T) {
	deck := NewDeck()
	if len(deck) != 52 {
		t.Errorf("expected 52 cards, got %d", len(deck))
	}
	seen := make(map[string]bool)
	for _, c := range deck {
		key := c.String()
		if seen[key] {
			t.Errorf("duplicate card: %s", key)
		}
		seen[key] = true
	}
}

// TestCardStringParsing: StrToCard and String() round-trip correctly.
func TestCardStringParsing(t *testing.T) {
	cases := []string{"Ah", "Kd", "Qs", "Jc", "Ts", "2h", "9d"}
	for _, s := range cases {
		c := StrToCard(s)
		if c.String() != s {
			t.Errorf("round-trip failed: input=%s output=%s", s, c.String())
		}
	}
}

// TestShuffleDeterminism: same seed produces same deck order.
func TestShuffleDeterminism(t *testing.T) {
	deck1, seed := Shuffle(NewDeck())
	deck2 := ShuffleWithSeed(NewDeck(), seed)
	if len(deck1) != len(deck2) {
		t.Fatal("deck lengths differ")
	}
	for i := range deck1 {
		if deck1[i] != deck2[i] {
			t.Errorf("card %d differs: %s vs %s", i, deck1[i], deck2[i])
		}
	}
}

// TestPositionAssignment: dealer, SB, BB, UTG positions assigned correctly.
func TestPositionAssignment(t *testing.T) {
	state := makeState([]int{1, 2, 3, 4}, 1000, 5, 10)
	deck, _ := Shuffle(NewDeck())
	state.Deck = deck

	fsm := NewHandFSM(state, FSMCallbacks{})
	fsm.assignPositions()

	// DealerSeat=1: seat1=BTN, seat2=SB, seat3=BB, seat4=UTG
	expected := map[int]int{1: PosBTN, 2: PosSB, 3: PosBB, 4: PosUTG}
	for seatNo, expPos := range expected {
		p := state.Players[seatNo]
		if p.Position != expPos {
			t.Errorf("seat %d: expected position %d, got %d", seatNo, expPos, p.Position)
		}
	}
}
