package game

import "testing"

// newTestGame は決定的なゲーム状態を組み立てる（シャッフルを使わない）
func newTestGame(hands ...[]*Card) *Game {
	players := make([]*Player, len(hands))
	for i, h := range hands {
		players[i] = &Player{UserID: int64(i + 1), Hand: h}
	}
	return &Game{
		Players:    players,
		FieldCards: []*Card{},
		Turn:       0,
	}
}

func TestPlay_WrongTurn(t *testing.T) {
	g := newTestGame(
		[]*Card{mkCard(1, SuitSpade, 3)},
		[]*Card{mkCard(2, SuitSpade, 4)},
	)
	// Turn は 0 (UserID 1) なのに UserID 2 が出そうとする
	if err := g.Play(2, []*Card{mkCard(2, SuitSpade, 4)}); err == nil {
		t.Error("他人のターンでの play はエラーになるべき")
	}
}

func TestPlay_CardNotInHand(t *testing.T) {
	g := newTestGame(
		[]*Card{mkCard(1, SuitSpade, 3)},
		[]*Card{mkCard(2, SuitSpade, 4)},
	)
	if err := g.Play(1, []*Card{mkCard(99, SuitHeart, 7)}); err == nil {
		t.Error("手札にないカードでの play はエラーになるべき")
	}
}

func TestPlay_BasicAdvancesTurn(t *testing.T) {
	g := newTestGame(
		[]*Card{mkCard(1, SuitSpade, 3), mkCard(2, SuitSpade, 5)},
		[]*Card{mkCard(3, SuitSpade, 4), mkCard(4, SuitSpade, 6)},
	)
	if err := g.Play(1, []*Card{mkCard(1, SuitSpade, 3)}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(g.FieldCards) != 1 || g.FieldCards[0].ID != 1 {
		t.Errorf("場に出したカードが反映されていない: %+v", g.FieldCards)
	}
	if len(g.Players[0].Hand) != 1 {
		t.Errorf("出したカードが手札から除かれていない: %+v", g.Players[0].Hand)
	}
	if g.Turn != 1 {
		t.Errorf("ターンが進んでいない: Turn = %d, want 1", g.Turn)
	}
}

func TestPlay_EightCutKeepsTurnAndClearsField(t *testing.T) {
	g := newTestGame(
		[]*Card{mkCard(1, SuitSpade, RankEight), mkCard(2, SuitSpade, 5)},
		[]*Card{mkCard(3, SuitSpade, 4), mkCard(4, SuitSpade, 6)},
	)
	if err := g.Play(1, []*Card{mkCard(1, SuitSpade, RankEight)}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(g.FieldCards) != 0 {
		t.Errorf("8切りで場が流れていない: %+v", g.FieldCards)
	}
	if g.Turn != 0 {
		t.Errorf("8切りでターンが継続していない: Turn = %d, want 0", g.Turn)
	}
}

func TestPlay_RevolutionOnFourCards(t *testing.T) {
	g := newTestGame(
		[]*Card{
			mkCard(1, SuitSpade, 3), mkCard(2, SuitHeart, 3),
			mkCard(3, SuitClub, 3), mkCard(4, SuitDiamond, 3),
			mkCard(5, SuitSpade, 5),
		},
		[]*Card{mkCard(6, SuitSpade, 4)},
	)
	four := []*Card{
		mkCard(1, SuitSpade, 3), mkCard(2, SuitHeart, 3),
		mkCard(3, SuitClub, 3), mkCard(4, SuitDiamond, 3),
	}
	if err := g.Play(1, four); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !g.IsRevolution {
		t.Error("4枚出しで革命が発生していない")
	}
}

func TestPlay_WinnerGetsRankAndGameFinishes(t *testing.T) {
	// 2人。UserID 1 が最後の1枚を出してあがる → ゲーム終了
	g := newTestGame(
		[]*Card{mkCard(1, SuitSpade, 3)},
		[]*Card{mkCard(2, SuitSpade, 4), mkCard(3, SuitSpade, 5)},
	)
	if err := g.Play(1, []*Card{mkCard(1, SuitSpade, 3)}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !g.IsFinished {
		t.Error("残り1人になったのにゲームが終了していない")
	}
	if len(g.Players[0].Hand) != 0 {
		t.Error("あがったのに手札が残っている")
	}
	if g.Players[0].Rank != 1 {
		t.Errorf("最初にあがった人の順位が1位でない: Rank = %d", g.Players[0].Rank)
	}
	if g.Players[1].Rank != 2 {
		t.Errorf("最後に残った人の順位が2位でない: Rank = %d", g.Players[1].Rank)
	}
}

func TestPlay_MiyakoOchi_Triggers(t *testing.T) {
	// 前局大富豪(player2)が、今局で player1 に先を越される → 都落ち
	g := newTestGame(
		[]*Card{mkCard(1, SuitSpade, 3)},                          // player1: あがる
		[]*Card{mkCard(2, SuitSpade, 4), mkCard(3, SuitSpade, 5)}, // player2: 前局大富豪
		[]*Card{mkCard(4, SuitSpade, 6), mkCard(5, SuitSpade, 7)}, // player3
	)
	g.Players[1].Rank = 1 // player2 は前局の大富豪

	if err := g.Play(1, []*Card{mkCard(1, SuitSpade, 3)}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(g.Players[1].Hand) != 0 {
		t.Error("都落ちした前局大富豪の手札が没収されていない")
	}
	if g.Players[1].Rank != 3 {
		t.Errorf("都落ちした人の最終順位 = %d, want 3 (最下位=大貧民)", g.Players[1].Rank)
	}
	if g.Players[0].Rank != 1 {
		t.Errorf("最初にあがった人の順位 = %d, want 1", g.Players[0].Rank)
	}
}

func TestPlay_MiyakoOchi_DefenseNoTrigger(t *testing.T) {
	// 前局大富豪(player1)が今局も最初にあがる → 防衛成功で都落ちなし
	g := newTestGame(
		[]*Card{mkCard(1, SuitSpade, 3)},                          // player1: 前局大富豪、あがる
		[]*Card{mkCard(2, SuitSpade, 4), mkCard(3, SuitSpade, 5)}, // player2
		[]*Card{mkCard(4, SuitSpade, 6), mkCard(5, SuitSpade, 7)}, // player3
	)
	g.Players[0].Rank = 1 // player1 が前局大富豪

	if err := g.Play(1, []*Card{mkCard(1, SuitSpade, 3)}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 誰の手札も都落ちで没収されていないこと
	if len(g.Players[1].Hand) == 0 || len(g.Players[2].Hand) == 0 {
		t.Error("防衛成功のはずが都落ちが発生している")
	}
	if g.IsFinished {
		t.Error("まだ2人が手札を持っているのにゲームが終了している")
	}
}

func TestPlay_MiyakoOchi_FirstGameNoTrigger(t *testing.T) {
	// 初回ゲーム(全員 Rank==0) では都落ちは発生しない
	g := newTestGame(
		[]*Card{mkCard(1, SuitSpade, 3)},
		[]*Card{mkCard(2, SuitSpade, 4), mkCard(3, SuitSpade, 5)},
		[]*Card{mkCard(4, SuitSpade, 6), mkCard(5, SuitSpade, 7)},
	)

	if err := g.Play(1, []*Card{mkCard(1, SuitSpade, 3)}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(g.Players[1].Hand) == 0 || len(g.Players[2].Hand) == 0 {
		t.Error("初回ゲームで都落ちが誤発生している")
	}
}

func TestPass_IncrementsAndAdvances(t *testing.T) {
	g := newTestGame(
		[]*Card{mkCard(1, SuitSpade, 3)},
		[]*Card{mkCard(2, SuitSpade, 4)},
		[]*Card{mkCard(3, SuitSpade, 5)},
	)
	// 場にカードがある状態にしておく
	g.FieldCards = []*Card{mkCard(9, SuitHeart, 7)}
	g.LastHandType = HandTypeSingle
	g.LastHandStrength = 5

	if err := g.Pass(1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if g.PassCount != 1 {
		t.Errorf("PassCount = %d, want 1", g.PassCount)
	}
	if g.Turn != 1 {
		t.Errorf("Turn = %d, want 1", g.Turn)
	}
}

func TestPass_WrongTurn(t *testing.T) {
	g := newTestGame(
		[]*Card{mkCard(1, SuitSpade, 3)},
		[]*Card{mkCard(2, SuitSpade, 4)},
	)
	if err := g.Pass(2); err == nil {
		t.Error("他人のターンでの pass はエラーになるべき")
	}
}

func TestRemovePlayer(t *testing.T) {
	g := newTestGame(
		[]*Card{mkCard(1, SuitSpade, 3), mkCard(2, SuitSpade, 5)},
		[]*Card{mkCard(3, SuitSpade, 4), mkCard(4, SuitSpade, 6)},
		[]*Card{mkCard(5, SuitSpade, 7), mkCard(6, SuitSpade, 8)},
	)
	g.RemovePlayer(2)

	for _, p := range g.Players {
		if p.UserID == 2 {
			t.Error("退出したプレイヤーが Players に残っている")
		}
	}
	if len(g.Players) != 2 {
		t.Errorf("Players 数 = %d, want 2", len(g.Players))
	}
}
