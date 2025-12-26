package user

import (
	"context"
	"fmt"

	"github.com/ne241099/daifugo-server/repository"
)

type DeleteUserUseCase interface {
	Execute(ctx context.Context, id int64) error
}

type DeleteUserInteractor struct {
	UserRepository repository.UserRepository
}

func (uc *DeleteUserInteractor) Execute(ctx context.Context, id int64) error {
	_, err := uc.UserRepository.GetUser(ctx, id)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// 削除実行
	if err := uc.UserRepository.DeleteUser(ctx, id); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}
