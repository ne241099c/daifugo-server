package game

import (
	"context"
	"testing"

	"github.com/ne241099/daifugo-server/model"
)

// startedRoom は2人でゲーム開始済みの部屋を用意する
func startedRoom(t *testing.T, repo *fakeRoomRepo) *model.Room {
	t.Helper()
	r := model.NewRoom("test", 1)
	r.MemberIDs = []int64{1, 2}
	_ = repo.SaveRoom(context.Background(), r)

	got, err := (&StartGameInteractor{RoomRepository: repo}).Execute(context.Background(), r.ID)
	if err != nil {
		t.Fatalf("failed to start game: %v", err)
	}
	return got
}

func TestPlayCard_Success(t *testing.T) {
	repo := newFakeRoomRepo()
	r := startedRoom(t, repo)

	// 手番プレイヤー（Turn==0）の手札から1枚を出す。場が空なので単騎は必ず有効。
	cur := r.Game.Players[r.Game.Turn]
	before := len(cur.Hand)
	card := cur.Hand[0]

	uc := &PlayCardInteractor{RoomRepository: repo}
	got, err := uc.Execute(context.Background(), r.ID, cur.UserID, []int{card.ID})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 出したプレイヤーの手札が1枚減っていること
	played := got.Game.Players[0]
	if played.UserID == cur.UserID && len(played.Hand) != before-1 {
		t.Errorf("hand size = %d, want %d", len(played.Hand), before-1)
	}
}

func TestPlayCard_GameNotStarted(t *testing.T) {
	repo := newFakeRoomRepo()
	r := model.NewRoom("test", 1)
	r.MemberIDs = []int64{1, 2}
	_ = repo.SaveRoom(context.Background(), r)

	uc := &PlayCardInteractor{RoomRepository: repo}
	if _, err := uc.Execute(context.Background(), r.ID, 1, []int{1}); err == nil {
		t.Error("expected error when game is not started")
	}
}

func TestPlayCard_PlayerNotInGame(t *testing.T) {
	repo := newFakeRoomRepo()
	r := startedRoom(t, repo)

	uc := &PlayCardInteractor{RoomRepository: repo}
	if _, err := uc.Execute(context.Background(), r.ID, 999, []int{1}); err == nil {
		t.Error("expected error when user is not a player in the game")
	}
}

func TestPlayCard_CardNotInHand(t *testing.T) {
	repo := newFakeRoomRepo()
	r := startedRoom(t, repo)

	cur := r.Game.Players[r.Game.Turn]
	uc := &PlayCardInteractor{RoomRepository: repo}
	// ID 9999 のカードは誰の手札にも存在しない
	if _, err := uc.Execute(context.Background(), r.ID, cur.UserID, []int{9999}); err == nil {
		t.Error("expected error when a card is not in the player's hand")
	}
}
