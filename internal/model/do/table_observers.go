// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// TableObservers is the golang structure of table table_observers for DAO operations like Where/Data.
type TableObservers struct {
	g.Meta    `orm:"table:table_observers, do:true"`
	Id        any         //
	SessionId any         //
	UserId    any         //
	Status    any         // 1=watching 2=left
	JoinedAt  *gtime.Time //
	LeftAt    *gtime.Time //
}
