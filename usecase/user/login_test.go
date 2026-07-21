package user

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/ne241099/daifugo-server/model"
)

// makeUser はパスワードをハッシュ化した既存ユーザーを作る
func makeUser(t *testing.T, email, password string) *model.User {
	t.Helper()
	u := &model.User{}
	if err := u.Create(model.CreateUserParam{
		Email:     email,
		Password:  password,
		Name:      "tester",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}); err != nil {
		t.Fatalf("failed to create user: %v", err)
	}
	return u
}

func TestLogin_Success(t *testing.T) {
	repo := newFakeUserRepo()
	repo.seed(makeUser(t, "a@example.com", "correct-password"))
	tokenGen := &fakeTokenGenerator{}

	uc := &LoginInteractor{UserRepository: repo, TokenGenerator: tokenGen}

	token, u, err := uc.Execute(context.Background(), "a@example.com", "correct-password")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if token != "signed-token" {
		t.Errorf("token = %q, want %q", token, "signed-token")
	}
	// トークンバージョンがインクリメントされ、発行に反映されていること
	if u.TokenVersion != 1 {
		t.Errorf("TokenVersion = %d, want 1", u.TokenVersion)
	}
	if tokenGen.lastVersion != 1 || tokenGen.lastUserID != u.ID {
		t.Errorf("token generated with userID=%d ver=%d, want userID=%d ver=1", tokenGen.lastUserID, tokenGen.lastVersion, u.ID)
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	repo := newFakeUserRepo()
	repo.seed(makeUser(t, "a@example.com", "correct-password"))
	uc := &LoginInteractor{UserRepository: repo, TokenGenerator: &fakeTokenGenerator{}}

	_, _, err := uc.Execute(context.Background(), "a@example.com", "wrong-password")
	if err == nil {
		t.Fatal("expected error for wrong password")
	}
	if !strings.Contains(err.Error(), "invalid email or password") {
		t.Errorf("error = %q, want to contain 'invalid email or password'", err.Error())
	}
}

func TestLogin_UnknownEmail(t *testing.T) {
	repo := newFakeUserRepo()
	uc := &LoginInteractor{UserRepository: repo, TokenGenerator: &fakeTokenGenerator{}}

	_, _, err := uc.Execute(context.Background(), "nobody@example.com", "whatever")
	if err == nil {
		t.Fatal("expected error for unknown email")
	}
	// セキュリティ上、存在しないメールでも同じメッセージであること
	if !strings.Contains(err.Error(), "invalid email or password") {
		t.Errorf("error = %q, want to contain 'invalid email or password'", err.Error())
	}
}
