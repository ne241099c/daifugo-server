package game

import (
	"math/rand"
	"time"
)

type Deck []*Card

func NewDeck(jokerCount int) Deck {
	rand.Seed(time.Now().UnixNano()) // シード初期化を追加（念の為）
	d := make(Deck, 0, 52+jokerCount)
	idCounter := 1

	for s := SuitSpade; s <= SuitClub; s++ {
		for r := 1; r <= 13; r++ {
			d = append(d, NewCard(idCounter, s, Rank(r)))
			idCounter++
		}
	}

	for i := 0; i < jokerCount; i++ {
		d = append(d, NewCard(idCounter, SuitJoker, 0))
		idCounter++
	}
	return d
}

func (d Deck) Shuffle() {
	rand.Shuffle(len(d), func(i, j int) {
		d[i], d[j] = d[j], d[i]
	})
}

func (d Deck) Deal(numPlayers int) [][]*Card {
	hands := make([][]*Card, numPlayers)
	for i, card := range d {
		playerIndex := i % numPlayers
		hands[playerIndex] = append(hands[playerIndex], card)
	}
	return hands
}
