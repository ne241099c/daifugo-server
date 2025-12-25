package game

type Player struct {
	UserID int64
	Hand   []*Card
}

type Game struct {
	Players []*Player // 参加者（順番に並んでいる）
	Turn    int       // 現在のターンのプレイヤーのインデックス
	Field   []*Card   // 場に出ているカード
}

func NewGame(userIDs []int64) *Game {
	deck := NewDeck()
	deck.Shuffle()

	// 人数分にカードを配る
	hands := deck.Deal(len(userIDs))

	// Player構造体を作成して手札を持たせる
	players := make([]*Player, len(userIDs))
	for i, uid := range userIDs {
		players[i] = &Player{
			UserID: uid,
			Hand:   hands[i],
		}
	}

	// Turnは0番目の人からスタート
	return &Game{
		Players: players,
		Turn:    0,
		Field:   []*Card{},
	}
}
