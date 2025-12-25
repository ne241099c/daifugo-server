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
