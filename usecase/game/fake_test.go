package game

import (
	"context"

	"github.com/ne241099/daifugo-server/model"
	"github.com/ne241099/daifugo-server/repository"
)

// fakeRoomRepo はテスト用のインメモリ RoomRepository
type fakeRoomRepo struct {
	rooms map[int64]*model.Room
	next  int64
}

func newFakeRoomRepo() *fakeRoomRepo {
	return &fakeRoomRepo{rooms: make(map[int64]*model.Room), next: 1}
}

func (f *fakeRoomRepo) SaveRoom(ctx context.Context, r *model.Room) error {
	if r.ID == 0 {
		r.ID = f.next
		f.next++
	}
	f.rooms[r.ID] = r
	return nil
}

func (f *fakeRoomRepo) UpdateRoom(ctx context.Context, r *model.Room) error {
	return f.SaveRoom(ctx, r)
}

func (f *fakeRoomRepo) DeleteRoom(ctx context.Context, id int64) error {
	delete(f.rooms, id)
	return nil
}

func (f *fakeRoomRepo) ListRooms(ctx context.Context) ([]*model.Room, error) {
	out := make([]*model.Room, 0, len(f.rooms))
	for _, r := range f.rooms {
		out = append(out, r)
	}
	return out, nil
}

func (f *fakeRoomRepo) GetRoomByID(ctx context.Context, id int64) (*model.Room, error) {
	r, ok := f.rooms[id]
	if !ok {
		return nil, repository.ErrEntityNotFound
	}
	return r, nil
}

// WithLock はテストでは逐次実行のため no-op で十分。
func (f *fakeRoomRepo) WithLock(id int64) func() { return func() {} }
