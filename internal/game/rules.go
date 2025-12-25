package game

import (
	"errors"
	"sort"
)

type HandType int

const (
	HandTypeInvalid  HandType = iota // 無効
	HandTypeSingle                   // 単騎
	HandTypePair                     // ペア
	HandTypeSequence                 // 階段
)

const (
	RankAce   = 1
	RankTwo   = 2
	RankEight = 8
	RankJack  = 11
	RankThree = 3
)

func AnalyzeHand(cards []*Card, isRev bool) (HandType, int, error) {
	count := len(cards)
	if count == 0 {
		return HandTypeInvalid, 0, errors.New("カードが選択されていません")
	}

	// 1. 単騎
	if count == 1 {
		return HandTypeSingle, GetStrength(cards[0], isRev), nil
	}

	// 2. ペア判定
	if isPair(cards) {
		// ペアの強さは構成カードの強さ（ジョーカー以外）
		str := calculateGroupStrength(cards, isRev)
		return HandTypePair, str, nil
	}

	// 3. 階段判定
	if isSequence(cards) {
		// 階段の強さは「一番弱いカード」または「一番強いカード」で定義できますが、
		// ここでは比較しやすいように「一番強いカードの強さ」を返します
		str := calculateSequenceStrength(cards, isRev)
		return HandTypeSequence, str, nil
	}

	return HandTypeInvalid, 0, errors.New("役として成立していません")
}

func GetStrength(c *Card, isRev bool) int {
	if c.Suit == SuitJoker {
		return 99
	}

	strength := 0

	switch c.Rank {
	case 1: // Ace
		strength = 14
	case 2: // Two
		strength = 15
	default:
		strength = int(c.Rank)
	}

	if isRev {
		strength *= -1
	}

	return strength
}

func IsStronger(candidate, target *Card, isRev bool) bool {
	if target == nil {
		return true
	}

	myStrength := GetStrength(candidate, isRev)
	targetStrength := GetStrength(target, isRev)

	return myStrength > targetStrength
}

func isPair(cards []*Card) bool {
	if len(cards) < 2 {
		return false
	}

	var baseRank Rank = -1
	for _, c := range cards {
		if c.Suit == SuitJoker {
			continue
		}
		if baseRank == -1 {
			baseRank = c.Rank
		} else if baseRank != c.Rank {
			return false
		}
	}
	return true
}

func isSequence(cards []*Card) bool {
	if len(cards) < 3 {
		return false
	}

	// ジョーカーを除いたカードを抽出
	var normals []*Card
	jokers := 0
	for _, c := range cards {
		if c.Suit == SuitJoker {
			jokers++
		} else {
			normals = append(normals, c)
		}
	}

	// ジョーカーのみは階段ではない
	if len(normals) == 0 {
		return false
	}

	// ランク順にソート
	sort.Slice(normals, func(i, j int) bool {
		return normals[i].Rank < normals[j].Rank
	})

	// 連続チェック
	for i := 0; i < len(normals)-1; i++ {
		diff := int(normals[i+1].Rank) - int(normals[i].Rank)

		// 同じ数字は階段不可
		if diff == 0 {
			return false
		}

		// 差が1なら連続
		// 差が2以上なら、その分ジョーカーが必要
		neededJokers := diff - 1
		if neededJokers > 0 {
			jokers -= neededJokers
			if jokers < 0 {
				return false // ジョーカー不足
			}
		}
	}

	return true
}

func calculateGroupStrength(cards []*Card, isRev bool) int {
	for _, c := range cards {
		if c.Suit != SuitJoker {
			return GetStrength(c, isRev)
		}
	}
	// ジョーカーのみの場合
	return GetStrength(cards[0], isRev)
}

func calculateSequenceStrength(cards []*Card, isRev bool) int {
	// 階段の強さ判定は一番強いカードを基準にする
	maxStr := -999
	for _, c := range cards {
		s := GetStrength(c, isRev)
		if s > maxStr {
			maxStr = s
		}
	}
	return maxStr
}

// 特殊効果判定用
func ContainsEight(cards []*Card) bool {
	for _, c := range cards {
		if c.Suit != SuitJoker && c.Rank == RankEight {
			return true
		}
	}
	return false
}

func ContainsJack(cards []*Card) bool {
	for _, c := range cards {
		if c.Suit != SuitJoker && c.Rank == RankJack {
			return true
		}
	}
	return false
}
