package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ne241099/daifugo-server/model"
)

type MySQLUserRepository struct {
	db *sql.DB
}

func NewMySQLUserRepository(db *sql.DB) *MySQLUserRepository {
	return &MySQLUserRepository{db: db}
}

// SaveUser はユーザーを新規作成または更新する
func (r *MySQLUserRepository) SaveUser(ctx context.Context, u *model.User) error {
	// IDが0の場合は新規作成、それ以外は更新とみなす簡易実装
	if u.ID == 0 {
		return r.create(ctx, u)
	}
	// 更新処理は今回省略（必要なら UPDATE 文を書く）
	return nil
}

func (r *MySQLUserRepository) create(ctx context.Context, u *model.User) error {
	query := `
		INSERT INTO users (name, email, password_hash, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`
	res, err := r.db.ExecContext(ctx, query, u.Name, u.Email, u.HashedPassword, u.CreatedAt, u.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}

	// 発行されたIDを取得してモデルにセット
	id, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}
	u.ID = id
	return nil
}

// GetUser はIDでユーザーを取得する
func (r *MySQLUserRepository) GetUser(ctx context.Context, id int64) (*model.User, error) {
	query := `
		SELECT id, name, email, password_hash, created_at, updated_at
		FROM users WHERE id = ?
	`
	row := r.db.QueryRowContext(ctx, query, id)

	var u model.User
	if err := row.Scan(&u.ID, &u.Name, &u.Email, &u.HashedPassword, &u.CreatedAt, &u.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("failed to get user: %w", err) // 共通のエラー型があればそれを使う
		}
		return nil, fmt.Errorf("failed to scan user: %w", err)
	}
	return &u, nil
}

// GetUserByEmail はEmailでユーザーを取得する
func (r *MySQLUserRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `
		SELECT id, name, email, password_hash, created_at, updated_at
		FROM users WHERE email = ?
	`
	row := r.db.QueryRowContext(ctx, query, email)

	var u model.User
	if err := row.Scan(&u.ID, &u.Name, &u.Email, &u.HashedPassword, &u.CreatedAt, &u.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("failed to get user: %w", err)
		}
		return nil, fmt.Errorf("failed to scan user: %w", err)
	}
	return &u, nil
}

// DeleteUser はIDでユーザーを削除する
func (r *MySQLUserRepository) DeleteUser(ctx context.Context, id int64) error {
	query := "DELETE FROM users WHERE id = ?"
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}
