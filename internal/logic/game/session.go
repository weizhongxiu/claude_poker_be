package game

import (
	"context"
	"fmt"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	"claude-test/internal/game"
	"claude-test/internal/model/entity"
	"claude-test/utility/ws"
)

// StartGameEngine starts the engine for a session and begins the first hand.
func StartGameEngine(ctx context.Context, sessionID int64) error {
	var session entity.RoomSessions
	if e := g.DB().Model("room_sessions").Where("id", sessionID).Scan(&session); e != nil || session.Id == 0 {
		return gerror.New("场次不存在")
	}

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
		Where("ts.table_id", session.TableId).Where("ts.status", 1).
		Scan(&seats); e != nil {
		return e
	}
	if len(seats) < 2 {
		return gerror.New("玩家不足，无法开局")
	}

	cfg := game.TableConfig{
		TableID:    int64(session.TableId),
		SessionID:  sessionID,
		SmallBlind: session.SmallBlind,
		BigBlind:   session.BigBlind,
		MaxSeats:   9,
	}

	cb := game.FSMCallbacks{
		OnAction:      onAction,
		OnStageChange: onStageChange,
		OnHandEnd:     makeOnHandEnd(sessionID),
	}

	eng := game.GlobalEngine()
	eng.StartTable(cfg, cb)

	for _, s := range seats {
		_ = eng.AddPlayer(cfg.TableID, game.PlayerState{
			UserID:   s.UserID,
			Nickname: s.Nickname,
			Avatar:   s.Avatar,
			SeatNo:   s.SeatNo,
			Chips:    s.Chips,
			Status:   game.PlayerActive,
		})
	}

	return startNextHand(ctx, cfg.TableID, sessionID, 1)
}

func startNextHand(ctx context.Context, tableID, sessionID int64, handIdx int) error {
	handNo := fmt.Sprintf("H%d%d", sessionID, time.Now().UnixMilli()%10000000)

	var session entity.RoomSessions
	_ = g.DB().Model("room_sessions").Where("id", sessionID).Scan(&session)

	result, e := g.DB().Model("games").Data(g.Map{
		"table_id":    tableID,
		"session_id":  sessionID,
		"hand_no":     handNo,
		"small_blind": session.SmallBlind,
		"big_blind":   session.BigBlind,
		"stage":       game.StageBlinds,
		"status":      1,
		"started_at":  time.Now(),
	}).Insert()
	if e != nil {
		return e
	}
	gameID, _ := result.LastInsertId()

	return game.GlobalEngine().StartHand(tableID, gameID, handNo, handIdx)
}

// EndSession ends a running session and settles all players.
func EndSession(ctx context.Context, sessionID int64, reason int) error {
	var session entity.RoomSessions
	if e := g.DB().Model("room_sessions").Where("id", sessionID).Scan(&session); e != nil || session.Id == 0 {
		return gerror.New("场次不存在")
	}
	if session.Status != 1 {
		return gerror.New("场次已结束")
	}

	now := time.Now()
	_, e := g.DB().Model("room_sessions").Where("id", sessionID).
		Data(g.Map{"status": 2, "ended_at": now, "end_reason": reason}).Update()
	if e != nil {
		return e
	}

	game.GlobalEngine().StopTable(int64(session.TableId))

	if err := settleSession(ctx, sessionID, int64(session.TableId)); err != nil {
		return err
	}

	ws.GlobalHub().BroadcastTable(int64(session.TableId), ws.MsgTypeSessionEnd, g.Map{
		"session_id": sessionID,
		"reason":     reason,
	})
	return nil
}

func settleSession(ctx context.Context, sessionID, tableID int64) error {
	type playerRow struct {
		UserID     int64 `orm:"user_id"`
		SeatNo     int   `orm:"seat_no"`
		TotalBuyin int64 `orm:"total_buyin"`
		TotalHands int   `orm:"total_hands"`
		Chips      int64 `orm:"chips"`
	}
	var players []*playerRow
	if e := g.DB().Model("session_players sp").
		LeftJoin("table_seats ts", fmt.Sprintf("ts.user_id = sp.user_id AND ts.table_id = %d", tableID)).
		Fields("sp.user_id, sp.total_buyin, sp.total_hands, ts.chips, ts.seat_no").
		Where("sp.session_id", sessionID).
		Scan(&players); e != nil {
		return e
	}

	for rank, p := range players {
		chipsFinal := p.Chips
		profit := chipsFinal - p.TotalBuyin

		_, _ = g.DB().Model("session_players").
			Where("session_id", sessionID).Where("user_id", p.UserID).
			Data(g.Map{"chips_final": chipsFinal, "result": profit, "rank": rank + 1}).Update()

		if chipsFinal > 0 {
			_, _ = g.DB().Model("user_wallets").Where("user_id", p.UserID).
				Data(g.Map{
					"chips":        gdb.Raw(fmt.Sprintf("chips + %d", chipsFinal)),
					"frozen_chips": gdb.Raw(fmt.Sprintf("GREATEST(frozen_chips - %d, 0)", p.TotalBuyin)),
					"version":      gdb.Raw("version + 1"),
				}).Update()
		}

		txType, amount := 3, profit
		if profit < 0 {
			txType, amount = 4, -profit
		}
		_, _ = g.DB().Model("chip_transactions").Data(g.Map{
			"user_id": p.UserID,
			"type":    txType,
			"amount":  amount,
			"ref_id":  sessionID,
			"remark":  fmt.Sprintf("session %d settle", sessionID),
		}).Insert()

		_, _ = g.DB().Model("table_seats").
			Where("user_id", p.UserID).Where("table_id", tableID).
			Data(g.Map{"status": 2}).Update()
	}

	_, _ = g.DB().Model("tables").Where("id", tableID).
		Data(g.Map{"status": 3, "ended_at": time.Now()}).Update()
	return nil
}

func onAction(state *game.GameState, action game.PlayerAction) {
	ws.GlobalHub().BroadcastTable(state.TableID, ws.MsgTypeActionResult, g.Map{
		"game_id": state.GameID,
		"seat":    action.SeatNo,
		"action":  action.Action,
		"amount":  action.Amount,
		"pot":     state.Pot,
	})
}

func onStageChange(state *game.GameState) {
	ws.GlobalHub().BroadcastTable(state.TableID, ws.MsgTypeGameState, buildGameStateMsg(state))
}

func makeOnHandEnd(sessionID int64) func(result game.HandEndResult) {
	return func(result game.HandEndResult) {
		ctx := context.Background()
		if err := OnHandEnd(ctx, sessionID, result); err != nil {
			g.Log().Errorf(ctx, "OnHandEnd error: %v", err)
		}
	}
}

func buildGameStateMsg(state *game.GameState) g.Map {
	players := make([]g.Map, 0)
	for _, p := range state.Players {
		players = append(players, g.Map{
			"seat":     p.SeatNo,
			"nickname": p.Nickname,
			"chips":    p.Chips,
			"bet":      p.Bet,
			"status":   p.Status,
		})
	}
	return g.Map{
		"game_id":         state.GameID,
		"stage":           state.Stage,
		"community_cards": game.CardsToStr(state.CommunityCards),
		"pot":             state.Pot,
		"current_seat":    state.CurrentSeat,
		"deadline":        state.ActionDeadline,
		"hand_index":      state.HandIndex,
		"players":         players,
	}
}
