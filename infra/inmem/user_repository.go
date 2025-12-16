package inmem

import (
	"context"
	"errors"
	"sync"

	"github.com/ne241099/daifugo-server/model"
	"github.com/ne241099/daifugo-server/repository"
)

var _ repository.UserRepository = &InmemUserRepository{}

type InmemUserRepository struct {
	mtx    sync.RWMutex
	data   map[int64]model.User
	number int64
}

func NewInmemUserRepository() *InmemUserRepository {
	return &InmemUserRepository{
		data: make(map[int64]model.User),
	}
}

func (r *InmemUserRepository) GetUser(ctx context.Context, id int64) (*model.User, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	user, ok := r.data[id]
	if !ok {
		return nil, errors.Join(repository.ErrEntityNotFound)
	}
	return &user, nil
}

func (r *InmemUserRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	for _, user := range r.data {
		if user.Email == email {
			return &user, nil
		}
	}
	return nil, errors.Join(repository.ErrEntityNotFound)
}

func (r *InmemUserRepository) SaveUser(ctx context.Context, user *model.User) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	if user.ID == 0 {
		r.number++
		user.ID = r.number
	}
	r.data[user.ID] = *user
	return nil
}
