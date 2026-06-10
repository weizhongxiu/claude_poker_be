package stats

import (
	"context"
	"fmt"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"

	v1 "claude-test/api/stats/v1"
)

// GetOverview returns career statistics for a user.
func GetOverview(ctx context.Context, userID int64, req *v1.OverviewReq) (*v1.OverviewRes, error) {
	// Try pre-aggregated for today
	if req.StatType == 1 && req.DateFrom == "" {
		today := time.Now().Format("2006-01-02")
		type statRow struct {
			TotalSessions int     `orm:"total_sessions"`
			TotalHands    int     `orm:"total_hands"`
			TotalProfit   int64   `orm:"total_profit"`
			TotalBuyin    int64   `orm:"total_buyin"`
			TotalFlow     int64   `orm:"total_flow"`
			BiggestPot    int64   `orm:"biggest_pot"`
			VPIP          float64 `orm:"vpip"`
			PFR           float64 `orm:"pfr"`
			WTSD          float64 `orm:"wtsd"`
		}
		var stat statRow
		err := g.DB().Model("user_stats").
			Where("user_id", userID).Where("game_type", req.GameType).
			Where("stat_type", 1).Where("stat_date", today).
			Scan(&stat)
		if err == nil && stat.TotalSessions > 0 {
			return &v1.OverviewRes{
				TotalSessions: stat.TotalSessions,
				TotalHands:    stat.TotalHands,
				TotalProfit:   stat.TotalProfit,
				TotalBuyin:    stat.TotalBuyin,
				TotalFlow:     stat.TotalFlow,
				BiggestPot:    stat.BiggestPot,
				VPIP:          stat.VPIP,
				PFR:           stat.PFR,
				WTSD:          stat.WTSD,
			}, nil
		}
	}

	dateFrom, dateTo := buildDateRange(req)

	type aggRow struct {
		TotalSessions int     `orm:"total_sessions"`
		TotalHands    int     `orm:"total_hands"`
		TotalProfit   int64   `orm:"total_profit"`
		TotalBuyin    int64   `orm:"total_buyin"`
		TotalFlow     int64   `orm:"total_flow"`
		AvgVPIP       float64 `orm:"avg_vpip"`
	}
	var agg aggRow

	m := g.DB().Model("session_players sp").
		LeftJoin("room_sessions rs", "rs.id = sp.session_id").
		Fields(`COUNT(DISTINCT rs.id) AS total_sessions,
			SUM(sp.total_hands) AS total_hands,
			SUM(sp.result)      AS total_profit,
			SUM(sp.total_buyin) AS total_buyin,
			SUM(rs.total_flow)  AS total_flow,
			AVG(sp.vpip)        AS avg_vpip`).
		Where("sp.user_id", userID).Where("rs.status", 2)

	if req.GameType > 0 {
		m = m.Where("rs.game_type", req.GameType)
	}
	if dateFrom != "" {
		m = m.WhereGTE("rs.started_at", dateFrom)
	}
	if dateTo != "" {
		m = m.WhereLTE("rs.started_at", dateTo)
	}
	if err := m.Scan(&agg); err != nil {
		return nil, err
	}

	// Biggest pot won by this user
	biggestPotVal, _ := g.DB().Model("pot_winner_details pwd").
		LeftJoin("games g", "g.id = pwd.game_id").
		LeftJoin("room_sessions rs", "rs.id = g.session_id").
		Fields("COALESCE(MAX(pwd.amount), 0)").
		Where("pwd.user_id", userID).Where("rs.status", 2).
		Value()
	biggestPot := biggestPotVal.Int64()

	// PFR
	type pfrRow struct {
		TotalPFR  int `orm:"total_pfr"`
		TotalHand int `orm:"total_hand"`
	}
	var pfr pfrRow
	_ = g.DB().Model("game_players gp").
		LeftJoin("games g", "g.id = gp.game_id").
		LeftJoin("room_sessions rs", "rs.id = g.session_id").
		Fields("SUM(gp.is_pfr) AS total_pfr, COUNT(*) AS total_hand").
		Where("gp.user_id", userID).Where("rs.status", 2).
		Scan(&pfr)
	pfrPct := 0.0
	if pfr.TotalHand > 0 {
		pfrPct = float64(pfr.TotalPFR) / float64(pfr.TotalHand) * 100
	}

	// WTSD
	type wtsdRow struct {
		WentToSD  int `orm:"went_to_sd"`
		TotalHand int `orm:"total_hand"`
	}
	var wtsd wtsdRow
	_ = g.DB().Model("game_players gp").
		LeftJoin("games g", "g.id = gp.game_id").
		LeftJoin("room_sessions rs", "rs.id = g.session_id").
		Fields("SUM(gp.went_to_sd) AS went_to_sd, COUNT(*) AS total_hand").
		Where("gp.user_id", userID).Where("rs.status", 2).
		Scan(&wtsd)
	wtsdPct := 0.0
	if wtsd.TotalHand > 0 {
		wtsdPct = float64(wtsd.WentToSD) / float64(wtsd.TotalHand) * 100
	}

	return &v1.OverviewRes{
		TotalSessions: agg.TotalSessions,
		TotalHands:    agg.TotalHands,
		TotalProfit:   agg.TotalProfit,
		TotalBuyin:    agg.TotalBuyin,
		TotalFlow:     agg.TotalFlow,
		BiggestPot:    biggestPot,
		VPIP:          agg.AvgVPIP,
		PFR:           pfrPct,
		WTSD:          wtsdPct,
	}, nil
}

func buildDateRange(req *v1.OverviewReq) (from, to string) {
	now := time.Now()
	switch req.StatType {
	case 1:
		from = now.Format("2006-01-02") + " 00:00:00"
		to = now.Format("2006-01-02") + " 23:59:59"
	case 2:
		weekday := int(now.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		monday := now.AddDate(0, 0, -(weekday - 1))
		from = monday.Format("2006-01-02") + " 00:00:00"
		to = now.Format("2006-01-02") + " 23:59:59"
	case 3:
		from = now.Format("2006-01") + "-01 00:00:00"
		to = now.Format("2006-01-02") + " 23:59:59"
	case 4:
		from = req.DateFrom + " 00:00:00"
		to = req.DateTo + " 23:59:59"
	}
	return
}

// UpsertDailyStats writes/updates the pre-aggregated user_stats row for today.
func UpsertDailyStats(ctx context.Context, userID int64, gameType int, profit, buyin, flow int64, vpip, pfr, wtsd float64) {
	today := time.Now().Format("2006-01-02")
	_, _ = g.DB().Model("user_stats").
		Data(g.Map{
			"user_id":        userID,
			"game_type":      gameType,
			"stat_type":      1,
			"stat_date":      today,
			"total_sessions": 1,
			"total_profit":   profit,
			"total_buyin":    buyin,
			"total_flow":     flow,
			"vpip":           vpip,
			"pfr":            pfr,
			"wtsd":           wtsd,
		}).
		OnDuplicate(g.Map{
			"total_sessions": gdb.Raw("total_sessions + 1"),
			"total_profit":   gdb.Raw(fmt.Sprintf("total_profit + %d", profit)),
			"total_buyin":    gdb.Raw(fmt.Sprintf("total_buyin + %d", buyin)),
			"total_flow":     gdb.Raw(fmt.Sprintf("total_flow + %d", flow)),
			"vpip":           vpip,
			"pfr":            pfr,
			"wtsd":           wtsd,
		}).
		Insert()
}
