package game

import (
	"errors"
	"fmt"
)

// Player はゲーム中のプレイヤー状態
type Player struct {
	UserID int64
	Hand   []*Card
	Rank   int
}

// Game はゲーム全体の進行状態
type Game struct {
	Players         []*Player
	FinishedPlayers []*Player // あがった人のリスト
	FieldCards      []*Card   // 場に出ているカード
	Turn            int       // 現在のターンのプレイヤーインデックス

	// ルール状態
	IsRevolution bool
	Is11Back     bool
	PassCount    int // 全員パス判定用
}

// NewGame は新しいゲームを開始
func NewGame(userIDs []int64) *Game {
	deck := NewDeck()
	deck.Shuffle()
	hands := deck.Deal(len(userIDs))

	players := make([]*Player, len(userIDs))
	for i, uid := range userIDs {
		players[i] = &Player{
			UserID: uid,
			Hand:   hands[i],
			Rank:   0,
		}
	}

	return &Game{
		Players:    players,
		Turn:       0,
		FieldCards: []*Card{},
	}
}

// Play は指定されたプレイヤーがカードを出す処理を行う
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

	// ルール判定
	if err := g.validatePlay(cards); err != nil {
		return err
	}

	// 手札から削除
	player.RemoveCards(cards)

	// 場に出す
	g.FieldCards = cards
	g.PassCount = 0 // パスカウントリセット

	// 革命判定
	if len(cards) >= 4 {
		g.IsRevolution = !g.IsRevolution
		fmt.Printf("革命発生! (Rev: %v)\n", g.IsRevolution)
	}

	// 11バック判定
	if ContainsJack(cards) {
		g.Is11Back = true
		fmt.Println("11バック有効")
	}

	// あがり判定
	if len(player.Hand) == 0 {
		g.handleWin(player)
	}

	// 8切り判定
	if ContainsEight(cards) {
		fmt.Println("8切り! 場を流して同じ人のターン")
		g.ClearField()
		// ターンを進めずにリターン
		return nil
	}

	// 次のターンへ
	g.advanceTurn()
	return nil
}

// Pass はプレイヤーがパスをする処理を行う
func (g *Game) Pass(userID int64) error {
	player := g.Players[g.Turn]
	if player.UserID != userID {
		return errors.New("あなたのターンではありません")
	}

	// 場に何もないのにパスはできない
	if len(g.FieldCards) == 0 {
		return errors.New("自分の番から開始の場合はパスできません")
	}

	g.PassCount++

	// 参加中の人数 - 1 人がパスしたら場が流れる
	activeCount := g.getActivePlayerCount()
	if g.PassCount >= activeCount-1 {
		fmt.Println("全員パス -> 場が流れました")
		g.ClearField()
	}

	g.advanceTurn()
	return nil
}

// バリデーション
func (g *Game) validatePlay(cards []*Card) error {
	// 自分の出すカードの解析
	isRev := g.IsRevolution != g.Is11Back
	myType, myStr, err := AnalyzeHand(cards, isRev)
	if err != nil {
		return err
	}

	// 場にカードがない場合 -> なんでも出せる
	if len(g.FieldCards) == 0 {
		return nil
	}

	// 場のカードの解析
	fieldType, fieldStr, _ := AnalyzeHand(g.FieldCards, isRev)

	// 役の種類が一致しているか
	if myType != fieldType {
		return errors.New("場のカードと役の種類が一致しません")
	}

	// 枚数が一致しているか
	if len(cards) != len(g.FieldCards) {
		return errors.New("場のカードと枚数が一致しません")
	}

	// 強さが上回っているか
	// AnalyzeHandは革命考慮済みの強さを返すので、単純比較でOK
	if myStr <= fieldStr {
		return errors.New("場のカードより弱いです")
	}

	// 4. (特殊) 縛り判定などを入れるならここ

	return nil
}

// ターンを進める
func (g *Game) advanceTurn() {
	// 無限ループ防止のため最大人数分トライ
	for i := 0; i < len(g.Players); i++ {
		g.Turn++
		if g.Turn >= len(g.Players) {
			g.Turn = 0
		}

		// まだあがっていない人ならターン確定
		if len(g.Players[g.Turn].Hand) > 0 {
			break
		}
	}
}

// 場を流す
func (g *Game) ClearField() {
	g.FieldCards = []*Card{}
	g.PassCount = 0
	g.Is11Back = false // 場が流れたら11バック終了
}

func (g *Game) handleWin(player *Player) {
	g.FinishedPlayers = append(g.FinishedPlayers, player)
	player.Rank = len(g.FinishedPlayers)
	fmt.Printf("Player %d finished at rank %d\n", player.UserID, player.Rank)
}

func (g *Game) getActivePlayerCount() int {
	count := 0
	for _, p := range g.Players {
		if len(p.Hand) > 0 {
			count++
		}
	}
	return count
}

// カードを持っているか確認
func (p *Player) HasCards(cards []*Card) bool {
	for _, target := range cards {
		found := false
		for _, handCard := range p.Hand {
			if handCard.ID == target.ID {
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

// カードを捨てる
func (p *Player) RemoveCards(targets []*Card) {
	var newHand []*Card

	for _, handCard := range p.Hand {
		isTarget := false
		for _, t := range targets {
			if handCard.ID == t.ID {
				isTarget = true
				break
			}
		}

		if !isTarget {
			newHand = append(newHand, handCard)
		}
	}
	p.Hand = newHand
}
