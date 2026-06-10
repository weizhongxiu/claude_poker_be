// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// RoomSessionsDao is the data access object for the table room_sessions.
type RoomSessionsDao struct {
	table    string              // table is the underlying table name of the DAO.
	group    string              // group is the database configuration group name of the current DAO.
	columns  RoomSessionsColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler  // handlers for customized model modification.
}

// RoomSessionsColumns defines and stores column names for the table room_sessions.
type RoomSessionsColumns struct {
	Id             string //
	TableId        string //
	SessionNo      string // Session number (unique)
	CreatorId      string //
	GameType       string // 1=holdem 2=shortdeck 3=plo 4=sng 5=chinese
	SmallBlind     string //
	BigBlind       string //
	TotalHands     string // Total hands played
	TotalFlow      string // Sum of all pots
	TotalBuyin     string // Total buy-in of all players
	MaxPot         string // Biggest pot in session
	AvgPot         string // Average pot size
	PlayerCount    string //
	SpectatorCount string //
	Duration       string // Actual duration in hours
	Status         string // 1=running 2=ended
	EndReason      string // 0=running 1=timeout 2=manual 3=all_left
	StartedAt      string //
	EndedAt        string //
}

// roomSessionsColumns holds the columns for the table room_sessions.
var roomSessionsColumns = RoomSessionsColumns{
	Id:             "id",
	TableId:        "table_id",
	SessionNo:      "session_no",
	CreatorId:      "creator_id",
	GameType:       "game_type",
	SmallBlind:     "small_blind",
	BigBlind:       "big_blind",
	TotalHands:     "total_hands",
	TotalFlow:      "total_flow",
	TotalBuyin:     "total_buyin",
	MaxPot:         "max_pot",
	AvgPot:         "avg_pot",
	PlayerCount:    "player_count",
	SpectatorCount: "spectator_count",
	Duration:       "duration",
	Status:         "status",
	EndReason:      "end_reason",
	StartedAt:      "started_at",
	EndedAt:        "ended_at",
}

// NewRoomSessionsDao creates and returns a new DAO object for table data access.
func NewRoomSessionsDao(handlers ...gdb.ModelHandler) *RoomSessionsDao {
	return &RoomSessionsDao{
		group:    "default",
		table:    "room_sessions",
		columns:  roomSessionsColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *RoomSessionsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *RoomSessionsDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *RoomSessionsDao) Columns() RoomSessionsColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *RoomSessionsDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *RoomSessionsDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, handler := range dao.handlers {
		model = handler(model)
	}
	return model.Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rolls back the transaction and returns the error if function f returns a non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note: Do not commit or roll back the transaction in function f,
// as it is automatically handled by this function.
func (dao *RoomSessionsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
