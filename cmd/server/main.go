package main

import (
	"github.com/ne241099/daifugo-server/graph"
	"github.com/ne241099/daifugo-server/infra/inmem"
	"github.com/ne241099/daifugo-server/infra/mysql"
	"github.com/ne241099/daifugo-server/internal/auth"
	"github.com/ne241099/daifugo-server/internal/config"
	internalMiddleware "github.com/ne241099/daifugo-server/internal/middleware"
	"github.com/ne241099/daifugo-server/internal/server"
	"github.com/ne241099/daifugo-server/internal/sse"
	"github.com/ne241099/daifugo-server/usecase/game"
	"github.com/ne241099/daifugo-server/usecase/room"
	"github.com/ne241099/daifugo-server/usecase/user"
)

func main() {
	cfg := config.Load()

	db, err := mysql.NewDB(cfg)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// リポジトリ初期化
	userRepo := mysql.NewMySQLUserRepository(db)
	roomRepo := inmem.NewInmemRoomRepository()

	// Configから読み込んだ秘密鍵を使用する
	authenticator := auth.NewJWTAuthenticator(cfg.JWTSecret)
	authMiddleware := internalMiddleware.NewAuthMiddleware(authenticator)

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
		ListUsersUseCase: &user.ListUsersInteractor{
			UserRepository: userRepo,
		},
		DeleteUserUseCase: &user.DeleteUserInteractor{
			UserRepository: userRepo,
		},
		LoginUseCase: &user.LoginInteractor{
			UserRepository: userRepo,
			Authenticator:  authenticator,
		},
		CreateRoomUseCase: &room.CreateRoomInteractor{
			RoomRepository: roomRepo,
		},
		JoinRoomUseCase: &room.JoinRoomInteractor{
			RoomRepository: roomRepo,
		},
		LeaveRoomUseCase: &room.LeaveRoomInteractor{
			RoomRepository: roomRepo,
		},
		ListRoomsUseCase: &room.ListRoomsInteractor{
			RoomRepository: roomRepo,
		},
		GetRoomUseCase: &room.GetRoomInteractor{
			RoomRepository: roomRepo,
		},
		StartGameUseCase: &game.StartGameInteractor{
			RoomRepository: roomRepo,
		},
		RestartGameUseCase: &game.RestartGameInteractor{
			RoomRepository: roomRepo,
		},
		PlayCardUseCase: &game.PlayCardInteractor{
			RoomRepository: roomRepo,
		},
		PassUseCase: &game.PassInteractor{
			RoomRepository: roomRepo,
		},
	}
	// サーバー作成
	srv := server.New(resolver, hub, authMiddleware)

	// サーバー起動
	srv.Logger.Fatal(srv.Start(":" + cfg.Port))
}
