package inmem

import (
	"context"
	"sync"
	"testing"

	"github.com/ne241099/daifugo-server/model"
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
