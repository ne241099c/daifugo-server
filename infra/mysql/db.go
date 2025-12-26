package mysql

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/ne241099/daifugo-server/internal/config"
)

// NewDB はデータベース接続を確立して返します
func NewDB(cfg *config.Config) (*sql.DB, error) {
	// DSN (Data Source Name) を取得
	dsn := cfg.DSN()

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %w", err)
	}

	// 接続確認
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping db: %w", err)
	}

	// 接続プールの設定
	db.SetMaxOpenConns(25)                 // 最大同時接続数
	db.SetMaxIdleConns(25)                 // アイドル状態の最大接続数
	db.SetConnMaxLifetime(5 * time.Minute) // 接続の最大寿命

	return db, nil
}
