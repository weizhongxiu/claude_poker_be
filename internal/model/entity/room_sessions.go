// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// RoomSessions is the golang structure for table room_sessions.
type RoomSessions struct {
	Id             uint64      `json:"id"             orm:"id"              description:""`                                           //
	TableId        uint64      `json:"tableId"        orm:"table_id"        description:""`                                           //
	SessionNo      string      `json:"sessionNo"      orm:"session_no"      description:"Session number (unique)"`                    // Session number (unique)
	CreatorId      uint64      `json:"creatorId"      orm:"creator_id"      description:""`                                           //
	GameType       int         `json:"gameType"       orm:"game_type"       description:"1=holdem 2=shortdeck 3=plo 4=sng 5=chinese"` // 1=holdem 2=shortdeck 3=plo 4=sng 5=chinese
	SmallBlind     int64       `json:"smallBlind"     orm:"small_blind"     description:""`                                           //
	BigBlind       int64       `json:"bigBlind"       orm:"big_blind"       description:""`                                           //
	TotalHands     int         `json:"totalHands"     orm:"total_hands"     description:"Total hands played"`                         // Total hands played
	TotalFlow      int64       `json:"totalFlow"      orm:"total_flow"      description:"Sum of all pots"`                            // Sum of all pots
	TotalBuyin     int64       `json:"totalBuyin"     orm:"total_buyin"     description:"Total buy-in of all players"`                // Total buy-in of all players
	MaxPot         int64       `json:"maxPot"         orm:"max_pot"         description:"Biggest pot in session"`                     // Biggest pot in session
	AvgPot         int64       `json:"avgPot"         orm:"avg_pot"         description:"Average pot size"`                           // Average pot size
	PlayerCount    int         `json:"playerCount"    orm:"player_count"    description:""`                                           //
	SpectatorCount int         `json:"spectatorCount" orm:"spectator_count" description:""`                                           //
	Duration       float64     `json:"duration"       orm:"duration"        description:"Actual duration in hours"`                   // Actual duration in hours
	Status         int         `json:"status"         orm:"status"          description:"1=running 2=ended"`                          // 1=running 2=ended
	EndReason      int         `json:"endReason"      orm:"end_reason"      description:"0=running 1=timeout 2=manual 3=all_left"`    // 0=running 1=timeout 2=manual 3=all_left
	StartedAt      *gtime.Time `json:"startedAt"      orm:"started_at"      description:""`                                           //
	EndedAt        *gtime.Time `json:"endedAt"        orm:"ended_at"        description:""`                                           //
}
