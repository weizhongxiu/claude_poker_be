// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SessionPlayers is the golang structure for table session_players.
type SessionPlayers struct {
	Id          uint64      `json:"id"          orm:"id"           description:""`                                        //
	SessionId   uint64      `json:"sessionId"   orm:"session_id"   description:""`                                        //
	UserId      uint64      `json:"userId"      orm:"user_id"      description:""`                                        //
	SeatNo      int         `json:"seatNo"      orm:"seat_no"      description:""`                                        //
	TotalHands  int         `json:"totalHands"  orm:"total_hands"  description:"Hands participated"`                      // Hands participated
	TotalBuyin  int64       `json:"totalBuyin"  orm:"total_buyin"  description:"Cumulative buy-in"`                       // Cumulative buy-in
	ChipsFinal  int64       `json:"chipsFinal"  orm:"chips_final"  description:"Final chips when session ended"`          // Final chips when session ended
	Result      int64       `json:"result"      orm:"result"       description:"Profit/loss = chips_final - total_buyin"` // Profit/loss = chips_final - total_buyin
	Vpip        float64     `json:"vpip"        orm:"vpip"         description:"VPIP% in this session"`                   // VPIP% in this session
	WinRate     float64     `json:"winRate"     orm:"win_rate"     description:"Win rate% in this session"`               // Win rate% in this session
	ActivityPts int         `json:"activityPts" orm:"activity_pts" description:"Activity points earned"`                  // Activity points earned
	IsMvp       int         `json:"isMvp"       orm:"is_mvp"       description:""`                                        //
	Rank        int         `json:"rank"        orm:"rank"         description:"Rank by profit desc"`                     // Rank by profit desc
	JoinedAt    *gtime.Time `json:"joinedAt"    orm:"joined_at"    description:""`                                        //
	LeftAt      *gtime.Time `json:"leftAt"      orm:"left_at"      description:""`                                        //
}
