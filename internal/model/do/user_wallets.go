// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// UserWallets is the golang structure of table user_wallets for DAO operations like Where/Data.
type UserWallets struct {
	g.Meta      `orm:"table:user_wallets, do:true"`
	Id          any         //
	UserId      any         //
	Chips       any         // Chips balance
	Gold        any         // Gold balance
	Diamonds    any         // Diamond balance
	FrozenChips any         // Chips frozen on table
	Version     any         // Optimistic lock version
	CreatedAt   *gtime.Time //
	UpdatedAt   *gtime.Time //
}
