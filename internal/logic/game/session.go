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
		OnDeal:        onDeal,
	}

	eng := game.GlobalEngine()
	eng.StartTable(cfg, cb)

	// Detect which seats are bots (phone prefix "bot_") and register them.
	var botIDs []int64
	for _, s := range seats {
		type phoneRow struct{ Phone string `orm:"phone"` }
		var pr phoneRow
		_ = g.DB().Model("users").Fields("phone").Where("id", s.UserID).Scan(&pr)
		if len(pr.Phone) > 4 && pr.Phone[:4] == "bot_" {
			botIDs = append(botIDs, s.UserID)
		}
		_ = eng.AddPlayer(cfg.TableID, game.PlayerState{
			UserID:   s.UserID,
			Nickname: s.Nickname,
			Avatar:   s.Avatar,
			SeatNo:   s.SeatNo,
			Chips:    s.Chips,
			Status:   game.PlayerActive,
		})
	}
	if len(botIDs) > 0 {
		RegisterBots(cfg.TableID, botIDs)
	}

	ws.GlobalHub().BroadcastTable(cfg.TableID, ws.MsgTypeSessionStarted, g.Map{
		"session_id": sessionID,
	})

	return startNextHand(ctx, cfg.TableID, sessionID, 1)
}

func startNextHand(ctx context.Context, tableID, sessionID int64, handIdx int) error {
	handNo := fmt.Sprintf("H%d%d", sessionID, time.Now().UnixMilli()%10000000)

	var session entity.RoomSessions
	_ = g.DB().Model("room_sessions").Where("id", sessionID).Scan(&session)

	// Check session duration (default 30 minutes if not configured).
	if session.StartedAt != nil {
		durationHours := session.Duration
		if durationHours <= 0 {
			durationHours = 0.5 // default 30 minutes
		}
		expireAt := session.StartedAt.Time.Add(time.Duration(float64(time.Hour) * durationHours))
		g.Log().Infof(ctx, "[duration] session=%d startedAt=%v expireAt=%v now=%v durationH=%.2f",
			sessionID, session.StartedAt.Time, expireAt, time.Now(), durationHours)
		if time.Now().After(expireAt) {
			g.Log().Infof(ctx, "[duration] session=%d EXPIRED, ending", sessionID)
			go func() {
				if e := EndSession(ctx, sessionID, 1); e != nil {
					g.Log().Errorf(ctx, "EndSession (time expired) error: %v", e)
				}
			}()
			return nil
		}
	}

	result, e := g.DB().Model("games").Data(g.Map{
		"table_id":    tableID,
		"session_id":  sessionID,
		"hand_no":     handNo,
		"small_blind": session.SmallBlind,
		"big_blind":   session.BigBlind,
		"dealer_seat": 0,
		"stage":       game.StageBlinds,
		"status":      1,
		"started_at":  time.Now(),
	}).Insert()
	if e != nil {
		return e
	}
	gameID, _ := result.LastInsertId()

	// Re-sync every seated player's chip count into the engine before each hand.
	// This catches any desync between engine memory and DB (e.g. new joiners,
	// rebuys that happened just after the last UpdatePlayerChips call).
	type seatRow struct {
		UserID   int64  `orm:"user_id"`
		SeatNo   int    `orm:"seat_no"`
		Chips    int64  `orm:"chips"`
		Nickname string `orm:"nickname"`
		Avatar   string `orm:"avatar"`
	}
	var activeSeats []*seatRow
	_ = g.DB().Model("table_seats ts").
		LeftJoin("users u", "u.id = ts.user_id").
		Fields("ts.user_id, ts.seat_no, ts.chips, u.nickname, u.avatar").
		Where("ts.table_id", tableID).Where("ts.status", 1).Scan(&activeSeats)

	eng := game.GlobalEngine()
	for _, s := range activeSeats {
		// Update chip count for existing engine players (corrects any desync).
		eng.UpdatePlayerChips(tableID, s.UserID, s.Chips)
		// Also ensure the player exists in the engine (catches late joiners).
		eng.EnsurePlayer(tableID, game.PlayerState{
			UserID:   s.UserID,
			Nickname: s.Nickname,
			Avatar:   s.Avatar,
			SeatNo:   s.SeatNo,
			Chips:    s.Chips,
			Status:   game.PlayerActive,
		})

		_, _ = g.DB().Model("game_players").Data(g.Map{
			"game_id":     gameID,
			"user_id":     s.UserID,
			"seat_no":     s.SeatNo,
			"chips_start": s.Chips,
		}).Insert()
	}

	return eng.StartHand(tableID, gameID, handNo, handIdx)
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
		ChipsFinal int64 `orm:"chips_final"` // updated after each hand
		SeatChips  int64 `orm:"seat_chips"`  // current table seat chips (may be 0 if seat cleared)
	}
	var players []*playerRow
	if e := g.DB().Model("session_players sp").
		LeftJoin("table_seats ts", fmt.Sprintf("ts.user_id = sp.user_id AND ts.table_id = %d", tableID)).
		Fields("sp.user_id, sp.total_buyin, sp.total_hands, sp.chips_final, COALESCE(ts.chips,0) AS seat_chips, ts.seat_no").
		Where("sp.session_id", sessionID).
		Scan(&players); e != nil {
		return e
	}

	for rank, p := range players {
		var chipsFinal int64
		if p.SeatChips > 0 {
			chipsFinal = p.SeatChips // live seat chips most accurate
		} else if p.TotalHands > 0 {
			chipsFinal = p.ChipsFinal // fallback: last hand ending chips
		} else {
			chipsFinal = p.TotalBuyin // 0 hands: return buyin in full
		}
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

func onDeal(state *game.GameState) {
	// Send each player their private hole cards as a string array
	for _, p := range state.Players {
		if len(p.HoleCards) == 0 {
			continue
		}
		cards := make([]string, len(p.HoleCards))
		for i, c := range p.HoleCards {
			cards[i] = c.String()
		}
		g.Log().Infof(nil, "[onDeal] sending hole_cards=%v to userID=%d seat=%d", cards, p.UserID, p.SeatNo)
		ws.GlobalHub().SendToUser(state.TableID, p.UserID, ws.MsgTypeDeal, g.Map{
			"hole_cards": cards,
		})
	}
}

func onAction(state *game.GameState, action game.PlayerAction) {
	// Persist to game_actions
	_, _ = g.DB().Model("game_actions").Data(g.Map{
		"game_id":    state.GameID,
		"user_id":    action.UserID,
		"seat_no":    action.SeatNo,
		"stage":      state.Stage,
		"action":     action.Action,
		"amount":     action.Amount,
		"pot_after":  state.Pot,
		"action_seq": state.ActionSeq,
	}).Insert()

	ws.GlobalHub().BroadcastTable(state.TableID, ws.MsgTypeActionResult, g.Map{
		"game_id":      state.GameID,
		"seat":         action.SeatNo,
		"action":       action.Action,
		"amount":       action.Amount,
		"pot":          state.Pot,
		"current_seat": state.CurrentSeat,
		"deadline":     state.ActionDeadline,
	})
	MaybeTriggerBot(state)
}

func onStageChange(state *game.GameState) {
	ws.GlobalHub().BroadcastTable(state.TableID, ws.MsgTypeGameState, buildGameStateMsg(state))
	MaybeTriggerBot(state)
}

func makeOnHandEnd(sessionID int64) func(result game.HandEndResult) {
	return func(result game.HandEndResult) {
		ctx := context.Background()
		if err := OnHandEnd(ctx, sessionID, result); err != nil {
			g.Log().Errorf(ctx, "OnHandEnd error: %v", err)
		}
	}
}

// BuildGameStateForClient returns a ws-ready map for the given game state.
func BuildGameStateForClient(state *game.GameState) g.Map {
	return buildGameStateMsg(state)
}

func buildGameStateMsg(state *game.GameState) g.Map {
	players := make([]g.Map, 0)
	for _, p := range state.Players {
		players = append(players, g.Map{
			"seat":     p.SeatNo,
			"nickname": p.Nickname,
			"avatar":   p.Avatar,
			"chips":    p.Chips,
			"bet":      p.Bet,
			"status":   p.Status,
		})
	}
	// community_cards as string array so frontend can iterate directly
	community := make([]string, len(state.CommunityCards))
	for i, c := range state.CommunityCards {
		community[i] = c.String()
	}
	return g.Map{
		"game_id":         state.GameID,
		"stage":           state.Stage,
		"community_cards": community,
		"pot":             state.Pot,
		"current_seat":    state.CurrentSeat,
		"deadline":        state.ActionDeadline,
		"hand_index":      state.HandIndex,
		"dealer_seat":     state.DealerSeat,
		"small_blind":     state.SmallBlind,
		"big_blind":       state.BigBlind,
		"players":         players,
	}
}
