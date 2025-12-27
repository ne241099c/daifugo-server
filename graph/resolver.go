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
	SignUpUseCase      user.SignUpUseCase
	LoginUseCase       user.LoginUseCase
	GetUserUseCase     user.GetUserUseCase
	ListUsersUseCase   user.ListUsersUseCase
	DeleteUserUseCase  user.DeleteUserUseCase
	CreateRoomUseCase  room.CreateRoomUseCase
	JoinRoomUseCase    room.JoinRoomUseCase
	LeaveRoomUseCase   room.LeaveRoomUseCase
	ListRoomsUseCase   room.ListRoomsUseCase
	GetRoomUseCase     room.GetRoomUseCase
	StartGameUseCase   *game.StartGameInteractor
	RestartGameUseCase *game.RestartGameInteractor
	PlayCardUseCase    *game.PlayCardInteractor
	PassUseCase        *game.PassInteractor
}
