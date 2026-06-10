// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// UserStats is the golang structure of table user_stats for DAO operations like Where/Data.
type UserStats struct {
	g.Meta        `orm:"table:user_stats, do:true"`
	Id            any         //
	UserId        any         //
	GameType      any         // 0=all 1=holdem 2=shortdeck 3=plo 4=sng 5=chinese
	StatType      any         // 1=day 2=week 3=month 4=custom
	StatDate      *gtime.Time // Stat start date, NULL=all-time
	StatEnd       *gtime.Time // Stat end date (stat_type=4 only)
	TotalGames    any         // Total hands
	TotalHands    any         // Total hands participated
	TotalSessions any         // Total sessions
	TotalWins     any         // Winning hands
	TotalProfit   any         // Total profit/loss
	TotalBuyin    any         // Total buy-in
	TotalFlow     any         // Total flow
	BiggestPot    any         // Biggest pot won
	Vpip          any         // VPIP%
	Pfr           any         // PFR%
	Wtsd          any         // Went to showdown%
	UpdatedAt     *gtime.Time //
}
