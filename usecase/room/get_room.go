package room

import (
	"context"

	"github.com/ne241099/daifugo-server/model"
	"github.com/ne241099/daifugo-server/repository"
)

type GetRoomUseCase interface {
	Execute(ctx context.Context, roomID int64) (*model.Room, error)
}

var _ GetRoomUseCase = &GetRoomInteractor{}

type GetRoomInteractor struct {
	RoomRepository repository.RoomRepository
}

func (uc *GetRoomInteractor) Execute(ctx context.Context, roomID int64) (*model.Room, error) {
	room, err := uc.RoomRepository.GetRoomByID(ctx, roomID)
	if err != nil {
		return nil, err
	}

	return room, nil
}
