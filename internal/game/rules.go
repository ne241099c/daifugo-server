package game

import (
	"errors"
	"fmt"
)

type HandType int

const (
	HandTypeInvalid HandType = iota
	HandTypeSingle
	HandTypePair
	HandTypeSequence
)

// AnalyzeHand は手札の役と強さを判定します
func AnalyzeHand(cards []*Card, isRev bool) (HandType, int, error) {
	count := len(cards)
	if count == 0 {
		return HandTypeInvalid, 0, errors.New("カードが選択されていません")
	}

	// 1. 単騎
	if count == 1 {
		return HandTypeSingle, GetStrength(cards[0], isRev), nil
	}

	// 2. ペア
	if IsPair(cards) {
		str := calculateGroupStrength(cards, isRev)
		return HandTypePair, str, nil
	}

	// 3. 階段
	if IsSequence(cards) {
		str := calculateSequenceStrength(cards, isRev)
		return HandTypeSequence, str, nil
	}

	return HandTypeInvalid, 0, errors.New("役として成立していません")
}

func ValidatePlay(fieldCards []*Card, fieldType HandType, fieldStrength int, playCards []*Card, playType HandType, playStrength int) error {
	// 場に何もないならOK
	if len(fieldCards) == 0 {
		return nil
	}

	// スペ3返し（ジョーカー単騎(最強)に対して、スペードの3）
	// fieldStrengthが最強(99)かつ単騎の場合
	if fieldType == HandTypeSingle && len(fieldCards) == 1 && fieldCards[0].Suit == SuitJoker {
		if len(playCards) == 1 && playCards[0].Suit == SuitSpade && playCards[0].Rank == RankThree {
			return nil // スペ3返し成功
		}
	}

	// 役の種類一致
	if fieldType != playType {
		return fmt.Errorf("場のカードと役の種類が違います")
	}

	// 枚数一致
	if len(fieldCards) != len(playCards) {
		return fmt.Errorf("場のカードと枚数が違います")
	}

	// 強さ比較（同値不可）
	if playStrength <= fieldStrength {
		return fmt.Errorf("場のカードより弱いです")
	}

	return nil
}

func GetStrength(c *Card, isRev bool) int {
	if c.Suit == SuitJoker {
		return 99
	}
	strength := 0
	switch c.Rank {
	case RankAce:
		strength = 13
	case RankTwo:
		strength = 14
	default:
		strength = int(c.Rank) - 2
	}
	if isRev {
		// 簡易反転: 14(2) -> 1, 1(3) -> 14
		// 通常: 3=1, ... 2=14 (range 1-14)
		// 革命: 3=14, ... 2=1
		strength = 15 - strength
	}
	return strength
}

func calculateGroupStrength(cards []*Card, isRev bool) int {
	for _, c := range cards {
		if c.Suit != SuitJoker {
			return GetStrength(c, isRev)
		}
	}
	return 99 // All Jokers
}

func calculateSequenceStrength(cards []*Card, isRev bool) int {
	// 階段の強さは「一番強いカード」とする（元のロジック準拠）
	// ただし革命中は「数字が小さいほうが強い」ので注意が必要だが、
	// GetStrengthが反転しているので、ここでのMax計算は「StrengthのMax」を取ればよい
	maxStr := -1
	for _, c := range cards {
		if c.Suit != SuitJoker {
			s := GetStrength(c, isRev)
			if s > maxStr {
				maxStr = s
			}
		}
	}
	// ジョーカー補正は複雑だが、基本は通常カードのStrengthに依存
	return maxStr
}
