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
	tableNo = fmt.Sprintf("%06d", time.Now().UnixMilli()%1000000)

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
		"duration":     table.Duration, // planned duration in hours (for expiry check)
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

	// Create session_players records for all seated players
	type seatRow struct {
		UserID int64 `orm:"user_id"`
		SeatNo int   `orm:"seat_no"`
		Chips  int64 `orm:"chips"`
	}
	var seatedPlayers []*seatRow
	if e2 := g.DB().Model("table_seats").Fields("user_id,seat_no,chips").
		Where("table_id", tableID).Where("status", 1).Scan(&seatedPlayers); e2 == nil {
		for _, sp := range seatedPlayers {
			// Get total buyin for this player
			buyinVal, _ := g.DB().Model("buyin_records").
				Fields("COALESCE(SUM(amount),0)").
				Where("user_id", sp.UserID).Where("status", 2).Value()
			_, _ = g.DB().Model("session_players").Data(g.Map{
				"session_id":  sessionID,
				"user_id":     sp.UserID,
				"seat_no":     sp.SeatNo,
				"total_buyin": buyinVal.Int64(),
				"chips_final": buyinVal.Int64(),
				"total_hands": 0,
				"result":      0,
			}).Insert()
		}
	}
	return
}

// TableSeatRow holds per-seat display info.
type TableSeatRow struct {
	SeatNo   int    `orm:"seat_no"   json:"seat_no"`
	UserID   int64  `orm:"user_id"   json:"user_id"`
	Nickname string `orm:"nickname"  json:"nickname"`
	Avatar   string `orm:"avatar"    json:"avatar"`
	Chips    int64  `orm:"chips"     json:"chips"`
}

// TableInfoResult bundles everything the table-info endpoint needs.
type TableInfoResult struct {
	Table         *entity.Tables
	Seats         []*TableSeatRow
	SessionID     int64
	StartedAt     string
	SessionStatus int // 0=none, 1=running, 2=ended
}

// TableInfo returns table config + active seats with user info.
func TableInfo(ctx context.Context, tableID int64) (*TableInfoResult, error) {
	table, err := GetTable(ctx, tableID)
	if err != nil {
		return nil, err
	}

	var seats []*TableSeatRow
	_ = g.DB().Model("table_seats ts").
		LeftJoin("users u", "u.id = ts.user_id").
		Fields("ts.seat_no, ts.user_id, u.nickname, u.avatar, ts.chips").
		Where("ts.table_id", tableID).Where("ts.status", 1).
		Scan(&seats)

	var sessionID int64
	var startedAt string
	var sessionStatus int
	type sessionRow struct {
		ID        int64  `orm:"id"`
		StartedAt string `orm:"started_at"`
		Status    int    `orm:"status"`
	}
	var sess sessionRow
	_ = g.DB().Model("room_sessions").Fields("id,started_at,status").
		Where("table_id", tableID).WhereIn("status", g.Slice{1, 2}).OrderDesc("id").Scan(&sess)
	if sess.ID > 0 {
		sessionID = sess.ID
		startedAt = sess.StartedAt
		sessionStatus = sess.Status
	}

	return &TableInfoResult{Table: table, Seats: seats, SessionID: sessionID, StartedAt: startedAt, SessionStatus: sessionStatus}, nil
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
