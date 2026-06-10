// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// UserWallets is the golang structure for table user_wallets.
type UserWallets struct {
	Id          uint64      `json:"id"          orm:"id"           description:""`                        //
	UserId      uint64      `json:"userId"      orm:"user_id"      description:""`                        //
	Chips       int64       `json:"chips"       orm:"chips"        description:"Chips balance"`           // Chips balance
	Gold        int64       `json:"gold"        orm:"gold"         description:"Gold balance"`            // Gold balance
	Diamonds    int         `json:"diamonds"    orm:"diamonds"     description:"Diamond balance"`         // Diamond balance
	FrozenChips int64       `json:"frozenChips" orm:"frozen_chips" description:"Chips frozen on table"`   // Chips frozen on table
	Version     int         `json:"version"     orm:"version"      description:"Optimistic lock version"` // Optimistic lock version
	CreatedAt   *gtime.Time `json:"createdAt"   orm:"created_at"   description:""`                        //
	UpdatedAt   *gtime.Time `json:"updatedAt"   orm:"updated_at"   description:""`                        //
}
