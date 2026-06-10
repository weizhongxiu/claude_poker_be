package game

import "sort"

// Hand rank constants.
// Internal strength value: larger = stronger.
// NOTE: Opposite to rules.md display rank (rules.md: 1=strongest Royal Flush, 10=weakest High Card).
// Stored in game_players.hand_rank.
const (
	HandHighCard      = 1  // High Card       rules.md rank: 10 (weakest)
	HandOnePair       = 2  // One Pair        rules.md rank: 9
	HandTwoPair       = 3  // Two Pair        rules.md rank: 8
	HandThreeOfAKind  = 4  // Three of a Kind rules.md rank: 7
	HandStraight      = 5  // Straight        rules.md rank: 6
	HandFlush         = 6  // Flush           rules.md rank: 5
	HandFullHouse     = 7  // Full House      rules.md rank: 4
	HandFourOfAKind   = 8  // Four of a Kind  rules.md rank: 3
	HandStraightFlush = 9  // Straight Flush  rules.md rank: 2
	HandRoyalFlush    = 10 // Royal Flush     rules.md rank: 1 (strongest)
)

var handRankDesc = map[int]string{
	HandHighCard:      "High Card",
	HandOnePair:       "One Pair",
	HandTwoPair:       "Two Pair",
	HandThreeOfAKind:  "Three of a Kind",
	HandStraight:      "Straight",
	HandFlush:         "Flush",
	HandFullHouse:     "Full House",
	HandFourOfAKind:   "Four of a Kind",
	HandStraightFlush: "Straight Flush",
	HandRoyalFlush:    "Royal Flush",
}

// HandResult stores the evaluation result for a 5-card hand.
type HandResult struct {
	Rank     int    // Internal strength value 1-10
	Desc     string // e.g. "Straight Flush"
	Cards    []Card // Best 5 cards
	Tiebreak []int  // Tiebreak values for same-rank comparison (descending)
}

// EvalBest5 selects the best 5-card hand from hole cards + board (7 cards total).
// Returns HandResult with the highest strength.
func EvalBest5(hole []Card, board []Card) HandResult {
	all := append(append([]Card{}, hole...), board...)
	n := len(all)
	if n < 5 {
		// Not enough cards, evaluate as is
		return evalFive(all)
	}

	var best HandResult
	// C(n,5) combinations
	for i := 0; i < n-4; i++ {
		for j := i + 1; j < n-3; j++ {
			for k := j + 1; k < n-2; k++ {
				for l := k + 1; l < n-1; l++ {
					for m := l + 1; m < n; m++ {
						hand := []Card{all[i], all[j], all[k], all[l], all[m]}
						result := evalFive(hand)
						if best.Rank == 0 || compareResults(result, best) > 0 {
							best = result
						}
					}
				}
			}
		}
	}
	return best
}

// CompareHands returns 1 if a > b, -1 if a < b, 0 if equal (split pot).
func CompareHands(a, b HandResult) int {
	return compareResults(a, b)
}

func compareResults(a, b HandResult) int {
	if a.Rank != b.Rank {
		if a.Rank > b.Rank {
			return 1
		}
		return -1
	}
	// Same rank: compare tiebreak values
	for i := 0; i < len(a.Tiebreak) && i < len(b.Tiebreak); i++ {
		if a.Tiebreak[i] > b.Tiebreak[i] {
			return 1
		}
		if a.Tiebreak[i] < b.Tiebreak[i] {
			return -1
		}
	}
	return 0 // Split pot
}

// evalFive evaluates exactly 5 cards.
func evalFive(cards []Card) HandResult {
	c := make([]Card, len(cards))
	copy(c, cards)

	// Sort descending by rank
	sort.Slice(c, func(i, j int) bool {
		return c[i].Rank > c[j].Rank
	})

	flush := isFlush(c)
	straight, highCard := isStraight(c)

	if flush && straight {
		if highCard == RankAce {
			return HandResult{
				Rank:     HandRoyalFlush,
				Desc:     handRankDesc[HandRoyalFlush],
				Cards:    c,
				Tiebreak: []int{highCard},
			}
		}
		return HandResult{
			Rank:     HandStraightFlush,
			Desc:     handRankDesc[HandStraightFlush],
			Cards:    c,
			Tiebreak: []int{highCard},
		}
	}

	groups := groupByRank(c)

	if r, ok := hasFourOfAKind(groups); ok {
		kicker := kickersExcept(c, r, 1)
		return HandResult{
			Rank:     HandFourOfAKind,
			Desc:     handRankDesc[HandFourOfAKind],
			Cards:    c,
			Tiebreak: append([]int{r}, kicker...),
		}
	}

	if trip, pair, ok := hasFullHouse(groups); ok {
		return HandResult{
			Rank:     HandFullHouse,
			Desc:     handRankDesc[HandFullHouse],
			Cards:    c,
			Tiebreak: []int{trip, pair},
		}
	}

	if flush {
		tiebreak := ranksDesc(c)
		return HandResult{
			Rank:     HandFlush,
			Desc:     handRankDesc[HandFlush],
			Cards:    c,
			Tiebreak: tiebreak,
		}
	}

	if straight {
		return HandResult{
			Rank:     HandStraight,
			Desc:     handRankDesc[HandStraight],
			Cards:    c,
			Tiebreak: []int{highCard},
		}
	}

	if r, ok := hasThreeOfAKind(groups); ok {
		kicker := kickersExcept(c, r, 2)
		return HandResult{
			Rank:     HandThreeOfAKind,
			Desc:     handRankDesc[HandThreeOfAKind],
			Cards:    c,
			Tiebreak: append([]int{r}, kicker...),
		}
	}

	if high, low, ok := hasTwoPair(groups); ok {
		kicker := kickersExcept2(c, high, low, 1)
		return HandResult{
			Rank:     HandTwoPair,
			Desc:     handRankDesc[HandTwoPair],
			Cards:    c,
			Tiebreak: append([]int{high, low}, kicker...),
		}
	}

	if r, ok := hasOnePair(groups); ok {
		kicker := kickersExcept(c, r, 3)
		return HandResult{
			Rank:     HandOnePair,
			Desc:     handRankDesc[HandOnePair],
			Cards:    c,
			Tiebreak: append([]int{r}, kicker...),
		}
	}

	// High Card
	return HandResult{
		Rank:     HandHighCard,
		Desc:     handRankDesc[HandHighCard],
		Cards:    c,
		Tiebreak: ranksDesc(c),
	}
}

// ---- helpers ----

func isFlush(cards []Card) bool {
	suit := cards[0].Suit
	for _, c := range cards[1:] {
		if c.Suit != suit {
			return false
		}
	}
	return true
}

// isStraight checks for a straight (including A-low: A-2-3-4-5).
// Returns (true, highCard) or (false, 0).
func isStraight(cards []Card) (bool, int) {
	ranks := make([]int, len(cards))
	for i, c := range cards {
		ranks[i] = c.Rank
	}
	sort.Sort(sort.Reverse(sort.IntSlice(ranks)))

	// Normal straight
	ok := true
	for i := 1; i < len(ranks); i++ {
		if ranks[i] != ranks[i-1]-1 {
			ok = false
			break
		}
	}
	if ok {
		return true, ranks[0]
	}

	// Ace-low straight: A-2-3-4-5 (wheel)
	if ranks[0] == RankAce && ranks[1] == 5 && ranks[2] == 4 && ranks[3] == 3 && ranks[4] == 2 {
		return true, 5 // high card is 5 for wheel
	}

	return false, 0
}

// groupByRank returns rank → count map.
func groupByRank(cards []Card) map[int]int {
	g := make(map[int]int)
	for _, c := range cards {
		g[c.Rank]++
	}
	return g
}

func hasFourOfAKind(groups map[int]int) (int, bool) {
	for r, cnt := range groups {
		if cnt == 4 {
			return r, true
		}
	}
	return 0, false
}

func hasFullHouse(groups map[int]int) (int, int, bool) {
	var trip, pair int
	for r, cnt := range groups {
		if cnt == 3 {
			trip = r
		} else if cnt == 2 {
			pair = r
		}
	}
	if trip > 0 && pair > 0 {
		return trip, pair, true
	}
	return 0, 0, false
}

func hasThreeOfAKind(groups map[int]int) (int, bool) {
	for r, cnt := range groups {
		if cnt == 3 {
			return r, true
		}
	}
	return 0, false
}

func hasTwoPair(groups map[int]int) (int, int, bool) {
	var pairs []int
	for r, cnt := range groups {
		if cnt == 2 {
			pairs = append(pairs, r)
		}
	}
	if len(pairs) >= 2 {
		sort.Sort(sort.Reverse(sort.IntSlice(pairs)))
		return pairs[0], pairs[1], true
	}
	return 0, 0, false
}

func hasOnePair(groups map[int]int) (int, bool) {
	for r, cnt := range groups {
		if cnt == 2 {
			return r, true
		}
	}
	return 0, false
}

// ranksDesc returns card ranks sorted descending (for tiebreak).
func ranksDesc(cards []Card) []int {
	ranks := make([]int, len(cards))
	for i, c := range cards {
		ranks[i] = c.Rank
	}
	sort.Sort(sort.Reverse(sort.IntSlice(ranks)))
	return ranks
}

// kickersExcept returns the top n kicker ranks excluding the given rank.
func kickersExcept(cards []Card, exclude int, n int) []int {
	var kickers []int
	for _, c := range cards {
		if c.Rank != exclude {
			kickers = append(kickers, c.Rank)
		}
	}
	sort.Sort(sort.Reverse(sort.IntSlice(kickers)))
	if len(kickers) > n {
		kickers = kickers[:n]
	}
	return kickers
}

// kickersExcept2 returns top n kicker ranks excluding two given ranks.
func kickersExcept2(cards []Card, ex1, ex2 int, n int) []int {
	var kickers []int
	for _, c := range cards {
		if c.Rank != ex1 && c.Rank != ex2 {
			kickers = append(kickers, c.Rank)
		}
	}
	sort.Sort(sort.Reverse(sort.IntSlice(kickers)))
	if len(kickers) > n {
		kickers = kickers[:n]
	}
	return kickers
}
