// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// ClubMembers is the golang structure for table club_members.
type ClubMembers struct {
	Id        uint64      `json:"id"        orm:"id"         description:""`                         //
	ClubId    uint64      `json:"clubId"    orm:"club_id"    description:""`                         //
	UserId    uint64      `json:"userId"    orm:"user_id"    description:""`                         //
	Role      int         `json:"role"      orm:"role"       description:"1=owner 2=admin 3=member"` // 1=owner 2=admin 3=member
	Chips     int64       `json:"chips"     orm:"chips"      description:"Club chips"`               // Club chips
	Status    int         `json:"status"    orm:"status"     description:"1=active 2=banned"`        // 1=active 2=banned
	JoinedAt  *gtime.Time `json:"joinedAt"  orm:"joined_at"  description:""`                         //
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at" description:""`                         //
}
