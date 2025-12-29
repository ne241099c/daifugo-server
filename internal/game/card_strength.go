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
