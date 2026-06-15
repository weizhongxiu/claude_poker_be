package table

import (
	"context"
	"fmt"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

var botNames = []string{
	"海豚AI", "小鲨AI", "深海AI", "浪花AI",
	"珊瑚AI", "章鱼AI", "海龟AI", "金枪AI",
}

// AddBots adds `count` bot players to the table, filling available seats.
// It creates bot DB users on demand and returns their user IDs.
func AddBots(ctx context.Context, tableID int64, count int) ([]int64, error) {
	if count < 1 || count > 8 {
		return nil, gerror.New("机器人数量 1-8")
	}

	table, err := GetTable(ctx, tableID)
	if err != nil {
		return nil, err
	}
	if table.Status == 3 {
		return nil, gerror.New("牌局已结束")
	}

	// Find occupied seat numbers
	type seatRow struct {
		SeatNo int `orm:"seat_no"`
	}
	var occupied []*seatRow
	_ = g.DB().Model("table_seats").Fields("seat_no").
		Where("table_id", tableID).Where("status", 1).Scan(&occupied)
	occupiedSet := make(map[int]bool)
	for _, s := range occupied {
		occupiedSet[s.SeatNo] = true
	}

	// Find free seats
	var freeSeats []int
	for i := 1; i <= int(table.MaxSeats); i++ {
		if !occupiedSet[i] {
			freeSeats = append(freeSeats, i)
		}
	}
	if len(freeSeats) < count {
		count = len(freeSeats)
	}
	if count == 0 {
		return nil, gerror.New("没有空余座位")
	}

	buyin := table.MinBuyin * 10 // bots buy in with 10× min, capped at MaxBuyin
	if table.MaxBuyin > 0 && buyin > table.MaxBuyin {
		buyin = table.MaxBuyin
	}

	var botIDs []int64
	for i := 0; i < count; i++ {
		seatNo := freeSeats[i]
		botName := botNames[i%len(botNames)]
		phone := fmt.Sprintf("bot_%s_%d", botName, tableID)

		// Fetch or create bot user
		type userRow struct {
			ID int64 `orm:"id"`
		}
		var u userRow
		_ = g.DB().Model("users").Fields("id").Where("phone", phone).Scan(&u)
		if u.ID == 0 {
			avatar := fmt.Sprintf("https://api.dicebear.com/7.x/bottts/svg?seed=%d_%d", tableID, i)
			uid := fmt.Sprintf("BOT%d%02d", tableID, i)
			res, e := g.DB().Model("users").Data(g.Map{
				"uid":        uid,
				"phone":      phone,
				"password":   "bot_no_login",
				"nickname":   botName,
				"avatar":     avatar,
				"created_at": time.Now(),
			}).Insert()
			if e != nil {
				continue
			}
			u.ID, _ = res.LastInsertId()
		}
		if u.ID == 0 {
			continue
		}

		// Ensure bot has a wallet with enough chips
		type walletRow struct{ ID int64 `orm:"id"` }
		var wl walletRow
		_ = g.DB().Model("user_wallets").Fields("id").Where("user_id", u.ID).Scan(&wl)
		if wl.ID == 0 {
			_, _ = g.DB().Model("user_wallets").Data(g.Map{
				"user_id": u.ID,
				"chips":   buyin * 100,
			}).Insert()
		} else {
			// Make sure balance covers buyin
			_, _ = g.DB().Model("user_wallets").Where("id", wl.ID).
				Data(g.Map{"chips": buyin * 100}).Update()
		}

		// Insert or update seat (reuse upsert logic)
		err2 := TakeSeat(ctx, u.ID, tableID, seatNo, buyin)
		if err2 != nil {
			continue
		}
		botIDs = append(botIDs, u.ID)
	}

	if len(botIDs) == 0 {
		return nil, gerror.New("添加机器人失败")
	}
	return botIDs, nil
}
