package game

// Stage constants for game progression.
const (
	StageBlinds   = 0 // Forced blinds / ante / straddle posting
	StagePreFlop  = 1
	StageFlop     = 2
	StageTurn     = 3
	StageRiver    = 4
	StageShowdown = 5
)

// Player status constants.
const (
	PlayerActive = 1
	PlayerFolded = 2
	PlayerAllIn  = 3
)

// Action type constants (maps to game_actions.action).
const (
	ActionFold     = 1 // Fold
	ActionCheck    = 2 // Check
	ActionCall     = 3 // Call
	ActionRaise    = 4 // Raise (after someone already bet)
	ActionAllIn    = 5 // All-in
	ActionBet      = 6 // Bet (first aggressor in a round)
	ActionBlind    = 7 // Blind post (SB/BB forced, stage=0)
	ActionAnte     = 8 // Ante post (forced, stage=0)
	ActionStraddle = 9 // Straddle (UTG optional 2×BB, stage=0)
)

// Position constants (maps to game_players.position).
const (
	PosBTN  = 0 // Button / Dealer
	PosSB   = 1 // Small Blind
	PosBB   = 2 // Big Blind
	PosUTG  = 3 // Under the Gun
	PosUTG1 = 4 // UTG+1
	PosMP   = 5 // Middle Position
	PosHJ   = 6 // Hijack
	PosCO   = 7 // Cutoff
)

// PlayerState holds the in-memory state of one player at the table.
type PlayerState struct {
	UserID    int64
	Nickname  string
	Avatar    string
	SeatNo    int
	Position  int    // PosBTN / PosSB / ...
	Chips     int64  // Chips currently on the table
	Bet       int64  // Chips bet in the current betting round
	TotalBet  int64  // Total voluntary bet this hand
	ForcedBet int64  // Forced bet this hand (blind + ante)
	Status    int    // PlayerActive / PlayerFolded / PlayerAllIn
	HoleCards []Card // Private hole cards
	FoldStage int    // Stage when folded (0 = not folded)
	IsVPIP    bool   // Voluntarily put chips in preflop
	IsPFR     bool   // Preflop raise
	WentToSD  bool   // Went to showdown
}

// GameState holds the full in-memory state of one hand.
type GameState struct {
	TableID        int64
	SessionID      int64
	GameID         int64
	HandNo         string
	ShuffleSeed    []byte
	Stage          int
	DealerSeat     int
	Pot            int64
	SidePots       []Pot
	CommunityCards []Card
	Deck           []Card   // Remaining deck (cards dealt are removed)
	Players        map[int]*PlayerState // key = seatNo
	SeatOrder      []int    // Clockwise seat order
	CurrentSeat    int      // Which seat is to act
	ActionDeadline int64    // Unix ms deadline for current action
	ActionSeq      int      // Global action counter for this hand
	HandIndex      int      // Which hand in the session (1-based)
	RunTwiceUsed   bool     // Run-twice actually executed
	RunTwiceBoard2 []Card   // Second board (Turn+River) if run-twice
	IsStraddled    bool     // Straddle was posted
	StraddleSeat   int      // Who straddled
	SmallBlind     int64
	BigBlind       int64
	Ante           int64
}

// PlayerAction is sent by a player to the FSM.
type PlayerAction struct {
	UserID  int64
	SeatNo  int
	Action  int   // ActionFold / ActionCheck / ...
	Amount  int64 // For bet/raise/call/allin
	TimedAt int64 // Unix ms when action was received
}

// HandEndResult is the settled outcome of a hand, passed to logic layer for DB writes.
type HandEndResult struct {
	GameID         int64
	HandNo         string
	ShuffleSeed    []byte
	CommunityCards []Card
	RunTwiceUsed   bool
	RunTwiceBoard2 []Card
	IsSplitPot     bool
	DurationMs     int
	Players        []*PlayerEndState
	Pots           []PotResult
	Snapshots      []StageSnapshot // One per stage played
}

// PlayerEndState is the per-player result of a hand.
type PlayerEndState struct {
	UserID      int64
	SeatNo      int
	Position    int
	HoleCards   []Card
	HandResult  HandResult
	ForcedBet   int64
	TotalBet    int64
	ChipsStart  int64
	ChipsEnd    int64
	Result      int64 // ChipsEnd - ChipsStart (positive = won)
	IsWinner    bool
	FoldStage   int
	IsVPIP      bool
	IsPFR       bool
	WentToSD    bool
	IsShowCard  bool
	FoldedEarly bool // Folded before showdown
}

// PotResult is the final distribution of one pot.
type PotResult struct {
	PotType  int    // 1=main 2=side
	PotIndex int    // Side pot index
	Amount   int64
	Winners  []PotShare
	WinRank  int    // Winning hand_rank value
	WinDesc  string // Winning hand description
}
