package game

import (
	"context"
	"fmt"

	"github.com/ne241099/daifugo-server/model"
	"github.com/ne241099/daifugo-server/repository"
)

type RestartGameUseCase interface {
	Execute(ctx context.Context, roomID int64) (*model.Room, error)
}

var _ RestartGameUseCase = &RestartGameInteractor{}

type RestartGameInteractor struct {
	RoomRepository repository.RoomRepository
}

func (uc *RestartGameInteractor) Execute(ctx context.Context, roomID int64) (*model.Room, error) {
	room, err := uc.RoomRepository.GetRoomByID(ctx, roomID)
	if err != nil {
		return nil, err
	}
	room.Mu.Lock()
	defer room.Mu.Unlock()

	if room.Game == nil {
		return nil, fmt.Errorf("game is not started")
	}

	if room.PrevRanks == nil {
		room.PrevRanks = make(map[int64]int)
	}

	// 前回の順位を保存
	for k := range room.PrevRanks {
		delete(room.PrevRanks, k)
	}
	for _, p := range room.Game.FinishedPlayers {
		room.PrevRanks[p.UserID] = p.Rank
	}
	// 途中退場したプレイヤーの順位も保存
	for _, p := range room.Game.Players {
		if p.Rank > 0 {
			room.PrevRanks[p.UserID] = p.Rank
		}
	}

	room.Game = nil

	if err := uc.RoomRepository.SaveRoom(ctx, room); err != nil {
		return nil, err
	}

	return room, nil
}
