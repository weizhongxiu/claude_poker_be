// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// TableSeatsDao is the data access object for the table table_seats.
type TableSeatsDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  TableSeatsColumns  // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// TableSeatsColumns defines and stores column names for the table table_seats.
type TableSeatsColumns struct {
	Id           string //
	TableId      string //
	UserId       string //
	SeatNo       string // Seat number 1-10
	Chips        string // Chips on table
	Status       string // 1=seated 2=left
	IsSittingOut string // Sitting out flag
	JoinedAt     string //
	UpdatedAt    string //
}

// tableSeatsColumns holds the columns for the table table_seats.
var tableSeatsColumns = TableSeatsColumns{
	Id:           "id",
	TableId:      "table_id",
	UserId:       "user_id",
	SeatNo:       "seat_no",
	Chips:        "chips",
	Status:       "status",
	IsSittingOut: "is_sitting_out",
	JoinedAt:     "joined_at",
	UpdatedAt:    "updated_at",
}

// NewTableSeatsDao creates and returns a new DAO object for table data access.
func NewTableSeatsDao(handlers ...gdb.ModelHandler) *TableSeatsDao {
	return &TableSeatsDao{
		group:    "default",
		table:    "table_seats",
		columns:  tableSeatsColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *TableSeatsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *TableSeatsDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *TableSeatsDao) Columns() TableSeatsColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *TableSeatsDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *TableSeatsDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *TableSeatsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
