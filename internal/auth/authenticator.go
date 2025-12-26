package auth

import "context"

// Authenticator はトークンを検証してユーザーIDを特定する責務を持ちます
type Authenticator interface {
	VerifyToken(ctx context.Context, token string) (int64, error)
}
