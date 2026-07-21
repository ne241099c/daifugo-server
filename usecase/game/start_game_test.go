package game

import (
	"context"
	"testing"

	"github.com/ne241099/daifugo-server/model"
)

func TestStartGame_Success(t *testing.T) {
	repo := newFakeRoomRepo()
	r := model.NewRoom("test", 1)
	r.MemberIDs = []int64{1, 2}
	_ = repo.SaveRoom(context.Background(), r)

	uc := &StartGameInteractor{RoomRepository: repo}
	got, err := uc.Execute(context.Background(), r.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Game == nil {
		t.Fatal("game should be created")
	}
	if len(got.Game.Players) != 2 {
		t.Errorf("players = %d, want 2", len(got.Game.Players))
	}
}

func TestStartGame_TooFewPlayers(t *testing.T) {
	repo := newFakeRoomRepo()
	r := model.NewRoom("test", 1) // members=[1] のみ
	_ = repo.SaveRoom(context.Background(), r)

	uc := &StartGameInteractor{RoomRepository: repo}
	if _, err := uc.Execute(context.Background(), r.ID); err == nil {
		t.Error("expected error when starting with fewer than 2 players")
	}
}

func TestStartGame_AlreadyStarted(t *testing.T) {
	repo := newFakeRoomRepo()
	r := model.NewRoom("test", 1)
	r.MemberIDs = []int64{1, 2}
	r.StartGame() // すでにゲーム開始済み
	_ = repo.SaveRoom(context.Background(), r)

	uc := &StartGameInteractor{RoomRepository: repo}
	if _, err := uc.Execute(context.Background(), r.ID); err == nil {
		t.Error("expected error when the game is already started")
	}
}

func TestStartGame_RoomNotFound(t *testing.T) {
	repo := newFakeRoomRepo()
	uc := &StartGameInteractor{RoomRepository: repo}
	if _, err := uc.Execute(context.Background(), 999); err == nil {
		t.Error("expected error for non-existent room")
	}
}
