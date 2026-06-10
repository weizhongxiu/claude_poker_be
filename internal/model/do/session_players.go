// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SessionPlayers is the golang structure of table session_players for DAO operations like Where/Data.
type SessionPlayers struct {
	g.Meta      `orm:"table:session_players, do:true"`
	Id          any         //
	SessionId   any         //
	UserId      any         //
	SeatNo      any         //
	TotalHands  any         // Hands participated
	TotalBuyin  any         // Cumulative buy-in
	ChipsFinal  any         // Final chips when session ended
	Result      any         // Profit/loss = chips_final - total_buyin
	Vpip        any         // VPIP% in this session
	WinRate     any         // Win rate% in this session
	ActivityPts any         // Activity points earned
	IsMvp       any         //
	Rank        any         // Rank by profit desc
	JoinedAt    *gtime.Time //
	LeftAt      *gtime.Time //
}
