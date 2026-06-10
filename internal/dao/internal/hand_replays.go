// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// HandReplaysDao is the data access object for the table hand_replays.
type HandReplaysDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  HandReplaysColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// HandReplaysColumns defines and stores column names for the table hand_replays.
type HandReplaysColumns struct {
	Id             string //
	GameId         string //
	HandIndex      string // Hand index in session (for replay progress)
	Stage          string // 1=preflop 2=flop 3=turn 4=river 5=showdown
	CommunityCards string // Board at this stage
	Pot            string //
	PlayersState   string // Player state snapshot [{seat,chips,bet,status,hole_cards}]
	ActionSeqStart string //
	ActionSeqEnd   string //
}

// handReplaysColumns holds the columns for the table hand_replays.
var handReplaysColumns = HandReplaysColumns{
	Id:             "id",
	GameId:         "game_id",
	HandIndex:      "hand_index",
	Stage:          "stage",
	CommunityCards: "community_cards",
	Pot:            "pot",
	PlayersState:   "players_state",
	ActionSeqStart: "action_seq_start",
	ActionSeqEnd:   "action_seq_end",
}

// NewHandReplaysDao creates and returns a new DAO object for table data access.
func NewHandReplaysDao(handlers ...gdb.ModelHandler) *HandReplaysDao {
	return &HandReplaysDao{
		group:    "default",
		table:    "hand_replays",
		columns:  handReplaysColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *HandReplaysDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *HandReplaysDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *HandReplaysDao) Columns() HandReplaysColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *HandReplaysDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *HandReplaysDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *HandReplaysDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
