package inmem

import (
	"context"
	"encoding/json"
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

	// 新規作成の場合はIDを割り当て
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

	rooms := make([]*model.Room, 0, len(r.data))
	for _, room := range r.data {
		rooms = append(rooms, r.safeCopy(room))
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

func (r *InmemRoomRepository) safeCopy(src *model.Room) *model.Room {
	src.Mu.Lock()
	defer src.Mu.Unlock()

	return r.deepCopy(src)
}

func (r *InmemRoomRepository) deepCopy(src *model.Room) *model.Room {
	if src == nil {
		return nil
	}

	// JSONに変換
	b, err := json.Marshal(src)
	if err != nil {
		return nil // 本来はエラーハンドリングすべき
	}

	// 新しい構造体に書き戻す
	dst := &model.Room{}
	if err := json.Unmarshal(b, dst); err != nil {
		return nil
	}

	return dst
}
