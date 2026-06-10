// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// ChipTransactions is the golang structure for table chip_transactions.
type ChipTransactions struct {
	Id        uint64      `json:"id"        orm:"id"         description:""`                                                     //
	UserId    uint64      `json:"userId"    orm:"user_id"    description:""`                                                     //
	Type      int         `json:"type"      orm:"type"       description:"1=recharge 2=withdraw 3=win 4=lose 5=buyin 6=cashout"` // 1=recharge 2=withdraw 3=win 4=lose 5=buyin 6=cashout
	Amount    int64       `json:"amount"    orm:"amount"     description:"Change amount (positive=add negative=sub)"`            // Change amount (positive=add negative=sub)
	Balance   int64       `json:"balance"   orm:"balance"    description:"Balance after change"`                                 // Balance after change
	RefId     uint64      `json:"refId"     orm:"ref_id"     description:"Related id (game_id etc)"`                             // Related id (game_id etc)
	Remark    string      `json:"remark"    orm:"remark"     description:""`                                                     //
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:""`                                                     //
}
