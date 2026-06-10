// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// GamesDao is the data access object for the table games.
type GamesDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  GamesColumns       // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// GamesColumns defines and stores column names for the table games.
type GamesColumns struct {
	Id             string //
	TableId        string //
	SessionId      string // Parent session id
	HandNo         string // Hand number (unique)
	ShuffleSeed    string // Shuffle seed, published after hand ends
	DealerSeat     string // Dealer button seat
	SmallBlind     string //
	BigBlind       string //
	Ante           string //
	Pot            string // Total pot
	CommunityCards string // Board cards e.g. Ah Kd Qc Jh Ts
	IsSplitPot     string // Split pot flag
	RunTwiceUsed   string // Run-twice actually executed
	RunTwiceBoard2 string // Second board for run-twice (Turn+River)
	Stage          string // 0=blinds 1=preflop 2=flop 3=turn 4=river 5=showdown
	Status         string // 1=running 2=ended
	StartedAt      string //
	EndedAt        string //
	DurationMs     string // Hand duration in milliseconds
}

// gamesColumns holds the columns for the table games.
var gamesColumns = GamesColumns{
	Id:             "id",
	TableId:        "table_id",
	SessionId:      "session_id",
	HandNo:         "hand_no",
	ShuffleSeed:    "shuffle_seed",
	DealerSeat:     "dealer_seat",
	SmallBlind:     "small_blind",
	BigBlind:       "big_blind",
	Ante:           "ante",
	Pot:            "pot",
	CommunityCards: "community_cards",
	IsSplitPot:     "is_split_pot",
	RunTwiceUsed:   "run_twice_used",
	RunTwiceBoard2: "run_twice_board2",
	Stage:          "stage",
	Status:         "status",
	StartedAt:      "started_at",
	EndedAt:        "ended_at",
	DurationMs:     "duration_ms",
}

// NewGamesDao creates and returns a new DAO object for table data access.
func NewGamesDao(handlers ...gdb.ModelHandler) *GamesDao {
	return &GamesDao{
		group:    "default",
		table:    "games",
		columns:  gamesColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *GamesDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *GamesDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *GamesDao) Columns() GamesColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *GamesDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *GamesDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *GamesDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
