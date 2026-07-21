package graph

import (
	"github.com/ne241099/daifugo-server/internal/sse"
	"github.com/ne241099/daifugo-server/usecase/game"
	"github.com/ne241099/daifugo-server/usecase/room"
	"github.com/ne241099/daifugo-server/usecase/user"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require
// here.

type Resolver struct {
	Hub                *sse.Hub
	CreateGuestUseCase user.CreateGuestUseCase
	RenameUserUseCase  user.RenameUserUseCase
	GetUserUseCase     user.GetUserUseCase
	ListUsersUseCase   user.ListUsersUseCase
	CreateRoomUseCase  room.CreateRoomUseCase
	JoinRoomUseCase    room.JoinRoomUseCase
	AddBotUseCase      room.AddBotUseCase
	LeaveRoomUseCase   room.LeaveRoomUseCase
	ListRoomsUseCase   room.ListRoomsUseCase
	GetRoomUseCase     room.GetRoomUseCase
	StartGameUseCase   *game.StartGameInteractor
	RestartGameUseCase *game.RestartGameInteractor
	PlayCardUseCase    *game.PlayCardInteractor
	PassUseCase        *game.PassInteractor
}
