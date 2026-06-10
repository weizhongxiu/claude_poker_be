package stats

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	v1 "claude-test/api/stats/v1"
)

// ListSessions returns paginated session history for a user.
func ListSessions(ctx context.Context, userID int64, req *v1.SessionListReq) (list []v1.SessionItem, total int, err error) {
	m := g.DB().Model("session_players sp").
		LeftJoin("room_sessions rs", "rs.id = sp.session_id").
		Fields(`
			sp.session_id, rs.session_no, rs.table_id, rs.game_type,
			rs.small_blind, rs.big_blind, rs.total_hands, rs.total_buyin,
			sp.result, sp.total_buyin as my_buyin,
			rs.started_at, rs.ended_at, rs.duration
		`).
		Where("sp.user_id", userID).
		Where("rs.status", 2)

	if req.GameType > 0 {
		m = m.Where("rs.game_type", req.GameType)
	}
	if req.DateFrom != "" {
		m = m.WhereGTE("rs.started_at", req.DateFrom)
	}
	if req.DateTo != "" {
		m = m.WhereLTE("rs.started_at", req.DateTo+" 23:59:59")
	}

	total, err = m.Count()
	if err != nil {
		return
	}

	type row struct {
		SessionID  uint64      `json:"session_id"`
		SessionNo  string      `json:"session_no"`
		TableID    uint64      `json:"table_id"`
		GameType   int         `json:"game_type"`
		SmallBlind int64       `json:"small_blind"`
		BigBlind   int64       `json:"big_blind"`
		TotalHands int         `json:"total_hands"`
		TotalBuyin int64       `json:"total_buyin"`
		Result     int64       `json:"result"`
		StartedAt  *gtime.Time `json:"started_at"`
		EndedAt    *gtime.Time `json:"ended_at"`
		Duration   float64     `json:"duration"`
	}

	var rows []*row
	err = m.OrderDesc("rs.started_at").Page(req.Page, req.PageSize).Scan(&rows)
	if err != nil {
		return
	}

	for _, r := range rows {
		item := v1.SessionItem{
			SessionID:  int64(r.SessionID),
			SessionNo:  r.SessionNo,
			TableID:    int64(r.TableID),
			GameType:   r.GameType,
			SmallBlind: r.SmallBlind,
			BigBlind:   r.BigBlind,
			TotalHands: r.TotalHands,
			TotalBuyin: r.TotalBuyin,
			Result:     r.Result,
			Duration:   r.Duration,
		}
		if r.StartedAt != nil {
			item.StartedAt = r.StartedAt.Format("Y-m-d H:i:s")
		}
		if r.EndedAt != nil {
			item.EndedAt = r.EndedAt.Format("Y-m-d H:i:s")
		}
		list = append(list, item)
	}
	return
}

// GetSessionDetail returns full settlement detail for one session.
func GetSessionDetail(ctx context.Context, userID, sessionID int64) (*v1.SessionDetailRes, error) {
	type sessionRow struct {
		Id         uint64      `json:"id"`
		SessionNo  string      `json:"session_no"`
		GameType   int         `json:"game_type"`
		SmallBlind int64       `json:"small_blind"`
		BigBlind   int64       `json:"big_blind"`
		TotalHands int         `json:"total_hands"`
		TotalFlow  int64       `json:"total_flow"`
		TotalBuyin int64       `json:"total_buyin"`
		AvgPot     int64       `json:"avg_pot"`
		MaxPot     int64       `json:"max_pot"`
		Duration   float64     `json:"duration"`
		StartedAt  *gtime.Time `json:"started_at"`
		EndedAt    *gtime.Time `json:"ended_at"`
	}
	var session sessionRow
	if e := g.DB().Model("room_sessions").Where("id", sessionID).Scan(&session); e != nil || session.Id == 0 {
		return nil, gerror.New("场次不存在")
	}

	// Only allow participants to view
	count, _ := g.DB().Model("session_players").
		Where("session_id", sessionID).
		Where("user_id", userID).
		Count()
	if count == 0 {
		return nil, gerror.New("无权查看该场次")
	}

	type playerRow struct {
		UserID      int64   `json:"user_id"`
		Nickname    string  `json:"nickname"`
		Avatar      string  `json:"avatar"`
		SeatNo      int     `json:"seat_no"`
		TotalBuyin  int64   `json:"total_buyin"`
		ChipsFinal  int64   `json:"chips_final"`
		Result      int64   `json:"result"`
		VPIP        float64 `json:"vpip"`
		WinRate     float64 `json:"win_rate"`
		TotalHands  int     `json:"total_hands"`
		ActivityPts int     `json:"activity_pts"`
		IsMVP       int     `json:"is_mvp"`
		Rank        int     `json:"rank"`
	}
	var players []*playerRow
	_ = g.DB().Model("session_players sp").
		LeftJoin("users u", "u.id = sp.user_id").
		Fields(`
			sp.user_id, u.nickname, u.avatar, sp.seat_no,
			sp.total_buyin, sp.chips_final, sp.result,
			sp.vpip, sp.win_rate, sp.total_hands,
			sp.activity_pts, sp.is_mvp, sp.rank
		`).
		Where("sp.session_id", sessionID).
		OrderAsc("sp.rank").
		Scan(&players)

	res := &v1.SessionDetailRes{
		SessionID:  sessionID,
		SessionNo:  session.SessionNo,
		GameType:   session.GameType,
		SmallBlind: session.SmallBlind,
		BigBlind:   session.BigBlind,
		TotalHands: session.TotalHands,
		TotalFlow:  session.TotalFlow,
		TotalBuyin: session.TotalBuyin,
		AvgPot:     session.AvgPot,
		MaxPot:     session.MaxPot,
		Duration:   session.Duration,
	}
	if session.StartedAt != nil {
		res.StartedAt = session.StartedAt.Format("Y-m-d H:i:s")
	}
	if session.EndedAt != nil {
		res.EndedAt = session.EndedAt.Format("Y-m-d H:i:s")
	}

	for _, p := range players {
		res.Players = append(res.Players, v1.SessionPlayer{
			UserID:      p.UserID,
			Nickname:    p.Nickname,
			Avatar:      p.Avatar,
			SeatNo:      p.SeatNo,
			TotalBuyin:  p.TotalBuyin,
			ChipsFinal:  p.ChipsFinal,
			Result:      p.Result,
			VPIP:        p.VPIP,
			WinRate:     p.WinRate,
			TotalHands:  p.TotalHands,
			ActivityPts: p.ActivityPts,
			IsMVP:       p.IsMVP == 1,
			Rank:        p.Rank,
		})
	}
	return res, nil
}
