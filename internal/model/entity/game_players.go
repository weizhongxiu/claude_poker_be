// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

// GamePlayers is the golang structure for table game_players.
type GamePlayers struct {
	Id           uint64 `json:"id"           orm:"id"             description:""`                                                                                                    //
	GameId       uint64 `json:"gameId"       orm:"game_id"        description:""`                                                                                                    //
	UserId       uint64 `json:"userId"       orm:"user_id"        description:""`                                                                                                    //
	SeatNo       int    `json:"seatNo"       orm:"seat_no"        description:""`                                                                                                    //
	Position     int    `json:"position"     orm:"position"       description:"0=BTN 1=SB 2=BB 3=UTG 4=UTG+1 5=MP 6=HJ 7=CO"`                                                        // 0=BTN 1=SB 2=BB 3=UTG 4=UTG+1 5=MP 6=HJ 7=CO
	HoleCards    string `json:"holeCards"    orm:"hole_cards"     description:"Hole cards e.g. Ah Kd"`                                                                               // Hole cards e.g. Ah Kd
	ForcedBet    int64  `json:"forcedBet"    orm:"forced_bet"     description:"Forced bet (blind+ante), separate from voluntary"`                                                    // Forced bet (blind+ante), separate from voluntary
	ChipsStart   int64  `json:"chipsStart"   orm:"chips_start"    description:"Chips at hand start"`                                                                                 // Chips at hand start
	ChipsEnd     int64  `json:"chipsEnd"     orm:"chips_end"      description:"Chips at hand end"`                                                                                   // Chips at hand end
	TotalBet     int64  `json:"totalBet"     orm:"total_bet"      description:"Total voluntary bet this hand"`                                                                       // Total voluntary bet this hand
	Result       int64  `json:"result"       orm:"result"         description:"Profit/loss this hand"`                                                                               // Profit/loss this hand
	BestHand     string `json:"bestHand"     orm:"best_hand"      description:"Best 5-card combo description"`                                                                       // Best 5-card combo description
	HandRank     int    `json:"handRank"     orm:"hand_rank"      description:"Hand strength 1(HighCard,weakest)-10(RoyalFlush,strongest). NOTE: opposite to rules.md display rank"` // Hand strength 1(HighCard,weakest)-10(RoyalFlush,strongest). NOTE: opposite to rules.md display rank
	HandRankDesc string `json:"handRankDesc" orm:"hand_rank_desc" description:"Hand description: High Card / One Pair / ... / Royal Flush"`                                          // Hand description: High Card / One Pair / ... / Royal Flush
	IsWinner     int    `json:"isWinner"     orm:"is_winner"      description:""`                                                                                                    //
	FoldStage    int    `json:"foldStage"    orm:"fold_stage"     description:"Stage when folded (0-4)"`                                                                             // Stage when folded (0-4)
	IsVpip       int    `json:"isVpip"       orm:"is_vpip"        description:"Voluntarily put chips in pot preflop"`                                                                // Voluntarily put chips in pot preflop
	IsPfr        int    `json:"isPfr"        orm:"is_pfr"         description:"Preflop raise flag"`                                                                                  // Preflop raise flag
	WentToSd     int    `json:"wentToSd"     orm:"went_to_sd"     description:"Went to showdown flag"`                                                                               // Went to showdown flag
	IsShowCard   int    `json:"isShowCard"   orm:"is_show_card"   description:"0=muck 1=show cards at showdown"`                                                                     // 0=muck 1=show cards at showdown
	Status       int    `json:"status"       orm:"status"         description:"1=active 2=folded 3=allin"`                                                                           // 1=active 2=folded 3=allin
}
