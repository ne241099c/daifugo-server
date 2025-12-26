package middleware

import (
	"net/http"
	"strings"

	"github.com/ne241099/daifugo-server/internal/auth"
)

type AuthMiddleware struct {
	authenticator auth.Authenticator
}

// コンストラクタで Authenticator を受け取る
func NewAuthMiddleware(authenticator auth.Authenticator) *AuthMiddleware {
	return &AuthMiddleware{
		authenticator: authenticator,
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

		// インターフェース経由で検証を実行
		uid, err := m.authenticator.VerifyToken(r.Context(), token)
		if err != nil {
			// 検証失敗時は未認証扱いとして通す
			next.ServeHTTP(w, r)
			return
		}

		// Contextに埋め込む
		ctx := auth.WithUserID(r.Context(), uid)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
