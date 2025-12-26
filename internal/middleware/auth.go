package middleware

import (
	"net/http"
	"strconv"

	"github.com/ne241099/daifugo-server/internal/auth"
)

// AuthMiddleware はリクエストヘッダー等からユーザーを特定し、Contextに埋め込む
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// ヘッダーからユーザーIDを取得
		// 本来は Cookie や Authorization ヘッダーのトークンを検証します
		uidStr := r.Header.Get("X-Showcase-User")

		// デフォルト値
		if uidStr == "" {
			uidStr = "1" // 開発中は便宜上 1 とする
		}

		uid, err := strconv.ParseInt(uidStr, 10, 64)
		if err != nil {
			// IDが不正な場合は401エラーなどを返すか、無視してゲスト扱いにする
			// ここでは無視してそのままnextを呼び出す（GetUserIDでエラーになる）
			next.ServeHTTP(w, r)
			return
		}

		// Contextに埋め込む
		ctx := auth.WithUserID(r.Context(), uid)

		// 新しいContextを持ったRequestで次の処理へ
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
