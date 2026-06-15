package game

import "math/rand"

// BotDecide returns the action a bot should take given the current game state.
// It uses a simplified strategy: pre-flop hand strength + pot odds.
func BotDecide(state *GameState, seatNo int) PlayerAction {
	p, ok := state.Players[seatNo]
	if !ok {
		return PlayerAction{SeatNo: seatNo, Action: ActionFold}
	}

	maxBet := maxBetInRound(state)
	callAmt := maxBet - p.Bet
	if callAmt < 0 {
		callAmt = 0
	}

	strength := preflopStrength(p.HoleCards)
	// Post-flop: upgrade strength with made hand
	if len(state.CommunityCards) >= 3 {
		res := EvalBest5(p.HoleCards, state.CommunityCards)
		if res.Rank >= HandTwoPair {
			strength += 3
		} else if res.Rank == HandOnePair {
			strength += 1
		}
	}

	potOdds := 0.0
	if state.Pot > 0 && callAmt > 0 {
		potOdds = float64(callAmt) / float64(state.Pot+callAmt)
	}

	r := rand.Float64()

	return decide(p, seatNo, state, callAmt, strength, potOdds, r)
}

func decide(p *PlayerState, seatNo int, state *GameState, callAmt int64, strength int, potOdds float64, r float64) PlayerAction {
	act := func(action int, amount int64) PlayerAction {
		return PlayerAction{SeatNo: seatNo, UserID: p.UserID, Action: action, Amount: amount}
	}

	bb := state.BigBlind
	if bb == 0 {
		bb = 2
	}

	switch {
	case strength >= 8: // premium: AA KK QQ AKs
		if callAmt == 0 {
			bet := clampBet(state.Pot/2+bb, bb, p.Chips)
			return act(ActionBet, bet)
		}
		if r < 0.65 {
			raise := clampBet(callAmt*3, callAmt+bb, p.Chips)
			return act(ActionRaise, raise)
		}
		return act(ActionCall, callAmt)

	case strength >= 5: // playable: JJ TT 99 AQ AJ KQ suited connectors
		if callAmt == 0 {
			if r < 0.35 {
				bet := clampBet(state.Pot/3+bb, bb, p.Chips)
				return act(ActionBet, bet)
			}
			return act(ActionCheck, 0)
		}
		if potOdds < 0.35 || r < 0.55 {
			return act(ActionCall, callAmt)
		}
		return act(ActionFold, 0)

	default: // weak
		if callAmt == 0 {
			return act(ActionCheck, 0)
		}
		if r < 0.12 && potOdds < 0.2 {
			return act(ActionCall, callAmt)
		}
		return act(ActionFold, 0)
	}
}

// preflopStrength returns 0-10 based on hole cards alone.
func preflopStrength(cards []Card) int {
	if len(cards) < 2 {
		return 0
	}
	c1, c2 := cards[0], cards[1]
	r1, r2 := c1.Rank, c2.Rank
	if r1 < r2 {
		r1, r2 = r2, r1 // r1 >= r2
	}
	suited := c1.Suit == c2.Suit
	isPair := r1 == r2

	if isPair {
		switch {
		case r1 >= RankAce:
			return 10
		case r1 >= RankKing:
			return 9
		case r1 >= RankQueen:
			return 8
		case r1 >= RankJack:
			return 7
		case r1 >= RankTen:
			return 6
		case r1 >= RankEight:
			return 5
		default:
			return 3
		}
	}

	// High-card hands
	score := 0
	if r1 == RankAce {
		score += 4
	} else if r1 == RankKing {
		score += 3
	} else if r1 == RankQueen {
		score += 2
	} else if r1 == RankJack {
		score += 1
	}
	gap := r1 - r2
	if gap == 1 {
		score += 3
	} else if gap == 2 {
		score += 2
	} else if gap <= 3 {
		score += 1
	}
	if suited {
		score += 1
	}
	return score
}

func maxBetInRound(state *GameState) int64 {
	var max int64
	for _, p := range state.Players {
		if p.Bet > max {
			max = p.Bet
		}
	}
	return max
}

func clampBet(amount, min, max int64) int64 {
	if amount < min {
		amount = min
	}
	if amount > max {
		amount = max
	}
	return amount
}
