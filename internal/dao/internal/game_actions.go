// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// GameActionsDao is the data access object for the table game_actions.
type GameActionsDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  GameActionsColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// GameActionsColumns defines and stores column names for the table game_actions.
type GameActionsColumns struct {
	Id        string //
	GameId    string //
	UserId    string //
	SeatNo    string //
	Stage     string // 0=blinds 1=preflop 2=flop 3=turn 4=river
	Action    string // 1=fold 2=check 3=call 4=raise 5=allin 6=bet 7=blind_post 8=ante_post 9=straddle
	Amount    string //
	PotAfter  string // Pot size after this action
	ActionSeq string // Action sequence in this hand
	ActionMs  string // Decision time in milliseconds
	CreatedAt string //
}

// gameActionsColumns holds the columns for the table game_actions.
var gameActionsColumns = GameActionsColumns{
	Id:        "id",
	GameId:    "game_id",
	UserId:    "user_id",
	SeatNo:    "seat_no",
	Stage:     "stage",
	Action:    "action",
	Amount:    "amount",
	PotAfter:  "pot_after",
	ActionSeq: "action_seq",
	ActionMs:  "action_ms",
	CreatedAt: "created_at",
}

// NewGameActionsDao creates and returns a new DAO object for table data access.
func NewGameActionsDao(handlers ...gdb.ModelHandler) *GameActionsDao {
	return &GameActionsDao{
		group:    "default",
		table:    "game_actions",
		columns:  gameActionsColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *GameActionsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *GameActionsDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *GameActionsDao) Columns() GameActionsColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *GameActionsDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *GameActionsDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *GameActionsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
