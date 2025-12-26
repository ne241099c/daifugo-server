package game

import (
	"context"
	"fmt"

	"github.com/ne241099/daifugo-server/internal/game"
	"github.com/ne241099/daifugo-server/model"
	"github.com/ne241099/daifugo-server/repository"
)

type PlayCardUseCase interface {
	Execute(ctx context.Context, roomID int64, userID int64, cardIDs []int) (*model.Room, error)
}

var _ PlayCardUseCase = &PlayCardInteractor{}

type PlayCardInteractor struct {
	RoomRepository repository.RoomRepository
}

func (i *PlayCardInteractor) Execute(ctx context.Context, roomID int64, userID int64, cardIDs []int) (*model.Room, error) {
	room, err := i.RoomRepository.GetRoomByID(ctx, roomID)
	if err != nil {
		return nil, fmt.Errorf("room not found: %w", err)
	}

	if room.Game == nil {
		return nil, fmt.Errorf("game not started")
	}

	// IDリストから判定用カード生成
	var targetCards []*game.Card
	for _, cid := range cardIDs {
		targetCards = append(targetCards, &game.Card{ID: cid})
	}

	// ロジック実行
	if err := room.Game.Play(userID, targetCards); err != nil {
		return nil, err
	}

	return room, nil
}
