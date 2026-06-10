// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Clubs is the golang structure for table clubs.
type Clubs struct {
	Id           uint64      `json:"id"           orm:"id"           description:""`                     //
	ClubNo       string      `json:"clubNo"       orm:"club_no"      description:"Club number"`          // Club number
	Name         string      `json:"name"         orm:"name"         description:""`                     //
	Logo         string      `json:"logo"         orm:"logo"         description:""`                     //
	OwnerId      uint64      `json:"ownerId"      orm:"owner_id"     description:"Owner user_id"`        // Owner user_id
	Announcement string      `json:"announcement" orm:"announcement" description:"Announcement"`         // Announcement
	MemberCount  int         `json:"memberCount"  orm:"member_count" description:""`                     //
	MaxMembers   int         `json:"maxMembers"   orm:"max_members"  description:""`                     //
	Status       int         `json:"status"       orm:"status"       description:"1=active 2=dissolved"` // 1=active 2=dissolved
	CreatedAt    *gtime.Time `json:"createdAt"    orm:"created_at"   description:""`                     //
	UpdatedAt    *gtime.Time `json:"updatedAt"    orm:"updated_at"   description:""`                     //
	DeletedAt    *gtime.Time `json:"deletedAt"    orm:"deleted_at"   description:""`                     //
}
