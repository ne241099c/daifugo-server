package model

import (
	"time"

	"github.com/ne241099/daifugo-server/internal/game"
)

type Room struct {
	ID        int64
	Name      string
	OwnerID   int64
	MemberIDs []int64
	Game      *game.Game
	CreatedAt time.Time
	UpdatedAt time.Time
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
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
