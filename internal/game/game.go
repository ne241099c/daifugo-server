package game

import (
	"errors"
	"fmt"
)

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

	FieldCards []*Card

	// 直前の役情報
	LastHandType     HandType
	LastHandStrength int

	LastPlayerID int64

	Turn         int
	IsRevolution bool
	PassCount    int

	IsFinished bool
}

func NewGame(memberIDs []int64) *Game {
	// 初期化処理
	deck := NewDeck(2)
	deck.Shuffle()
	hands := deck.Deal(len(memberIDs))

	players := make([]*Player, len(memberIDs))
	for i, uid := range memberIDs {
		players[i] = &Player{
			UserID: uid,
			Hand:   hands[i],
			Name:   fmt.Sprintf("User%d", uid),
			Rank:   0,
		}
	}

	return &Game{
		Players:    players,
		FieldCards: []*Card{},
		Turn:       0,
	}
}

// Play カードを出す
func (g *Game) Play(userID int64, cards []*Card) error {
	player := g.Players[g.Turn]

	// ターンの確認
	if player.UserID != userID {
		return errors.New("あなたのターンではありません")
	}

	// 手札所有チェック
	if !player.HasCards(cards) {
		return errors.New("持っていないカードが含まれています")
	}

	// 役の解析
	hType, strength, err := AnalyzeHand(cards, g.IsRevolution)
	if err != nil {
		return err
	}

	// ルール判定
	// 場が流れているならチェック不要
	if err := ValidatePlay(g.FieldCards, g.LastHandType, g.LastHandStrength, cards, hType, strength); err != nil {
		return err
	}

	player.RemoveCards(cards)

	// 場の更新
	g.FieldCards = cards
	g.LastHandType = hType
	g.LastHandStrength = strength
	g.LastPlayerID = userID
	g.PassCount = 0

	// 特殊効果
	// 革命
	if len(cards) >= 4 {
		g.IsRevolution = !g.IsRevolution
	}

	// 8切り判定
	is8giri := false
	for _, c := range cards {
		if c.Rank == RankEight {
			is8giri = true
			break
		}
	}

	// あがり判定
	if len(player.Hand) == 0 {
		g.handleWin(player)
		// あがった場合は8切りでもターンを進める
		g.advanceTurn()
		return nil
	}

	// ターン進行
	if is8giri {
		// 場を流して同じ人のターン
		g.clearTable()
		// ターンは進めない
	} else {
		g.advanceTurn()
	}

	if len(player.Hand) == 0 {
		g.handleWin(player)

		// ゲーム終了判定
		if g.IsFinished {
			return nil // 終了
		}

		// あがった場合は次の人へ
		g.advanceTurn()
		return nil
	}

	// 8切りならターン継続、それ以外なら次へ
	if is8giri {
		g.clearTable()
	} else {
		g.advanceTurn()
	}

	return nil
}

func (g *Game) Reset() *Game {
	// デッキの再生成とシャッフル
	deck := NewDeck(2)
	deck.Shuffle()

	// カードを配る
	hands := deck.Deal(len(g.Players))

	// プレイヤー状態のリセット
	for i, p := range g.Players {
		p.Hand = hands[i]
	}

	// ゲーム状態の初期化
	g.FinishedPlayers = []*Player{}
	g.FieldCards = []*Card{}
	g.LastHandType = HandTypeInvalid
	g.LastHandStrength = 0
	g.LastPlayerID = 0
	g.IsRevolution = false
	g.PassCount = 0
	g.IsFinished = false

	// ターンの決定
	// 大富豪から開始
	g.Turn = 0
	for i, p := range g.Players {
		if p.Rank == 1 {
			g.Turn = i
			break
		}
	}

	return g
}

// Pass
func (g *Game) Pass(userID int64) error {
	player := g.Players[g.Turn]
	if player.UserID != userID {
		return errors.New("あなたのターンではありません")
	}

	g.PassCount++
	g.advanceTurn()

	// 全員パス判定 (プレイ人数 - 1)
	activeCount := g.getActivePlayerCount()
	if activeCount > 0 && g.PassCount >= activeCount-1 {
		g.clearTable()
		// パスで流れたら、親（最後にカードを出した人）のターンにする
		// ただし今回は簡易的に、現在手番の人を親としてスタートさせる
	}

	return nil
}

func (g *Game) clearTable() {
	g.FieldCards = []*Card{}
	g.LastHandType = HandTypeInvalid
	g.LastHandStrength = 0
	g.PassCount = 0
}

func (g *Game) advanceTurn() {
	startTurn := g.Turn
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

	activeCount := g.getActivePlayerCount()

	if activeCount == 1 {
		g.IsFinished = true

		// 残った1人を敗者として確定
		for _, loser := range g.Players {
			if len(loser.Hand) > 0 {
				g.FinishedPlayers = append(g.FinishedPlayers, loser)
				loser.Rank = len(g.FinishedPlayers) // 最下位
				break
			}
		}
	}
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
