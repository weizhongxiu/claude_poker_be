package game

import "sort"

// Pot represents a single pot (main or side).
type Pot struct {
	Amount        int64 // Total chips in this pot
	EligibleSeats []int // Seat numbers eligible to win this pot
}

// PotShare is one winner's share of a pot.
type PotShare struct {
	SeatNo int
	Amount int64
}

// CalcPots calculates main pot and side pots given each player's total bet.
//
// Parameters:
//   - bets: map[seatNo]totalBetAmount (includes all-in amounts)
//   - foldedSeats: seats that folded (they contributed to pots but cannot win)
//
// Returns ordered pots: pots[0] = main pot, pots[1..] = side pots.
func CalcPots(bets map[int]int64, foldedSeats []int) ([]Pot, map[int]int64) {
	foldedSet := make(map[int]bool)
	for _, s := range foldedSeats {
		foldedSet[s] = true
	}

	// Collect all seats (active + folded) that put chips in
	type seatBet struct {
		seat int
		bet  int64
	}
	var allBets []seatBet
	for seat, bet := range bets {
		if bet > 0 {
			allBets = append(allBets, seatBet{seat, bet})
		}
	}
	sort.Slice(allBets, func(i, j int) bool {
		return allBets[i].bet < allBets[j].bet
	})

	remaining := make(map[int]int64)
	for _, sb := range allBets {
		remaining[sb.seat] = sb.bet
	}

	var pots []Pot
	refunds := make(map[int]int64) // chips to return to players

	processed := int64(0)
	for _, sb := range allBets {
		level := sb.bet - processed
		if level <= 0 {
			continue
		}

		// Determine who participates in this pot level
		var potAmount int64
		var eligible []int
		for seat, rem := range remaining {
			contribution := min64(rem, level)
			potAmount += contribution
			remaining[seat] -= contribution
			if !foldedSet[seat] {
				eligible = append(eligible, seat)
			}
		}

		// Clean up zero remaining
		for seat, rem := range remaining {
			if rem == 0 {
				delete(remaining, seat)
			}
		}

		if potAmount > 0 {
			sort.Ints(eligible)
			pots = append(pots, Pot{Amount: potAmount, EligibleSeats: eligible})
		}
		processed = sb.bet
	}

	// Any leftover in remaining = over-bets from callers vs allin player, return them
	for seat, rem := range remaining {
		if rem > 0 {
			refunds[seat] += rem
		}
	}

	// Merge consecutive pots with identical eligible sets (simplify)
	return mergePots(pots), refunds
}

func mergePots(pots []Pot) []Pot {
	if len(pots) == 0 {
		return pots
	}
	result := []Pot{pots[0]}
	for i := 1; i < len(pots); i++ {
		last := &result[len(result)-1]
		if seatSetsEqual(last.EligibleSeats, pots[i].EligibleSeats) {
			last.Amount += pots[i].Amount
		} else {
			result = append(result, pots[i])
		}
	}
	return result
}

func seatSetsEqual(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// SplitPot divides a pot amount among multiple winners.
//
// Odd chips rule: the player whose seat is closest to the left of the dealer
// (dealerLeftFirstSeat = seat number of that player) receives the extra chip(s).
//
// Returns []PotShare, one entry per winner.
func SplitPot(amount int64, winnerSeats []int, dealerLeftFirstSeat int) []PotShare {
	if len(winnerSeats) == 0 {
		return nil
	}
	if len(winnerSeats) == 1 {
		return []PotShare{{SeatNo: winnerSeats[0], Amount: amount}}
	}

	base := amount / int64(len(winnerSeats))
	remainder := int(amount % int64(len(winnerSeats)))

	shares := make([]PotShare, len(winnerSeats))
	for i, seat := range winnerSeats {
		shares[i] = PotShare{SeatNo: seat, Amount: base}
	}

	// Distribute remainder chips to winners in seat order starting from dealerLeftFirstSeat
	if remainder > 0 {
		// Sort winners by proximity to dealerLeftFirstSeat (closest first)
		sorted := seatsByProximity(winnerSeats, dealerLeftFirstSeat)
		for i := 0; i < remainder && i < len(sorted); i++ {
			for j := range shares {
				if shares[j].SeatNo == sorted[i] {
					shares[j].Amount++
					break
				}
			}
		}
	}

	return shares
}

// seatsByProximity returns seats ordered by clockwise proximity to startSeat.
func seatsByProximity(seats []int, startSeat int) []int {
	maxSeat := 10
	result := make([]int, len(seats))
	copy(result, seats)
	sort.Slice(result, func(i, j int) bool {
		di := (result[i] - startSeat + maxSeat) % maxSeat
		dj := (result[j] - startSeat + maxSeat) % maxSeat
		return di < dj
	})
	return result
}

func min64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}
