// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// SessionPlayersDao is the data access object for the table session_players.
type SessionPlayersDao struct {
	table    string                // table is the underlying table name of the DAO.
	group    string                // group is the database configuration group name of the current DAO.
	columns  SessionPlayersColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler    // handlers for customized model modification.
}

// SessionPlayersColumns defines and stores column names for the table session_players.
type SessionPlayersColumns struct {
	Id          string //
	SessionId   string //
	UserId      string //
	SeatNo      string //
	TotalHands  string // Hands participated
	TotalBuyin  string // Cumulative buy-in
	ChipsFinal  string // Final chips when session ended
	Result      string // Profit/loss = chips_final - total_buyin
	Vpip        string // VPIP% in this session
	WinRate     string // Win rate% in this session
	ActivityPts string // Activity points earned
	IsMvp       string //
	Rank        string // Rank by profit desc
	JoinedAt    string //
	LeftAt      string //
}

// sessionPlayersColumns holds the columns for the table session_players.
var sessionPlayersColumns = SessionPlayersColumns{
	Id:          "id",
	SessionId:   "session_id",
	UserId:      "user_id",
	SeatNo:      "seat_no",
	TotalHands:  "total_hands",
	TotalBuyin:  "total_buyin",
	ChipsFinal:  "chips_final",
	Result:      "result",
	Vpip:        "vpip",
	WinRate:     "win_rate",
	ActivityPts: "activity_pts",
	IsMvp:       "is_mvp",
	Rank:        "rank",
	JoinedAt:    "joined_at",
	LeftAt:      "left_at",
}

// NewSessionPlayersDao creates and returns a new DAO object for table data access.
func NewSessionPlayersDao(handlers ...gdb.ModelHandler) *SessionPlayersDao {
	return &SessionPlayersDao{
		group:    "default",
		table:    "session_players",
		columns:  sessionPlayersColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *SessionPlayersDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *SessionPlayersDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *SessionPlayersDao) Columns() SessionPlayersColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *SessionPlayersDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *SessionPlayersDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *SessionPlayersDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
