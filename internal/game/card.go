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

const (
	RankAce   = 1
	RankTwo   = 2
	RankThree = 3
	RankEight = 8
	RankJack  = 11
	RankQueen = 12
	RankKing  = 13
)

type Card struct {
	ID   int  `json:"id"`
	Suit Suit `json:"suit"`
	Rank Rank `json:"rank"`
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
	case RankAce:
		rankStr = "A"
	case RankJack:
		rankStr = "J"
	case RankQueen:
		rankStr = "Q"
	case RankKing:
		rankStr = "K"
	}
	return fmt.Sprintf("%s%s", c.Suit.String(), rankStr)
}
