package game

import (
	"context"
	"fmt"

	"github.com/ne241099/daifugo-server/model"
	"github.com/ne241099/daifugo-server/repository"
)

type PassUseCase interface {
	Execute(ctx context.Context, roomID int64, userID int64) (*model.Room, error)
}

var _ PassUseCase = &PassInteractor{}

type PassInteractor struct {
	RoomRepository repository.RoomRepository
}

func (uc *PassInteractor) Execute(ctx context.Context, roomID int64, userID int64) (*model.Room, error) {
	room, err := uc.RoomRepository.GetRoomByID(ctx, roomID)
	if err != nil {
		return nil, fmt.Errorf("room not found: %w", err)
	}

	if room.Game == nil {
		return nil, fmt.Errorf("game not started")
	}

	if err := room.Game.Pass(userID); err != nil {
		return nil, err
	}

	return room, nil
}
