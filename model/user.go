package model

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             int64
	Email          string
	HashedPassword string `jason:"-"`
	Name           string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type CreateUserParam struct {
	Email    string
	Password string
	Name     string
}

func (u *User) Create(params CreateUserParam) error {
	hp, err := hashPassword(params.Password)
	if err != nil {
		return errors.Join(err)
	}
	u.Email = params.Email
	u.HashedPassword = hp
	u.Name = params.Name

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
