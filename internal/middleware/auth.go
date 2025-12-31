package middleware

import (
	"net/http"
	"strings"

	"github.com/ne241099/daifugo-server/internal/auth"
	"github.com/ne241099/daifugo-server/repository"
)

type AuthMiddleware struct {
	authenticator auth.Authenticator
	userRepo      repository.UserRepository
}

// コンストラクタで Authenticator を受け取る
func NewAuthMiddleware(authenticator auth.Authenticator, userRepo repository.UserRepository) *AuthMiddleware {
	return &AuthMiddleware{
		authenticator: authenticator,
		userRepo:      userRepo,
	}
}

func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Authorizationヘッダーの取得
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			next.ServeHTTP(w, r)
			return
		}

		// "Bearer <token>" 形式の解析
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			next.ServeHTTP(w, r)
			return
		}
		token := parts[1]

		uid, tokenVer, err := m.authenticator.VerifyToken(r.Context(), token)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		user, err := m.userRepo.GetUser(r.Context(), uid)
		if err != nil {
			// ユーザーが存在しない
			http.Error(w, "User not found", http.StatusUnauthorized)
			return
		}

		if user.TokenVersion != tokenVer {
			http.Error(w, "Session expired (Logged in on another device)", http.StatusUnauthorized)
			return
		}

		// Contextに埋め込む
		ctx := auth.WithUserID(r.Context(), uid)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
