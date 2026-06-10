// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Tournaments is the golang structure for table tournaments.
type Tournaments struct {
	Id             uint64      `json:"id"             orm:"id"              description:""`                                //
	ClubId         uint64      `json:"clubId"         orm:"club_id"         description:""`                                //
	Name           string      `json:"name"           orm:"name"            description:""`                                //
	Type           int         `json:"type"           orm:"type"            description:"1=MTT 2=SNG 3=Bounty"`            // 1=MTT 2=SNG 3=Bounty
	Buyin          int64       `json:"buyin"          orm:"buyin"           description:""`                                //
	Fee            int64       `json:"fee"            orm:"fee"             description:""`                                //
	StartingChips  int64       `json:"startingChips"  orm:"starting_chips"  description:""`                                //
	MaxPlayers     int         `json:"maxPlayers"     orm:"max_players"     description:"0=unlimited"`                     // 0=unlimited
	CurrentPlayers int         `json:"currentPlayers" orm:"current_players" description:""`                                //
	PrizePool      int64       `json:"prizePool"      orm:"prize_pool"      description:""`                                //
	Status         int         `json:"status"         orm:"status"          description:"1=registering 2=running 3=ended"` // 1=registering 2=running 3=ended
	RegisterStart  *gtime.Time `json:"registerStart"  orm:"register_start"  description:""`                                //
	RegisterEnd    *gtime.Time `json:"registerEnd"    orm:"register_end"    description:""`                                //
	StartedAt      *gtime.Time `json:"startedAt"      orm:"started_at"      description:""`                                //
	EndedAt        *gtime.Time `json:"endedAt"        orm:"ended_at"        description:""`                                //
	CreatedAt      *gtime.Time `json:"createdAt"      orm:"created_at"      description:""`                                //
}
