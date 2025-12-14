package repository

import (
	"context"

	"github.com/ne241099/daifugo-server/model"
)

type UserRepository interface {
	GetUser(ctx context.Context, id int64) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	SaveUser(ctx context.Context, user *model.User) error
}
