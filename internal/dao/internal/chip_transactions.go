// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ChipTransactionsDao is the data access object for the table chip_transactions.
type ChipTransactionsDao struct {
	table    string                  // table is the underlying table name of the DAO.
	group    string                  // group is the database configuration group name of the current DAO.
	columns  ChipTransactionsColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler      // handlers for customized model modification.
}

// ChipTransactionsColumns defines and stores column names for the table chip_transactions.
type ChipTransactionsColumns struct {
	Id        string //
	UserId    string //
	Type      string // 1=recharge 2=withdraw 3=win 4=lose 5=buyin 6=cashout
	Amount    string // Change amount (positive=add negative=sub)
	Balance   string // Balance after change
	RefId     string // Related id (game_id etc)
	Remark    string //
	CreatedAt string //
}

// chipTransactionsColumns holds the columns for the table chip_transactions.
var chipTransactionsColumns = ChipTransactionsColumns{
	Id:        "id",
	UserId:    "user_id",
	Type:      "type",
	Amount:    "amount",
	Balance:   "balance",
	RefId:     "ref_id",
	Remark:    "remark",
	CreatedAt: "created_at",
}

// NewChipTransactionsDao creates and returns a new DAO object for table data access.
func NewChipTransactionsDao(handlers ...gdb.ModelHandler) *ChipTransactionsDao {
	return &ChipTransactionsDao{
		group:    "default",
		table:    "chip_transactions",
		columns:  chipTransactionsColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *ChipTransactionsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *ChipTransactionsDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *ChipTransactionsDao) Columns() ChipTransactionsColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *ChipTransactionsDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *ChipTransactionsDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *ChipTransactionsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
