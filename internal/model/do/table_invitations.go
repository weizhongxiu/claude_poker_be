// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// TableInvitations is the golang structure of table table_invitations for DAO operations like Where/Data.
type TableInvitations struct {
	g.Meta    `orm:"table:table_invitations, do:true"`
	Id        any         //
	TableId   any         // Table id (before session)
	SessionId any         // Session id (after start), 0=pre-session
	InviterId any         //
	InviteeId any         //
	Status    any         // 1=pending 2=accepted 3=rejected 4=expired
	ExpiredAt *gtime.Time //
	CreatedAt *gtime.Time //
}
