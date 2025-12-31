package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ne241099/daifugo-server/model"
	"github.com/ne241099/daifugo-server/repository"
)

type MySQLUserRepository struct {
	db *sql.DB
}

func NewMySQLUserRepository(db *sql.DB) *MySQLUserRepository {
	return &MySQLUserRepository{db: db}
}

// SaveUser はユーザーを新規作成または更新する
func (r *MySQLUserRepository) SaveUser(ctx context.Context, u *model.User) error {
	if u.ID == 0 {
		return r.create(ctx, u)
	}
	return r.update(ctx, u)
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

	id, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}
	u.ID = id
	return nil
}

func (r *MySQLUserRepository) update(ctx context.Context, u *model.User) error {
	query := `
		UPDATE users 
		SET name = ?, email = ?, password_hash = ?, updated_at = ?
		WHERE id = ?
	`
	_, err := r.db.ExecContext(ctx, query, u.Name, u.Email, u.HashedPassword, u.UpdatedAt, u.ID)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

// GetUser はIDでユーザーを取得する
func (r *MySQLUserRepository) GetUser(ctx context.Context, id int64) (*model.User, error) {
	fmt.Printf("デバッグ: MySQL GetUser id=%d 開始\n", id)

	query := `
		SELECT id, name, email, password_hash, created_at, updated_at
		FROM users WHERE id = ?
	`
	row := r.db.QueryRowContext(ctx, query, id)

	var u model.User
	if err := row.Scan(&u.ID, &u.Name, &u.Email, &u.HashedPassword, &u.CreatedAt, &u.UpdatedAt); err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrEntityNotFound
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
			return nil, repository.ErrEntityNotFound
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

func (r *MySQLUserRepository) ListUsers(ctx context.Context) ([]*model.User, error) {
	{
		query := `
		SELECT id, name, email, password_hash, created_at, updated_at
		FROM users
	`
		rows, err := r.db.QueryContext(ctx, query)
		if err != nil {
			return nil, fmt.Errorf("failed to query users: %w", err)
		}
		defer rows.Close()

		var users []*model.User
		for rows.Next() {
			var u model.User
			if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.HashedPassword, &u.CreatedAt, &u.UpdatedAt); err != nil {
				return nil, fmt.Errorf("failed to scan user: %w", err)
			}
			users = append(users, &u)
		}

		if err := rows.Err(); err != nil {
			return nil, fmt.Errorf("rows iteration error: %w", err)
		}

		return users, nil
	}
}

func (r *MySQLUserRepository) IncrementTokenVersion(ctx context.Context, userID int64) (int, error) {
	query := `UPDATE users SET token_version = token_version + 1 WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return 0, err
	}

	// 更新後の値を取得
	var newVersion int
	err = r.db.QueryRowContext(ctx, "SELECT token_version FROM users WHERE id = ?", userID).Scan(&newVersion)
	if err != nil {
		return 0, err
	}
	return newVersion, nil
}
