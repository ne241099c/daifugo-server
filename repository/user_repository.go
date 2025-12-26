package repository

import (
	"context"

	"github.com/ne241099/daifugo-server/model"
)

type UserRepository interface {
	// GetUser は、ユーザを取得する
	GetUser(ctx context.Context, id int64) (*model.User, error)
	// GetUserByEmail は、e-mailでユーザを取得する
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	// SaveUser は、ユーザを保存する
	SaveUser(ctx context.Context, user *model.User) error
	// DeleteUser は、ユーザを削除する
	DeleteUser(ctx context.Context, id int64) error
}
