package user

import (
	"context"
	"errors"
	"time"

	"github.com/ne241099/daifugo-server/model"
	"github.com/ne241099/daifugo-server/repository"
	"github.com/ne241099/daifugo-server/usecase"
)

// SignUpInput はサインアップのユースケース入力
// プレゼンテーション層の型に依存しないよう、ユースケース層で定義する
type SignUpInput struct {
	Name     string
	Email    string
	Password string
}

type SignUpUseCase interface {
	Execute(ctx context.Context, input SignUpInput) (*model.User, error)
}

var _ SignUpUseCase = &SignUpInteractor{}

type SignUpInteractor struct {
	UserRepository repository.UserRepository
}

func (uc *SignUpInteractor) Execute(ctx context.Context, input SignUpInput) (*model.User, error) {
	// 重複チェック
	u, err := uc.UserRepository.GetUserByEmail(ctx, input.Email)
	if u != nil {
		return nil, usecase.ErrDuplicateEntity
	}

	if !errors.Is(err, repository.ErrEntityNotFound) {
		return nil, err
	}
	now := time.Now()

	// ユーザ作成
	user := &model.User{}
	if err := user.Create(model.CreateUserParam{
		Email:     input.Email,
		Password:  input.Password,
		Name:      input.Name,
		CreatedAt: now,
		UpdatedAt: now,
	}); err != nil {
		return nil, err
	}

	// ユーザ保存
	if err := uc.UserRepository.SaveUser(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}
