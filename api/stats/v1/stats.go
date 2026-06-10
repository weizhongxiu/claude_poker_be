package v1

import "github.com/gogf/gf/v2/frame/g"

// --- Session list ---

type SessionListReq struct {
	g.Meta   `path:"/stats/sessions" method:"get" tags:"统计" summary:"历史牌局列表"`
	GameType int    `json:"game_type"` // 0=all
	DateFrom string `json:"date_from"` // yyyy-mm-dd
	DateTo   string `json:"date_to"`
	Page     int    `json:"page"      d:"1"`
	PageSize int    `json:"page_size" d:"20"`
}

type SessionListRes struct {
	List  []SessionItem `json:"list"`
	Total int           `json:"total"`
}

type SessionItem struct {
	SessionID   int64   `json:"session_id"`
	SessionNo   string  `json:"session_no"`
	TableID     int64   `json:"table_id"`
	GameType    int     `json:"game_type"`
	SmallBlind  int64   `json:"small_blind"`
	BigBlind    int64   `json:"big_blind"`
	TotalHands  int     `json:"total_hands"`
	TotalBuyin  int64   `json:"total_buyin"`
	Result      int64   `json:"result"`   // my profit/loss
	StartedAt   string  `json:"started_at"`
	EndedAt     string  `json:"ended_at"`
	Duration    float64 `json:"duration"`
}

// --- Session detail ---

type SessionDetailReq struct {
	g.Meta `path:"/stats/sessions/{id}" method:"get" tags:"统计" summary:"牌局结算详情"`
	ID     int64 `json:"id" in:"path"`
}

type SessionDetailRes struct {
	SessionID    int64          `json:"session_id"`
	SessionNo    string         `json:"session_no"`
	GameType     int            `json:"game_type"`
	SmallBlind   int64          `json:"small_blind"`
	BigBlind     int64          `json:"big_blind"`
	TotalHands   int            `json:"total_hands"`
	TotalFlow    int64          `json:"total_flow"`
	TotalBuyin   int64          `json:"total_buyin"`
	AvgPot       int64          `json:"avg_pot"`
	MaxPot       int64          `json:"max_pot"`
	Duration     float64        `json:"duration"`
	StartedAt    string         `json:"started_at"`
	EndedAt      string         `json:"ended_at"`
	Players      []SessionPlayer `json:"players"`
}

type SessionPlayer struct {
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
	IsMVP       bool    `json:"is_mvp"`
	Rank        int     `json:"rank"`
}

// --- Hand list ---

type HandListReq struct {
	g.Meta    `path:"/stats/hands" method:"get" tags:"统计" summary:"牌谱列表"`
	SessionID int64 `json:"session_id"` // 0=all
	Favorites int   `json:"favorites"`  // 1=only favorites
	Page      int   `json:"page"      d:"1"`
	PageSize  int   `json:"page_size" d:"20"`
}

type HandListRes struct {
	List  []HandItem `json:"list"`
	Total int        `json:"total"`
}

type HandItem struct {
	GameID         int64  `json:"game_id"`
	HandNo         string `json:"hand_no"`
	SessionID      int64  `json:"session_id"`
	SmallBlind     int64  `json:"small_blind"`
	BigBlind       int64  `json:"big_blind"`
	Pot            int64  `json:"pot"`
	CommunityCards string `json:"community_cards"`
	HoleCards      string `json:"hole_cards"`      // my hole cards
	HandRank       int    `json:"hand_rank"`
	HandRankDesc   string `json:"hand_rank_desc"`
	Result         int64  `json:"result"`
	IsFavorite     bool   `json:"is_favorite"`
	StartedAt      string `json:"started_at"`
}

// --- Hand replay ---

type HandReplayReq struct {
	g.Meta `path:"/stats/hands/{id}/replay" method:"get" tags:"统计" summary:"手牌回放"`
	ID     int64 `json:"id" in:"path"`
}

type HandReplayRes struct {
	GameID      int64          `json:"game_id"`
	HandNo      string         `json:"hand_no"`
	ShuffleSeed string         `json:"shuffle_seed"` // revealed after hand ends
	SmallBlind  int64          `json:"small_blind"`
	BigBlind    int64          `json:"big_blind"`
	DealerSeat  int            `json:"dealer_seat"`
	Stages      []ReplayStage  `json:"stages"`
}

type ReplayStage struct {
	Stage          int             `json:"stage"`
	CommunityCards string          `json:"community_cards"`
	Pot            int64           `json:"pot"`
	PlayersState   interface{}     `json:"players_state"`
	Actions        []ReplayAction  `json:"actions"`
}

type ReplayAction struct {
	Seq      int    `json:"seq"`
	SeatNo   int    `json:"seat_no"`
	UserID   int64  `json:"user_id"`
	Nickname string `json:"nickname"`
	Action   int    `json:"action"`
	Amount   int64  `json:"amount"`
	PotAfter int64  `json:"pot_after"`
}

// --- Favorite ---

type FavoriteReq struct {
	g.Meta `path:"/stats/hands/{id}/favorite" method:"post" tags:"统计" summary:"收藏/取消收藏手牌"`
	ID     int64 `json:"id" in:"path"`
}

type FavoriteRes struct {
	IsFavorite bool `json:"is_favorite"`
}

// --- Overview ---

type OverviewReq struct {
	g.Meta   `path:"/stats/overview" method:"get" tags:"统计" summary:"生涯总览"`
	GameType int    `json:"game_type"` // 0=all
	StatType int    `json:"stat_type" d:"1"` // 1=day 2=week 3=month 4=custom
	DateFrom string `json:"date_from"`
	DateTo   string `json:"date_to"`
}

type OverviewRes struct {
	TotalSessions int     `json:"total_sessions"`
	TotalHands    int     `json:"total_hands"`
	TotalProfit   int64   `json:"total_profit"`
	TotalBuyin    int64   `json:"total_buyin"`
	TotalFlow     int64   `json:"total_flow"`
	BiggestPot    int64   `json:"biggest_pot"`
	VPIP          float64 `json:"vpip"`
	PFR           float64 `json:"pfr"`
	WTSD          float64 `json:"wtsd"`
}
