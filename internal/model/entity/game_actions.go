// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// GameActions is the golang structure for table game_actions.
type GameActions struct {
	Id        uint64      `json:"id"        orm:"id"         description:""`                                                                                //
	GameId    uint64      `json:"gameId"    orm:"game_id"    description:""`                                                                                //
	UserId    uint64      `json:"userId"    orm:"user_id"    description:""`                                                                                //
	SeatNo    int         `json:"seatNo"    orm:"seat_no"    description:""`                                                                                //
	Stage     int         `json:"stage"     orm:"stage"      description:"0=blinds 1=preflop 2=flop 3=turn 4=river"`                                        // 0=blinds 1=preflop 2=flop 3=turn 4=river
	Action    int         `json:"action"    orm:"action"     description:"1=fold 2=check 3=call 4=raise 5=allin 6=bet 7=blind_post 8=ante_post 9=straddle"` // 1=fold 2=check 3=call 4=raise 5=allin 6=bet 7=blind_post 8=ante_post 9=straddle
	Amount    int64       `json:"amount"    orm:"amount"     description:""`                                                                                //
	PotAfter  int64       `json:"potAfter"  orm:"pot_after"  description:"Pot size after this action"`                                                      // Pot size after this action
	ActionSeq int         `json:"actionSeq" orm:"action_seq" description:"Action sequence in this hand"`                                                    // Action sequence in this hand
	ActionMs  int         `json:"actionMs"  orm:"action_ms"  description:"Decision time in milliseconds"`                                                   // Decision time in milliseconds
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:""`                                                                                //
}
