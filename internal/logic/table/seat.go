package table

import (
	"context"
	"fmt"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

const seatLockTTL = 10 // seconds

// TakeSeat seats a player at a table.
func TakeSeat(ctx context.Context, userID, tableID int64, seatNo int, buyIn int64) error {
	table, err := GetTable(ctx, tableID)
	if err != nil {
		return err
	}
	if buyIn < table.MinBuyin || buyIn > table.MaxBuyin {
		return gerror.Newf("买入金额须在 %d ~ %d 之间", table.MinBuyin, table.MaxBuyin)
	}

	lockKey := fmt.Sprintf("seat_lock:%d:%d", tableID, seatNo)
	lockVal := fmt.Sprintf("%d", userID)
	set, e := g.Redis().SetNX(ctx, lockKey, lockVal)
	if e != nil {
		return e
	}
	if !set {
		return gerror.New("该座位已被占用，请稍后重试")
	}
	_, _ = g.Redis().Expire(ctx, lockKey, seatLockTTL)

	count, e := g.DB().Model("table_seats").
		Where("table_id", tableID).Where("seat_no", seatNo).Where("status", 1).Count()
	if e != nil || count > 0 {
		_, _ = g.Redis().Del(ctx, lockKey)
		return gerror.New("该座位已有玩家")
	}

	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		type walletRow struct {
			ID      int64 `orm:"id"`
			Chips   int64 `orm:"chips"`
			Version int   `orm:"version"`
		}
		var w walletRow
		if e2 := tx.Model("user_wallets").Fields("id,chips,version").
			Where("user_id", userID).Scan(&w); e2 != nil {
			return e2
		}
		if w.ID == 0 {
			return gerror.New("钱包不存在")
		}
		if w.Chips < buyIn {
			return gerror.New("筹码不足")
		}

		if table.MaxBuyinTotal > 0 {
			cumVal, _ := tx.Model("buyin_records").Fields("COALESCE(SUM(amount),0)").
				Where("user_id", userID).Where("status", 2).Value()
			if cumVal.Int64()+buyIn > table.MaxBuyinTotal {
				return gerror.Newf("累计买入超过限制 %d", table.MaxBuyinTotal)
			}
		}

		res, e2 := tx.Model("user_wallets").
			Where("id", w.ID).Where("version", w.Version).Where("chips >= ?", buyIn).
			Data(g.Map{
				"chips":        gdb.Raw(fmt.Sprintf("chips - %d", buyIn)),
				"frozen_chips": gdb.Raw(fmt.Sprintf("frozen_chips + %d", buyIn)),
				"version":      gdb.Raw("version + 1"),
			}).Update()
		if e2 != nil {
			return e2
		}
		if rows, _ := res.RowsAffected(); rows == 0 {
			return gerror.New("筹码变更失败，请重试")
		}

		if _, e2 = tx.Model("table_seats").Data(g.Map{
			"table_id":  tableID,
			"user_id":   userID,
			"seat_no":   seatNo,
			"chips":     buyIn,
			"status":    1,
			"joined_at": time.Now(),
		}).Insert(); e2 != nil {
			return e2
		}

		if _, e2 = tx.Model("buyin_records").Data(g.Map{
			"session_id": 0,
			"user_id":    userID,
			"amount":     buyIn,
			"type":       1,
			"status":     2,
		}).Insert(); e2 != nil {
			return e2
		}

		_, _ = tx.Model("tables").Where("id", tableID).
			Data(g.Map{"current_players": gdb.Raw("current_players + 1")}).Update()
		return nil
	})

	if err != nil {
		_, _ = g.Redis().Del(ctx, lockKey)
	}
	return err
}

// LeaveSeat removes a player from their seat and refunds chips.
func LeaveSeat(ctx context.Context, userID, tableID int64) error {
	var seat struct {
		ID     int64 `orm:"id"`
		SeatNo int   `orm:"seat_no"`
		Chips  int64 `orm:"chips"`
	}
	if e := g.DB().Model("table_seats").Fields("id,seat_no,chips").
		Where("table_id", tableID).Where("user_id", userID).Where("status", 1).
		Scan(&seat); e != nil || seat.ID == 0 {
		return gerror.New("您未在该桌就座")
	}

	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if _, e := tx.Model("table_seats").Where("id", seat.ID).
			Data(g.Map{"status": 2}).Update(); e != nil {
			return e
		}
		if seat.Chips > 0 {
			if _, e := tx.Model("user_wallets").Where("user_id", userID).
				Data(g.Map{
					"chips":        gdb.Raw(fmt.Sprintf("chips + %d", seat.Chips)),
					"frozen_chips": gdb.Raw(fmt.Sprintf("GREATEST(frozen_chips - %d, 0)", seat.Chips)),
					"version":      gdb.Raw("version + 1"),
				}).Update(); e != nil {
				return e
			}
		}
		_, _ = tx.Model("tables").Where("id", tableID).
			Data(g.Map{"current_players": gdb.Raw("GREATEST(current_players - 1, 0)")}).Update()

		lockKey := fmt.Sprintf("seat_lock:%d:%d", tableID, seat.SeatNo)
		_, _ = g.Redis().Del(ctx, lockKey)
		return nil
	})
}
