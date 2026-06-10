// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

// PotDistributions is the golang structure for table pot_distributions.
type PotDistributions struct {
	Id          uint64 `json:"id"          orm:"id"           description:""`                                                                         //
	GameId      uint64 `json:"gameId"      orm:"game_id"      description:""`                                                                         //
	PotType     int    `json:"potType"     orm:"pot_type"     description:"1=main 2=side"`                                                            // 1=main 2=side
	PotIndex    int    `json:"potIndex"    orm:"pot_index"    description:"Side pot index"`                                                           // Side pot index
	Amount      int64  `json:"amount"      orm:"amount"       description:"Total pot amount"`                                                         // Total pot amount
	WinnerIds   string `json:"winnerIds"   orm:"winner_ids"   description:"Winner user_ids (display only), see pot_winner_details for exact amounts"` // Winner user_ids (display only), see pot_winner_details for exact amounts
	WinnerCount int    `json:"winnerCount" orm:"winner_count" description:"Number of winners (>1 = split pot)"`                                       // Number of winners (>1 = split pot)
	WinReason   string `json:"winReason"   orm:"win_reason"   description:"Winning hand description (display)"`                                       // Winning hand description (display)
	WinRank     int    `json:"winRank"     orm:"win_rank"     description:"Winning hand strength value (1-10)"`                                       // Winning hand strength value (1-10)
}
