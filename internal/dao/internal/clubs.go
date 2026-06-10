// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ClubsDao is the data access object for the table clubs.
type ClubsDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  ClubsColumns       // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// ClubsColumns defines and stores column names for the table clubs.
type ClubsColumns struct {
	Id           string //
	ClubNo       string // Club number
	Name         string //
	Logo         string //
	OwnerId      string // Owner user_id
	Announcement string // Announcement
	MemberCount  string //
	MaxMembers   string //
	Status       string // 1=active 2=dissolved
	CreatedAt    string //
	UpdatedAt    string //
	DeletedAt    string //
}

// clubsColumns holds the columns for the table clubs.
var clubsColumns = ClubsColumns{
	Id:           "id",
	ClubNo:       "club_no",
	Name:         "name",
	Logo:         "logo",
	OwnerId:      "owner_id",
	Announcement: "announcement",
	MemberCount:  "member_count",
	MaxMembers:   "max_members",
	Status:       "status",
	CreatedAt:    "created_at",
	UpdatedAt:    "updated_at",
	DeletedAt:    "deleted_at",
}

// NewClubsDao creates and returns a new DAO object for table data access.
func NewClubsDao(handlers ...gdb.ModelHandler) *ClubsDao {
	return &ClubsDao{
		group:    "default",
		table:    "clubs",
		columns:  clubsColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *ClubsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *ClubsDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *ClubsDao) Columns() ClubsColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *ClubsDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *ClubsDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *ClubsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
