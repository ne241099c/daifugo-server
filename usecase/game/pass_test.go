package game

import (
	"context"
	"testing"

	"github.com/ne241099/daifugo-server/model"
)

func TestPass_Success(t *testing.T) {
	repo := newFakeRoomRepo()
	r := startedRoom(t, repo)

	turnBefore := r.Game.Turn
	cur := r.Game.Players[turnBefore]

	uc := &PassInteractor{RoomRepository: repo}
	got, err := uc.Execute(context.Background(), r.ID, cur.UserID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// パスすると手番が次のプレイヤーへ移る
	// （2人対戦では全員パス扱いで場が流れ PassCount は 0 に戻るため、
	//   観測点はターンの移動にする）
	if got.Game.Turn == turnBefore {
		t.Errorf("Turn = %d, want it to advance from %d", got.Game.Turn, turnBefore)
	}
}

func TestPass_GameNotStarted(t *testing.T) {
	repo := newFakeRoomRepo()
	r := model.NewRoom("test", 1)
	r.MemberIDs = []int64{1, 2}
	_ = repo.SaveRoom(context.Background(), r)

	uc := &PassInteractor{RoomRepository: repo}
	if _, err := uc.Execute(context.Background(), r.ID, 1); err == nil {
		t.Error("expected error when game is not started")
	}
}

func TestPass_WrongTurn(t *testing.T) {
	repo := newFakeRoomRepo()
	r := startedRoom(t, repo)

	// 手番でないプレイヤーを選ぶ
	notCurrent := r.Game.Players[1-r.Game.Turn]
	uc := &PassInteractor{RoomRepository: repo}
	if _, err := uc.Execute(context.Background(), r.ID, notCurrent.UserID); err == nil {
		t.Error("expected error when passing out of turn")
	}
}
