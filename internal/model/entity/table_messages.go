// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// TableMessages is the golang structure for table table_messages.
type TableMessages struct {
	Id        uint64      `json:"id"        orm:"id"         description:""`                                 //
	SessionId uint64      `json:"sessionId" orm:"session_id" description:""`                                 //
	UserId    uint64      `json:"userId"    orm:"user_id"    description:"0=system"`                         // 0=system
	Type      int         `json:"type"      orm:"type"       description:"1=text 2=emoji 3=system 4=invite"` // 1=text 2=emoji 3=system 4=invite
	Content   string      `json:"content"   orm:"content"    description:""`                                 //
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:""`                                 //
}
