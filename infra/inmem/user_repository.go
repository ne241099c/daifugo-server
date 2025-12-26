package inmem

import (
	"context"
	"sync"

	"github.com/ne241099/daifugo-server/model"
	"github.com/ne241099/daifugo-server/repository"
)

var _ repository.UserRepository = &InmemUserRepository{}

type InmemUserRepository struct {
	mtx    sync.RWMutex
	data   map[int64]model.User
	emails map[string]int64
	number int64
}

func NewInmemUserRepository() *InmemUserRepository {
	return &InmemUserRepository{
		data:   make(map[int64]model.User),
		emails: make(map[string]int64),
	}
}

func (r *InmemUserRepository) GetUser(ctx context.Context, id int64) (*model.User, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	user, ok := r.data[id]
	if !ok {
		return nil, repository.ErrEntityNotFound
	}
	return &user, nil
}

func (r *InmemUserRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	id, ok := r.emails[email]
	if !ok {
		return nil, repository.ErrEntityNotFound
	}

	user, ok := r.data[id]
	if !ok {
		// インデックスにあるのにデータがない
		return nil, repository.ErrEntityNotFound
	}
	return &user, nil
}

func (r *InmemUserRepository) SaveUser(ctx context.Context, user *model.User) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	r.mtx.Lock()
	defer r.mtx.Unlock()

	if user.ID == 0 {
		// 新規作成
		r.number++
		user.ID = r.number
	} else {
		// 更新の
		// Emailが変更されていたら、古いインデックスを削除
		if oldUser, ok := r.data[user.ID]; ok {
			if oldUser.Email != user.Email {
				delete(r.emails, oldUser.Email)
			}
		}
	}

	// 保存
	r.data[user.ID] = *user
	r.emails[user.Email] = user.ID
	return nil
}

func (r *InmemUserRepository) DeleteUser(ctx context.Context, id int64) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	user, ok := r.data[id]
	if !ok {
		// 存在しない場合はエラーを返すのが親切
		return repository.ErrEntityNotFound
	}

	// Emailインデックスから削除
	delete(r.emails, user.Email)

	// データ本体から削除
	delete(r.data, id)

	return nil
}
