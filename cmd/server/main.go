package main

import (
	"github.com/ne241099/daifugo-server/graph"
	"github.com/ne241099/daifugo-server/infra/inmem"
	"github.com/ne241099/daifugo-server/internal/server"
	"github.com/ne241099/daifugo-server/internal/sse"
	"github.com/ne241099/daifugo-server/usecase/room"
	"github.com/ne241099/daifugo-server/usecase/user"
)

func main() {
	// リポジトリ初期化
	userRepo := inmem.NewInmemUserRepository()
	roomRepo := inmem.NewInmemRoomRepository()

	// SSE Hub 作成
	hub := sse.NewHub()

	// Resolver 作成
	resolver := &graph.Resolver{
		Hub: hub,
		SignUpUseCase: &user.SignUpInteractor{
			UserRepository: userRepo,
		},
		GetUserUseCase: &user.GetUserInteractor{
			UserRepository: userRepo,
		},
		CreateRoomUseCase: &room.CreateRoomInteractor{
			RoomRepository: roomRepo,
		},
		JoinRoomUseCase: &room.JoinRoomInteractor{
			RoomRepository: roomRepo,
		},
		ListRoomsUseCase: &room.ListRoomsInteractor{
			RoomRepository: roomRepo,
		},
	}
	// サーバー作成
	srv := server.New(resolver, hub)

	// サーバー起動
	srv.Logger.Fatal(srv.Start(":8080"))
}
