package user

import (
	"context"
	"errors"

	gqlmodel "github.com/ne241099/daifugo-server/graph/model"
	"github.com/ne241099/daifugo-server/model"
	"github.com/ne241099/daifugo-server/repository"
	"github.com/ne241099/daifugo-server/usecase"
)

type SignUpUseCase interface {
	Execute(context.Context, gqlmodel.SignUpInput) (*model.User, error)
}

var _ SignUpUseCase = &SignUpInteractor{}

type SignUpInteractor struct {
	UserRepository repository.UserRepository
}

func (uc *SignUpInteractor) Execute(ctx context.Context, input gqlmodel.SignUpInput) (*model.User, error) {
	// 重複チェック
	u, err := uc.UserRepository.GetUserByEmail(ctx, input.Email)
	if u != nil {
		return nil, errors.Join(usecase.ErrDuplicateEntity)
	}

	if !errors.Is(err, repository.ErrEntityNotFound) {
		return nil, errors.Join(err)
	}

	// ユーザ作成
	user := &model.User{}
	if err := user.Create(model.CreateUserParam{
		Email:    input.Email,
		Password: input.Password,
		Name:     input.Name,
	}); err != nil {
		return nil, errors.Join(err)
	}

	// ユーザ保存
	if err := uc.UserRepository.SaveUser(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}
