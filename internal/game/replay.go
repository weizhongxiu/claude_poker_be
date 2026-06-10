package game

// PlayerSnapshot is one player's state at a specific stage snapshot.
type PlayerSnapshot struct {
	SeatNo    int    `json:"seat_no"`
	UserID    int64  `json:"user_id"`
	Nickname  string `json:"nickname"`
	Avatar    string `json:"avatar"`
	Chips     int64  `json:"chips"`
	Bet       int64  `json:"bet"`
	Status    int    `json:"status"` // 1=active 2=folded 3=allin
	HoleCards string `json:"hole_cards,omitempty"` // Only visible after showdown
}

// StageSnapshot captures the full table state at one stage boundary.
type StageSnapshot struct {
	Stage          int              `json:"stage"`
	CommunityCards string           `json:"community_cards"` // Space-separated e.g. "Ah Kd Qc"
	Pot            int64            `json:"pot"`
	PlayersState   []PlayerSnapshot `json:"players_state"`
	ActionSeqStart int              `json:"action_seq_start"`
	ActionSeqEnd   int              `json:"action_seq_end"`
}

// BuildStageSnapshot creates a snapshot of the current game state for a given stage.
func BuildStageSnapshot(state *GameState, stage, actionSeqStart, actionSeqEnd int) StageSnapshot {
	players := make([]PlayerSnapshot, 0, len(state.Players))
	for _, p := range state.Players {
		ps := PlayerSnapshot{
			SeatNo:   p.SeatNo,
			UserID:   p.UserID,
			Nickname: p.Nickname,
			Avatar:   p.Avatar,
			Chips:    p.Chips,
			Bet:      p.Bet,
			Status:   p.Status,
		}
		// Only reveal hole cards at showdown
		if stage == StageShowdown && p.Status != PlayerFolded && len(p.HoleCards) > 0 {
			ps.HoleCards = CardsToStr(p.HoleCards)
		}
		players = append(players, ps)
	}

	return StageSnapshot{
		Stage:          stage,
		CommunityCards: CardsToStr(state.CommunityCards),
		Pot:            state.Pot,
		PlayersState:   players,
		ActionSeqStart: actionSeqStart,
		ActionSeqEnd:   actionSeqEnd,
	}
}
