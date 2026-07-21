package ai

import (
	"sort"
	"strconv"
	"strings"

	"github.com/ne241099/daifugo-server/internal/game"
)

// FeatureNames は特徴量の並び（Python 側の FEATURE_NAMES と完全一致させること）。
var FeatureNames = []string{
	// --- 手ごとに変わる特徴 ---
	"is_pass", "played", "hand_after", "jokers_after", "twos_after", "strong_after",
	"uses_joker", "move_has_two", "is_8giri", "is_revolution",
	"type_single", "type_pair", "type_seq", "move_strength",
	// --- 状態特徴 ---
	"hand_size", "field_strength", "is_lead", "eff_rev",
	"min_opp_hand", "opp_le2", "n_active_opp", "sum_opp_hand", "discard_size",
	// --- 状態 × 手 の交互作用 ---
	"minopp_x_str", "oppLe2_x_str", "oppLe2_x_strongmove", "oppLe2_x_pass",
	"minopp_x_pass", "discard_x_str",
}

// effectiveRev は現在の実効的な革命状態（革命 XOR 11バック）を返す。
func effectiveRev(g *game.Game) bool {
	is11 := false
	for _, c := range g.FieldCards {
		if c.Rank == game.RankJack {
			is11 = true
			break
		}
	}
	return g.IsRevolution != is11
}

// legalMoves は手札から出せる手（カードの組）をすべて列挙する。
// Python の legal_moves と同じ（単騎・ペア系・階段を候補にして役成立＆場に出せるものを残す）。
func legalMoves(hand []*game.Card, field []*game.Card, isRev bool) [][]*game.Card {
	var jokers []*game.Card
	for _, c := range hand {
		if c.Suit == game.SuitJoker {
			jokers = append(jokers, c)
		}
	}

	var cands [][]*game.Card
	for _, c := range hand { // 単騎
		cands = append(cands, []*game.Card{c})
	}
	cands = append(cands, sameRankGroups(hand, jokers)...) // ペア系
	cands = append(cands, sequences(hand, jokers)...)      // 階段

	var fType game.HandType
	var fStr int
	if len(field) > 0 {
		fType, fStr, _ = game.AnalyzeHand(field, isRev)
	}

	var out [][]*game.Card
	seen := map[string]bool{}
	for _, cs := range cands {
		pType, pStr, err := game.AnalyzeHand(cs, isRev)
		if err != nil || pType == game.HandTypeInvalid {
			continue
		}
		if game.ValidatePlay(field, fType, fStr, cs, pType, pStr) != nil {
			continue
		}
		key := cardsKey(cs)
		if !seen[key] {
			seen[key] = true
			out = append(out, cs)
		}
	}
	return out
}

func sameRankGroups(hand, jokers []*game.Card) [][]*game.Card {
	byRank := map[int][]*game.Card{}
	order := []int{}
	for _, c := range hand {
		if c.Suit == game.SuitJoker {
			continue
		}
		r := int(c.Rank)
		if _, ok := byRank[r]; !ok {
			order = append(order, r)
		}
		byRank[r] = append(byRank[r], c)
	}
	var out [][]*game.Card
	for _, r := range order {
		cs := byRank[r]
		for k := 2; k <= len(cs); k++ { // 純粋なペア/トリオ/…
			out = append(out, append([]*game.Card{}, cs[:k]...))
		}
		for jk := 1; jk <= len(jokers); jk++ { // ジョーカーで枚数を伸ばす
			if len(cs) >= 1 {
				g := append([]*game.Card{}, cs...)
				g = append(g, jokers[:jk]...)
				out = append(out, g)
			}
		}
	}
	return out
}

func sequences(hand, jokers []*game.Card) [][]*game.Card {
	bySuit := map[game.Suit]map[int]*game.Card{}
	var suitOrder []game.Suit
	for _, c := range hand {
		if c.Suit == game.SuitJoker {
			continue
		}
		if _, ok := bySuit[c.Suit]; !ok {
			bySuit[c.Suit] = map[int]*game.Card{}
			suitOrder = append(suitOrder, c.Suit)
		}
		bySuit[c.Suit][int(c.Rank)] = c
	}
	nj := len(jokers)
	var out [][]*game.Card
	for _, suit := range suitOrder {
		ranks := bySuit[suit]
		for start := 1; start <= 13; start++ {
			for length := 3; length <= 13; length++ {
				end := start + length - 1
				if end > 13 {
					break
				}
				var cards []*game.Card
				need, present := 0, 0
				for r := start; r <= end; r++ {
					if c, ok := ranks[r]; ok {
						cards = append(cards, c)
						present++
					} else {
						need++
					}
				}
				if present > 0 && need <= nj {
					g := append([]*game.Card{}, cards...)
					g = append(g, jokers[:need]...)
					out = append(out, g)
				}
			}
		}
	}
	return out
}

func cardsKey(cs []*game.Card) string {
	ids := make([]int, len(cs))
	for i, c := range cs {
		ids[i] = c.ID
	}
	sort.Ints(ids)
	parts := make([]string, len(ids))
	for i, id := range ids {
		parts[i] = strconv.Itoa(id)
	}
	return strings.Join(parts, ",")
}

// extractFeatures は「局面 g で playerIdx が move を出す（isPass ならパス）」ときの
// 29 次元の特徴ベクトルを返す。Python の extract_features と完全一致させること。
func extractFeatures(g *game.Game, playerIdx int, move []*game.Card, isPass bool) []float64 {
	hand := g.Players[playerIdx].Hand
	rev := effectiveRev(g)

	mType := game.HandTypeSingle // パス時は Python 実装に合わせて Single 扱い
	mStr := 0
	usesJoker, moveHasTwo, is8, isRevMove, played := 0, 0, 0, 0, 0
	var remaining []*game.Card
	if isPass {
		remaining = hand
	} else {
		ids := map[int]bool{}
		for _, c := range move {
			ids[c.ID] = true
		}
		for _, c := range hand {
			if !ids[c.ID] {
				remaining = append(remaining, c)
			}
		}
		mType, mStr, _ = game.AnalyzeHand(move, rev)
		for _, c := range move {
			if c.Suit == game.SuitJoker {
				usesJoker = 1
			}
			if c.Rank == game.RankTwo {
				moveHasTwo = 1
			}
			if c.Rank == game.RankEight {
				is8 = 1
			}
		}
		if len(move) >= 4 {
			isRevMove = 1
		}
		played = len(move)
	}

	jokersAfter, twosAfter := 0, 0
	for _, c := range remaining {
		if c.Suit == game.SuitJoker {
			jokersAfter++
		}
		if c.Rank == game.RankTwo {
			twosAfter++
		}
	}

	fieldStr := 0
	if len(g.FieldCards) > 0 {
		_, fieldStr, _ = game.AnalyzeHand(g.FieldCards, rev)
	}

	minOpp, oppLe2, nOpp, sumOpp := 0, 0, 0, 0
	firstOpp := true
	for i, p := range g.Players {
		if i == playerIdx || len(p.Hand) == 0 {
			continue
		}
		s := len(p.Hand)
		nOpp++
		sumOpp += s
		if s <= 2 {
			oppLe2++
		}
		if firstOpp || s < minOpp {
			minOpp = s
			firstOpp = false
		}
	}

	clamp := func(s int) float64 {
		if s > 15 {
			s = 15
		}
		return float64(s) / 15.0
	}
	mstrN := clamp(mStr)
	fieldN := clamp(fieldStr)
	strongMove := usesJoker + moveHasTwo
	isLead := 0
	if len(g.FieldCards) == 0 {
		isLead = 1
	}
	passF := 0.0
	if isPass {
		passF = 1.0
	}
	discard := float64(len(g.DiscardPile))

	return []float64{
		passF,
		float64(played),
		float64(len(hand) - played),
		float64(jokersAfter),
		float64(twosAfter),
		float64(jokersAfter + twosAfter),
		float64(usesJoker),
		float64(moveHasTwo),
		float64(is8),
		float64(isRevMove),
		boolf(mType == game.HandTypeSingle),
		boolf(mType == game.HandTypePair),
		boolf(mType == game.HandTypeSequence),
		mstrN,
		float64(len(hand)),
		fieldN,
		float64(isLead),
		boolf(rev),
		float64(minOpp),
		float64(oppLe2),
		float64(nOpp),
		float64(sumOpp),
		discard,
		float64(minOpp) * mstrN,
		float64(oppLe2) * mstrN,
		float64(oppLe2) * float64(strongMove),
		float64(oppLe2) * passF,
		float64(minOpp) * passF,
		discard * mstrN,
	}
}

func boolf(b bool) float64 {
	if b {
		return 1
	}
	return 0
}
