package mysql

// import (
// 	"context"
// 	"database/sql"
// 	"errors"

// 	"github.com/ne241099/daifugo-server/model"
// 	"github.com/ne241099/daifugo-server/repository"
// )

// var _ repository.UserRepository = &MySQLUserRepository{}

// type MySQLUserRepository struct {
// 	db *sql.DB
// }

// func NewMySQLUserRepository() *MySQLUserRepository {
// 	return &MySQLUserRepository{
// 		data: make(map[int64]*model.User),
// 	}
// }

// func (r *MySQLUserRepository) GetUser(ctx context.Context, id int64) (*model.User, error) {
// 	r.mtx.RLock()
// 	defer r.mtx.RUnlock()

// 	user, ok := r.data[id]
// 	if !ok {
// 		return nil, errors.Join(repository.ErrEntityNotFound)
// 	}
// 	return user, nil
// }

// func (r *MySQLUserRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
// 	r.mtx.RLock()
// 	defer r.mtx.RUnlock()

// 	for _, user := range r.data {
// 		if user.Email == email {
// 			return user, nil
// 		}
// 	}
// 	return nil, errors.Join(repository.ErrEntityNotFound)
// }

// func (r *MySQLUserRepository) SaveUser(ctx context.Context, user *model.User) error {
// 	r.mtx.Lock()
// 	defer r.mtx.Unlock()

// 	if user.ID == 0 {
// 		r.number++
// 		user.ID = r.number
// 	}
// 	r.data[user.ID] = user
// 	return nil
// }
