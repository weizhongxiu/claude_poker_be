// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// TablesDao is the data access object for the table tables.
type TablesDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  TablesColumns      // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// TablesColumns defines and stores column names for the table tables.
type TablesColumns struct {
	Id                string //
	TableNo           string // Table number
	ClubId            string // Club id, NULL=public
	Name              string //
	HasPassword       string // Password protected
	Password          string //
	GameType          string // 1=holdem 2=shortdeck 3=plo 4=sng 5=chinese
	BlindType         string // 1=fixed 2=increasing
	SmallBlind        string //
	BigBlind          string //
	Ante              string //
	StraddleEnabled   string // Allow straddle (2x BB)
	MinBuyin          string //
	MaxBuyin          string //
	MaxBuyinTotal     string // Cumulative max buy-in, 0=unlimited
	Duration          string // Session duration in hours
	RunTwice          string // Feature switch: allow run-twice
	LowWaterInsurance string // Low water insurance
	CritGameplay      string // Critical hit gameplay
	ActivityPoints    string // Activity points enabled
	AutoRebuy         string // Auto rebuy/cashout
	BuyinApproval     string // Rebuy needs admin approval
	DelayShowCard     string // Delay show card
	RandomSeat        string // Random seat assignment
	SpectatorMute     string // Mute spectators
	GpsIpRestrict     string // GPS and IP restriction
	FullTableStart    string // Start only when full
	MaxSeats          string // Max seats 2-10 (standard 6-9)
	CurrentPlayers    string //
	Tag               string // Custom tag
	CreatorId         string //
	Status            string // 1=waiting 2=playing 3=closed
	CreatedAt         string //
	UpdatedAt         string //
	EndedAt           string //
}

// tablesColumns holds the columns for the table tables.
var tablesColumns = TablesColumns{
	Id:                "id",
	TableNo:           "table_no",
	ClubId:            "club_id",
	Name:              "name",
	HasPassword:       "has_password",
	Password:          "password",
	GameType:          "game_type",
	BlindType:         "blind_type",
	SmallBlind:        "small_blind",
	BigBlind:          "big_blind",
	Ante:              "ante",
	StraddleEnabled:   "straddle_enabled",
	MinBuyin:          "min_buyin",
	MaxBuyin:          "max_buyin",
	MaxBuyinTotal:     "max_buyin_total",
	Duration:          "duration",
	RunTwice:          "run_twice",
	LowWaterInsurance: "low_water_insurance",
	CritGameplay:      "crit_gameplay",
	ActivityPoints:    "activity_points",
	AutoRebuy:         "auto_rebuy",
	BuyinApproval:     "buyin_approval",
	DelayShowCard:     "delay_show_card",
	RandomSeat:        "random_seat",
	SpectatorMute:     "spectator_mute",
	GpsIpRestrict:     "gps_ip_restrict",
	FullTableStart:    "full_table_start",
	MaxSeats:          "max_seats",
	CurrentPlayers:    "current_players",
	Tag:               "tag",
	CreatorId:         "creator_id",
	Status:            "status",
	CreatedAt:         "created_at",
	UpdatedAt:         "updated_at",
	EndedAt:           "ended_at",
}

// NewTablesDao creates and returns a new DAO object for table data access.
func NewTablesDao(handlers ...gdb.ModelHandler) *TablesDao {
	return &TablesDao{
		group:    "default",
		table:    "tables",
		columns:  tablesColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *TablesDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *TablesDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *TablesDao) Columns() TablesColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *TablesDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *TablesDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *TablesDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
