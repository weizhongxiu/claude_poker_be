// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Tournaments is the golang structure of table tournaments for DAO operations like Where/Data.
type Tournaments struct {
	g.Meta         `orm:"table:tournaments, do:true"`
	Id             any         //
	ClubId         any         //
	Name           any         //
	Type           any         // 1=MTT 2=SNG 3=Bounty
	Buyin          any         //
	Fee            any         //
	StartingChips  any         //
	MaxPlayers     any         // 0=unlimited
	CurrentPlayers any         //
	PrizePool      any         //
	Status         any         // 1=registering 2=running 3=ended
	RegisterStart  *gtime.Time //
	RegisterEnd    *gtime.Time //
	StartedAt      *gtime.Time //
	EndedAt        *gtime.Time //
	CreatedAt      *gtime.Time //
}
