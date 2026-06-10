// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// BuyinRecordsDao is the data access object for the table buyin_records.
type BuyinRecordsDao struct {
	table    string              // table is the underlying table name of the DAO.
	group    string              // group is the database configuration group name of the current DAO.
	columns  BuyinRecordsColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler  // handlers for customized model modification.
}

// BuyinRecordsColumns defines and stores column names for the table buyin_records.
type BuyinRecordsColumns struct {
	Id         string //
	SessionId  string //
	UserId     string //
	Amount     string // Buy-in amount
	Type       string // 1=initial 2=rebuy 3=cashout
	Status     string // 1=pending 2=approved 3=rejected
	ApprovedBy string // Admin user_id who approved
	ApprovedAt string //
	Remark     string //
	CreatedAt  string //
}

// buyinRecordsColumns holds the columns for the table buyin_records.
var buyinRecordsColumns = BuyinRecordsColumns{
	Id:         "id",
	SessionId:  "session_id",
	UserId:     "user_id",
	Amount:     "amount",
	Type:       "type",
	Status:     "status",
	ApprovedBy: "approved_by",
	ApprovedAt: "approved_at",
	Remark:     "remark",
	CreatedAt:  "created_at",
}

// NewBuyinRecordsDao creates and returns a new DAO object for table data access.
func NewBuyinRecordsDao(handlers ...gdb.ModelHandler) *BuyinRecordsDao {
	return &BuyinRecordsDao{
		group:    "default",
		table:    "buyin_records",
		columns:  buyinRecordsColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *BuyinRecordsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *BuyinRecordsDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *BuyinRecordsDao) Columns() BuyinRecordsColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *BuyinRecordsDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *BuyinRecordsDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *BuyinRecordsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
