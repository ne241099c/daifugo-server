package main

import (
	"log"
	"time"

	"github.com/ne241099/daifugo-server/graph"
	"github.com/ne241099/daifugo-server/infra/inmem"
	"github.com/ne241099/daifugo-server/internal/config"
	"github.com/ne241099/daifugo-server/internal/game/ai"
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

	// 学習済み SVM-AI（CPU 用）を読み込む。失敗しても CPU 無しで起動する。
	aiModel, err := ai.Default()
	if err != nil {
		log.Printf("AIモデルの読み込みに失敗しました。CPU対戦は無効になります: %v", err)
		aiModel = nil
	}

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
		AddBotUseCase: &room.AddBotInteractor{
			RoomRepository: roomRepo,
			UserRepository: userRepo,
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
			AIModel:        aiModel,
		},
		RestartGameUseCase: &game.RestartGameInteractor{
			RoomRepository: roomRepo,
		},
		PlayCardUseCase: &game.PlayCardInteractor{
			RoomRepository: roomRepo,
			AIModel:        aiModel,
		},
		PassUseCase: &game.PassInteractor{
			RoomRepository: roomRepo,
			AIModel:        aiModel,
		},
	}
	// サーバー作成
	srv := server.New(resolver, hub, authMiddleware)

	// サーバー起動
	srv.Logger.Fatal(srv.Start(":" + cfg.Port))
}
