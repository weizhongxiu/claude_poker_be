package stats

import (
	"context"
	"encoding/json"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	v1 "claude-test/api/stats/v1"
)

// ListHands returns the hand list for a user, optionally filtered by session or favorites.
func ListHands(ctx context.Context, userID int64, req *v1.HandListReq) (list []v1.HandItem, total int, err error) {
	base := g.DB().Model("game_players gp").
		LeftJoin("games g", "g.id = gp.game_id").
		Where("gp.user_id", userID).
		Where("g.status", 2)

	if req.SessionID > 0 {
		base = base.Where("g.session_id", req.SessionID)
	}

	// Favorites filter
	if req.Favorites == 1 {
		base = base.Where("EXISTS (SELECT 1 FROM hand_favorites hf WHERE hf.game_id=g.id AND hf.user_id=?)", userID)
	}

	total, err = base.Count()
	if err != nil {
		return
	}

	m := base.Fields(`
		g.id as game_id, g.hand_no, g.session_id,
		g.small_blind, g.big_blind, g.pot,
		g.community_cards, g.started_at,
		gp.hole_cards, gp.hand_rank, gp.hand_rank_desc, gp.result
	`)

	type row struct {
		GameID         uint64      `json:"game_id"`
		HandNo         string      `json:"hand_no"`
		SessionID      uint64      `json:"session_id"`
		SmallBlind     int64       `json:"small_blind"`
		BigBlind       int64       `json:"big_blind"`
		Pot            int64       `json:"pot"`
		CommunityCards string      `json:"community_cards"`
		StartedAt      *gtime.Time `json:"started_at"`
		HoleCards      string      `json:"hole_cards"`
		HandRank       int         `json:"hand_rank"`
		HandRankDesc   string      `json:"hand_rank_desc"`
		Result         int64       `json:"result"`
	}
	var rows []*row
	err = m.OrderDesc("g.id").Page(req.Page, req.PageSize).Scan(&rows)
	if err != nil {
		return
	}

	// Fetch favorites set for this user
	type favRow struct{ GameId uint64 }
	var favs []*favRow
	_ = g.DB().Model("hand_favorites").Fields("game_id").Where("user_id", userID).Scan(&favs)
	favSet := make(map[uint64]bool)
	for _, f := range favs {
		favSet[f.GameId] = true
	}

	for _, r := range rows {
		item := v1.HandItem{
			GameID:         int64(r.GameID),
			HandNo:         r.HandNo,
			SessionID:      int64(r.SessionID),
			SmallBlind:     r.SmallBlind,
			BigBlind:       r.BigBlind,
			Pot:            r.Pot,
			CommunityCards: r.CommunityCards,
			HoleCards:      r.HoleCards,
			HandRank:       r.HandRank,
			HandRankDesc:   r.HandRankDesc,
			Result:         r.Result,
			IsFavorite:     favSet[r.GameID],
		}
		if r.StartedAt != nil {
			item.StartedAt = r.StartedAt.Format("Y-m-d H:i:s")
		}
		list = append(list, item)
	}
	return
}

// GetHandReplay returns full replay data for one hand.
func GetHandReplay(ctx context.Context, userID, gameID int64) (*v1.HandReplayRes, error) {
	// Must be a participant or observer
	count, _ := g.DB().Model("game_players").
		Where("game_id", gameID).
		Where("user_id", userID).
		Count()
	if count == 0 {
		// Check if observer (was in same session)
		count2, _ := g.DB().Model("table_observers tob").
			LeftJoin("games g", "g.session_id = tob.session_id").
			Where("g.id", gameID).
			Where("tob.user_id", userID).
			Count()
		if count2 == 0 {
			return nil, gerror.New("无权查看该手牌回放")
		}
	}

	type gameRow struct {
		Id             uint64 `json:"id"`
		HandNo         string `json:"hand_no"`
		ShuffleSeed    string `json:"shuffle_seed"`
		SmallBlind     int64  `json:"small_blind"`
		BigBlind       int64  `json:"big_blind"`
		DealerSeat     int    `json:"dealer_seat"`
	}
	var g2 gameRow
	if e := g.DB().Model("games").Where("id", gameID).Scan(&g2); e != nil || g2.Id == 0 {
		return nil, gerror.New("手牌不存在")
	}

	// Load stage snapshots
	type replayRow struct {
		Stage          int    `json:"stage"`
		CommunityCards string `json:"community_cards"`
		Pot            int64  `json:"pot"`
		PlayersState   string `json:"players_state"`
		ActionSeqStart int    `json:"action_seq_start"`
		ActionSeqEnd   int    `json:"action_seq_end"`
	}
	var snapshots []*replayRow
	_ = g.DB().Model("hand_replays").
		Where("game_id", gameID).
		OrderAsc("stage").
		Scan(&snapshots)

	// Load all actions for this hand
	type actionRow struct {
		ActionSeq int    `json:"action_seq"`
		SeatNo    int    `json:"seat_no"`
		UserID    uint64 `json:"user_id"`
		Nickname  string `json:"nickname"`
		Action    int    `json:"action"`
		Amount    int64  `json:"amount"`
		PotAfter  int64  `json:"pot_after"`
		Stage     int    `json:"stage"`
	}
	var actions []*actionRow
	_ = g.DB().Model("game_actions ga").
		LeftJoin("users u", "u.id = ga.user_id").
		Fields("ga.action_seq, ga.seat_no, ga.user_id, u.nickname, ga.action, ga.amount, ga.pot_after, ga.stage").
		Where("ga.game_id", gameID).
		OrderAsc("ga.action_seq").
		Scan(&actions)

	// Group actions by stage
	actionsByStage := make(map[int][]v1.ReplayAction)
	for _, a := range actions {
		actionsByStage[a.Stage] = append(actionsByStage[a.Stage], v1.ReplayAction{
			Seq:      a.ActionSeq,
			SeatNo:   a.SeatNo,
			UserID:   int64(a.UserID),
			Nickname: a.Nickname,
			Action:   a.Action,
			Amount:   a.Amount,
			PotAfter: a.PotAfter,
		})
	}

	stages := make([]v1.ReplayStage, 0, len(snapshots))
	for _, snap := range snapshots {
		var playersState interface{}
		_ = json.Unmarshal([]byte(snap.PlayersState), &playersState)

		stages = append(stages, v1.ReplayStage{
			Stage:          snap.Stage,
			CommunityCards: snap.CommunityCards,
			Pot:            snap.Pot,
			PlayersState:   playersState,
			Actions:        actionsByStage[snap.Stage],
		})
	}

	return &v1.HandReplayRes{
		GameID:      gameID,
		HandNo:      g2.HandNo,
		ShuffleSeed: g2.ShuffleSeed,
		SmallBlind:  g2.SmallBlind,
		BigBlind:    g2.BigBlind,
		DealerSeat:  g2.DealerSeat,
		Stages:      stages,
	}, nil
}

// ToggleFavorite adds or removes a hand from favorites. Returns new state.
func ToggleFavorite(ctx context.Context, userID, gameID int64) (bool, error) {
	count, _ := g.DB().Model("hand_favorites").
		Where("user_id", userID).
		Where("game_id", gameID).
		Count()

	if count > 0 {
		_, err := g.DB().Model("hand_favorites").
			Where("user_id", userID).
			Where("game_id", gameID).
			Delete()
		return false, err
	}
	_, err := g.DB().Model("hand_favorites").Data(g.Map{
		"user_id": userID,
		"game_id": gameID,
	}).Insert()
	return true, err
}
