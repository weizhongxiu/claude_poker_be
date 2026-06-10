// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// TableObservers is the golang structure for table table_observers.
type TableObservers struct {
	Id        uint64      `json:"id"        orm:"id"         description:""`                  //
	SessionId uint64      `json:"sessionId" orm:"session_id" description:""`                  //
	UserId    uint64      `json:"userId"    orm:"user_id"    description:""`                  //
	Status    int         `json:"status"    orm:"status"     description:"1=watching 2=left"` // 1=watching 2=left
	JoinedAt  *gtime.Time `json:"joinedAt"  orm:"joined_at"  description:""`                  //
	LeftAt    *gtime.Time `json:"leftAt"    orm:"left_at"    description:""`                  //
}
