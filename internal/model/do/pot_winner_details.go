// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// PotWinnerDetails is the golang structure of table pot_winner_details for DAO operations like Where/Data.
type PotWinnerDetails struct {
	g.Meta         `orm:"table:pot_winner_details, do:true"`
	Id             any //
	DistributionId any // pot_distributions.id
	GameId         any // Redundant for query
	UserId         any //
	Amount         any // Exact amount this winner receives
	IsSplit        any // Part of a split pot
}
