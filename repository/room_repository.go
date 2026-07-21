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
	// WithLock は、指定した部屋IDに対する排他ロックを取得し、解放用の関数を返す。
	// 「取得→変更→保存」を直列化するために使う。呼び出し側は戻り値を必ず defer で呼ぶこと。
	// （並行制御はドメインではなく永続化層の責務であるため、ここに置く）
	WithLock(id int64) (release func())
}
