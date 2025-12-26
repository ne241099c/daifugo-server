package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/ne241099/daifugo-server/internal/auth"
	"github.com/ne241099/daifugo-server/model"
	"github.com/ne241099/daifugo-server/repository"
	"golang.org/x/crypto/bcrypt"
)

type LoginUseCase interface {
	Execute(ctx context.Context, email, password string) (string, *model.User, error)
}

var _ LoginUseCase = &LoginInteractor{}

type LoginInteractor struct {
	// UserRepository ユーザリポジトリ
	UserRepository repository.UserRepository
	// Authenticator 認証サービス
	Authenticator auth.Authenticator
}

func (uc *LoginInteractor) Execute(ctx context.Context, email, password string) (string, *model.User, error) {
	u, err := uc.UserRepository.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrEntityNotFound) {
			return "", nil, errors.New("invalid email or password") // セキュリティのため詳細は伏せる
		}
		return "", nil, fmt.Errorf("failed to get user: %w", err)
	}

	// パスワードの検証
	if err := bcrypt.CompareHashAndPassword([]byte(u.HashedPassword), []byte(password)); err != nil {
		return "", nil, errors.New("invalid email or password")
	}

	// トークンの生成
	token, err := uc.Authenticator.CreateToken(ctx, u.ID)
	if err != nil {
		return "", nil, fmt.Errorf("failed to create token: %w", err)
	}

	return token, u, nil
}
