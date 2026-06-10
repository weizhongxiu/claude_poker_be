// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// HandFavorites is the golang structure for table hand_favorites.
type HandFavorites struct {
	Id        uint64      `json:"id"        orm:"id"         description:""` //
	UserId    uint64      `json:"userId"    orm:"user_id"    description:""` //
	GameId    uint64      `json:"gameId"    orm:"game_id"    description:""` //
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:""` //
}
