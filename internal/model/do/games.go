// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Games is the golang structure of table games for DAO operations like Where/Data.
type Games struct {
	g.Meta         `orm:"table:games, do:true"`
	Id             any         //
	TableId        any         //
	SessionId      any         // Parent session id
	HandNo         any         // Hand number (unique)
	ShuffleSeed    any         // Shuffle seed, published after hand ends
	DealerSeat     any         // Dealer button seat
	SmallBlind     any         //
	BigBlind       any         //
	Ante           any         //
	Pot            any         // Total pot
	CommunityCards any         // Board cards e.g. Ah Kd Qc Jh Ts
	IsSplitPot     any         // Split pot flag
	RunTwiceUsed   any         // Run-twice actually executed
	RunTwiceBoard2 any         // Second board for run-twice (Turn+River)
	Stage          any         // 0=blinds 1=preflop 2=flop 3=turn 4=river 5=showdown
	Status         any         // 1=running 2=ended
	StartedAt      *gtime.Time //
	EndedAt        *gtime.Time //
	DurationMs     any         // Hand duration in milliseconds
}
