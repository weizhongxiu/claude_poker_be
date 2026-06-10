// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// TableInvitationsDao is the data access object for the table table_invitations.
type TableInvitationsDao struct {
	table    string                  // table is the underlying table name of the DAO.
	group    string                  // group is the database configuration group name of the current DAO.
	columns  TableInvitationsColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler      // handlers for customized model modification.
}

// TableInvitationsColumns defines and stores column names for the table table_invitations.
type TableInvitationsColumns struct {
	Id        string //
	TableId   string // Table id (before session)
	SessionId string // Session id (after start), 0=pre-session
	InviterId string //
	InviteeId string //
	Status    string // 1=pending 2=accepted 3=rejected 4=expired
	ExpiredAt string //
	CreatedAt string //
}

// tableInvitationsColumns holds the columns for the table table_invitations.
var tableInvitationsColumns = TableInvitationsColumns{
	Id:        "id",
	TableId:   "table_id",
	SessionId: "session_id",
	InviterId: "inviter_id",
	InviteeId: "invitee_id",
	Status:    "status",
	ExpiredAt: "expired_at",
	CreatedAt: "created_at",
}

// NewTableInvitationsDao creates and returns a new DAO object for table data access.
func NewTableInvitationsDao(handlers ...gdb.ModelHandler) *TableInvitationsDao {
	return &TableInvitationsDao{
		group:    "default",
		table:    "table_invitations",
		columns:  tableInvitationsColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *TableInvitationsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *TableInvitationsDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *TableInvitationsDao) Columns() TableInvitationsColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *TableInvitationsDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *TableInvitationsDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *TableInvitationsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
