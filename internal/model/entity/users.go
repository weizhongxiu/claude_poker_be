// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Users is the golang structure for table users.
type Users struct {
	Id          uint64      `json:"id"          orm:"id"            description:""`                          //
	Uid         string      `json:"uid"         orm:"uid"           description:"User unique number"`        // User unique number
	Nickname    string      `json:"nickname"    orm:"nickname"      description:""`                          //
	Avatar      string      `json:"avatar"      orm:"avatar"        description:""`                          //
	Phone       string      `json:"phone"       orm:"phone"         description:""`                          //
	Password    string      `json:"password"    orm:"password"      description:""`                          //
	Gender      int         `json:"gender"      orm:"gender"        description:"0=unknown 1=male 2=female"` // 0=unknown 1=male 2=female
	Status      int         `json:"status"      orm:"status"        description:"1=normal 2=banned"`         // 1=normal 2=banned
	LastLoginAt *gtime.Time `json:"lastLoginAt" orm:"last_login_at" description:""`                          //
	CreatedAt   *gtime.Time `json:"createdAt"   orm:"created_at"    description:""`                          //
	UpdatedAt   *gtime.Time `json:"updatedAt"   orm:"updated_at"    description:""`                          //
	DeletedAt   *gtime.Time `json:"deletedAt"   orm:"deleted_at"    description:""`                          //
}
