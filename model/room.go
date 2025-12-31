package model

import (
	"sync"
	"time"

	"github.com/ne241099/daifugo-server/internal/game"
)

type Room struct {
	ID        int64         `json:"id"`
	Name      string        `json:"name"`
	OwnerID   int64         `json:"owner_id"`
	MemberIDs []int64       `json:"member_ids"`
	Game      *game.Game    `json:"game"`
	PrevRanks map[int64]int `json:"prev_ranks"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
	Mu        sync.Mutex    `json:"-"`
}

func (r *Room) IsFull() bool {
	return len(r.MemberIDs) >= 4
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
