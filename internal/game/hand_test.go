package game

import "testing"

func TestIsPair(t *testing.T) {
	tests := []struct {
		name  string
		cards []*Card
		want  bool
	}{
		{"1枚はペア不成立", []*Card{mkCard(1, SuitSpade, 3)}, false},
		{"同ランク2枚", []*Card{mkCard(1, SuitSpade, 3), mkCard(2, SuitHeart, 3)}, true},
		{"同ランク3枚", []*Card{mkCard(1, SuitSpade, 3), mkCard(2, SuitHeart, 3), mkCard(3, SuitClub, 3)}, true},
		{"異ランクは不成立", []*Card{mkCard(1, SuitSpade, 3), mkCard(2, SuitHeart, 4)}, false},
		{"ジョーカー混じりの同ランク", []*Card{mkCard(1, SuitSpade, 3), mkCard(2, SuitHeart, 3), mkCard(3, SuitJoker, 0)}, true},
	}
	for _, tt := range tests {
		if got := IsPair(tt.cards); got != tt.want {
			t.Errorf("%s: IsPair() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestIsSequence(t *testing.T) {
	tests := []struct {
		name  string
		cards []*Card
		want  bool
	}{
		{"2枚は階段不成立", []*Card{mkCard(1, SuitSpade, 3), mkCard(2, SuitSpade, 4)}, false},
		{"同スート連続3枚", []*Card{mkCard(1, SuitSpade, 3), mkCard(2, SuitSpade, 4), mkCard(3, SuitSpade, 5)}, true},
		{"スート不一致は不成立", []*Card{mkCard(1, SuitSpade, 3), mkCard(2, SuitHeart, 4), mkCard(3, SuitSpade, 5)}, false},
		{"間が空くとジョーカーなしでは不成立", []*Card{mkCard(1, SuitSpade, 3), mkCard(2, SuitSpade, 4), mkCard(3, SuitSpade, 6)}, false},
		{"ジョーカーで間を埋めれば成立", []*Card{mkCard(1, SuitSpade, 3), mkCard(2, SuitSpade, 4), mkCard(3, SuitSpade, 6), mkCard(4, SuitJoker, 0)}, true},
		{"同ランク重複は不成立", []*Card{mkCard(1, SuitSpade, 3), mkCard(2, SuitSpade, 3), mkCard(3, SuitSpade, 4)}, false},
	}
	for _, tt := range tests {
		if got := IsSequence(tt.cards); got != tt.want {
			t.Errorf("%s: IsSequence() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
