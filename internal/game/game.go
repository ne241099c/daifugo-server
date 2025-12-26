package game

import (
	"errors"
	"fmt"
)

// Player 構造体
// システム整合性のため ID -> UserID に変更
type Player struct {
	UserID int64   `json:"user_id"`
	Hand   []*Card `json:"hand"`
	Name   string  `json:"name"`
	Rank   int     `json:"rank"`
}

// HasCards 手札チェック
func (p *Player) HasCards(cards []*Card) bool {
	for _, target := range cards {
		found := false
		for _, hand := range p.Hand {
			if hand.ID == target.ID {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func (p *Player) RemoveCards(cards []*Card) {
	var newHand []*Card
	for _, hand := range p.Hand {
		keep := true
		for _, target := range cards {
			if hand.ID == target.ID {
				keep = false
				break
			}
		}
		if keep {
			newHand = append(newHand, hand)
		}
	}
	p.Hand = newHand
}

// Game 構造体
type Game struct {
	Players         []*Player
	FinishedPlayers []*Player

	// システム整合性のため TableCards -> FieldCards に変更
	FieldCards []*Card

	// 直前の役情報（判定用）
	LastHandType     HandType
	LastHandStrength int

	LastPlayerID int64

	// システム整合性のため TurnIndex -> Turn に変更
	Turn         int
	IsRevolution bool
	PassCount    int
}

func NewGame(memberIDs []int64) *Game {
	// 初期化処理
	deck := NewDeck(2)
	deck.Shuffle()
	hands := deck.Deal(len(memberIDs))

	players := make([]*Player, len(memberIDs))
	for i, uid := range memberIDs {
		players[i] = &Player{
			UserID: uid, // ID -> UserID
			Hand:   hands[i],
			Name:   fmt.Sprintf("User%d", uid),
			Rank:   0,
		}
	}

	return &Game{
		Players:    players,
		FieldCards: []*Card{}, // TableCards -> FieldCards
		Turn:       0,         // TurnIndex -> Turn
	}
}

// Play カードを出す
func (g *Game) Play(userID int64, cards []*Card) error {
	player := g.Players[g.Turn] // TurnIndex -> Turn

	// 1. ターンの確認
	if player.UserID != userID { // ID -> UserID
		return errors.New("あなたのターンではありません")
	}

	// 2. 手札所有チェック
	if !player.HasCards(cards) {
		return errors.New("持っていないカードが含まれています")
	}

	// 3. 役の解析
	hType, strength, err := AnalyzeHand(cards, g.IsRevolution)
	if err != nil {
		return err
	}

	// 4. ルール判定 (ValidatePlay)
	// 場が流れている(FieldCardsが空)ならチェック不要
	// TableCards -> FieldCards
	if err := ValidatePlay(g.FieldCards, g.LastHandType, g.LastHandStrength, cards, hType, strength); err != nil {
		return err
	}

	// --- 実行 ---
	player.RemoveCards(cards)

	// 場の更新
	g.FieldCards = cards // TableCards -> FieldCards
	g.LastHandType = hType
	g.LastHandStrength = strength
	g.LastPlayerID = userID
	g.PassCount = 0

	// 5. 特殊効果
	// 革命 (4枚以上)
	if len(cards) >= 4 {
		g.IsRevolution = !g.IsRevolution
	}

	// 8切り判定
	is8giri := false
	for _, c := range cards {
		if c.Rank == RankEight { // RankEightはcard.go/rules.goの定義に依存
			is8giri = true
			break
		}
	}

	// 6. あがり判定
	if len(player.Hand) == 0 {
		g.handleWin(player)
		// あがった場合は8切りでもターンを進める
		g.advanceTurn()
		return nil
	}

	// 7. ターン進行
	if is8giri {
		// 場を流して同じ人のターン
		g.clearTable()
		// ターンは進めない
	} else {
		g.advanceTurn()
	}

	return nil
}

// Pass パス
func (g *Game) Pass(userID int64) error {
	player := g.Players[g.Turn]  // TurnIndex -> Turn
	if player.UserID != userID { // ID -> UserID
		return errors.New("あなたのターンではありません")
	}

	g.PassCount++
	g.advanceTurn()

	// 全員パス判定 (プレイ人数 - 1)
	activeCount := g.getActivePlayerCount()
	if activeCount > 0 && g.PassCount >= activeCount-1 {
		g.clearTable()
		// パスで流れたら、親（最後にカードを出した人）のターンにする
		// ただし今回は簡易的に、現在手番の人（advanceTurn後の人）を親としてスタートさせる
	}

	return nil
}

func (g *Game) clearTable() {
	g.FieldCards = []*Card{} // TableCards -> FieldCards
	g.LastHandType = HandTypeInvalid
	g.LastHandStrength = 0
	g.PassCount = 0
}

func (g *Game) advanceTurn() {
	startTurn := g.Turn // TurnIndex -> Turn
	for {
		g.Turn++
		if g.Turn >= len(g.Players) {
			g.Turn = 0
		}
		// まだあがっていない人ならOK
		if len(g.Players[g.Turn].Hand) > 0 {
			break
		}
		// 一周したら終了
		if g.Turn == startTurn {
			break
		}
	}
}

func (g *Game) handleWin(p *Player) {
	g.FinishedPlayers = append(g.FinishedPlayers, p)
	p.Rank = len(g.FinishedPlayers)
}

func (g *Game) getActivePlayerCount() int {
	c := 0
	for _, p := range g.Players {
		if len(p.Hand) > 0 {
			c++
		}
	}
	return c
}
