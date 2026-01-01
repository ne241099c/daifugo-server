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

func (uc *PlayCardInteractor) Execute(ctx context.Context, roomID int64, userID int64, cardIDs []int) (*model.Room, error) {
	room, err := uc.RoomRepository.GetRoomByID(ctx, roomID)
	if err != nil {
		return nil, fmt.Errorf("room not found: %w", err)
	}
	room.Mu.Lock()
	defer room.Mu.Unlock()

	if room.Game == nil {
		return nil, fmt.Errorf("game not started")
	}

	// プレイヤーを特定
	var targetPlayer *game.Player
	for _, p := range room.Game.Players {
		if p.UserID == userID {
			targetPlayer = p
			break
		}
	}
	if targetPlayer == nil {
		return nil, fmt.Errorf("player not found in this game")
	}

	// 手札から指定されたカードを取得
	var targetCards []*game.Card
	for _, cid := range cardIDs {
		found := false
		for _, handCard := range targetPlayer.Hand {
			if handCard.ID == cid {
				targetCards = append(targetCards, handCard)
				found = true
				break
			}
		}
		if !found {
			return nil, fmt.Errorf("card %d not found in player's hand", cid)
		}
	}

	// ロジック実行
	if err := room.Game.Play(userID, targetCards); err != nil {
		return nil, err
	}

	room.Game.ProcessBots()

	if err := uc.RoomRepository.SaveRoom(ctx, room); err != nil {
		return nil, err
	}

	return room, nil
}
