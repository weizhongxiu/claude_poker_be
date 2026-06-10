// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// TableMessages is the golang structure of table table_messages for DAO operations like Where/Data.
type TableMessages struct {
	g.Meta    `orm:"table:table_messages, do:true"`
	Id        any         //
	SessionId any         //
	UserId    any         // 0=system
	Type      any         // 1=text 2=emoji 3=system 4=invite
	Content   any         //
	CreatedAt *gtime.Time //
}
