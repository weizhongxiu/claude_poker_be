package table

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

// RebuyRequest creates a rebuy request. Auto-approved when buyin_approval=0.
func RebuyRequest(ctx context.Context, userID, sessionID int64, amount int64) (recordID int64, status int, err error) {
	var session entity.RoomSessions
	if e := g.DB().Model("room_sessions").Where("id", sessionID).Scan(&session); e != nil || session.Id == 0 {
		err = gerror.New("场次不存在")
		return
	}
	table, e := GetTable(ctx, int64(session.TableId))
	if e != nil {
		err = e
		return
	}
	if amount < table.MinBuyin || amount > table.MaxBuyin {
		err = gerror.Newf("补码金额须在 %d ~ %d 之间", table.MinBuyin, table.MaxBuyin)
		return
	}
	if table.MaxBuyinTotal > 0 {
		cumVal, _ := g.DB().Model("buyin_records").Fields("COALESCE(SUM(amount),0)").
			Where("session_id", sessionID).Where("user_id", userID).Where("status", 2).Value()
		if cumVal.Int64()+amount > table.MaxBuyinTotal {
			err = gerror.Newf("累计买入超过上限 %d", table.MaxBuyinTotal)
			return
		}
	}

	initStatus := 1
	if table.BuyinApproval == 0 {
		initStatus = 2
	}
	result, e := g.DB().Model("buyin_records").Data(g.Map{
		"session_id": sessionID,
		"user_id":    userID,
		"amount":     amount,
		"type":       2,
		"status":     initStatus,
	}).Insert()
	if e != nil {
		err = e
		return
	}
	recordID, err = result.LastInsertId()
	if err != nil {
		return
	}
	status = initStatus
	if initStatus == 2 {
		err = executeRebuy(ctx, userID, sessionID, amount)
	}
	return
}

// ApproveBuyin approves or rejects a pending buyin record.
func ApproveBuyin(ctx context.Context, adminUserID, recordID int64, approve bool) error {
	var record entity.BuyinRecords
	if e := g.DB().Model("buyin_records").Where("id", recordID).Scan(&record); e != nil || record.Id == 0 {
		return gerror.New("记录不存在")
	}
	if record.Status != 1 {
		return gerror.New("该记录已处理")
	}
	if !approve {
		_, e := g.DB().Model("buyin_records").Where("id", recordID).Data(g.Map{"status": 3}).Update()
		return e
	}
	_, e := g.DB().Model("buyin_records").Where("id", recordID).
		Data(g.Map{"status": 2, "approved_by": adminUserID, "approved_at": time.Now()}).Update()
	if e != nil {
		return e
	}
	return executeRebuy(ctx, int64(record.UserId), int64(record.SessionId), record.Amount)
}

func executeRebuy(ctx context.Context, userID, sessionID, amount int64) error {
	type seatRow struct {
		ID      int64 `orm:"id"`
		TableID int64 `orm:"table_id"`
		SeatNo  int   `orm:"seat_no"`
		Chips   int64 `orm:"chips"`
	}
	var seat seatRow
	if e := g.DB().Model("table_seats ts").
		LeftJoin("room_sessions rs", "rs.table_id = ts.table_id").
		Fields("ts.id, ts.table_id, ts.seat_no, ts.chips").
		Where("rs.id", sessionID).Where("ts.user_id", userID).Where("ts.status", 1).
		Scan(&seat); e != nil || seat.ID == 0 {
		return gerror.New("玩家未在座，无法补码")
	}

	err := g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		res, e := tx.Model("user_wallets").
			Where("user_id", userID).Where(fmt.Sprintf("chips >= %d", amount)).
			Data(g.Map{
				"chips":        gdb.Raw(fmt.Sprintf("chips - %d", amount)),
				"frozen_chips": gdb.Raw(fmt.Sprintf("frozen_chips + %d", amount)),
				"version":      gdb.Raw("version + 1"),
			}).Update()
		if e != nil {
			return e
		}
		if rows, _ := res.RowsAffected(); rows == 0 {
			return gerror.New("筹码不足，无法补码")
		}
		if _, e = tx.Model("table_seats").Where("id", seat.ID).
			Data(g.Map{"chips": gdb.Raw(fmt.Sprintf("chips + %d", amount))}).Update(); e != nil {
			return e
		}
		_, _ = tx.Model("session_players").
			Where("session_id", sessionID).Where("user_id", userID).
			Data(g.Map{"total_buyin": gdb.Raw(fmt.Sprintf("total_buyin + %d", amount))}).Update()
		// Also tag any unlinked buyin_records row (session_id=0) to this session
		_, _ = tx.Model("buyin_records").
			Where("user_id", userID).Where("session_id", 0).Where("status", 2).
			OrderDesc("id").Limit(1).
			Data(g.Map{"session_id": sessionID}).Update()
		return nil
	})
	if err != nil {
		return err
	}
	newChips := seat.Chips + amount
	// Sync new chip count into engine's in-memory state so the next hand sees it immediately
	game.GlobalEngine().UpdatePlayerChips(seat.TableID, userID, newChips)
	// Broadcast chip update so all clients refresh the chip display immediately
	ws.GlobalHub().BroadcastTable(seat.TableID, ws.MsgTypeChipUpdate, g.Map{
		"seat_no": seat.SeatNo,
		"chips":   newChips,
	})
	return nil
}

// GetTableRank returns real-time ranking for a table.
func GetTableRank(ctx context.Context, tableID int64) (interface{}, error) {
	var session entity.RoomSessions
	if e := g.DB().Model("room_sessions").
		Where("table_id", tableID).Where("status", 1).OrderDesc("id").
		Scan(&session); e != nil || session.Id == 0 {
		return g.Map{"session_id": 0, "players": []interface{}{}}, nil
	}

	type rankRow struct {
		UserID      int64   `json:"user_id"`
		Nickname    string  `json:"nickname"`
		Avatar      string  `json:"avatar"`
		TotalBuyin  int64   `json:"total_buyin"`
		Result      int64   `json:"result"`
		VPIP        float64 `json:"vpip"`
		WinRate     float64 `json:"win_rate"`
		ActivityPts int     `json:"activity_pts"`
		IsMVP       int     `json:"is_mvp"`
		Rank        int     `json:"rank"`
	}
	var players []*rankRow
	_ = g.DB().Model("session_players sp").
		LeftJoin("users u", "u.id = sp.user_id").
		Fields("sp.user_id, u.nickname, u.avatar, sp.total_buyin, sp.result, sp.vpip, sp.win_rate, sp.activity_pts, sp.is_mvp, sp.rank").
		Where("sp.session_id", session.Id).OrderAsc("sp.rank").
		Scan(&players)

	return g.Map{
		"session_id":  session.Id,
		"total_hands": session.TotalHands,
		"total_flow":  session.TotalFlow,
		"total_buyin": session.TotalBuyin,
		"avg_pot":     session.AvgPot,
		"players":     players,
	}, nil
}
