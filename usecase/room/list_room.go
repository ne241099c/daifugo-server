package room

import (
	"context"

	"github.com/ne241099/daifugo-server/model"
	"github.com/ne241099/daifugo-server/repository"
)

type ListRoomsUseCase interface {
	Execute(ctx context.Context) ([]*model.Room, error)
}

var _ ListRoomsUseCase = &ListRoomsInteractor{}

type ListRoomsInteractor struct {
	RoomRepository repository.RoomRepository
}

func (uc *ListRoomsInteractor) Execute(ctx context.Context) ([]*model.Room, error) {
	rooms, err := uc.RoomRepository.ListRooms(ctx)
	if err != nil {
		return nil, err
	}

	return rooms, nil
}
