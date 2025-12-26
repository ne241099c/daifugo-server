package auth

import (
	"context"
	"errors"
)

type contextKey string

const userIDKey contextKey = "userID"

// GetUserID はContextからユーザーIDを取得する
func GetUserID(ctx context.Context) (int64, error) {
	v := ctx.Value(userIDKey)
	id, ok := v.(int64)
	if !ok {
		return 0, errors.New("user id not found in context")
	}
	return id, nil
}

// WithUserID はContextにユーザーIDをセットした新しいContextを返す
func WithUserID(ctx context.Context, userID int64) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}
