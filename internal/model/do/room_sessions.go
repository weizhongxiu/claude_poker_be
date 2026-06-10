// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// RoomSessions is the golang structure of table room_sessions for DAO operations like Where/Data.
type RoomSessions struct {
	g.Meta         `orm:"table:room_sessions, do:true"`
	Id             any         //
	TableId        any         //
	SessionNo      any         // Session number (unique)
	CreatorId      any         //
	GameType       any         // 1=holdem 2=shortdeck 3=plo 4=sng 5=chinese
	SmallBlind     any         //
	BigBlind       any         //
	TotalHands     any         // Total hands played
	TotalFlow      any         // Sum of all pots
	TotalBuyin     any         // Total buy-in of all players
	MaxPot         any         // Biggest pot in session
	AvgPot         any         // Average pot size
	PlayerCount    any         //
	SpectatorCount any         //
	Duration       any         // Actual duration in hours
	Status         any         // 1=running 2=ended
	EndReason      any         // 0=running 1=timeout 2=manual 3=all_left
	StartedAt      *gtime.Time //
	EndedAt        *gtime.Time //
}
