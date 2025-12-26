package game

import (
	"context"
	"fmt"

	"github.com/ne241099/daifugo-server/model"
	"github.com/ne241099/daifugo-server/repository"
)

type StartGameUseCase interface {
	Execute(ctx context.Context, roomID int64) (*model.Room, error)
}

var _ StartGameUseCase = &StartGameInteractor{}

type StartGameInteractor struct {
	RoomRepository repository.RoomRepository
}

func (uc *StartGameInteractor) Execute(ctx context.Context, roomID int64) (*model.Room, error) {
	room, err := uc.RoomRepository.GetRoomByID(ctx, roomID)
	if err != nil {
		return nil, fmt.Errorf("room not found: %w", err)
	}

	if room.Game != nil {
		return nil, fmt.Errorf("game already started")
	}

	if len(room.MemberIDs) < 2 {
		return nil, fmt.Errorf("at least 2 players are required")
	}

	room.StartGame()

	// 本来はここで roomRepo.Update(room) などを呼んで保存します

	return room, nil
}
