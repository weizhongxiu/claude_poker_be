// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

// HandReplays is the golang structure for table hand_replays.
type HandReplays struct {
	Id             uint64 `json:"id"             orm:"id"               description:""`                                                           //
	GameId         uint64 `json:"gameId"         orm:"game_id"          description:""`                                                           //
	HandIndex      int    `json:"handIndex"      orm:"hand_index"       description:"Hand index in session (for replay progress)"`                // Hand index in session (for replay progress)
	Stage          int    `json:"stage"          orm:"stage"            description:"1=preflop 2=flop 3=turn 4=river 5=showdown"`                 // 1=preflop 2=flop 3=turn 4=river 5=showdown
	CommunityCards string `json:"communityCards" orm:"community_cards"  description:"Board at this stage"`                                        // Board at this stage
	Pot            int64  `json:"pot"            orm:"pot"              description:""`                                                           //
	PlayersState   string `json:"playersState"   orm:"players_state"    description:"Player state snapshot [{seat,chips,bet,status,hole_cards}]"` // Player state snapshot [{seat,chips,bet,status,hole_cards}]
	ActionSeqStart int    `json:"actionSeqStart" orm:"action_seq_start" description:""`                                                           //
	ActionSeqEnd   int    `json:"actionSeqEnd"   orm:"action_seq_end"   description:""`                                                           //
}
