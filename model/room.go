package model

import (
	"time"

	"github.com/ne241099/daifugo-server/internal/game"
)

// Room はドメインエンティティ。並行制御（ロック）は永続化の関心事なので
// ここには持たせず、RoomRepository 側で部屋ごとに排他制御する。
type Room struct {
	ID        int64   `json:"id"`
	Name      string  `json:"name"`
	OwnerID   int64   `json:"owner_id"`
	MemberIDs []int64 `json:"member_ids"`
	// BotIDs は MemberIDs のうち CPU（SVM-AI）が操作するプレイヤーの userID。
	BotIDs    []int64       `json:"bot_ids"`
	Game      *game.Game    `json:"game"`
	PrevRanks map[int64]int `json:"prev_ranks"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

func (r *Room) IsFull() bool {
	return len(r.MemberIDs) >= 4
}

// IsBot は userID が CPU（SVM-AI）プレイヤーかどうかを返す。
func (r *Room) IsBot(userID int64) bool {
	for _, id := range r.BotIDs {
		if id == userID {
			return true
		}
	}
	return false
}

// HasHumanMember は CPU(Bot) 以外の人間メンバーが1人でもいるかを返す。
// これが false の部屋（CPUのみ／空）は存在させ続ける意味がない。
func (r *Room) HasHumanMember() bool {
	for _, mid := range r.MemberIDs {
		if !r.IsBot(mid) {
			return true
		}
	}
	return false
}

func (r *Room) StartGame() {
	r.Game = game.NewGame(r.MemberIDs)
}

func (r *Room) RestartGame() {
	r.Game = r.Game.Reset()
}

func NewRoom(name string, ownerID int64) *Room {
	return &Room{
		Name:      name,
		OwnerID:   ownerID,
		MemberIDs: []int64{ownerID},
		PrevRanks: make(map[int64]int),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
