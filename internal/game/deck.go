package game

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"strings"
)

// Suit constants
const (
	Spades   = 0 // ♠
	Hearts   = 1 // ♥
	Diamonds = 2 // ♦
	Clubs    = 3 // ♣
)

// Rank constants: 2-14, Ace=14
const (
	RankTwo   = 2
	RankThree = 3
	RankFour  = 4
	RankFive  = 5
	RankSix   = 6
	RankSeven = 7
	RankEight = 8
	RankNine  = 9
	RankTen   = 10
	RankJack  = 11
	RankQueen = 12
	RankKing  = 13
	RankAce   = 14
)

var rankStr = map[int]string{
	2: "2", 3: "3", 4: "4", 5: "5", 6: "6",
	7: "7", 8: "8", 9: "9", 10: "T",
	11: "J", 12: "Q", 13: "K", 14: "A",
}

var suitStr = map[int]string{
	Spades: "s", Hearts: "h", Diamonds: "d", Clubs: "c",
}

var strRank = map[string]int{
	"2": 2, "3": 3, "4": 4, "5": 5, "6": 6,
	"7": 7, "8": 8, "9": 9, "T": 10,
	"J": 11, "Q": 12, "K": 13, "A": 14,
}

var strSuit = map[string]int{
	"s": Spades, "h": Hearts, "d": Diamonds, "c": Clubs,
}

// Card represents a single playing card.
type Card struct {
	Rank int // 2-14 (Ace=14)
	Suit int // 0=♠ 1=♥ 2=♦ 3=♣
}

// String returns the card notation, e.g. "Ah", "Kd", "2c", "Ts"
func (c Card) String() string {
	return rankStr[c.Rank] + suitStr[c.Suit]
}

// CardToStr converts a Card to its string notation.
func CardToStr(c Card) string {
	return c.String()
}

// StrToCard parses a card string like "Ah", "Kd", "Ts", "2c".
// Returns zero Card{} on invalid input.
func StrToCard(s string) Card {
	if len(s) != 2 {
		return Card{}
	}
	r, okR := strRank[string(s[0])]
	su, okS := strSuit[string(s[1])]
	if !okR || !okS {
		return Card{}
	}
	return Card{Rank: r, Suit: su}
}

// CardsToStr converts a slice of cards to a space-separated string.
func CardsToStr(cards []Card) string {
	parts := make([]string, len(cards))
	for i, c := range cards {
		parts[i] = c.String()
	}
	return strings.Join(parts, " ")
}

// StrToCards parses a space-separated card string like "Ah Kd Qc".
func StrToCards(s string) []Card {
	if s == "" {
		return nil
	}
	parts := strings.Fields(s)
	cards := make([]Card, 0, len(parts))
	for _, p := range parts {
		c := StrToCard(p)
		if c.Rank != 0 {
			cards = append(cards, c)
		}
	}
	return cards
}

// NewDeck returns a fresh 52-card deck in suit-rank order.
func NewDeck() []Card {
	deck := make([]Card, 0, 52)
	for suit := 0; suit <= 3; suit++ {
		for rank := 2; rank <= 14; rank++ {
			deck = append(deck, Card{Rank: rank, Suit: suit})
		}
	}
	return deck
}

// Shuffle performs a Fisher-Yates shuffle using crypto/rand for true randomness.
// It returns a new shuffled slice (input is not mutated).
func Shuffle(deck []Card) ([]Card, []byte) {
	result := make([]Card, len(deck))
	copy(result, deck)

	// Generate seed for reproducibility / fairness audit
	seed := make([]byte, 32)
	_, _ = rand.Read(seed)

	shuffleWithSeed(result, seed)
	return result, seed
}

// ShuffleWithSeed shuffles using a provided seed (for replay / verification).
func ShuffleWithSeed(deck []Card, seed []byte) []Card {
	result := make([]Card, len(deck))
	copy(result, deck)
	shuffleWithSeed(result, seed)
	return result
}

func shuffleWithSeed(cards []Card, seed []byte) {
	n := len(cards)
	// Build a simple PRNG from seed (XOR-shift seeded from crypto/rand seed)
	var state [4]uint64
	for i := 0; i < 4; i++ {
		if i*8+8 <= len(seed) {
			state[i] = binary.LittleEndian.Uint64(seed[i*8:])
		}
	}
	if state[0] == 0 {
		state[0] = 1
	}

	next := func() uint64 {
		// xoshiro256** algorithm
		result := rotl(state[1]*5, 7) * 9
		t := state[1] << 17
		state[2] ^= state[0]
		state[3] ^= state[1]
		state[1] ^= state[2]
		state[0] ^= state[3]
		state[2] ^= t
		state[3] = rotl(state[3], 45)
		return result
	}

	for i := n - 1; i > 0; i-- {
		j := int(next() % uint64(i+1))
		cards[i], cards[j] = cards[j], cards[i]
	}
}

func rotl(x uint64, k int) uint64 {
	return (x << k) | (x >> (64 - k))
}

// SeedToHex converts a seed to a hex string for storage.
func SeedToHex(seed []byte) string {
	return fmt.Sprintf("%x", seed)
}
