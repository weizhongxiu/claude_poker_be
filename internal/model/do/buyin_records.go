// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// BuyinRecords is the golang structure of table buyin_records for DAO operations like Where/Data.
type BuyinRecords struct {
	g.Meta     `orm:"table:buyin_records, do:true"`
	Id         any         //
	SessionId  any         //
	UserId     any         //
	Amount     any         // Buy-in amount
	Type       any         // 1=initial 2=rebuy 3=cashout
	Status     any         // 1=pending 2=approved 3=rejected
	ApprovedBy any         // Admin user_id who approved
	ApprovedAt *gtime.Time //
	Remark     any         //
	CreatedAt  *gtime.Time //
}
