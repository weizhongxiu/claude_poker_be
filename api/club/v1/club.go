package v1

import "github.com/gogf/gf/v2/frame/g"

// --- Create Club ---

type CreateClubReq struct {
	g.Meta `path:"/club/create" method:"post" tags:"俱乐部" summary:"创建俱乐部"`
	Name   string `json:"name" v:"required|max-length:100#俱乐部名称必填|名称过长"`
	Logo   string `json:"logo"`
}

type CreateClubRes struct {
	ClubID int64  `json:"club_id"`
	ClubNo string `json:"club_no"`
}

// --- Join Club ---

type JoinClubReq struct {
	g.Meta `path:"/club/join" method:"post" tags:"俱乐部" summary:"加入俱乐部"`
	ClubNo string `json:"club_no" v:"required#俱乐部编号必填"`
}

type JoinClubRes struct {
	ClubID int64 `json:"club_id"`
}

// --- Member List ---

type MemberListReq struct {
	g.Meta   `path:"/club/{id}/members" method:"get" tags:"俱乐部" summary:"成员列表"`
	ID       int64 `json:"id" in:"path"`
	Page     int   `json:"page"      d:"1"`
	PageSize int   `json:"page_size" d:"20"`
}

type MemberListRes struct {
	List  []ClubMember `json:"list"`
	Total int          `json:"total"`
}

type ClubMember struct {
	UserID   int64  `json:"user_id"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Role     int    `json:"role"` // 1=owner 2=admin 3=member
	Chips    int64  `json:"chips"`
	JoinedAt string `json:"joined_at"`
}

// --- Club Tables ---

type ClubTablesReq struct {
	g.Meta   `path:"/club/{id}/tables" method:"get" tags:"俱乐部" summary:"俱乐部桌列表"`
	ID       int64 `json:"id" in:"path"`
	Page     int   `json:"page"      d:"1"`
	PageSize int   `json:"page_size" d:"20"`
}

type ClubTablesRes struct {
	List  []ClubTableItem `json:"list"`
	Total int             `json:"total"`
}

type ClubTableItem struct {
	TableID        int64  `json:"table_id"`
	TableNo        string `json:"table_no"`
	Name           string `json:"name"`
	GameType       int    `json:"game_type"`
	SmallBlind     int64  `json:"small_blind"`
	BigBlind       int64  `json:"big_blind"`
	CurrentPlayers int    `json:"current_players"`
	MaxSeats       int    `json:"max_seats"`
	Status         int    `json:"status"`
}

// --- Club Info ---

type ClubInfoReq struct {
	g.Meta `path:"/club/{id}" method:"get" tags:"俱乐部" summary:"俱乐部详情"`
	ID     int64 `json:"id" in:"path"`
}

type ClubInfoRes struct {
	ClubID      int64  `json:"club_id"`
	ClubNo      string `json:"club_no"`
	Name        string `json:"name"`
	Logo        string `json:"logo"`
	MemberCount int    `json:"member_count"`
	Announcement string `json:"announcement"`
	MyRole      int    `json:"my_role"` // 0=not member
	MyChips     int64  `json:"my_chips"`
}
