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
	room.Mu.Lock()
	defer room.Mu.Unlock()

	if room.Game != nil {
		return nil, fmt.Errorf("game already started")
	}

	if len(room.MemberIDs) < 2 {
		return nil, fmt.Errorf("at least 2 players are required")
	}

	room.StartGame()

	if len(room.PrevRanks) > 0 {
		restoredCount := 0
		// 前回の順位を復元
		for _, p := range room.Game.Players {
			if rank, ok := room.PrevRanks[p.UserID]; ok {
				p.Rank = rank
				restoredCount++
			}
		}

		// 順位がついている人がいれば、Reset() を呼んで手札交換を実行させる
		if restoredCount > 0 {
			room.Game = room.Game.Reset()
		}
	}
	if err := uc.RoomRepository.SaveRoom(ctx, room); err != nil {
		return nil, err
	}

	return room, nil
}
