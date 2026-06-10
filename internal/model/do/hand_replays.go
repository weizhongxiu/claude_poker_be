// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// HandReplays is the golang structure of table hand_replays for DAO operations like Where/Data.
type HandReplays struct {
	g.Meta         `orm:"table:hand_replays, do:true"`
	Id             any //
	GameId         any //
	HandIndex      any // Hand index in session (for replay progress)
	Stage          any // 1=preflop 2=flop 3=turn 4=river 5=showdown
	CommunityCards any // Board at this stage
	Pot            any //
	PlayersState   any // Player state snapshot [{seat,chips,bet,status,hole_cards}]
	ActionSeqStart any //
	ActionSeqEnd   any //
}
