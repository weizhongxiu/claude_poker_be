package game

import "testing"

// 重现用户报告的牌力误判 bug
func TestBug_TwoCardsFlushStraight(t *testing.T) {
	// 2张连续同花色 → 不应该是同花顺
	h := evalFive(cards("5h 4h"))
	if h.Rank == HandStraightFlush || h.Rank == HandFlush || h.Rank == HandStraight {
		t.Errorf("5h 4h (2 cards) wrongly evaluated as %s (rank=%d), expected High Card", h.Desc, h.Rank)
	}
}

func TestBug_TwoCardsOffsuit_Straight(t *testing.T) {
	// 2张连续不同花色 → 不应该是顺子
	h := evalFive(cards("Ah Kd"))
	if h.Rank == HandStraight || h.Rank == HandStraightFlush {
		t.Errorf("Ah Kd (2 cards) wrongly evaluated as %s, expected High Card", h.Desc)
	}
}

func TestBug_OneCard_NotFlushStraight(t *testing.T) {
	// 1张牌 → 不应该是任何花型
	h := evalFive(cards("Ah"))
	if h.Rank != HandHighCard {
		t.Errorf("Ah (1 card) wrongly evaluated as %s, expected High Card", h.Desc)
	}
}

func TestBug_TwoDifferentCards_NotPair(t *testing.T) {
	// 2张不同点数的牌 → 不应该是对子
	h := evalFive(cards("Th 3d"))
	if h.Rank == HandOnePair {
		t.Errorf("Th 3d (2 different cards) wrongly evaluated as One Pair, expected High Card")
	}
}

func TestBug_EvalBest5_FoldWin_NoBoard(t *testing.T) {
	// fold win 场景：只有2张底牌，0张公共牌
	// 连续同花色的牌不应误判为同花顺
	hole := cards("Ks Qs")
	board := []Card{}
	res := EvalBest5(hole, board)
	if res.Rank >= HandStraight {
		t.Errorf("fold-win (2 hole cards, no board): wrongly got %s (rank=%d), expected High Card", res.Desc, res.Rank)
	}
}
