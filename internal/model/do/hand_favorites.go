// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// HandFavorites is the golang structure of table hand_favorites for DAO operations like Where/Data.
type HandFavorites struct {
	g.Meta    `orm:"table:hand_favorites, do:true"`
	Id        any         //
	UserId    any         //
	GameId    any         //
	CreatedAt *gtime.Time //
}
