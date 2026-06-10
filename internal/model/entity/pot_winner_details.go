// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

// PotWinnerDetails is the golang structure for table pot_winner_details.
type PotWinnerDetails struct {
	Id             uint64 `json:"id"             orm:"id"              description:""`                                  //
	DistributionId uint64 `json:"distributionId" orm:"distribution_id" description:"pot_distributions.id"`              // pot_distributions.id
	GameId         uint64 `json:"gameId"         orm:"game_id"         description:"Redundant for query"`               // Redundant for query
	UserId         uint64 `json:"userId"         orm:"user_id"         description:""`                                  //
	Amount         int64  `json:"amount"         orm:"amount"          description:"Exact amount this winner receives"` // Exact amount this winner receives
	IsSplit        int    `json:"isSplit"        orm:"is_split"        description:"Part of a split pot"`               // Part of a split pot
}
