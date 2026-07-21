package user

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/ne241099/daifugo-server/model"
	"github.com/ne241099/daifugo-server/repository"
)

type RenameUserUseCase interface {
	Execute(ctx context.Context, userID int64, name string) (*model.User, error)
}

var _ RenameUserUseCase = &RenameUserInteractor{}

type RenameUserInteractor struct {
	UserRepository repository.UserRepository
}

func (uc *RenameUserInteractor) Execute(ctx context.Context, userID int64, name string) (*model.User, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, errors.New("name is required")
	}

	user, err := uc.UserRepository.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	user.Name = name
	user.UpdatedAt = time.Now()

	if err := uc.UserRepository.SaveUser(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}
