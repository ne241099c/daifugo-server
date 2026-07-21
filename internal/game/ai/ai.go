package ai

import (
	"fmt"
	"math"

	"github.com/ne241099/daifugo-server/internal/game"
)

// ChooseMove は userID のプレイヤーの手番として、SVM のスコアが最大になる手を選ぶ。
// pass=true のときは cards は nil（パスを選んだ）。
func (m *Model) ChooseMove(g *game.Game, userID int64) (cards []*game.Card, pass bool) {
	idx := -1
	for i, p := range g.Players {
		if p.UserID == userID {
			idx = i
			break
		}
	}
	if idx < 0 {
		return nil, true
	}

	rev := effectiveRev(g)
	moves := legalMoves(g.Players[idx].Hand, g.FieldCards, rev)

	type action struct {
		cards  []*game.Card
		isPass bool
	}
	acts := make([]action, 0, len(moves)+1)
	for _, mv := range moves {
		acts = append(acts, action{cards: mv})
	}
	// パスは「場がある」または「出せる手が無い」ときだけ候補に入れる。
	if len(g.FieldCards) > 0 || len(moves) == 0 {
		acts = append(acts, action{isPass: true})
	}
	if len(acts) == 0 {
		return nil, true
	}

	best, bestScore := 0, math.Inf(-1)
	for i, a := range acts {
		sc := m.score(extractFeatures(g, idx, a.cards, a.isPass))
		if sc > bestScore {
			bestScore = sc
			best = i
		}
	}
	return acts[best].cards, acts[best].isPass
}

// AutoPlay は userID の手番を SVM に判断させ、その手をゲームに適用する。
// CPU の手番でこれを呼べば、学習済み AI として着手する。
func (m *Model) AutoPlay(g *game.Game, userID int64) error {
	cards, pass := m.ChooseMove(g, userID)
	if pass {
		return g.Pass(userID)
	}
	if len(cards) == 0 {
		return fmt.Errorf("ai: 出す手が決まりませんでした")
	}
	return g.Play(userID, cards)
}

// RunBotTurns は、現在の手番プレイヤーが Bot である間、AI に自動で着手させ続ける。
// isBot(userID) が false（＝人間の手番）になるか、ゲームが終了したら止まる。
//
// 使い方: 人間が PlayCard / Pass / StartGame した後にこれを呼ぶと、
// 続く CPU の手番が自動で進む。
func (m *Model) RunBotTurns(g *game.Game, isBot func(userID int64) bool) error {
	guard := 0
	for !g.IsFinished && guard < 10000 {
		guard++
		uid := g.Players[g.Turn].UserID
		if !isBot(uid) {
			return nil
		}
		if err := m.AutoPlay(g, uid); err != nil {
			return err
		}
	}
	return nil
}
