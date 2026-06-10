// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// GamePlayersDao is the data access object for the table game_players.
type GamePlayersDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  GamePlayersColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// GamePlayersColumns defines and stores column names for the table game_players.
type GamePlayersColumns struct {
	Id           string //
	GameId       string //
	UserId       string //
	SeatNo       string //
	Position     string // 0=BTN 1=SB 2=BB 3=UTG 4=UTG+1 5=MP 6=HJ 7=CO
	HoleCards    string // Hole cards e.g. Ah Kd
	ForcedBet    string // Forced bet (blind+ante), separate from voluntary
	ChipsStart   string // Chips at hand start
	ChipsEnd     string // Chips at hand end
	TotalBet     string // Total voluntary bet this hand
	Result       string // Profit/loss this hand
	BestHand     string // Best 5-card combo description
	HandRank     string // Hand strength 1(HighCard,weakest)-10(RoyalFlush,strongest). NOTE: opposite to rules.md display rank
	HandRankDesc string // Hand description: High Card / One Pair / ... / Royal Flush
	IsWinner     string //
	FoldStage    string // Stage when folded (0-4)
	IsVpip       string // Voluntarily put chips in pot preflop
	IsPfr        string // Preflop raise flag
	WentToSd     string // Went to showdown flag
	IsShowCard   string // 0=muck 1=show cards at showdown
	Status       string // 1=active 2=folded 3=allin
}

// gamePlayersColumns holds the columns for the table game_players.
var gamePlayersColumns = GamePlayersColumns{
	Id:           "id",
	GameId:       "game_id",
	UserId:       "user_id",
	SeatNo:       "seat_no",
	Position:     "position",
	HoleCards:    "hole_cards",
	ForcedBet:    "forced_bet",
	ChipsStart:   "chips_start",
	ChipsEnd:     "chips_end",
	TotalBet:     "total_bet",
	Result:       "result",
	BestHand:     "best_hand",
	HandRank:     "hand_rank",
	HandRankDesc: "hand_rank_desc",
	IsWinner:     "is_winner",
	FoldStage:    "fold_stage",
	IsVpip:       "is_vpip",
	IsPfr:        "is_pfr",
	WentToSd:     "went_to_sd",
	IsShowCard:   "is_show_card",
	Status:       "status",
}

// NewGamePlayersDao creates and returns a new DAO object for table data access.
func NewGamePlayersDao(handlers ...gdb.ModelHandler) *GamePlayersDao {
	return &GamePlayersDao{
		group:    "default",
		table:    "game_players",
		columns:  gamePlayersColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *GamePlayersDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *GamePlayersDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *GamePlayersDao) Columns() GamePlayersColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *GamePlayersDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *GamePlayersDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *GamePlayersDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
