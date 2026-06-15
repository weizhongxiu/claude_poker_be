package game

import (
	"math/rand"
	"sync"
	"time"

	"claude-test/internal/game"
)

// botManager tracks which seats are bots on each table.
type botManager struct {
	mu      sync.RWMutex
	tables  map[int64]map[int64]bool // tableID → set of bot userIDs
}

var globalBotMgr = &botManager{
	tables: make(map[int64]map[int64]bool),
}

// RegisterBots records that the given user IDs are bots on tableID.
// Exported so the controller layer can call it directly.
func RegisterBots(tableID int64, botUserIDs []int64) {
	globalBotMgr.mu.Lock()
	defer globalBotMgr.mu.Unlock()
	set := make(map[int64]bool, len(botUserIDs))
	for _, id := range botUserIDs {
		set[id] = true
	}
	globalBotMgr.tables[tableID] = set
}

// UnregisterBots removes bot records for a table.
func UnregisterBots(tableID int64) {
	globalBotMgr.mu.Lock()
	delete(globalBotMgr.tables, tableID)
	globalBotMgr.mu.Unlock()
}

// IsBot returns true if userID is a registered bot on tableID.
func IsBot(tableID, userID int64) bool {
	globalBotMgr.mu.RLock()
	defer globalBotMgr.mu.RUnlock()
	if m, ok := globalBotMgr.tables[tableID]; ok {
		return m[userID]
	}
	return false
}

// MaybeTriggerBot is called after an action or stage change.
// It waits a tiny bit for the FSM to update CurrentSeat, then checks if the
// next player to act is a bot and schedules their action.
func MaybeTriggerBot(state *game.GameState) {
	if state == nil || state.Stage == 0 {
		return
	}
	tableID := state.TableID
	gameID := state.GameID

	go func() {
		// 5ms is enough for the FSM goroutine to update CurrentSeat after this callback returns.
		time.Sleep(5 * time.Millisecond)
		triggerBotForTable(tableID, gameID)
	}()
}

// triggerBotForTable reads fresh state and, if the current seat is a bot, schedules their action.
func triggerBotForTable(tableID, gameID int64) {
	eng := game.GlobalEngine()
	state := eng.GetState(tableID)
	if state == nil || state.GameID != gameID || state.Stage == 0 || state.CurrentSeat < 0 {
		return
	}
	p, ok := state.Players[state.CurrentSeat]
	if !ok || !IsBot(tableID, p.UserID) {
		return
	}

	seatNo := state.CurrentSeat
	stateCopy := *state

	go func() {
		// Simulate thinking time: 600ms–2s
		delay := time.Duration(600+rand.Intn(1400)) * time.Millisecond
		time.Sleep(delay)

		action := game.BotDecide(&stateCopy, seatNo)
		action.TimedAt = time.Now().UnixMilli()

		// Verify it's still this bot's turn
		cur := eng.GetState(tableID)
		if cur == nil || cur.GameID != gameID || cur.CurrentSeat != seatNo {
			return
		}
		_ = eng.SubmitAction(tableID, action)
	}()
}
