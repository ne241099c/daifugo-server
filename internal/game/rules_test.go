package game

import "testing"

func TestGetStrength(t *testing.T) {
	tests := []struct {
		name  string
		card  *Card
		isRev bool
		want  int
	}{
		{"3は最弱", mkCard(1, SuitSpade, RankThree), false, 1},
		{"Kは11", mkCard(2, SuitSpade, RankKing), false, 11},
		{"Aは13", mkCard(3, SuitSpade, RankAce), false, 13},
		{"2は最強14", mkCard(4, SuitSpade, RankTwo), false, 14},
		{"ジョーカーは99", mkCard(5, SuitJoker, 0), false, 99},
		{"革命時: 3が強くなる", mkCard(6, SuitSpade, RankThree), true, 14},
		{"革命時: 2が弱くなる", mkCard(7, SuitSpade, RankTwo), true, 1},
		{"革命時: Aは2", mkCard(8, SuitSpade, RankAce), true, 2},
		{"革命時でもジョーカーは99", mkCard(9, SuitJoker, 0), true, 99},
	}
	for _, tt := range tests {
		if got := GetStrength(tt.card, tt.isRev); got != tt.want {
			t.Errorf("%s: GetStrength() = %d, want %d", tt.name, got, tt.want)
		}
	}
}

func TestAnalyzeHand(t *testing.T) {
	tests := []struct {
		name     string
		cards    []*Card
		isRev    bool
		wantType HandType
		wantStr  int
		wantErr  bool
	}{
		{"空はエラー", []*Card{}, false, HandTypeInvalid, 0, true},
		{"単騎", []*Card{mkCard(1, SuitSpade, 5)}, false, HandTypeSingle, 3, false},
		{"ペア", []*Card{mkCard(1, SuitSpade, 5), mkCard(2, SuitHeart, 5)}, false, HandTypePair, 3, false},
		{"階段", []*Card{mkCard(1, SuitSpade, 3), mkCard(2, SuitSpade, 4), mkCard(3, SuitSpade, 5)}, false, HandTypeSequence, 3, false},
		{"役不成立はエラー", []*Card{mkCard(1, SuitSpade, 3), mkCard(2, SuitHeart, 5)}, false, HandTypeInvalid, 0, true},
	}
	for _, tt := range tests {
		gotType, gotStr, err := AnalyzeHand(tt.cards, tt.isRev)
		if tt.wantErr {
			if err == nil {
				t.Errorf("%s: expected error, got nil", tt.name)
			}
			continue
		}
		if err != nil {
			t.Errorf("%s: unexpected error: %v", tt.name, err)
			continue
		}
		if gotType != tt.wantType {
			t.Errorf("%s: type = %v, want %v", tt.name, gotType, tt.wantType)
		}
		if gotStr != tt.wantStr {
			t.Errorf("%s: strength = %d, want %d", tt.name, gotStr, tt.wantStr)
		}
	}
}

func TestValidatePlay(t *testing.T) {
	single := func(c *Card) []*Card { return []*Card{c} }

	t.Run("場が空なら何でも出せる", func(t *testing.T) {
		err := ValidatePlay(nil, HandTypeInvalid, 0, single(mkCard(1, SuitSpade, 3)), HandTypeSingle, 1)
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}
	})

	t.Run("より強ければ出せる", func(t *testing.T) {
		field := single(mkCard(1, SuitSpade, 5)) // strength 3
		play := single(mkCard(2, SuitHeart, 7))  // strength 5
		err := ValidatePlay(field, HandTypeSingle, 3, play, HandTypeSingle, 5)
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}
	})

	t.Run("同値は出せない", func(t *testing.T) {
		field := single(mkCard(1, SuitSpade, 5))
		play := single(mkCard(2, SuitHeart, 5))
		err := ValidatePlay(field, HandTypeSingle, 3, play, HandTypeSingle, 3)
		if err == nil {
			t.Error("expected error for equal strength")
		}
	})

	t.Run("弱いと出せない", func(t *testing.T) {
		field := single(mkCard(1, SuitSpade, 7))
		play := single(mkCard(2, SuitHeart, 5))
		err := ValidatePlay(field, HandTypeSingle, 5, play, HandTypeSingle, 3)
		if err == nil {
			t.Error("expected error for weaker card")
		}
	})

	t.Run("役の種類が違うと出せない", func(t *testing.T) {
		field := single(mkCard(1, SuitSpade, 5))
		play := []*Card{mkCard(2, SuitHeart, 7), mkCard(3, SuitClub, 7)}
		err := ValidatePlay(field, HandTypeSingle, 3, play, HandTypePair, 5)
		if err == nil {
			t.Error("expected error for type mismatch")
		}
	})

	t.Run("枚数が違うと出せない", func(t *testing.T) {
		field := []*Card{mkCard(1, SuitSpade, 5), mkCard(2, SuitHeart, 5)}
		play := []*Card{mkCard(3, SuitClub, 7), mkCard(4, SuitDiamond, 7), mkCard(5, SuitSpade, 7)}
		err := ValidatePlay(field, HandTypePair, 3, play, HandTypePair, 5)
		if err == nil {
			t.Error("expected error for count mismatch")
		}
	})

	t.Run("スペ3返し: ジョーカー単騎にスペード3で返せる", func(t *testing.T) {
		field := single(mkCard(1, SuitJoker, 0))
		play := single(mkCard(2, SuitSpade, RankThree))
		err := ValidatePlay(field, HandTypeSingle, 99, play, HandTypeSingle, 1)
		if err != nil {
			t.Errorf("スペ3返しは成功するはず, got %v", err)
		}
	})
}
