package game

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"

	"claude-test/internal/game"
	"claude-test/utility/ws"
)

// OnHandEnd persists hand results to DB and broadcasts via WebSocket.
func OnHandEnd(ctx context.Context, sessionID int64, result game.HandEndResult) error {
	now := time.Now()
	communityStr := game.CardsToStr(result.CommunityCards)
	seedHex := game.SeedToHex(result.ShuffleSeed)

	isSplit, runTwiceUsed := 0, 0
	if result.IsSplitPot {
		isSplit = 1
	}
	if result.RunTwiceUsed {
		runTwiceUsed = 1
	}
	board2Str := game.CardsToStr(result.RunTwiceBoard2)

	stage := game.StageRiver
	for _, p := range result.Players {
		if p.WentToSD {
			stage = game.StageShowdown
			break
		}
	}

	// 1. Update games row
	_, _ = g.DB().Model("games").Where("id", result.GameID).Data(g.Map{
		"community_cards":  communityStr,
		"shuffle_seed":     seedHex,
		"pot":              totalPot(result.Pots),
		"is_split_pot":     isSplit,
		"run_twice_used":   runTwiceUsed,
		"run_twice_board2": board2Str,
		"stage":            stage,
		"status":           2,
		"ended_at":         now,
		"duration_ms":      result.DurationMs,
	}).Update()

	// 2. Update game_players
	for _, p := range result.Players {
		isWinner, isShowCard, isVPIP, isPFR, wentToSD := 0, 0, 0, 0, 0
		if p.IsWinner {
			isWinner = 1
		}
		if p.IsShowCard {
			isShowCard = 1
		}
		if p.IsVPIP {
			isVPIP = 1
		}
		if p.IsPFR {
			isPFR = 1
		}
		if p.WentToSD {
			wentToSD = 1
		}
		handRank, handRankDesc, bestHand := 0, "", ""
		if p.HandResult.Rank > 0 {
			handRank = p.HandResult.Rank
			handRankDesc = p.HandResult.Desc
			bestHand = game.CardsToStr(p.HandResult.Cards)
		}
		_, _ = g.DB().Model("game_players").
			Where("game_id", result.GameID).Where("user_id", p.UserID).
			Data(g.Map{
				"hole_cards":     game.CardsToStr(p.HoleCards),
				"chips_end":      p.ChipsEnd,
				"total_bet":      p.TotalBet,
				"forced_bet":     p.ForcedBet,
				"result":         p.Result,
				"best_hand":      bestHand,
				"hand_rank":      handRank,
				"hand_rank_desc": handRankDesc,
				"is_winner":      isWinner,
				"fold_stage":     p.FoldStage,
				"is_vpip":        isVPIP,
				"is_pfr":         isPFR,
				"went_to_sd":     wentToSD,
				"is_show_card":   isShowCard,
				"position":       p.Position,
				"status":         playerStatus(p),
			}).Update()
		_, _ = g.DB().Model("table_seats").
			Where("user_id", p.UserID).
			Data(g.Map{"chips": p.ChipsEnd}).Update()
	}

	// 3. pot_distributions + pot_winner_details
	for _, pot := range result.Pots {
		potRes, e := g.DB().Model("pot_distributions").Data(g.Map{
			"game_id":      result.GameID,
			"pot_type":     pot.PotType,
			"pot_index":    pot.PotIndex,
			"amount":       pot.Amount,
			"winner_ids":   winnerIDStr(pot.Winners),
			"winner_count": len(pot.Winners),
			"win_reason":   pot.WinDesc,
			"win_rank":     pot.WinRank,
		}).Insert()
		if e != nil {
			continue
		}
		distID, _ := potRes.LastInsertId()
		for _, share := range pot.Winners {
			isSplitRow := 0
			if len(pot.Winners) > 1 {
				isSplitRow = 1
			}
			_, _ = g.DB().Model("pot_winner_details").Data(g.Map{
				"distribution_id": distID,
				"game_id":         result.GameID,
				"user_id":         seatToUserID(result.Players, share.SeatNo),
				"amount":          share.Amount,
				"is_split":        isSplitRow,
			}).Insert()
		}
	}

	// 4. hand_replays
	for _, snap := range result.Snapshots {
		playersJSON, _ := json.Marshal(snap.PlayersState)
		_, _ = g.DB().Model("hand_replays").Data(g.Map{
			"game_id":          result.GameID,
			"stage":            snap.Stage,
			"community_cards":  snap.CommunityCards,
			"pot":              snap.Pot,
			"players_state":    string(playersJSON),
			"action_seq_start": snap.ActionSeqStart,
			"action_seq_end":   snap.ActionSeqEnd,
		}).Insert()
	}

	// 5. session_players running stats
	potTotal := totalPot(result.Pots)
	for _, p := range result.Players {
		vpipVal, winVal := 0, 0
		if p.IsVPIP {
			vpipVal = 1
		}
		if p.IsWinner {
			winVal = 1
		}
		_, _ = g.DB().Model("session_players").
			Where("session_id", sessionID).Where("user_id", p.UserID).
			Data(g.Map{
				"total_hands": gdb.Raw("total_hands + 1"),
				"result":      gdb.Raw(fmt.Sprintf("result + %d", p.Result)),
				"vpip":        gdb.Raw(fmt.Sprintf("ROUND((vpip * total_hands + %d) / (total_hands + 1), 2)", vpipVal*100)),
				"win_rate":    gdb.Raw(fmt.Sprintf("ROUND((win_rate * total_hands + %d) / (total_hands + 1), 2)", winVal*100)),
			}).Update()
	}

	// 6. room_sessions aggregate
	_, _ = g.DB().Model("room_sessions").Where("id", sessionID).
		Data(g.Map{
			"total_hands": gdb.Raw("total_hands + 1"),
			"total_flow":  gdb.Raw(fmt.Sprintf("total_flow + %d", potTotal)),
			"avg_pot":     gdb.Raw("IF(total_hands > 0, ROUND(total_flow / total_hands), 0)"),
			"max_pot":     gdb.Raw(fmt.Sprintf("GREATEST(max_pot, %d)", potTotal)),
		}).Update()

	// 7. Broadcast hand result
	broadcastHandResult(sessionID, result)

	// 8. Start next hand after delay
	tableIDVal, _ := g.DB().Model("room_sessions").Fields("table_id").Where("id", sessionID).Value()
	tableID := tableIDVal.Int64()
	if tableID > 0 {
		go func() {
			time.Sleep(3 * time.Second)
			statusVal, _ := g.DB().Model("room_sessions").Fields("status").Where("id", sessionID).Value()
			if statusVal.Int() != 1 {
				return
			}
			totalHandsVal, _ := g.DB().Model("room_sessions").Fields("total_hands").Where("id", sessionID).Value()
			nextIdx := totalHandsVal.Int()
			ctx2 := context.Background()
			if e := startNextHand(ctx2, tableID, sessionID, nextIdx); e != nil {
				g.Log().Errorf(ctx2, "startNextHand error: %v", e)
			}
		}()
	}
	return nil
}

func broadcastHandResult(sessionID int64, result game.HandEndResult) {
	tableIDVal, _ := g.DB().Model("room_sessions").Fields("table_id").Where("id", sessionID).Value()
	tableID := tableIDVal.Int64()
	if tableID == 0 {
		return
	}
	type winnerInfo struct {
		UserID    int64  `json:"user_id"`
		Amount    int64  `json:"amount"`
		HandRank  int    `json:"hand_rank"`
		HandDesc  string `json:"hand_desc"`
		ShowCards bool   `json:"show_cards"`
		HoleCards string `json:"hole_cards"`
	}
	type potMsg struct {
		Type    int          `json:"type"`
		Amount  int64        `json:"amount"`
		Winners []winnerInfo `json:"winners"`
	}
	pots := make([]potMsg, 0)
	for _, pot := range result.Pots {
		winners := make([]winnerInfo, 0)
		for _, share := range pot.Winners {
			uid := seatToUserID(result.Players, share.SeatNo)
			for _, p := range result.Players {
				if p.UserID == uid {
					winners = append(winners, winnerInfo{
						UserID: uid, Amount: share.Amount,
						HandRank: p.HandResult.Rank, HandDesc: p.HandResult.Desc,
						ShowCards: p.IsShowCard, HoleCards: game.CardsToStr(p.HoleCards),
					})
					break
				}
			}
		}
		pots = append(pots, potMsg{Type: pot.PotType, Amount: pot.Amount, Winners: winners})
	}
	ws.GlobalHub().BroadcastTable(tableID, ws.MsgTypeHandResult, g.Map{
		"game_id": result.GameID, "hand_no": result.HandNo,
		"is_split_pot": result.IsSplitPot, "run_twice_used": result.RunTwiceUsed,
		"pots": pots,
	})
}

func totalPot(pots []game.PotResult) int64 {
	var total int64
	for _, p := range pots {
		total += p.Amount
	}
	return total
}

func winnerIDStr(winners []game.PotShare) string {
	s := ""
	for i, w := range winners {
		if i > 0 {
			s += ","
		}
		s += fmt.Sprintf("%d", w.SeatNo)
	}
	return s
}

func seatToUserID(players []*game.PlayerEndState, seatNo int) int64 {
	for _, p := range players {
		if p.SeatNo == seatNo {
			return p.UserID
		}
	}
	return 0
}

func playerStatus(p *game.PlayerEndState) int {
	if p.FoldedEarly {
		return game.PlayerFolded
	}
	return game.PlayerActive
}
