package repository

import (
	"context"

	"github.com/ne241099/daifugo-server/graph/model"
)

type RoomRepository interface {
	// SaveRoom は、部屋を保存する
	SaveRoom(ctx context.Context, room *model.Room) error
	// ListRooms は、部屋一覧を取得する
	ListRooms(ctx context.Context) ([]*model.Room, error)
	// GetRoomByID は、IDから部屋を取得する
	GetRoomByID(ctx context.Context, id int64) (*model.Room, error)
}
