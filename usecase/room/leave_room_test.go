package room

import (
	"context"
	"errors"
	"testing"

	"github.com/ne241099/daifugo-server/model"
	"github.com/ne241099/daifugo-server/repository"
)

func TestLeaveRoom_RemovesMember(t *testing.T) {
	repo := newFakeRoomRepo()
	r := model.NewRoom("test", 1)
	r.MemberIDs = []int64{1, 2}
	_ = repo.SaveRoom(context.Background(), r)

	uc := &LeaveRoomInteractor{RoomRepository: repo}
	if err := uc.Execute(context.Background(), r.ID, 2); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, _ := repo.GetRoomByID(context.Background(), r.ID)
	if len(got.MemberIDs) != 1 || got.MemberIDs[0] != 1 {
		t.Errorf("members = %v, want [1]", got.MemberIDs)
	}
}

func TestLeaveRoom_LastMemberDeletesRoom(t *testing.T) {
	repo := newFakeRoomRepo()
	r := model.NewRoom("test", 1) // members=[1]
	_ = repo.SaveRoom(context.Background(), r)

	uc := &LeaveRoomInteractor{RoomRepository: repo}
	if err := uc.Execute(context.Background(), r.ID, 1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := repo.GetRoomByID(context.Background(), r.ID); !errors.Is(err, repository.ErrEntityNotFound) {
		t.Error("empty room should be deleted")
	}
}

func TestLeaveRoom_OwnerLeavesTransfersOwnership(t *testing.T) {
	repo := newFakeRoomRepo()
	r := model.NewRoom("test", 1)
	r.MemberIDs = []int64{1, 2}
	_ = repo.SaveRoom(context.Background(), r)

	uc := &LeaveRoomInteractor{RoomRepository: repo}
	if err := uc.Execute(context.Background(), r.ID, 1); err != nil { // owner leaves
		t.Fatalf("unexpected error: %v", err)
	}

	got, _ := repo.GetRoomByID(context.Background(), r.ID)
	if got.OwnerID != 2 {
		t.Errorf("ownerID = %d, want 2 (transferred to remaining member)", got.OwnerID)
	}
}

func TestLeaveRoom_NotAMember(t *testing.T) {
	repo := newFakeRoomRepo()
	r := model.NewRoom("test", 1)
	r.MemberIDs = []int64{1, 2}
	_ = repo.SaveRoom(context.Background(), r)

	uc := &LeaveRoomInteractor{RoomRepository: repo}
	if err := uc.Execute(context.Background(), r.ID, 99); err == nil {
		t.Error("expected error when a non-member leaves")
	}
}
