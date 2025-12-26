package game

import "sort"

func IsPair(cards []*Card) bool {
	if len(cards) < 2 {
		return false
	}
	baseRank := -1
	for _, c := range cards {
		if c.Suit == SuitJoker {
			continue
		}
		if baseRank == -1 {
			baseRank = int(c.Rank)
		} else {
			if baseRank != int(c.Rank) {
				return false
			}
		}
	}
	return true
}

func toSeqRank(c *Card) int {
	return int(c.Rank)
}

func IsSequence(cards []*Card) bool {
	if len(cards) < 3 {
		return false
	}
	var normalCards []*Card
	var baseSuit Suit = -1
	jokerCount := 0

	for _, c := range cards {
		if c.Suit == SuitJoker {
			jokerCount++
			continue
		}
		if baseSuit == -1 {
			baseSuit = c.Suit
		} else if baseSuit != c.Suit {
			return false // スート不一致
		}
		normalCards = append(normalCards, c)
	}

	// ジョーカーのみは階段不可
	if len(normalCards) == 0 {
		return false
	}

	sort.Slice(normalCards, func(i, j int) bool {
		return toSeqRank(normalCards[i]) < toSeqRank(normalCards[j])
	})

	for i := 0; i < len(normalCards)-1; i++ {
		current := toSeqRank(normalCards[i])
		next := toSeqRank(normalCards[i+1])
		diff := next - current

		if diff == 0 {
			return false // 同じランク
		}
		// 差が1なら連続。差が2以上ならその分ジョーカーが必要
		needed := diff - 1
		if needed > 0 {
			jokerCount -= needed
			if jokerCount < 0 {
				return false // ジョーカー不足
			}
		}
	}
	return true
}
