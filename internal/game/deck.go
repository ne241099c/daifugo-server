package game

import (
	"math/rand"
)

type Deck []*Card

func NewDeck() Deck {
	d := make(Deck, 0, 54)

	for s := SuitSpade; s <= SuitClub; s++ {
		for r := 1; r <= 13; r++ {
			d = append(d, NewCard(s, Rank(r)))
		}
	}

	for range 2 {
		d = append(d, NewCard(SuitJoker, 0))
	}
	return d
}

func (d Deck) Shuffle() {
	// ランダムに並び替える
	rand.Shuffle(len(d), func(i, j int) {
		d[i], d[j] = d[j], d[i]
	})
}

func (d Deck) Deal(numPlayers int) [][]*Card {
	hands := make([][]*Card, numPlayers)

	// 山札のカードを1枚ずつ配っていく
	for i, card := range d {
		playerIndex := i % numPlayers
		hands[playerIndex] = append(hands[playerIndex], card)
	}

	return hands
}
