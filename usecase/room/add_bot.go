package room

import (
	"context"
	"errors"
	"fmt"

	"github.com/ne241099/daifugo-server/model"
	"github.com/ne241099/daifugo-server/repository"
)

type AddBotUseCase interface {
	Execute(ctx context.Context, roomID int64) (*model.Room, error)
}

var _ AddBotUseCase = &AddBotInteractor{}

// AddBotInteractor は部屋に CPU(SVM-AI) プレイヤーを1人追加する。
// Bot 用のゲストユーザーを作成し、その userID を MemberIDs と BotIDs に登録する。
type AddBotInteractor struct {
	RoomRepository repository.RoomRepository
	UserRepository repository.UserRepository
}

func (uc *AddBotInteractor) Execute(ctx context.Context, roomID int64) (*model.Room, error) {
	release := uc.RoomRepository.WithLock(roomID)
	defer release()

	room, err := uc.RoomRepository.GetRoomByID(ctx, roomID)
	if err != nil {
		return nil, err
	}
	if room.Game != nil {
		return nil, errors.New("game already started")
	}
	if room.IsFull() {
		return nil, errors.New("room is full")
	}

	// Bot 用のゲストユーザーを作成（ID は SaveUser が採番）
	bot := &model.User{}
	bot.CreateGuest(fmt.Sprintf("CPU%d", len(room.MemberIDs)))
	if err := uc.UserRepository.SaveUser(ctx, bot); err != nil {
		return nil, err
	}

	room.MemberIDs = append(room.MemberIDs, bot.ID)
	room.BotIDs = append(room.BotIDs, bot.ID)

	if err := uc.RoomRepository.SaveRoom(ctx, room); err != nil {
		return nil, err
	}
	return room, nil
}
