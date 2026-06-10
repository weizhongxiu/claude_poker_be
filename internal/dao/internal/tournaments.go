// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// TournamentsDao is the data access object for the table tournaments.
type TournamentsDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  TournamentsColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// TournamentsColumns defines and stores column names for the table tournaments.
type TournamentsColumns struct {
	Id             string //
	ClubId         string //
	Name           string //
	Type           string // 1=MTT 2=SNG 3=Bounty
	Buyin          string //
	Fee            string //
	StartingChips  string //
	MaxPlayers     string // 0=unlimited
	CurrentPlayers string //
	PrizePool      string //
	Status         string // 1=registering 2=running 3=ended
	RegisterStart  string //
	RegisterEnd    string //
	StartedAt      string //
	EndedAt        string //
	CreatedAt      string //
}

// tournamentsColumns holds the columns for the table tournaments.
var tournamentsColumns = TournamentsColumns{
	Id:             "id",
	ClubId:         "club_id",
	Name:           "name",
	Type:           "type",
	Buyin:          "buyin",
	Fee:            "fee",
	StartingChips:  "starting_chips",
	MaxPlayers:     "max_players",
	CurrentPlayers: "current_players",
	PrizePool:      "prize_pool",
	Status:         "status",
	RegisterStart:  "register_start",
	RegisterEnd:    "register_end",
	StartedAt:      "started_at",
	EndedAt:        "ended_at",
	CreatedAt:      "created_at",
}

// NewTournamentsDao creates and returns a new DAO object for table data access.
func NewTournamentsDao(handlers ...gdb.ModelHandler) *TournamentsDao {
	return &TournamentsDao{
		group:    "default",
		table:    "tournaments",
		columns:  tournamentsColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *TournamentsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *TournamentsDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *TournamentsDao) Columns() TournamentsColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *TournamentsDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *TournamentsDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *TournamentsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
