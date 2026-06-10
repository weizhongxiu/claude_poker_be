// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// PotWinnerDetailsDao is the data access object for the table pot_winner_details.
type PotWinnerDetailsDao struct {
	table    string                  // table is the underlying table name of the DAO.
	group    string                  // group is the database configuration group name of the current DAO.
	columns  PotWinnerDetailsColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler      // handlers for customized model modification.
}

// PotWinnerDetailsColumns defines and stores column names for the table pot_winner_details.
type PotWinnerDetailsColumns struct {
	Id             string //
	DistributionId string // pot_distributions.id
	GameId         string // Redundant for query
	UserId         string //
	Amount         string // Exact amount this winner receives
	IsSplit        string // Part of a split pot
}

// potWinnerDetailsColumns holds the columns for the table pot_winner_details.
var potWinnerDetailsColumns = PotWinnerDetailsColumns{
	Id:             "id",
	DistributionId: "distribution_id",
	GameId:         "game_id",
	UserId:         "user_id",
	Amount:         "amount",
	IsSplit:        "is_split",
}

// NewPotWinnerDetailsDao creates and returns a new DAO object for table data access.
func NewPotWinnerDetailsDao(handlers ...gdb.ModelHandler) *PotWinnerDetailsDao {
	return &PotWinnerDetailsDao{
		group:    "default",
		table:    "pot_winner_details",
		columns:  potWinnerDetailsColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *PotWinnerDetailsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *PotWinnerDetailsDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *PotWinnerDetailsDao) Columns() PotWinnerDetailsColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *PotWinnerDetailsDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *PotWinnerDetailsDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *PotWinnerDetailsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
