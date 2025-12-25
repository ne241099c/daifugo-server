package model

import "time"

type Room struct {
	ID        int64
	Name      string
	OwnerID   int64
	MemberIDs []int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (r *Room) IsFull() bool {
	return len(r.MemberIDs) >= 4
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
