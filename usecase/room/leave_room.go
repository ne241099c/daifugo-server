package room

import (
	"context"
	"fmt"

	"github.com/ne241099/daifugo-server/repository"
)

type LeaveRoomUseCase interface {
	Execute(ctx context.Context, roomID int64, userID int64) error
}

var _ LeaveRoomUseCase = &LeaveRoomInteractor{}

type LeaveRoomInteractor struct {
	RoomRepository repository.RoomRepository
}

func (uc *LeaveRoomInteractor) Execute(ctx context.Context, roomID int64, userID int64) error {
	// 部屋情報を取得
	room, err := uc.RoomRepository.GetRoomByID(ctx, roomID)
	if err != nil {
		return fmt.Errorf("room not found: %w", err)
	}
	room.Mu.Lock()
	defer room.Mu.Unlock()

	// メンバーリストからユーザーを削除
	newMembers := make([]int64, 0, len(room.MemberIDs))
	found := false

	for _, mid := range room.MemberIDs {
		if mid == userID {
			found = true
			continue
		}
		newMembers = append(newMembers, mid)
	}

	if !found {
		return fmt.Errorf("user is not in the room")
	}
	room.MemberIDs = newMembers

	// 部屋が空になった場合は削除
	if len(room.MemberIDs) == 0 {
		if err := uc.RoomRepository.DeleteRoom(ctx, roomID); err != nil {
			return fmt.Errorf("failed to delete empty room: %w", err)
		}
		return nil
	}

	// オーナーが退出した場合は新しいオーナーを設定
	if room.OwnerID == userID && len(room.MemberIDs) > 0 {
		room.OwnerID = room.MemberIDs[0]
	}

	// 更新を保存
	if err := uc.RoomRepository.UpdateRoom(ctx, room); err != nil {
		return fmt.Errorf("failed to update room: %w", err)
	}

	return nil
}
