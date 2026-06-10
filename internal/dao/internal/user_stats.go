// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// UserStatsDao is the data access object for the table user_stats.
type UserStatsDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  UserStatsColumns   // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// UserStatsColumns defines and stores column names for the table user_stats.
type UserStatsColumns struct {
	Id            string //
	UserId        string //
	GameType      string // 0=all 1=holdem 2=shortdeck 3=plo 4=sng 5=chinese
	StatType      string // 1=day 2=week 3=month 4=custom
	StatDate      string // Stat start date, NULL=all-time
	StatEnd       string // Stat end date (stat_type=4 only)
	TotalGames    string // Total hands
	TotalHands    string // Total hands participated
	TotalSessions string // Total sessions
	TotalWins     string // Winning hands
	TotalProfit   string // Total profit/loss
	TotalBuyin    string // Total buy-in
	TotalFlow     string // Total flow
	BiggestPot    string // Biggest pot won
	Vpip          string // VPIP%
	Pfr           string // PFR%
	Wtsd          string // Went to showdown%
	UpdatedAt     string //
}

// userStatsColumns holds the columns for the table user_stats.
var userStatsColumns = UserStatsColumns{
	Id:            "id",
	UserId:        "user_id",
	GameType:      "game_type",
	StatType:      "stat_type",
	StatDate:      "stat_date",
	StatEnd:       "stat_end",
	TotalGames:    "total_games",
	TotalHands:    "total_hands",
	TotalSessions: "total_sessions",
	TotalWins:     "total_wins",
	TotalProfit:   "total_profit",
	TotalBuyin:    "total_buyin",
	TotalFlow:     "total_flow",
	BiggestPot:    "biggest_pot",
	Vpip:          "vpip",
	Pfr:           "pfr",
	Wtsd:          "wtsd",
	UpdatedAt:     "updated_at",
}

// NewUserStatsDao creates and returns a new DAO object for table data access.
func NewUserStatsDao(handlers ...gdb.ModelHandler) *UserStatsDao {
	return &UserStatsDao{
		group:    "default",
		table:    "user_stats",
		columns:  userStatsColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *UserStatsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *UserStatsDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *UserStatsDao) Columns() UserStatsColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *UserStatsDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *UserStatsDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *UserStatsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
