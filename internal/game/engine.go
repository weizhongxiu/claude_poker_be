package game

import (
	"context"
	"fmt"
	"sync"
)

// TableConfig holds the configuration needed to start a table.
type TableConfig struct {
	TableID    int64
	SessionID  int64
	SmallBlind int64
	BigBlind   int64
	Ante       int64
	RunTwice   bool
	MaxSeats   int
}

// TableRoom manages one active poker table (one goroutine per table).
type TableRoom struct {
	cfg      TableConfig
	actionCh chan PlayerAction
	players  map[int]*PlayerState // seatNo → player
	mu       sync.RWMutex
	cancel   context.CancelFunc
	fsm      *HandFSM
	handIdx  int
	callbacks FSMCallbacks
}

// Engine manages all active tables.
type Engine struct {
	tables map[int64]*TableRoom
	mu     sync.RWMutex
}

var globalEngine = &Engine{
	tables: make(map[int64]*TableRoom),
}

// GlobalEngine returns the singleton engine instance.
func GlobalEngine() *Engine {
	return globalEngine
}

// StartTable starts a table goroutine. Must be called once per session start.
func (e *Engine) StartTable(cfg TableConfig, cb FSMCallbacks) {
	e.mu.Lock()
	defer e.mu.Unlock()

	ctx, cancel := context.WithCancel(context.Background())
	room := &TableRoom{
		cfg:       cfg,
		actionCh:  make(chan PlayerAction, 32),
		players:   make(map[int]*PlayerState),
		cancel:    cancel,
		callbacks: cb,
	}
	e.tables[cfg.TableID] = room
	go room.run(ctx)
}

// StopTable stops the table goroutine.
func (e *Engine) StopTable(tableID int64) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if room, ok := e.tables[tableID]; ok {
		room.cancel()
		delete(e.tables, tableID)
	}
}

// AddPlayer adds a player to a table room (thread-safe).
func (e *Engine) AddPlayer(tableID int64, p PlayerState) error {
	e.mu.RLock()
	room, ok := e.tables[tableID]
	e.mu.RUnlock()
	if !ok {
		return fmt.Errorf("table %d not found", tableID)
	}
	room.mu.Lock()
	defer room.mu.Unlock()
	room.players[p.SeatNo] = &p
	return nil
}

// RemovePlayer removes a player from a table room (thread-safe).
func (e *Engine) RemovePlayer(tableID, userID int64) {
	e.mu.RLock()
	room, ok := e.tables[tableID]
	e.mu.RUnlock()
	if !ok {
		return
	}
	room.mu.Lock()
	defer room.mu.Unlock()
	for seatNo, p := range room.players {
		if p.UserID == userID {
			delete(room.players, seatNo)
			break
		}
	}
}

// EnsurePlayer adds a player only if they are not already in the room's player map.
// Used before each hand to guarantee late-joiners are included.
func (e *Engine) EnsurePlayer(tableID int64, p PlayerState) {
	e.mu.RLock()
	room, ok := e.tables[tableID]
	e.mu.RUnlock()
	if !ok {
		return
	}
	room.mu.Lock()
	defer room.mu.Unlock()
	if _, exists := room.players[p.SeatNo]; !exists {
		room.players[p.SeatNo] = &p
	}
}

// UpdatePlayerChips updates a player's chip count in the room's persistent player map.
// Called after each hand ends so the next hand starts with the correct chip count.
func (e *Engine) UpdatePlayerChips(tableID int64, userID int64, chips int64) {
	e.mu.RLock()
	room, ok := e.tables[tableID]
	e.mu.RUnlock()
	if !ok {
		return
	}
	room.mu.Lock()
	defer room.mu.Unlock()
	for _, p := range room.players {
		if p.UserID == userID {
			p.Chips = chips
			return
		}
	}
}

// SubmitAction forwards a player action to the correct table's FSM.
func (e *Engine) SubmitAction(tableID int64, action PlayerAction) error {
	e.mu.RLock()
	room, ok := e.tables[tableID]
	e.mu.RUnlock()
	if !ok {
		return fmt.Errorf("table %d not active", tableID)
	}
	if room.fsm == nil {
		return fmt.Errorf("hand not started")
	}
	room.fsm.SubmitAction(action)
	return nil
}

// GetState returns a snapshot of the current game state for a table, or nil if none.
func (e *Engine) GetState(tableID int64) *GameState {
	e.mu.RLock()
	room, ok := e.tables[tableID]
	e.mu.RUnlock()
	if !ok {
		return nil
	}
	room.mu.RLock()
	defer room.mu.RUnlock()
	if room.fsm == nil {
		return nil
	}
	s := *room.fsm.state // shallow copy
	return &s
}

// StartHand begins a new hand on the table.
func (e *Engine) StartHand(tableID int64, gameID int64, handNo string, handIdx int) error {
	e.mu.RLock()
	room, ok := e.tables[tableID]
	e.mu.RUnlock()
	if !ok {
		return fmt.Errorf("table %d not found", tableID)
	}

	room.mu.Lock()
	if len(room.players) < 2 {
		room.mu.Unlock()
		return fmt.Errorf("need at least 2 players")
	}

	// Copy players for this hand; skip players with no chips (awaiting rebuy)
	players := make(map[int]*PlayerState)
	for k, v := range room.players {
		if v.Chips <= 0 {
			continue
		}
		cp := *v
		cp.HoleCards = nil
		cp.Bet = 0
		cp.TotalBet = 0
		cp.ForcedBet = 0
		cp.Status = PlayerActive
		cp.IsVPIP = false
		cp.IsPFR = false
		cp.WentToSD = false
		cp.FoldStage = 0
		players[k] = &cp
	}
	room.mu.Unlock()

	// Poker action goes to the left of each seated player.
	// Table layout: 1=bottom, 2=bottom-right, 3=right, 4=top-right,
	// 5=top, 6=top-left, 7=left, 8=bottom-left, 9=center.
	// "Left" on screen from seat 1 is seat 8, then 7, 6... (screen counter-clockwise = real table clockwise).
	clockwiseSeats := []int{1, 8, 7, 6, 5, 4, 3, 2, 9}
	clockwisePos := make(map[int]int, len(clockwiseSeats))
	for i, s := range clockwiseSeats {
		clockwisePos[s] = i
	}
	seatOrder := make([]int, 0, len(players))
	for seatNo := range players {
		seatOrder = append(seatOrder, seatNo)
	}
	// Sort by clockwise position; unknown seats go last.
	for i := 1; i < len(seatOrder); i++ {
		for j := i; j > 0; j-- {
			pi := clockwisePos[seatOrder[j-1]]
			pj := clockwisePos[seatOrder[j]]
			if _, ok := clockwisePos[seatOrder[j-1]]; !ok {
				pi = 999
			}
			if _, ok := clockwisePos[seatOrder[j]]; !ok {
				pj = 999
			}
			if pi > pj {
				seatOrder[j-1], seatOrder[j] = seatOrder[j], seatOrder[j-1]
			} else {
				break
			}
		}
	}

	// Rotate dealer seat (simple: increment)
	dealerSeat := seatOrder[handIdx%len(seatOrder)]

	// Shuffle deck
	deck, seed := Shuffle(NewDeck())

	state := &GameState{
		TableID:     room.cfg.TableID,
		SessionID:   room.cfg.SessionID,
		GameID:      gameID,
		HandNo:      handNo,
		ShuffleSeed: seed,
		Stage:       StageBlinds,
		DealerSeat:  dealerSeat,
		Pot:         0,
		Deck:        deck,
		Players:     players,
		SeatOrder:   seatOrder,
		SmallBlind:  room.cfg.SmallBlind,
		BigBlind:    room.cfg.BigBlind,
		Ante:        room.cfg.Ante,
		HandIndex:   handIdx,
	}

	fsm := NewHandFSM(state, room.callbacks)
	room.fsm = fsm

	go fsm.Run(context.Background())
	return nil
}

func (r *TableRoom) run(ctx context.Context) {
	// The room goroutine listens for player join/leave events etc.
	// Hand lifecycle is managed by StartHand → FSM.
	<-ctx.Done()
}

// sortInts sorts a slice of ints ascending.
func sortInts(s []int) {
	for i := 1; i < len(s); i++ {
		for j := i; j > 0 && s[j] < s[j-1]; j-- {
			s[j], s[j-1] = s[j-1], s[j]
		}
	}
}
