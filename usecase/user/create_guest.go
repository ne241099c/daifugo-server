package user

import (
	"context"
	"errors"
	"strings"

	"github.com/ne241099/daifugo-server/model"
	"github.com/ne241099/daifugo-server/repository"
)

type CreateGuestUseCase interface {
	Execute(ctx context.Context, name string) (*model.User, error)
}

var _ CreateGuestUseCase = &CreateGuestInteractor{}

type CreateGuestInteractor struct {
	UserRepository repository.UserRepository
}

func (uc *CreateGuestInteractor) Execute(ctx context.Context, name string) (*model.User, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, errors.New("name is required")
	}

	user := &model.User{}
	user.CreateGuest(name)

	if err := uc.UserRepository.SaveUser(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}
