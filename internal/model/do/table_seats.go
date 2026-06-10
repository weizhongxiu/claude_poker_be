// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// TableSeats is the golang structure of table table_seats for DAO operations like Where/Data.
type TableSeats struct {
	g.Meta       `orm:"table:table_seats, do:true"`
	Id           any         //
	TableId      any         //
	UserId       any         //
	SeatNo       any         // Seat number 1-10
	Chips        any         // Chips on table
	Status       any         // 1=seated 2=left
	IsSittingOut any         // Sitting out flag
	JoinedAt     *gtime.Time //
	UpdatedAt    *gtime.Time //
}
