// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ClubMembersDao is the data access object for the table club_members.
type ClubMembersDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  ClubMembersColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// ClubMembersColumns defines and stores column names for the table club_members.
type ClubMembersColumns struct {
	Id        string //
	ClubId    string //
	UserId    string //
	Role      string // 1=owner 2=admin 3=member
	Chips     string // Club chips
	Status    string // 1=active 2=banned
	JoinedAt  string //
	UpdatedAt string //
}

// clubMembersColumns holds the columns for the table club_members.
var clubMembersColumns = ClubMembersColumns{
	Id:        "id",
	ClubId:    "club_id",
	UserId:    "user_id",
	Role:      "role",
	Chips:     "chips",
	Status:    "status",
	JoinedAt:  "joined_at",
	UpdatedAt: "updated_at",
}

// NewClubMembersDao creates and returns a new DAO object for table data access.
func NewClubMembersDao(handlers ...gdb.ModelHandler) *ClubMembersDao {
	return &ClubMembersDao{
		group:    "default",
		table:    "club_members",
		columns:  clubMembersColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *ClubMembersDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *ClubMembersDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *ClubMembersDao) Columns() ClubMembersColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *ClubMembersDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *ClubMembersDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *ClubMembersDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
