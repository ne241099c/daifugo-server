package graph

import (
	"github.com/ne241099/daifugo-server/internal/sse"
	"github.com/ne241099/daifugo-server/usecase/room"
	"github.com/ne241099/daifugo-server/usecase/user"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require
// here.

type Resolver struct {
	Hub               *sse.Hub
	SignUpUseCase     user.SignUpUseCase
	CreateRoomUseCase room.CreateRoomUseCase
	JoinRoomUseCase   room.JoinRoomUseCase
	ListRoomsUseCase  room.ListRoomsUseCase
}
