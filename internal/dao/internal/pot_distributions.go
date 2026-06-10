// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// PotDistributionsDao is the data access object for the table pot_distributions.
type PotDistributionsDao struct {
	table    string                  // table is the underlying table name of the DAO.
	group    string                  // group is the database configuration group name of the current DAO.
	columns  PotDistributionsColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler      // handlers for customized model modification.
}

// PotDistributionsColumns defines and stores column names for the table pot_distributions.
type PotDistributionsColumns struct {
	Id          string //
	GameId      string //
	PotType     string // 1=main 2=side
	PotIndex    string // Side pot index
	Amount      string // Total pot amount
	WinnerIds   string // Winner user_ids (display only), see pot_winner_details for exact amounts
	WinnerCount string // Number of winners (>1 = split pot)
	WinReason   string // Winning hand description (display)
	WinRank     string // Winning hand strength value (1-10)
}

// potDistributionsColumns holds the columns for the table pot_distributions.
var potDistributionsColumns = PotDistributionsColumns{
	Id:          "id",
	GameId:      "game_id",
	PotType:     "pot_type",
	PotIndex:    "pot_index",
	Amount:      "amount",
	WinnerIds:   "winner_ids",
	WinnerCount: "winner_count",
	WinReason:   "win_reason",
	WinRank:     "win_rank",
}

// NewPotDistributionsDao creates and returns a new DAO object for table data access.
func NewPotDistributionsDao(handlers ...gdb.ModelHandler) *PotDistributionsDao {
	return &PotDistributionsDao{
		group:    "default",
		table:    "pot_distributions",
		columns:  potDistributionsColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *PotDistributionsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *PotDistributionsDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *PotDistributionsDao) Columns() PotDistributionsColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *PotDistributionsDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *PotDistributionsDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *PotDistributionsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
