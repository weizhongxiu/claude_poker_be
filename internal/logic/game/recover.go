package game

import (
	"context"
	"fmt"
	"time"

	"github.com/gogf/gf/v2/frame/g"

	"claude-test/internal/game"
	"claude-test/utility/ws"
)

// RecoverActiveSessions is called on server startup to re-initialize any sessions
// that were left in "running" state after a restart (engine is in-memory and gets reset).
func RecoverActiveSessions(ctx context.Context) {
	type sessionRow struct {
		Id         int64 `orm:"id"`
		TableId    int64 `orm:"table_id"`
		SmallBlind int64 `orm:"small_blind"`
		BigBlind   int64 `orm:"big_blind"`
		TotalHands int   `orm:"total_hands"`
	}
	var sessions []*sessionRow
	_ = g.DB().Model("room_sessions").
		Where("status", 1).
		Scan(&sessions)

	for _, sess := range sessions {
		if err := recoverSession(ctx, sess.Id, sess.TableId, sess.SmallBlind, sess.BigBlind, sess.TotalHands); err != nil {
			g.Log().Warningf(ctx, "[recover] session %d failed: %v", sess.Id, err)
		} else {
			g.Log().Infof(ctx, "[recover] session %d table %d recovered", sess.Id, sess.TableId)
		}
	}
}

func recoverSession(ctx context.Context, sessionID, tableID, smallBlind, bigBlind int64, totalHands int) error {
	// Mark any incomplete hands (status=1) as abandoned
	_, _ = g.DB().Model("games").
		Where("session_id", sessionID).
		Where("status", 1).
		Data(g.Map{
			"status":   3, // abandoned
			"ended_at": time.Now(),
			"stage":    0,
		}).Update()

	// Reset started_at to now so the duration countdown restarts from the moment of recovery.
	// Without this, sessions started hours ago would expire immediately upon recovery.
	_, _ = g.DB().Model("room_sessions").Where("id", sessionID).
		Data(g.Map{"started_at": time.Now()}).Update()

	// Get seated players
	type seatRow struct {
		UserID   int64  `orm:"user_id"`
		SeatNo   int    `orm:"seat_no"`
		Chips    int64  `orm:"chips"`
		Nickname string `orm:"nickname"`
		Avatar   string `orm:"avatar"`
	}
	var seats []*seatRow
	if e := g.DB().Model("table_seats ts").
		LeftJoin("users u", "u.id = ts.user_id").
		Fields("ts.user_id, ts.seat_no, ts.chips, u.nickname, u.avatar").
		Where("ts.table_id", tableID).Where("ts.status", 1).
		Scan(&seats); e != nil {
		return e
	}
	if len(seats) < 2 {
		return fmt.Errorf("not enough players (%d) for table %d", len(seats), tableID)
	}

	cfg := game.TableConfig{
		TableID:    tableID,
		SessionID:  sessionID,
		SmallBlind: smallBlind,
		BigBlind:   bigBlind,
		MaxSeats:   9,
	}

	cb := game.FSMCallbacks{
		OnAction:      onAction,
		OnStageChange: onStageChange,
		OnHandEnd:     makeOnHandEnd(sessionID),
		OnDeal:        onDeal,
	}

	eng := game.GlobalEngine()
	eng.StartTable(cfg, cb)

	// Detect bots and register all players
	var botIDs []int64
	for _, s := range seats {
		type phoneRow struct{ Phone string `orm:"phone"` }
		var pr phoneRow
		_ = g.DB().Model("users").Fields("phone").Where("id", s.UserID).Scan(&pr)
		if len(pr.Phone) > 4 && pr.Phone[:4] == "bot_" {
			botIDs = append(botIDs, s.UserID)
		}
		_ = eng.AddPlayer(tableID, game.PlayerState{
			UserID:   s.UserID,
			Nickname: s.Nickname,
			Avatar:   s.Avatar,
			SeatNo:   s.SeatNo,
			Chips:    s.Chips,
			Status:   game.PlayerActive,
		})
	}
	if len(botIDs) > 0 {
		RegisterBots(tableID, botIDs)
	}

	// Brief delay then start the next hand
	go func() {
		time.Sleep(1 * time.Second)
		if e := startNextHand(ctx, tableID, sessionID, totalHands+1); e != nil {
			g.Log().Errorf(ctx, "[recover] startNextHand error: %v", e)
			return
		}
		// Notify all connected clients with fresh state
		state := eng.GetState(tableID)
		if state != nil {
			ws.GlobalHub().BroadcastTable(tableID, ws.MsgTypeGameState, BuildGameStateForClient(state))
		}
	}()

	return nil
}
