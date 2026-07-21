package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/ne241099/daifugo-server/model"
	"github.com/ne241099/daifugo-server/repository"
	"golang.org/x/crypto/bcrypt"
)

// TokenGenerator は認証トークンを発行するポート。
// 具体的な実装（JWT など）は infra/adapter 層が提供する。
type TokenGenerator interface {
	CreateToken(ctx context.Context, userID int64, version int) (string, error)
}

type LoginUseCase interface {
	Execute(ctx context.Context, email, password string) (string, *model.User, error)
}

var _ LoginUseCase = &LoginInteractor{}

type LoginInteractor struct {
	// UserRepository ユーザリポジトリ
	UserRepository repository.UserRepository
	// TokenGenerator トークン発行サービス
	TokenGenerator TokenGenerator
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

	newVersion, err := uc.UserRepository.IncrementTokenVersion(ctx, u.ID)
	if err != nil {
		return "", nil, fmt.Errorf("failed to increment token version: %w", err)
	}

	u.TokenVersion = newVersion

	// トークンの生成
	token, err := uc.TokenGenerator.CreateToken(ctx, u.ID, newVersion)
	if err != nil {
		return "", nil, fmt.Errorf("failed to create token: %w", err)
	}

	return token, u, nil
}
