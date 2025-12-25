package user

import (
	"context"

	"github.com/ne241099/daifugo-server/model"
	"github.com/ne241099/daifugo-server/repository"
)

type GetUserUseCase interface {
	Execute(context.Context, int64) (*model.User, error)
}

var _ GetUserUseCase = &GetUserInteractor{}

type GetUserInteractor struct {
	UserRepository repository.UserRepository
}

func NewGetUserUseCase(r repository.UserRepository) GetUserUseCase {
	return &GetUserInteractor{UserRepository: r}
}

func (u *GetUserInteractor) Execute(ctx context.Context, id int64) (*model.User, error) {
	return u.UserRepository.FindByID(ctx, id)
}
