package model

import (
	"errors"
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

func (u *User) Create(params CreateUserParam) error {
	hp, err := hashPassword(params.Password)
	if err != nil {
		return errors.Join(err)
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
		return "", errors.Join(err)
	}

	// バイト列を文字列に戻して返します
	return string(hashedBytes), nil

}
