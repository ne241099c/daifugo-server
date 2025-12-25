package game

import "fmt"

type Suit int

const (
	SuitSpade Suit = iota
	SuitHeart
	SuitDiamond
	SuitClub
	SuitJoker
)

func (s Suit) String() string {
	switch s {
	case SuitSpade:
		return "♠"
	case SuitHeart:
		return "♥"
	case SuitDiamond:
		return "♦"
	case SuitClub:
		return "♣"
	case SuitJoker:
		return "Joker"
	default:
		return "?"
	}
}

type Rank int

type Card struct {
	ID   int
	Suit Suit
	Rank Rank
}

func NewCard(id int, suit Suit, rank Rank) *Card {
	return &Card{
		ID:   id,
		Suit: suit,
		Rank: rank,
	}
}

func (c *Card) String() string {
	if c.Suit == SuitJoker {
		return "Joker"
	}

	rankStr := fmt.Sprintf("%d", c.Rank)
	switch c.Rank {
	case 1:
		rankStr = "A"
	case 11:
		rankStr = "J"
	case 12:
		rankStr = "Q"
	case 13:
		rankStr = "K"
	}
	return fmt.Sprintf("%s%s", c.Suit.String(), rankStr)
}
