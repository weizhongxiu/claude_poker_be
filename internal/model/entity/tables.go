// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Tables is the golang structure for table tables.
type Tables struct {
	Id                uint64      `json:"id"                orm:"id"                  description:""`                                           //
	TableNo           string      `json:"tableNo"           orm:"table_no"            description:"Table number"`                               // Table number
	ClubId            uint64      `json:"clubId"            orm:"club_id"             description:"Club id, NULL=public"`                       // Club id, NULL=public
	Name              string      `json:"name"              orm:"name"                description:""`                                           //
	HasPassword       int         `json:"hasPassword"       orm:"has_password"        description:"Password protected"`                         // Password protected
	Password          string      `json:"password"          orm:"password"            description:""`                                           //
	GameType          int         `json:"gameType"          orm:"game_type"           description:"1=holdem 2=shortdeck 3=plo 4=sng 5=chinese"` // 1=holdem 2=shortdeck 3=plo 4=sng 5=chinese
	BlindType         int         `json:"blindType"         orm:"blind_type"          description:"1=fixed 2=increasing"`                       // 1=fixed 2=increasing
	SmallBlind        int64       `json:"smallBlind"        orm:"small_blind"         description:""`                                           //
	BigBlind          int64       `json:"bigBlind"          orm:"big_blind"           description:""`                                           //
	Ante              int64       `json:"ante"              orm:"ante"                description:""`                                           //
	StraddleEnabled   int         `json:"straddleEnabled"   orm:"straddle_enabled"    description:"Allow straddle (2x BB)"`                     // Allow straddle (2x BB)
	MinBuyin          int64       `json:"minBuyin"          orm:"min_buyin"           description:""`                                           //
	MaxBuyin          int64       `json:"maxBuyin"          orm:"max_buyin"           description:""`                                           //
	MaxBuyinTotal     int64       `json:"maxBuyinTotal"     orm:"max_buyin_total"     description:"Cumulative max buy-in, 0=unlimited"`         // Cumulative max buy-in, 0=unlimited
	Duration          float64     `json:"duration"          orm:"duration"            description:"Session duration in hours"`                  // Session duration in hours
	RunTwice          int         `json:"runTwice"          orm:"run_twice"           description:"Feature switch: allow run-twice"`            // Feature switch: allow run-twice
	LowWaterInsurance int         `json:"lowWaterInsurance" orm:"low_water_insurance" description:"Low water insurance"`                        // Low water insurance
	CritGameplay      int         `json:"critGameplay"      orm:"crit_gameplay"       description:"Critical hit gameplay"`                      // Critical hit gameplay
	ActivityPoints    int         `json:"activityPoints"    orm:"activity_points"     description:"Activity points enabled"`                    // Activity points enabled
	AutoRebuy         int         `json:"autoRebuy"         orm:"auto_rebuy"          description:"Auto rebuy/cashout"`                         // Auto rebuy/cashout
	BuyinApproval     int         `json:"buyinApproval"     orm:"buyin_approval"      description:"Rebuy needs admin approval"`                 // Rebuy needs admin approval
	DelayShowCard     int         `json:"delayShowCard"     orm:"delay_show_card"     description:"Delay show card"`                            // Delay show card
	RandomSeat        int         `json:"randomSeat"        orm:"random_seat"         description:"Random seat assignment"`                     // Random seat assignment
	SpectatorMute     int         `json:"spectatorMute"     orm:"spectator_mute"      description:"Mute spectators"`                            // Mute spectators
	GpsIpRestrict     int         `json:"gpsIpRestrict"     orm:"gps_ip_restrict"     description:"GPS and IP restriction"`                     // GPS and IP restriction
	FullTableStart    int         `json:"fullTableStart"    orm:"full_table_start"    description:"Start only when full"`                       // Start only when full
	MaxSeats          int         `json:"maxSeats"          orm:"max_seats"           description:"Max seats 2-10 (standard 6-9)"`              // Max seats 2-10 (standard 6-9)
	CurrentPlayers    int         `json:"currentPlayers"    orm:"current_players"     description:""`                                           //
	Tag               string      `json:"tag"               orm:"tag"                 description:"Custom tag"`                                 // Custom tag
	CreatorId         uint64      `json:"creatorId"         orm:"creator_id"          description:""`                                           //
	Status            int         `json:"status"            orm:"status"              description:"1=waiting 2=playing 3=closed"`               // 1=waiting 2=playing 3=closed
	CreatedAt         *gtime.Time `json:"createdAt"         orm:"created_at"          description:""`                                           //
	UpdatedAt         *gtime.Time `json:"updatedAt"         orm:"updated_at"          description:""`                                           //
	EndedAt           *gtime.Time `json:"endedAt"           orm:"ended_at"            description:""`                                           //
}
