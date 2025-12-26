package auth

import (
	"context"
	"errors"
	"strconv"
	"strings"
)

// DummyAuthenticator は開発用の簡易認証実装
type DummyAuthenticator struct{}

func NewDummyAuthenticator() *DummyAuthenticator {
	return &DummyAuthenticator{}
}

func (a *DummyAuthenticator) VerifyToken(ctx context.Context, token string) (int64, error) {
	// "token:" というプレフィックスを除去してID化する
	if strings.HasPrefix(token, "token:") {
		idStr := strings.TrimPrefix(token, "token:")
		return strconv.ParseInt(idStr, 10, 64)
	}
	return 0, errors.New("invalid token")
}
