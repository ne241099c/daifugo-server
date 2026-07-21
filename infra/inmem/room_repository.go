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
	mtx   sync.RWMutex
	data  map[int64]*model.Room
	locks map[int64]*sync.Mutex // 部屋IDごとの排他ロック
	next  int64
}

func NewInmemRoomRepository() *InmemRoomRepository {
	return &InmemRoomRepository{
		data:  make(map[int64]*model.Room),
		locks: make(map[int64]*sync.Mutex),
		next:  1,
	}
}

// roomLock は指定IDのロックを取得（無ければ生成）して返す。
func (r *InmemRoomRepository) roomLock(id int64) *sync.Mutex {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	m, ok := r.locks[id]
	if !ok {
		m = &sync.Mutex{}
		r.locks[id] = m
	}
	return m
}

// WithLock は部屋ごとの排他ロックを取得し、解放関数を返す。
func (r *InmemRoomRepository) WithLock(id int64) func() {
	m := r.roomLock(id)
	m.Lock()
	return m.Unlock
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
	// まず対象の部屋ポインタをスナップショットする（r.mtx を長時間保持しない）
	r.mtx.RLock()
	originals := make([]*model.Room, 0, len(r.data))
	for _, room := range r.data {
		originals = append(originals, room)
	}
	r.mtx.RUnlock()

	// 各部屋を、その部屋のロックを取りながら安全にディープコピーする。
	// （Play/Pass 等の変更中に Game をマーシャルしないようにするため）
	rooms := make([]*model.Room, 0, len(originals))
	for _, original := range originals {
		release := r.WithLock(original.ID)
		safeCopy := jsonDeepCopy(original)
		release()

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

	// UpdatedAt は SaveRoom（r.mtx 保持下）でのみ書かれるため、
	// ここで r.mtx を保持したまま読めば整合的に判定できる。
	for id, room := range r.data {
		if room.UpdatedAt.Before(threshold) {
			delete(r.data, id)
			delete(r.locks, id)
			deletedCount++
		}
	}

	if deletedCount > 0 {
		fmt.Printf("Cleaned up %d rooms\n", deletedCount)
	}
}

// jsonDeepCopy は JSON を使った簡易ディープコピー。
// 呼び出し側が対象部屋のロックを保持している前提。
func jsonDeepCopy(src *model.Room) *model.Room {
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
