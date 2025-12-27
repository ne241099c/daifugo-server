package room

import (
	"context"
	"errors"

	"github.com/ne241099/daifugo-server/model"
	"github.com/ne241099/daifugo-server/repository"
)

type JoinRoomUseCase interface {
	Execute(ctx context.Context, roomID int64, userID int64) (*model.Room, error)
}

var _ JoinRoomUseCase = &JoinRoomInteractor{}

type JoinRoomInteractor struct {
	RoomRepository repository.RoomRepository
}

func (uc *JoinRoomInteractor) Execute(ctx context.Context, roomID int64, userID int64) (*model.Room, error) {
	room, err := uc.RoomRepository.GetRoomByID(ctx, roomID)
	if err != nil {
		return nil, err
	}
	room.Mu.Lock()
	defer room.Mu.Unlock()

	if room.IsFull() {
		return nil, errors.New("room is full")
	}

	for _, memberID := range room.MemberIDs {
		if memberID == userID {
			return room, nil
		}
	}

	room.MemberIDs = append(room.MemberIDs, userID)

	if err := uc.RoomRepository.SaveRoom(ctx, room); err != nil {
		return nil, err
	}

	return room, nil
}
