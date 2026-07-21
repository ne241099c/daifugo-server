package main

import (
	"time"

	"github.com/ne241099/daifugo-server/graph"
	"github.com/ne241099/daifugo-server/infra/inmem"
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

	// ログイン/永続化を廃止したため、ユーザーも部屋もインメモリで扱う（DB 不要）。
	userRepo := inmem.NewInmemUserRepository()
	roomRepo := inmem.NewInmemRoomRepository()

	// 定期クリーンアップ開始
	go func() {
		// 1時間に1回チェック
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()

		for range ticker.C {
			// 24時間以上触られていない部屋を削除
			roomRepo.CleanupRooms(24 * time.Hour)
		}
	}()

	// ゲストIDヘッダー(X-User-ID)でプレイヤーを識別する
	authMiddleware := internalMiddleware.NewAuthMiddleware(userRepo)

	// SSE Hub 作成
	hub := sse.NewHub()

	// Resolver 作成
	resolver := &graph.Resolver{
		Hub: hub,
		CreateGuestUseCase: &user.CreateGuestInteractor{
			UserRepository: userRepo,
		},
		RenameUserUseCase: &user.RenameUserInteractor{
			UserRepository: userRepo,
		},
		GetUserUseCase: &user.GetUserInteractor{
			UserRepository: userRepo,
		},
		ListUsersUseCase: &user.ListUsersInteractor{
			UserRepository: userRepo,
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
