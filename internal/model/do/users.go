// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Users is the golang structure of table users for DAO operations like Where/Data.
type Users struct {
	g.Meta      `orm:"table:users, do:true"`
	Id          any         //
	Uid         any         // User unique number
	Nickname    any         //
	Avatar      any         //
	Phone       any         //
	Password    any         //
	Gender      any         // 0=unknown 1=male 2=female
	Status      any         // 1=normal 2=banned
	LastLoginAt *gtime.Time //
	CreatedAt   *gtime.Time //
	UpdatedAt   *gtime.Time //
	DeletedAt   *gtime.Time //
}
