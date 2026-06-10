// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Clubs is the golang structure of table clubs for DAO operations like Where/Data.
type Clubs struct {
	g.Meta       `orm:"table:clubs, do:true"`
	Id           any         //
	ClubNo       any         // Club number
	Name         any         //
	Logo         any         //
	OwnerId      any         // Owner user_id
	Announcement any         // Announcement
	MemberCount  any         //
	MaxMembers   any         //
	Status       any         // 1=active 2=dissolved
	CreatedAt    *gtime.Time //
	UpdatedAt    *gtime.Time //
	DeletedAt    *gtime.Time //
}
