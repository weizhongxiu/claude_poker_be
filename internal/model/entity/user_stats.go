// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// UserStats is the golang structure for table user_stats.
type UserStats struct {
	Id            uint64      `json:"id"            orm:"id"             description:""`                                                 //
	UserId        uint64      `json:"userId"        orm:"user_id"        description:""`                                                 //
	GameType      int         `json:"gameType"      orm:"game_type"      description:"0=all 1=holdem 2=shortdeck 3=plo 4=sng 5=chinese"` // 0=all 1=holdem 2=shortdeck 3=plo 4=sng 5=chinese
	StatType      int         `json:"statType"      orm:"stat_type"      description:"1=day 2=week 3=month 4=custom"`                    // 1=day 2=week 3=month 4=custom
	StatDate      *gtime.Time `json:"statDate"      orm:"stat_date"      description:"Stat start date, NULL=all-time"`                   // Stat start date, NULL=all-time
	StatEnd       *gtime.Time `json:"statEnd"       orm:"stat_end"       description:"Stat end date (stat_type=4 only)"`                 // Stat end date (stat_type=4 only)
	TotalGames    int         `json:"totalGames"    orm:"total_games"    description:"Total hands"`                                      // Total hands
	TotalHands    int         `json:"totalHands"    orm:"total_hands"    description:"Total hands participated"`                         // Total hands participated
	TotalSessions int         `json:"totalSessions" orm:"total_sessions" description:"Total sessions"`                                   // Total sessions
	TotalWins     int         `json:"totalWins"     orm:"total_wins"     description:"Winning hands"`                                    // Winning hands
	TotalProfit   int64       `json:"totalProfit"   orm:"total_profit"   description:"Total profit/loss"`                                // Total profit/loss
	TotalBuyin    int64       `json:"totalBuyin"    orm:"total_buyin"    description:"Total buy-in"`                                     // Total buy-in
	TotalFlow     int64       `json:"totalFlow"     orm:"total_flow"     description:"Total flow"`                                       // Total flow
	BiggestPot    int64       `json:"biggestPot"    orm:"biggest_pot"    description:"Biggest pot won"`                                  // Biggest pot won
	Vpip          float64     `json:"vpip"          orm:"vpip"           description:"VPIP%"`                                            // VPIP%
	Pfr           float64     `json:"pfr"           orm:"pfr"            description:"PFR%"`                                             // PFR%
	Wtsd          float64     `json:"wtsd"          orm:"wtsd"           description:"Went to showdown%"`                                // Went to showdown%
	UpdatedAt     *gtime.Time `json:"updatedAt"     orm:"updated_at"     description:""`                                                 //
}
