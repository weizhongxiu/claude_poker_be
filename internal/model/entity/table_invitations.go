// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// TableInvitations is the golang structure for table table_invitations.
type TableInvitations struct {
	Id        uint64      `json:"id"        orm:"id"         description:""`                                          //
	TableId   uint64      `json:"tableId"   orm:"table_id"   description:"Table id (before session)"`                 // Table id (before session)
	SessionId uint64      `json:"sessionId" orm:"session_id" description:"Session id (after start), 0=pre-session"`   // Session id (after start), 0=pre-session
	InviterId uint64      `json:"inviterId" orm:"inviter_id" description:""`                                          //
	InviteeId uint64      `json:"inviteeId" orm:"invitee_id" description:""`                                          //
	Status    int         `json:"status"    orm:"status"     description:"1=pending 2=accepted 3=rejected 4=expired"` // 1=pending 2=accepted 3=rejected 4=expired
	ExpiredAt *gtime.Time `json:"expiredAt" orm:"expired_at" description:""`                                          //
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:""`                                          //
}
