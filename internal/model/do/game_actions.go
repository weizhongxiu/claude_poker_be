// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// GameActions is the golang structure of table game_actions for DAO operations like Where/Data.
type GameActions struct {
	g.Meta    `orm:"table:game_actions, do:true"`
	Id        any         //
	GameId    any         //
	UserId    any         //
	SeatNo    any         //
	Stage     any         // 0=blinds 1=preflop 2=flop 3=turn 4=river
	Action    any         // 1=fold 2=check 3=call 4=raise 5=allin 6=bet 7=blind_post 8=ante_post 9=straddle
	Amount    any         //
	PotAfter  any         // Pot size after this action
	ActionSeq any         // Action sequence in this hand
	ActionMs  any         // Decision time in milliseconds
	CreatedAt *gtime.Time //
}
