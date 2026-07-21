package user

import (
	"context"
	"errors"
	"testing"

	"github.com/ne241099/daifugo-server/usecase"
)

func TestSignUp_Success(t *testing.T) {
	repo := newFakeUserRepo()
	uc := &SignUpInteractor{UserRepository: repo}

	u, err := uc.Execute(context.Background(), SignUpInput{
		Name:     "new user",
		Email:    "new@example.com",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if u.ID == 0 {
		t.Error("saved user should have an assigned ID")
	}
	if u.HashedPassword == "" || u.HashedPassword == "password123" {
		t.Error("password should be hashed, not stored in plain text")
	}
	// 実際に保存されていること
	if _, err := repo.GetUserByEmail(context.Background(), "new@example.com"); err != nil {
		t.Errorf("user not persisted: %v", err)
	}
}

func TestSignUp_DuplicateEmail(t *testing.T) {
	repo := newFakeUserRepo()
	repo.seed(makeUser(t, "dup@example.com", "pw"))
	uc := &SignUpInteractor{UserRepository: repo}

	_, err := uc.Execute(context.Background(), SignUpInput{
		Name:     "second",
		Email:    "dup@example.com",
		Password: "password123",
	})
	if !errors.Is(err, usecase.ErrDuplicateEntity) {
		t.Errorf("error = %v, want ErrDuplicateEntity", err)
	}
}
