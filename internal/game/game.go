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
	Players          []*Player
	FinishedPlayers  []*Player
	MiyakoOchiPlayer *Player
	FieldCards       []*Card

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
	g.MiyakoOchiPlayer = nil

	// カード交換
	if g.Players[0].Rank > 0 {
		g.exchangeCards()
	}

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

func (g *Game) handleWin(winner *Player) {
	// 順位リストに追加
	g.FinishedPlayers = append(g.FinishedPlayers, winner)

	// 順位付け
	winner.Rank = len(g.FinishedPlayers)

	// 都落ち判定
	if len(g.FinishedPlayers) == 1 && winner.Rank != 1 {
		for _, p := range g.Players {
			if p.Rank == 1 && len(p.Hand) > 0 {
				// 都落ち発生！
				g.triggerMiyakoOchi(p)
				break
			}
		}
	}

	// ゲーム終了判定
	if g.getActivePlayerCount() == 1 {
		g.finishGame()
	}
}

func (g *Game) triggerMiyakoOchi(loser *Player) {
	// 手札を没収
	loser.Hand = []*Card{}

	// 一時退避
	g.MiyakoOchiPlayer = loser
	// ※ ここでクライアントに都落ち発生を通知するイベントを追加
}

func (g *Game) finishGame() {
	g.IsFinished = true

	// 残っているプレイヤーを探す
	var lastPlayer *Player
	for _, p := range g.Players {
		if len(p.Hand) > 0 {
			lastPlayer = p
			break
		}
	}

	// 残っていた人を追加
	if lastPlayer != nil {
		g.FinishedPlayers = append(g.FinishedPlayers, lastPlayer)
	}

	// 都落ちした人がいれば、リストの最後に追加
	if g.MiyakoOchiPlayer != nil {
		g.FinishedPlayers = append(g.FinishedPlayers, g.MiyakoOchiPlayer)
		g.MiyakoOchiPlayer = nil // リセット
	}

	// 次のゲームのために Rank を確定させる
	for i, p := range g.FinishedPlayers {
		p.Rank = i + 1 // 1位, 2位...
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

func (g *Game) exchangeCards() {
	// プレイヤーをランク順に取得するためのマップ
	rankMap := make(map[int]*Player)
	for _, p := range g.Players {
		rankMap[p.Rank] = p
		sortHandForExchange(p.Hand) // 手札を強さ順にソートしておく
	}

	playerCount := len(g.Players)
	if playerCount < 3 {
		return // 2人以下の場合は交換なし
	}

	// 大富豪<->大貧民
	daifugo := rankMap[1]
	daihinmin := rankMap[playerCount]

	if daifugo != nil && daihinmin != nil {
		// 大富豪
		giveLow := daifugo.Hand[:2]
		// 大貧民
		giveHigh := daihinmin.Hand[len(daihinmin.Hand)-2:]

		// 交換実行
		cardsFromDaifugo := make([]*Card, 2)
		copy(cardsFromDaifugo, giveLow)

		cardsFromDaihinmin := make([]*Card, 2)
		copy(cardsFromDaihinmin, giveHigh)

		// 手札から削除
		daifugo.Hand = daifugo.Hand[2:]
		daihinmin.Hand = daihinmin.Hand[:len(daihinmin.Hand)-2]

		// 手札に追加
		daifugo.Hand = append(daifugo.Hand, cardsFromDaihinmin...)
		daihinmin.Hand = append(daihinmin.Hand, cardsFromDaifugo...)
	}

	// 富豪<->貧民
	if playerCount >= 4 {
		fugo := rankMap[2]
		hinmin := rankMap[playerCount-1]

		if fugo != nil && hinmin != nil {
			giveLow := fugo.Hand[:1]
			giveHigh := hinmin.Hand[len(hinmin.Hand)-1:]

			cardFromFugo := giveLow[0]
			cardFromHinmin := giveHigh[0]

			fugo.Hand = fugo.Hand[1:]
			hinmin.Hand = hinmin.Hand[:len(hinmin.Hand)-1]

			fugo.Hand = append(fugo.Hand, cardFromHinmin)
			hinmin.Hand = append(hinmin.Hand, cardFromFugo)
		}
	}

	// 交換後の手札を再度ソート
	for _, p := range g.Players {
		sortHandForExchange(p.Hand)
	}
}
