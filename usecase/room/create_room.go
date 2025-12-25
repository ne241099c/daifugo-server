package room

import (
	"context"

	"github.com/ne241099/daifugo-server/model"
	"github.com/ne241099/daifugo-server/repository"
)

type CreateRoomUseCase interface {
	Execute(ctx context.Context, name string, ownerID int64) (*model.Room, error)
}

var _ CreateRoomUseCase = &CreateRoomInteractor{}

type CreateRoomInteractor struct {
	RoomRepository repository.RoomRepository
}

func (uc *CreateRoomInteractor) Execute(ctx context.Context, name string, ownerID int64) (*model.Room, error) {
	room := model.NewRoom(name, ownerID)

	if err := uc.RoomRepository.SaveRoom(ctx, room); err != nil {
		return nil, err
	}

	return room, nil
}
