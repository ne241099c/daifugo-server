package inmem

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/ne241099/daifugo-server/model"
	"github.com/ne241099/daifugo-server/repository"
)

// TestWithLock_SerializesMutations は WithLock が同一部屋への
// 「取得→変更→保存」を直列化することを確認する（-race で競合検出）。
func TestWithLock_SerializesMutations(t *testing.T) {
	repo := NewInmemRoomRepository()
	r := model.NewRoom("t", 1) // members=[1]
	_ = repo.SaveRoom(context.Background(), r)

	const n = 100
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			release := repo.WithLock(r.ID)
			defer release()

			got, _ := repo.GetRoomByID(context.Background(), r.ID)
			got.MemberIDs = append(got.MemberIDs, 99)
			_ = repo.SaveRoom(context.Background(), got)
		}()
	}
	wg.Wait()

	got, _ := repo.GetRoomByID(context.Background(), r.ID)
	if len(got.MemberIDs) != 1+n {
		t.Errorf("members = %d, want %d (mutations should be serialized)", len(got.MemberIDs), 1+n)
	}
}

// TestCleanupRooms_DeletesBotOnlyRoom は、人間がいない部屋（CPUのみ）が
// 放置期間に関わらず即削除され、人間がいる部屋は残ることを確認する。
func TestCleanupRooms_DeletesBotOnlyRoom(t *testing.T) {
	repo := NewInmemRoomRepository()

	// 人間がいる部屋（members=[1], bot なし）
	human := model.NewRoom("human", 1)
	_ = repo.SaveRoom(context.Background(), human)

	// CPU のみの部屋（member 2 は bot）
	botOnly := model.NewRoom("bot", 2)
	botOnly.BotIDs = []int64{2}
	_ = repo.SaveRoom(context.Background(), botOnly)

	// 放置しきい値を長く取っても、CPUのみ部屋は削除される
	repo.CleanupRooms(24 * time.Hour)

	if _, err := repo.GetRoomByID(context.Background(), botOnly.ID); !errors.Is(err, repository.ErrEntityNotFound) {
		t.Error("CPUのみの部屋は放置期間に関わらず削除されるべき")
	}
	if _, err := repo.GetRoomByID(context.Background(), human.ID); err != nil {
		t.Errorf("人間がいる部屋は残るべき: %v", err)
	}
}

// TestCleanupRooms_ConcurrentSafe は掃除と部屋の変更が並行しても
// レースしないことを確認する（-race で検出）。
func TestCleanupRooms_ConcurrentSafe(t *testing.T) {
	repo := NewInmemRoomRepository()
	r := model.NewRoom("t", 1)
	r.MemberIDs = []int64{1, 2} // 人間がいるので掃除では消えない
	_ = repo.SaveRoom(context.Background(), r)

	done := make(chan struct{})
	go func() {
		defer close(done)
		for i := 0; i < 300; i++ {
			release := repo.WithLock(r.ID)
			got, _ := repo.GetRoomByID(context.Background(), r.ID)
			got.MemberIDs = append(got.MemberIDs, 3)
			got.MemberIDs = got.MemberIDs[:2]
			_ = repo.SaveRoom(context.Background(), got)
			release()
		}
	}()

	for i := 0; i < 300; i++ {
		repo.CleanupRooms(24 * time.Hour)
	}
	<-done

	if _, err := repo.GetRoomByID(context.Background(), r.ID); err != nil {
		t.Errorf("人間がいる部屋は掃除で消えないはず: %v", err)
	}
}

// TestListRooms_ConcurrentWithMutation は ListRooms のディープコピーが
// ゲーム状態の変更と競合しないことを確認する（-race で検出）。
func TestListRooms_ConcurrentWithMutation(t *testing.T) {
	repo := NewInmemRoomRepository()
	r := model.NewRoom("t", 1)
	r.MemberIDs = []int64{1, 2}
	r.StartGame() // Game を生成して変更対象にする
	_ = repo.SaveRoom(context.Background(), r)

	done := make(chan struct{})
	go func() {
		defer close(done)
		for i := 0; i < 300; i++ {
			release := repo.WithLock(r.ID)
			got, _ := repo.GetRoomByID(context.Background(), r.ID)
			if got.Game != nil {
				got.Game.Turn++
			}
			_ = repo.SaveRoom(context.Background(), got)
			release()
		}
	}()

	for i := 0; i < 300; i++ {
		if _, err := repo.ListRooms(context.Background()); err != nil {
			t.Fatalf("ListRooms error: %v", err)
		}
	}
	<-done
}
