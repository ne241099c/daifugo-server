package game

import "testing"

// mkCard はテスト用のカード生成ヘルパー
func mkCard(id int, s Suit, r Rank) *Card {
	return &Card{ID: id, Suit: s, Rank: r}
}

func TestSuitString(t *testing.T) {
	tests := []struct {
		suit Suit
		want string
	}{
		{SuitSpade, "♠"},
		{SuitHeart, "♥"},
		{SuitDiamond, "♦"},
		{SuitClub, "♣"},
		{SuitJoker, "Joker"},
		{Suit(99), "?"},
	}
	for _, tt := range tests {
		if got := tt.suit.String(); got != tt.want {
			t.Errorf("Suit(%d).String() = %q, want %q", tt.suit, got, tt.want)
		}
	}
}

func TestCardString(t *testing.T) {
	tests := []struct {
		card *Card
		want string
	}{
		{mkCard(1, SuitJoker, 0), "Joker"},
		{mkCard(2, SuitSpade, RankAce), "♠A"},
		{mkCard(3, SuitSpade, RankJack), "♠J"},
		{mkCard(4, SuitSpade, RankQueen), "♠Q"},
		{mkCard(5, SuitSpade, RankKing), "♠K"},
		{mkCard(6, SuitHeart, 10), "♥10"},
		{mkCard(7, SuitClub, RankThree), "♣3"},
	}
	for _, tt := range tests {
		if got := tt.card.String(); got != tt.want {
			t.Errorf("Card%+v.String() = %q, want %q", tt.card, got, tt.want)
		}
	}
}

func TestCardStrength(t *testing.T) {
	tests := []struct {
		card *Card
		want int
	}{
		{mkCard(1, SuitJoker, 0), 99},
		{mkCard(2, SuitSpade, 1), 11},  // A
		{mkCard(3, SuitSpade, 2), 12},  // 2
		{mkCard(4, SuitSpade, 3), 0},   // 最弱
		{mkCard(5, SuitSpade, 13), 10}, // K
	}
	for _, tt := range tests {
		if got := cardStrength(tt.card); got != tt.want {
			t.Errorf("cardStrength(%s) = %d, want %d", tt.card, got, tt.want)
		}
	}
}

func TestHasRank(t *testing.T) {
	cards := []*Card{
		mkCard(1, SuitSpade, RankThree),
		mkCard(2, SuitHeart, RankEight),
		mkCard(3, SuitJoker, 0),
	}
	if !hasRank(cards, RankEight) {
		t.Error("hasRank should find RankEight")
	}
	if hasRank(cards, RankKing) {
		t.Error("hasRank should not find RankKing")
	}
	if hasRank(nil, RankThree) {
		t.Error("hasRank(nil) should be false")
	}
}

func TestSortHandForExchange(t *testing.T) {
	hand := []*Card{
		mkCard(1, SuitSpade, 2),  // strength 12
		mkCard(2, SuitSpade, 3),  // strength 0
		mkCard(3, SuitJoker, 0),  // strength 99
		mkCard(4, SuitSpade, 1),  // strength 11 (A)
		mkCard(5, SuitSpade, 10), // strength 7
	}
	sortHandForExchange(hand)

	wantOrder := []int{2, 5, 4, 1, 3} // ID order after ascending strength sort
	for i, id := range wantOrder {
		if hand[i].ID != id {
			t.Errorf("sorted[%d].ID = %d, want %d (order: %v)", i, hand[i].ID, id, wantOrder)
		}
	}
}
