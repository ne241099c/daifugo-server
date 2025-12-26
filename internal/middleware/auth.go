package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/ne241099/daifugo-server/internal/auth"
)

type AuthMiddleware struct {
	// 今回は外部依存がないためフィールドなし
}

func NewAuthMiddleware() *AuthMiddleware {
	return &AuthMiddleware{}
}

func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Authorizationヘッダーの取得
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			// ヘッダーがない場合は未認証としてそのまま通す
			next.ServeHTTP(w, r)
			return
		}

		// "Bearer <token>" 形式の解析
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			// 形式不正
			next.ServeHTTP(w, r)
			return
		}
		token := parts[1]

		// トークンの検証とユーザーIDの特定
		// 【本来の形式】
		// ここで JWT の Verify や、DB の Session Token 照合を行います。
		// 例: uid, err := m.authClient.VerifyIDToken(ctx, token)
		uid, err := m.verifyTokenDummy(token)
		if err != nil {
			// トークンが無効な場合
			// 厳密に 401 を返すか、未認証として通すかは要件によりますが、
			// ここでは未認証（Contextにセットしない）として通します。
			next.ServeHTTP(w, r)
			return
		}

		// Contextに埋め込む
		ctx := auth.WithUserID(r.Context(), uid)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *AuthMiddleware) verifyTokenDummy(token string) (int64, error) {
	// 簡易実装: "token:" というプレフィックスを除去してID化する
	// 実際にはこんなことはせず、暗号化されたトークンを検証します
	if strings.HasPrefix(token, "token:") {
		idStr := strings.TrimPrefix(token, "token:")
		return strconv.ParseInt(idStr, 10, 64)
	}
	// 検証失敗
	return 0, http.ErrNoCookie
}
