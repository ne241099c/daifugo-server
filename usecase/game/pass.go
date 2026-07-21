package game

import (
	"context"
	"fmt"

	"github.com/ne241099/daifugo-server/internal/game/ai"
	"github.com/ne241099/daifugo-server/model"
	"github.com/ne241099/daifugo-server/repository"
)

type PassUseCase interface {
	Execute(ctx context.Context, roomID int64, userID int64) (*model.Room, error)
}

var _ PassUseCase = &PassInteractor{}

type PassInteractor struct {
	RoomRepository repository.RoomRepository
	// AIModel は CPU の着手判断に使う学習済み SVM。nil の場合 CPU は動かない。
	AIModel *ai.Model
}

func (uc *PassInteractor) Execute(ctx context.Context, roomID int64, userID int64) (*model.Room, error) {
	release := uc.RoomRepository.WithLock(roomID)
	defer release()

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

	// 続く CPU（Bot）の手番を SVM-AI に自動で進めさせる。
	if uc.AIModel != nil {
		if err := uc.AIModel.RunBotTurns(room.Game, room.IsBot); err != nil {
			return nil, err
		}
	}

	if err := uc.RoomRepository.SaveRoom(ctx, room); err != nil {
		return nil, err
	}

	return room, nil
}
