package game

import "sort"

func cardStrength(c *Card) int {
	if c.Suit == SuitJoker {
		return 99
	}
	if c.Rank == 1 {
		return 11 // A
	}
	if c.Rank == 2 {
		return 12 // 2
	}
	return int(c.Rank) - 3
}

func sortHandForExchange(hand []*Card) {
	sort.Slice(hand, func(i, j int) bool {
		return cardStrength(hand[i]) < cardStrength(hand[j])
	})
}

// hasRank は指定したランクのカードが含まれているかを返す
func hasRank(cards []*Card, rank Rank) bool {
	for _, c := range cards {
		if c.Rank == rank {
			return true
		}
	}
	return false
}
