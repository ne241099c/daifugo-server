package user

import (
	"context"

	"github.com/ne241099/daifugo-server/model"
	"github.com/ne241099/daifugo-server/repository"
)

// fakeUserRepo はテスト用のインメモリ UserRepository
type fakeUserRepo struct {
	byID    map[int64]*model.User
	byEmail map[string]*model.User
	next    int64

	// エラー注入用
	getByEmailErr error
	saveErr       error
}

func newFakeUserRepo() *fakeUserRepo {
	return &fakeUserRepo{
		byID:    make(map[int64]*model.User),
		byEmail: make(map[string]*model.User),
		next:    1,
	}
}

// seed は既存ユーザーとしてリポジトリへ直接登録する
func (f *fakeUserRepo) seed(u *model.User) {
	if u.ID == 0 {
		u.ID = f.next
		f.next++
	}
	f.byID[u.ID] = u
	f.byEmail[u.Email] = u
}

func (f *fakeUserRepo) GetUser(ctx context.Context, id int64) (*model.User, error) {
	u, ok := f.byID[id]
	if !ok {
		return nil, repository.ErrEntityNotFound
	}
	return u, nil
}

func (f *fakeUserRepo) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	if f.getByEmailErr != nil {
		return nil, f.getByEmailErr
	}
	u, ok := f.byEmail[email]
	if !ok {
		return nil, repository.ErrEntityNotFound
	}
	return u, nil
}

func (f *fakeUserRepo) SaveUser(ctx context.Context, u *model.User) error {
	if f.saveErr != nil {
		return f.saveErr
	}
	f.seed(u)
	return nil
}

func (f *fakeUserRepo) DeleteUser(ctx context.Context, id int64) error {
	u, ok := f.byID[id]
	if !ok {
		return repository.ErrEntityNotFound
	}
	delete(f.byEmail, u.Email)
	delete(f.byID, id)
	return nil
}

func (f *fakeUserRepo) ListUsers(ctx context.Context) ([]*model.User, error) {
	out := make([]*model.User, 0, len(f.byID))
	for _, u := range f.byID {
		out = append(out, u)
	}
	return out, nil
}

func (f *fakeUserRepo) IncrementTokenVersion(ctx context.Context, userID int64) (int, error) {
	u, ok := f.byID[userID]
	if !ok {
		return 0, repository.ErrEntityNotFound
	}
	u.TokenVersion++
	return u.TokenVersion, nil
}

// fakeTokenGenerator はテスト用の TokenGenerator
type fakeTokenGenerator struct {
	lastUserID  int64
	lastVersion int
	err         error
}

func (f *fakeTokenGenerator) CreateToken(ctx context.Context, userID int64, version int) (string, error) {
	if f.err != nil {
		return "", f.err
	}
	f.lastUserID = userID
	f.lastVersion = version
	return "signed-token", nil
}
