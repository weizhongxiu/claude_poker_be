// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// PotDistributions is the golang structure of table pot_distributions for DAO operations like Where/Data.
type PotDistributions struct {
	g.Meta      `orm:"table:pot_distributions, do:true"`
	Id          any //
	GameId      any //
	PotType     any // 1=main 2=side
	PotIndex    any // Side pot index
	Amount      any // Total pot amount
	WinnerIds   any // Winner user_ids (display only), see pot_winner_details for exact amounts
	WinnerCount any // Number of winners (>1 = split pot)
	WinReason   any // Winning hand description (display)
	WinRank     any // Winning hand strength value (1-10)
}
