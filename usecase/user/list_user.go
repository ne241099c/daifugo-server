package user

import (
	"context"
	"fmt"

	"github.com/ne241099/daifugo-server/model"
	"github.com/ne241099/daifugo-server/repository"
)

type ListUsersUseCase interface {
	Execute(ctx context.Context) ([]*model.User, error)
}

var _ ListUsersUseCase = &ListUsersInteractor{}

type ListUsersInteractor struct {
	UserRepository repository.UserRepository
}

func (uc *ListUsersInteractor) Execute(ctx context.Context) ([]*model.User, error) {
	users, err := uc.UserRepository.ListUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	return users, nil
}
