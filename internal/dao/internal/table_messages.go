// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// TableMessagesDao is the data access object for the table table_messages.
type TableMessagesDao struct {
	table    string               // table is the underlying table name of the DAO.
	group    string               // group is the database configuration group name of the current DAO.
	columns  TableMessagesColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler   // handlers for customized model modification.
}

// TableMessagesColumns defines and stores column names for the table table_messages.
type TableMessagesColumns struct {
	Id        string //
	SessionId string //
	UserId    string // 0=system
	Type      string // 1=text 2=emoji 3=system 4=invite
	Content   string //
	CreatedAt string //
}

// tableMessagesColumns holds the columns for the table table_messages.
var tableMessagesColumns = TableMessagesColumns{
	Id:        "id",
	SessionId: "session_id",
	UserId:    "user_id",
	Type:      "type",
	Content:   "content",
	CreatedAt: "created_at",
}

// NewTableMessagesDao creates and returns a new DAO object for table data access.
func NewTableMessagesDao(handlers ...gdb.ModelHandler) *TableMessagesDao {
	return &TableMessagesDao{
		group:    "default",
		table:    "table_messages",
		columns:  tableMessagesColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *TableMessagesDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *TableMessagesDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *TableMessagesDao) Columns() TableMessagesColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *TableMessagesDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *TableMessagesDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *TableMessagesDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
