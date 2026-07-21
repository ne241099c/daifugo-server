package game

import "testing"

func TestNewDeck(t *testing.T) {
	tests := []struct {
		jokerCount int
		wantLen    int
	}{
		{0, 52},
		{1, 53},
		{2, 54},
	}
	for _, tt := range tests {
		d := NewDeck(tt.jokerCount)
		if len(d) != tt.wantLen {
			t.Errorf("NewDeck(%d) len = %d, want %d", tt.jokerCount, len(d), tt.wantLen)
		}

		// ジョーカー枚数の確認
		jokers := 0
		ids := make(map[int]bool)
		for _, c := range d {
			if c.Suit == SuitJoker {
				jokers++
			}
			if ids[c.ID] {
				t.Errorf("duplicate card ID %d in NewDeck(%d)", c.ID, tt.jokerCount)
			}
			ids[c.ID] = true
		}
		if jokers != tt.jokerCount {
			t.Errorf("NewDeck(%d) has %d jokers, want %d", tt.jokerCount, jokers, tt.jokerCount)
		}
	}
}

func TestDeckDeal(t *testing.T) {
	d := NewDeck(2) // 54 枚
	hands := d.Deal(3)

	if len(hands) != 3 {
		t.Fatalf("Deal(3) returned %d hands, want 3", len(hands))
	}

	total := 0
	seen := make(map[int]bool)
	for i, h := range hands {
		if len(h) != 18 {
			t.Errorf("hand[%d] has %d cards, want 18", i, len(h))
		}
		for _, c := range h {
			if seen[c.ID] {
				t.Errorf("card ID %d dealt to multiple players", c.ID)
			}
			seen[c.ID] = true
			total++
		}
	}
	if total != 54 {
		t.Errorf("total dealt cards = %d, want 54 (no card lost)", total)
	}
}
