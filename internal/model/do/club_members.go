// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ClubMembers is the golang structure of table club_members for DAO operations like Where/Data.
type ClubMembers struct {
	g.Meta    `orm:"table:club_members, do:true"`
	Id        any         //
	ClubId    any         //
	UserId    any         //
	Role      any         // 1=owner 2=admin 3=member
	Chips     any         // Club chips
	Status    any         // 1=active 2=banned
	JoinedAt  *gtime.Time //
	UpdatedAt *gtime.Time //
}
