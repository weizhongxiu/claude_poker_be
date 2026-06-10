package v1

import "github.com/gogf/gf/v2/frame/g"

// --- Create Table ---

type CreateTableReq struct {
	g.Meta `path:"/table/create" method:"post" tags:"牌桌" summary:"创建牌桌"`
	// Basic config
	Name         string  `json:"name"`
	HasPassword  int     `json:"has_password"`
	Password     string  `json:"password"`
	GameType     int     `json:"game_type"     v:"required|in:1,2,3,4,5#游戏类型必填|游戏类型不合法"` // 1=holdem
	BlindType    int     `json:"blind_type"    d:"1"`                                             // 1=fixed 2=increasing
	SmallBlind   int64   `json:"small_blind"   v:"required|min:1#小盲必填"`
	BigBlind     int64   `json:"big_blind"     v:"required|min:1#大盲必填"`
	Ante         int64   `json:"ante"`
	StraddleEnabled int  `json:"straddle_enabled"`
	MinBuyin     int64   `json:"min_buyin"     v:"required|min:1#最小买入必填"`
	MaxBuyin     int64   `json:"max_buyin"     v:"required|min:1#最大买入必填"`
	MaxBuyinTotal int64  `json:"max_buyin_total"` // 0 = unlimited
	Duration     float64 `json:"duration"      d:"2"` // hours
	MaxSeats     int     `json:"max_seats"     d:"9" v:"min:2|max:10"`
	Tag          string  `json:"tag"`
	ClubID       int64   `json:"club_id"` // 0 = public
	// Advanced settings
	RunTwice          int `json:"run_twice"`
	LowWaterInsurance int `json:"low_water_insurance"`
	CritGameplay      int `json:"crit_gameplay"`
	ActivityPoints    int `json:"activity_points"`
	AutoRebuy         int `json:"auto_rebuy"`
	BuyinApproval     int `json:"buyin_approval"`
	DelayShowCard     int `json:"delay_show_card"`
	RandomSeat        int `json:"random_seat"`
	SpectatorMute     int `json:"spectator_mute"`
	GpsIpRestrict     int `json:"gps_ip_restrict"`
	FullTableStart    int `json:"full_table_start"`
}

type CreateTableRes struct {
	TableID  int64  `json:"table_id"`
	TableNo  string `json:"table_no"`
}

// --- Join Table ---

type JoinTableReq struct {
	g.Meta   `path:"/table/join" method:"post" tags:"牌桌" summary:"加入牌桌(旁观)"`
	TableID  int64  `json:"table_id"  v:"required"`
	Password string `json:"password"`
}

type JoinTableRes struct {
	TableID int64 `json:"table_id"`
}

// --- Take Seat ---

type TakeSeatReq struct {
	g.Meta  `path:"/table/seat/take" method:"post" tags:"牌桌" summary:"入座"`
	TableID int64 `json:"table_id" v:"required"`
	SeatNo  int   `json:"seat_no"  v:"required|min:1|max:10"`
	BuyIn   int64 `json:"buyin"    v:"required|min:1#买入金额必填"`
}

type TakeSeatRes struct {
	SeatNo int   `json:"seat_no"`
	Chips  int64 `json:"chips"`
}

// --- Leave Seat ---

type LeaveSeatReq struct {
	g.Meta  `path:"/table/seat/leave" method:"post" tags:"牌桌" summary:"离座"`
	TableID int64 `json:"table_id" v:"required"`
}

type LeaveSeatRes struct{}

// --- Start Session ---

type StartSessionReq struct {
	g.Meta  `path:"/table/start" method:"post" tags:"牌桌" summary:"开局"`
	TableID int64 `json:"table_id" v:"required"`
}

type StartSessionRes struct {
	SessionID int64  `json:"session_id"`
	SessionNo string `json:"session_no"`
}

// --- Buy In (Rebuy) ---

type BuyInReq struct {
	g.Meta    `path:"/table/buyin" method:"post" tags:"牌桌" summary:"买入/补码"`
	SessionID int64 `json:"session_id" v:"required"`
	Amount    int64 `json:"amount"     v:"required|min:1"`
	Type      int   `json:"type"       d:"2"` // 1=initial 2=rebuy
}

type BuyInRes struct {
	RecordID int64 `json:"record_id"`
	Status   int   `json:"status"` // 1=pending 2=approved
}

// --- Rebuy Approve (admin) ---

type ApproveBuyInReq struct {
	g.Meta   `path:"/table/buyin/approve" method:"post" tags:"牌桌" summary:"审核补码"`
	RecordID int64 `json:"record_id" v:"required"`
	Approve  bool  `json:"approve"`
}

type ApproveBuyInRes struct{}

// --- Table Rank ---

type TableRankReq struct {
	g.Meta  `path:"/table/{id}/rank" method:"get" tags:"牌桌" summary:"实时排名"`
	ID      int64 `json:"id" in:"path"`
}

type TableRankRes struct {
	SessionID    int64         `json:"session_id"`
	TotalHands   int           `json:"total_hands"`
	TotalFlow    int64         `json:"total_flow"`
	TotalBuyin   int64         `json:"total_buyin"`
	AvgPot       int64         `json:"avg_pot"`
	Players      []RankPlayer  `json:"players"`
}

type RankPlayer struct {
	UserID      int64   `json:"user_id"`
	Nickname    string  `json:"nickname"`
	Avatar      string  `json:"avatar"`
	TotalBuyin  int64   `json:"total_buyin"`
	Result      int64   `json:"result"`
	VPIP        float64 `json:"vpip"`
	WinRate     float64 `json:"win_rate"`
	ActivityPts int     `json:"activity_pts"`
	IsMVP       bool    `json:"is_mvp"`
	Rank        int     `json:"rank"`
}

// --- Lobby ---

type LobbyTablesReq struct {
	g.Meta   `path:"/lobby/tables" method:"get" tags:"大厅" summary:"公开桌列表"`
	GameType int `json:"game_type"` // 0=all
	Page     int `json:"page"     d:"1"`
	PageSize int `json:"page_size" d:"20"`
}

type LobbyTablesRes struct {
	List  []TableInfo `json:"list"`
	Total int         `json:"total"`
}

type TableInfo struct {
	TableID        int64   `json:"table_id"`
	TableNo        string  `json:"table_no"`
	Name           string  `json:"name"`
	GameType       int     `json:"game_type"`
	SmallBlind     int64   `json:"small_blind"`
	BigBlind       int64   `json:"big_blind"`
	MinBuyin       int64   `json:"min_buyin"`
	MaxBuyin       int64   `json:"max_buyin"`
	HasPassword    bool    `json:"has_password"`
	CurrentPlayers int     `json:"current_players"`
	MaxSeats       int     `json:"max_seats"`
	Status         int     `json:"status"`
	CreatorID      int64   `json:"creator_id"`
}
