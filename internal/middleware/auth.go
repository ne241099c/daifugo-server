package middleware

import (
	"net/http"
	"strconv"

	"github.com/ne241099/daifugo-server/internal/auth"
	"github.com/ne241099/daifugo-server/repository"
)

// userIDHeader はゲストのユーザーIDを載せるリクエストヘッダー名。
// ログイン/JWT を廃止し、クライアントが保持するゲストIDをそのまま識別子として使う。
const userIDHeader = "X-User-ID"

type AuthMiddleware struct {
	userRepo repository.UserRepository
}

func NewAuthMiddleware(userRepo repository.UserRepository) *AuthMiddleware {
	return &AuthMiddleware{userRepo: userRepo}
}

func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw := r.Header.Get(userIDHeader)
		if raw == "" {
			next.ServeHTTP(w, r)
			return
		}

		uid, err := strconv.ParseInt(raw, 10, 64)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		// 存在しないIDは無視する（無認証で続行し、認証が要るリゾルバ側で弾く）
		if _, err := m.userRepo.GetUser(r.Context(), uid); err != nil {
			next.ServeHTTP(w, r)
			return
		}

		ctx := auth.WithUserID(r.Context(), uid)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
