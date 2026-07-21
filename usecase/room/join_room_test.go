package room

import (
	"context"
	"testing"

	"github.com/ne241099/daifugo-server/model"
)

func TestJoinRoom_AddsMember(t *testing.T) {
	repo := newFakeRoomRepo()
	r := model.NewRoom("test", 1) // owner=1, members=[1]
	_ = repo.SaveRoom(context.Background(), r)

	uc := &JoinRoomInteractor{RoomRepository: repo}
	joined, err := uc.Execute(context.Background(), r.ID, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(joined.MemberIDs) != 2 {
		t.Fatalf("members = %v, want 2 members", joined.MemberIDs)
	}
	if joined.MemberIDs[1] != 2 {
		t.Errorf("second member = %d, want 2", joined.MemberIDs[1])
	}
}

func TestJoinRoom_AlreadyMemberIsNoop(t *testing.T) {
	repo := newFakeRoomRepo()
	r := model.NewRoom("test", 1)
	_ = repo.SaveRoom(context.Background(), r)

	uc := &JoinRoomInteractor{RoomRepository: repo}
	joined, err := uc.Execute(context.Background(), r.ID, 1) // owner rejoins
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(joined.MemberIDs) != 1 {
		t.Errorf("members = %v, want unchanged (1 member)", joined.MemberIDs)
	}
}

func TestJoinRoom_FullRoom(t *testing.T) {
	repo := newFakeRoomRepo()
	r := model.NewRoom("test", 1)
	r.MemberIDs = []int64{1, 2, 3, 4} // 満室(4人)
	_ = repo.SaveRoom(context.Background(), r)

	uc := &JoinRoomInteractor{RoomRepository: repo}
	_, err := uc.Execute(context.Background(), r.ID, 5)
	if err == nil {
		t.Error("expected error when joining a full room")
	}
}

func TestJoinRoom_NotFound(t *testing.T) {
	repo := newFakeRoomRepo()
	uc := &JoinRoomInteractor{RoomRepository: repo}
	_, err := uc.Execute(context.Background(), 999, 2)
	if err == nil {
		t.Error("expected error for non-existent room")
	}
}
