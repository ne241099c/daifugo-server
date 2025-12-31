package inmem

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"sync"
	"time"

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
	return r.SaveRoom(ctx, room)
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

	// 新規作成の場合はIDを割り当て
	if room.ID == 0 {
		room.ID = r.next
		r.next++
	}

	room.UpdatedAt = time.Now()

	r.data[room.ID] = room
	return nil
}

func (r *InmemRoomRepository) ListRooms(ctx context.Context) ([]*model.Room, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	rooms := make([]*model.Room, 0, len(r.data))
	for _, original := range r.data {
		safeCopy := r.jsonDeepCopy(original)
		if safeCopy != nil {
			rooms = append(rooms, safeCopy)
		}
	}

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

func (r *InmemRoomRepository) CleanupRooms(expiration time.Duration) {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	threshold := time.Now().Add(-expiration)
	deletedCount := 0

	for id, room := range r.data {
		room.Mu.Lock()
		updatedAt := room.UpdatedAt
		room.Mu.Unlock()

		if updatedAt.Before(threshold) {
			delete(r.data, id)
			deletedCount++
		}
	}

	if deletedCount > 0 {
		fmt.Printf("Cleaned up %d rooms\n", deletedCount)
	}
}

// JSONを使った簡易DeepCopy
func (r *InmemRoomRepository) jsonDeepCopy(src *model.Room) *model.Room {
	src.Mu.Lock()
	defer src.Mu.Unlock()

	b, err := json.Marshal(src)
	if err != nil {
		return nil
	}
	var dst model.Room
	if err := json.Unmarshal(b, &dst); err != nil {
		return nil
	}
	return &dst
}
