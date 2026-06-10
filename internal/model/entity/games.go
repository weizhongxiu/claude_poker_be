// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Games is the golang structure for table games.
type Games struct {
	Id             uint64      `json:"id"             orm:"id"               description:""`                                                    //
	TableId        uint64      `json:"tableId"        orm:"table_id"         description:""`                                                    //
	SessionId      uint64      `json:"sessionId"      orm:"session_id"       description:"Parent session id"`                                   // Parent session id
	HandNo         string      `json:"handNo"         orm:"hand_no"          description:"Hand number (unique)"`                                // Hand number (unique)
	ShuffleSeed    string      `json:"shuffleSeed"    orm:"shuffle_seed"     description:"Shuffle seed, published after hand ends"`             // Shuffle seed, published after hand ends
	DealerSeat     int         `json:"dealerSeat"     orm:"dealer_seat"      description:"Dealer button seat"`                                  // Dealer button seat
	SmallBlind     int64       `json:"smallBlind"     orm:"small_blind"      description:""`                                                    //
	BigBlind       int64       `json:"bigBlind"       orm:"big_blind"        description:""`                                                    //
	Ante           int64       `json:"ante"           orm:"ante"             description:""`                                                    //
	Pot            int64       `json:"pot"            orm:"pot"              description:"Total pot"`                                           // Total pot
	CommunityCards string      `json:"communityCards" orm:"community_cards"  description:"Board cards e.g. Ah Kd Qc Jh Ts"`                     // Board cards e.g. Ah Kd Qc Jh Ts
	IsSplitPot     int         `json:"isSplitPot"     orm:"is_split_pot"     description:"Split pot flag"`                                      // Split pot flag
	RunTwiceUsed   int         `json:"runTwiceUsed"   orm:"run_twice_used"   description:"Run-twice actually executed"`                         // Run-twice actually executed
	RunTwiceBoard2 string      `json:"runTwiceBoard2" orm:"run_twice_board2" description:"Second board for run-twice (Turn+River)"`             // Second board for run-twice (Turn+River)
	Stage          int         `json:"stage"          orm:"stage"            description:"0=blinds 1=preflop 2=flop 3=turn 4=river 5=showdown"` // 0=blinds 1=preflop 2=flop 3=turn 4=river 5=showdown
	Status         int         `json:"status"         orm:"status"           description:"1=running 2=ended"`                                   // 1=running 2=ended
	StartedAt      *gtime.Time `json:"startedAt"      orm:"started_at"       description:""`                                                    //
	EndedAt        *gtime.Time `json:"endedAt"        orm:"ended_at"         description:""`                                                    //
	DurationMs     int         `json:"durationMs"     orm:"duration_ms"      description:"Hand duration in milliseconds"`                       // Hand duration in milliseconds
}
