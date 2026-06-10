package table

import (
	"context"
	"fmt"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	v1 "claude-test/api/table/v1"
	"claude-test/internal/model/entity"
)

// CreateTable creates a new table and returns its ID and table number.
func CreateTable(ctx context.Context, creatorID int64, req *v1.CreateTableReq) (tableID int64, tableNo string, err error) {
	// Generate unique table number
	tableNo = fmt.Sprintf("T%d", time.Now().UnixMilli()%1000000000)

	data := g.Map{
		"table_no":            tableNo,
		"club_id":             req.ClubID,
		"name":                req.Name,
		"has_password":        req.HasPassword,
		"password":            req.Password,
		"game_type":           req.GameType,
		"blind_type":          req.BlindType,
		"small_blind":         req.SmallBlind,
		"big_blind":           req.BigBlind,
		"ante":                req.Ante,
		"straddle_enabled":    req.StraddleEnabled,
		"min_buyin":           req.MinBuyin,
		"max_buyin":           req.MaxBuyin,
		"max_buyin_total":     req.MaxBuyinTotal,
		"duration":            req.Duration,
		"max_seats":           req.MaxSeats,
		"tag":                 req.Tag,
		"creator_id":          creatorID,
		"run_twice":           req.RunTwice,
		"low_water_insurance": req.LowWaterInsurance,
		"crit_gameplay":       req.CritGameplay,
		"activity_points":     req.ActivityPoints,
		"auto_rebuy":          req.AutoRebuy,
		"buyin_approval":      req.BuyinApproval,
		"delay_show_card":     req.DelayShowCard,
		"random_seat":         req.RandomSeat,
		"spectator_mute":      req.SpectatorMute,
		"gps_ip_restrict":     req.GpsIpRestrict,
		"full_table_start":    req.FullTableStart,
		"current_players":     0,
		"status":              1,
	}

	result, e := g.DB().Model("tables").Data(data).Insert()
	if e != nil {
		err = e
		return
	}
	tableID, err = result.LastInsertId()
	return
}

// GetTable returns a table by ID.
func GetTable(ctx context.Context, tableID int64) (*entity.Tables, error) {
	t := &entity.Tables{}
	err := g.DB().Model("tables").Where("id", tableID).Scan(t)
	if err != nil {
		return nil, err
	}
	if t.Id == 0 {
		return nil, gerror.New("牌桌不存在")
	}
	return t, nil
}

// StartSession creates a room_session and starts the game engine for a table.
func StartSession(ctx context.Context, tableID, creatorID int64) (sessionID int64, sessionNo string, err error) {
	table, err := GetTable(ctx, tableID)
	if err != nil {
		return
	}
	if table.Status != 1 {
		err = gerror.New("牌桌状态不允许开局")
		return
	}

	// Check player count
	count, e := g.DB().Model("table_seats").
		Where("table_id", tableID).
		Where("status", 1).
		Count()
	if e != nil {
		err = e
		return
	}
	if count < 2 {
		err = gerror.New("至少需要2名玩家才能开局")
		return
	}

	sessionNo = fmt.Sprintf("S%d%d", tableID, time.Now().UnixMilli()%1000000000)

	result, e := g.DB().Model("room_sessions").Data(g.Map{
		"table_id":     tableID,
		"session_no":   sessionNo,
		"creator_id":   creatorID,
		"game_type":    table.GameType,
		"small_blind":  table.SmallBlind,
		"big_blind":    table.BigBlind,
		"player_count": count,
		"status":       1,
		"started_at":   time.Now(),
	}).Insert()
	if e != nil {
		err = e
		return
	}
	sessionID, err = result.LastInsertId()
	if err != nil {
		return
	}

	// Update table status
	_, _ = g.DB().Model("tables").Where("id", tableID).Update(g.Map{"status": 2})
	return
}

// LobbyTables returns public table list.
func LobbyTables(ctx context.Context, gameType, page, pageSize int) (list []*entity.Tables, total int, err error) {
	m := g.DB().Model("tables").Where("status", g.Slice{1, 2})
	if gameType > 0 {
		m = m.Where("game_type", gameType)
	}
	m = m.Where("club_id", 0) // public tables only
	total, err = m.Count()
	if err != nil {
		return
	}
	err = m.Page(page, pageSize).OrderDesc("id").Scan(&list)
	return
}
