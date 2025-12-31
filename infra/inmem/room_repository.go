package inmem

import (
	"context"
	"sort"
	"sync"

	"github.com/ne241099/daifugo-server/internal/game"
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
	r.data[room.ID] = r.deepCopy(room)
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

	saved := r.deepCopy(room)
	r.data[room.ID] = saved

	room.ID = saved.ID
	return nil
}

func (r *InmemRoomRepository) ListRooms(ctx context.Context) ([]*model.Room, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	// マップからリストに詰め替え
	rooms := make([]*model.Room, 0, len(r.data))
	for _, room := range r.data {
		rooms = append(rooms, r.deepCopy(room))
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
	return r.deepCopy(room), nil
}

func (r *InmemRoomRepository) deepCopy(src *model.Room) *model.Room {
	if src == nil {
		return nil
	}
	dst := &model.Room{
		ID:        src.ID,
		Name:      src.Name,
		OwnerID:   src.OwnerID,
		CreatedAt: src.CreatedAt,
		UpdatedAt: src.UpdatedAt,
	}

	// MemberIDsを新しく作り直してコピー
	if src.MemberIDs != nil {
		dst.MemberIDs = make([]int64, len(src.MemberIDs))
		copy(dst.MemberIDs, src.MemberIDs)
	}

	if src.PrevRanks != nil {
		dst.PrevRanks = make(map[int64]int)
		for k, v := range src.PrevRanks {
			dst.PrevRanks[k] = v
		}
	} else {
		dst.PrevRanks = make(map[int64]int)
	}

	// Game構造体もポインタを含んでいるため再帰的にコピー
	if src.Game != nil {
		dstGame := *src.Game

		if src.Game.Players != nil {
			dstGame.Players = make([]*game.Player, len(src.Game.Players))
			for i, p := range src.Game.Players {
				if p != nil {
					pCopy := *p
					if p.Hand != nil {
						pCopy.Hand = make([]*game.Card, len(p.Hand))
						for j, c := range p.Hand {
							if c != nil {
								cCopy := *c
								pCopy.Hand[j] = &cCopy
							}
						}
					}
					dstGame.Players[i] = &pCopy
				}
			}
		}

		if src.Game.FieldCards != nil {
			dstGame.FieldCards = make([]*game.Card, len(src.Game.FieldCards))
			for i, c := range src.Game.FieldCards {
				if c != nil {
					cCopy := *c
					dstGame.FieldCards[i] = &cCopy
				}
			}
		}

		if src.Game.FinishedPlayers != nil {
			dstGame.FinishedPlayers = make([]*game.Player, len(src.Game.FinishedPlayers))
			for i, p := range src.Game.FinishedPlayers {
				if p != nil {
					pCopy := *p
					if p.Hand != nil {
						pCopy.Hand = make([]*game.Card, len(p.Hand))
						for j, c := range p.Hand {
							if c != nil {
								cCopy := *c
								pCopy.Hand[j] = &cCopy
							}
						}
					}
					dstGame.FinishedPlayers[i] = &pCopy
				}
			}
		}

		dst.Game = &dstGame
	}

	return dst
}
