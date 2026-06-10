// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// GamePlayers is the golang structure of table game_players for DAO operations like Where/Data.
type GamePlayers struct {
	g.Meta       `orm:"table:game_players, do:true"`
	Id           any //
	GameId       any //
	UserId       any //
	SeatNo       any //
	Position     any // 0=BTN 1=SB 2=BB 3=UTG 4=UTG+1 5=MP 6=HJ 7=CO
	HoleCards    any // Hole cards e.g. Ah Kd
	ForcedBet    any // Forced bet (blind+ante), separate from voluntary
	ChipsStart   any // Chips at hand start
	ChipsEnd     any // Chips at hand end
	TotalBet     any // Total voluntary bet this hand
	Result       any // Profit/loss this hand
	BestHand     any // Best 5-card combo description
	HandRank     any // Hand strength 1(HighCard,weakest)-10(RoyalFlush,strongest). NOTE: opposite to rules.md display rank
	HandRankDesc any // Hand description: High Card / One Pair / ... / Royal Flush
	IsWinner     any //
	FoldStage    any // Stage when folded (0-4)
	IsVpip       any // Voluntarily put chips in pot preflop
	IsPfr        any // Preflop raise flag
	WentToSd     any // Went to showdown flag
	IsShowCard   any // 0=muck 1=show cards at showdown
	Status       any // 1=active 2=folded 3=allin
}
