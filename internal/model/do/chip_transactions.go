// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ChipTransactions is the golang structure of table chip_transactions for DAO operations like Where/Data.
type ChipTransactions struct {
	g.Meta    `orm:"table:chip_transactions, do:true"`
	Id        any         //
	UserId    any         //
	Type      any         // 1=recharge 2=withdraw 3=win 4=lose 5=buyin 6=cashout
	Amount    any         // Change amount (positive=add negative=sub)
	Balance   any         // Balance after change
	RefId     any         // Related id (game_id etc)
	Remark    any         //
	CreatedAt *gtime.Time //
}
