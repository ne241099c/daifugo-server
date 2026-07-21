package model

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             int64     `json:"id"`
	Email          string    `json:"email"`
	HashedPassword string    `json:"-"`
	Name           string    `json:"name"`
	TokenVersion   int       `json:"token_version"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type CreateUserParam struct {
	Email     string
	Password  string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// CreateGuest は名前だけのゲストユーザーを初期化する。
// パスワード認証は行わないため、Email はリポジトリのインデックス衝突を避ける
// ためだけのランダム値を割り当てる。
func (u *User) CreateGuest(name string) {
	now := time.Now()
	u.Name = name
	u.Email = "guest-" + randomToken() + "@local"
	u.HashedPassword = ""
	u.CreatedAt = now
	u.UpdatedAt = now
}

func randomToken() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		// 乱数取得に失敗した場合でも一意性を保つためナノ秒を使う
		return hex.EncodeToString([]byte(time.Now().Format(time.RFC3339Nano)))
	}
	return hex.EncodeToString(b)
}

func (u *User) Create(params CreateUserParam) error {
	hp, err := hashPassword(params.Password)
	if err != nil {
		return err
	}
	u.Email = params.Email
	u.HashedPassword = hp
	u.Name = params.Name
	u.CreatedAt = params.CreatedAt
	u.UpdatedAt = params.UpdatedAt

	return nil
}

func hashPassword(password string) (string, error) {
	// クリプトでハッシュ化する
	// パスワードをバイト列に変換し、デフォルトのコストでハッシュ化します
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	// バイト列を文字列に戻して返します
	return string(hashedBytes), nil
}
