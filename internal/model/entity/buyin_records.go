// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// BuyinRecords is the golang structure for table buyin_records.
type BuyinRecords struct {
	Id         uint64      `json:"id"         orm:"id"          description:""`                                //
	SessionId  uint64      `json:"sessionId"  orm:"session_id"  description:""`                                //
	UserId     uint64      `json:"userId"     orm:"user_id"     description:""`                                //
	Amount     int64       `json:"amount"     orm:"amount"      description:"Buy-in amount"`                   // Buy-in amount
	Type       int         `json:"type"       orm:"type"        description:"1=initial 2=rebuy 3=cashout"`     // 1=initial 2=rebuy 3=cashout
	Status     int         `json:"status"     orm:"status"      description:"1=pending 2=approved 3=rejected"` // 1=pending 2=approved 3=rejected
	ApprovedBy uint64      `json:"approvedBy" orm:"approved_by" description:"Admin user_id who approved"`      // Admin user_id who approved
	ApprovedAt *gtime.Time `json:"approvedAt" orm:"approved_at" description:""`                                //
	Remark     string      `json:"remark"     orm:"remark"      description:""`                                //
	CreatedAt  *gtime.Time `json:"createdAt"  orm:"created_at"  description:""`                                //
}
