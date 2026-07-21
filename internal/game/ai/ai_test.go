package ai

import (
	"testing"

	"github.com/ne241099/daifugo-server/internal/game"
)

// AutoPlay が常に合法手を出し（Game.Play/Pass がエラーを返さない）、
// 2〜4人のゲームがきちんと終了することを確認する。
func TestAutoPlayFinishesGames(t *testing.T) {
	m, err := Default()
	if err != nil {
		t.Fatalf("モデル読み込み失敗: %v", err)
	}
	if len(m.W) != len(FeatureNames) {
		t.Fatalf("モデル次元 %d が FeatureNames %d と不一致", len(m.W), len(FeatureNames))
	}

	for _, n := range []int{2, 3, 4} {
		for trial := 0; trial < 300; trial++ {
			ids := make([]int64, n)
			for i := range ids {
				ids[i] = int64(i + 1)
			}
			g := game.NewGame(ids)

			guard := 0
			for !g.IsFinished && guard < 5000 {
				guard++
				uid := g.Players[g.Turn].UserID
				if err := m.AutoPlay(g, uid); err != nil {
					t.Fatalf("n=%d trial=%d: AutoPlay がエラー: %v", n, trial, err)
				}
			}
			if !g.IsFinished {
				t.Fatalf("n=%d trial=%d: ゲームが終了しなかった (guard=%d)", n, trial, guard)
			}
		}
	}
}

// RunBotTurns: 全員 Bot なら最後まで進み、人間の手番があればそこで止まる。
func TestRunBotTurns(t *testing.T) {
	m, err := Default()
	if err != nil {
		t.Fatal(err)
	}

	// 全員 Bot → ゲームが最後まで進む
	g := game.NewGame([]int64{1, 2, 3})
	if err := m.RunBotTurns(g, func(int64) bool { return true }); err != nil {
		t.Fatal(err)
	}
	if !g.IsFinished {
		t.Fatal("全員Botなのにゲームが終了しなかった")
	}

	// userID==1 が人間 → 1 の手番で止まる
	g2 := game.NewGame([]int64{1, 2, 3})
	if err := m.RunBotTurns(g2, func(uid int64) bool { return uid != 1 }); err != nil {
		t.Fatal(err)
	}
	if !g2.IsFinished && g2.Players[g2.Turn].UserID != 1 {
		t.Fatalf("人間(1)の手番で止まるはずが turn=%d", g2.Players[g2.Turn].UserID)
	}
}

// ChooseMove が返す手が legalMoves に含まれる（＝合法）ことを確認する。
func TestChooseMoveIsLegal(t *testing.T) {
	m, err := Default()
	if err != nil {
		t.Fatal(err)
	}
	g := game.NewGame([]int64{1, 2, 3})
	for step := 0; step < 200 && !g.IsFinished; step++ {
		idx := g.Turn
		uid := g.Players[idx].UserID
		cards, pass := m.ChooseMove(g, uid)
		if !pass {
			legal := legalMoves(g.Players[idx].Hand, g.FieldCards, effectiveRev(g))
			want := cardsKey(cards)
			found := false
			for _, mv := range legal {
				if cardsKey(mv) == want {
					found = true
					break
				}
			}
			if !found {
				t.Fatalf("step=%d: ChooseMove が非合法な手を返した", step)
			}
		}
		if err := m.AutoPlay(g, uid); err != nil {
			t.Fatalf("step=%d: %v", step, err)
		}
	}
}
