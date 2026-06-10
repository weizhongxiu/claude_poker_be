// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// TableSeats is the golang structure for table table_seats.
type TableSeats struct {
	Id           uint64      `json:"id"           orm:"id"             description:""`                 //
	TableId      uint64      `json:"tableId"      orm:"table_id"       description:""`                 //
	UserId       uint64      `json:"userId"       orm:"user_id"        description:""`                 //
	SeatNo       int         `json:"seatNo"       orm:"seat_no"        description:"Seat number 1-10"` // Seat number 1-10
	Chips        int64       `json:"chips"        orm:"chips"          description:"Chips on table"`   // Chips on table
	Status       int         `json:"status"       orm:"status"         description:"1=seated 2=left"`  // 1=seated 2=left
	IsSittingOut int         `json:"isSittingOut" orm:"is_sitting_out" description:"Sitting out flag"` // Sitting out flag
	JoinedAt     *gtime.Time `json:"joinedAt"     orm:"joined_at"      description:""`                 //
	UpdatedAt    *gtime.Time `json:"updatedAt"    orm:"updated_at"     description:""`                 //
}
