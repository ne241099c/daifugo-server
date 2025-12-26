package inmem

import (
	"context"
	"sort"
	"sync"

	"github.com/ne241099/daifugo-server/model"
	"github.com/ne241099/daifugo-server/repository"
)

var _ repository.RoomRepository = &InmemRoomRepository{}

type InmemRoomRepository struct {
	mtx  sync.RWMutex
	data map[int64]*model.Room
	next int64
}

func (r *InmemRoomRepository) UpdateRoom(ctx context.Context, room *model.Room) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	// マップの値を更新
	r.data[room.ID] = room
	return nil
}

func (r *InmemRoomRepository) DeleteRoom(ctx context.Context, id int64) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	delete(r.data, id)
	return nil
}

func NewInmemRoomRepository() *InmemRoomRepository {
	return &InmemRoomRepository{
		data: make(map[int64]*model.Room),
		next: 1,
	}
}

func (r *InmemRoomRepository) SaveRoom(ctx context.Context, room *model.Room) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	if room.ID == 0 {
		room.ID = r.next
		r.next++
	}

	r.data[room.ID] = room
	return nil
}

func (r *InmemRoomRepository) ListRooms(ctx context.Context) ([]*model.Room, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	// マップからリストに詰め替え
	rooms := make([]*model.Room, 0, len(r.data))
	for _, room := range r.data {
		rooms = append(rooms, room)
	}

	// ID順に並び替え
	sort.Slice(rooms, func(i, j int) bool {
		return rooms[i].ID < rooms[j].ID
	})

	return rooms, nil
}

func (r *InmemRoomRepository) GetRoomByID(ctx context.Context, id int64) (*model.Room, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	room, ok := r.data[id]
	if !ok {
		return nil, repository.ErrEntityNotFound
	}
	return room, nil
}
