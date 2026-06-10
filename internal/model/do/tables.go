// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Tables is the golang structure of table tables for DAO operations like Where/Data.
type Tables struct {
	g.Meta            `orm:"table:tables, do:true"`
	Id                any         //
	TableNo           any         // Table number
	ClubId            any         // Club id, NULL=public
	Name              any         //
	HasPassword       any         // Password protected
	Password          any         //
	GameType          any         // 1=holdem 2=shortdeck 3=plo 4=sng 5=chinese
	BlindType         any         // 1=fixed 2=increasing
	SmallBlind        any         //
	BigBlind          any         //
	Ante              any         //
	StraddleEnabled   any         // Allow straddle (2x BB)
	MinBuyin          any         //
	MaxBuyin          any         //
	MaxBuyinTotal     any         // Cumulative max buy-in, 0=unlimited
	Duration          any         // Session duration in hours
	RunTwice          any         // Feature switch: allow run-twice
	LowWaterInsurance any         // Low water insurance
	CritGameplay      any         // Critical hit gameplay
	ActivityPoints    any         // Activity points enabled
	AutoRebuy         any         // Auto rebuy/cashout
	BuyinApproval     any         // Rebuy needs admin approval
	DelayShowCard     any         // Delay show card
	RandomSeat        any         // Random seat assignment
	SpectatorMute     any         // Mute spectators
	GpsIpRestrict     any         // GPS and IP restriction
	FullTableStart    any         // Start only when full
	MaxSeats          any         // Max seats 2-10 (standard 6-9)
	CurrentPlayers    any         //
	Tag               any         // Custom tag
	CreatorId         any         //
	Status            any         // 1=waiting 2=playing 3=closed
	CreatedAt         *gtime.Time //
	UpdatedAt         *gtime.Time //
	EndedAt           *gtime.Time //
}
