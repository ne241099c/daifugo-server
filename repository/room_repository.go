package repository

import (
	"context"

	"github.com/ne241099/daifugo-server/model"
)

type RoomRepository interface {
	// SaveRoom は、部屋を保存する
	SaveRoom(ctx context.Context, room *model.Room) error
	// DeleteRoom は、部屋を削除する
	DeleteRoom(ctx context.Context, id int64) error
	// UpdateRoom は、部屋情報を更新する
	UpdateRoom(ctx context.Context, room *model.Room) error
	// ListRooms は、部屋一覧を取得する
	ListRooms(ctx context.Context) ([]*model.Room, error)
	// GetRoomByID は、IDから部屋を取得する
	GetRoomByID(ctx context.Context, id int64) (*model.Room, error)
}
