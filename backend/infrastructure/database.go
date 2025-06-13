package infrastructure

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3" // アンダースコア(_)でインポートするのは、このドライバの関数を直接は使わないが、
	// 内部的に "database/sql" に "sqlite3" ドライバを登録させるために必要だから
)

// ConnectDB はデータベース接続を初期化し、必要なテーブルを作成します
func ConnectDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./todo.db?_foreign_keys=on")
	// "sqlite3"ドライバを使って、カレントディレクトリの "todo.db" ファイルを開く（なければ新規作成）
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// データベース接続の確認
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// 各テーブルの作成SQL
	createTableSQLs := []string{
		`
		CREATE TABLE IF NOT EXISTS todos (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			content TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`,
		`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL UNIQUE, -- ユーザー名は重複不可
			password_hash TEXT NOT NULL
		);`,
		`
		CREATE TABLE IF NOT EXISTS user_sessions (
			user_id INTEGER PRIMARY KEY, -- ユーザーIDを主キーにすることで、1ユーザー1セッションを保証
			token_id TEXT NOT NULL,
			expires_at TIMESTAMP NOT NULL
		);`,
		`
		CREATE TABLE IF NOT EXISTS organizations (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			parent_id INTEGER, -- 親オーガナイゼーションのID (NULL許容)
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (parent_id) REFERENCES organizations(id) ON DELETE SET NULL
		);`,
		`
		CREATE TABLE IF NOT EXISTS view_group_companies (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,
		`
		CREATE TABLE IF NOT EXISTS shared_products (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			description TEXT,
			price REAL NOT NULL,
			sku TEXT UNIQUE,
			organization_id INTEGER NOT NULL, -- 親オーガナイゼーションを識別
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (organization_id) REFERENCES organizations(id) ON DELETE CASCADE
		);`,
		`
		CREATE TABLE IF NOT EXISTS projects (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			organization_id INTEGER NOT NULL, -- 所属オーガナイゼーションを識別
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			deleted_at DATETIME, -- 論理削除用
			FOREIGN KEY (organization_id) REFERENCES organizations(id) ON DELETE CASCADE
		);`,
		`
		CREATE TABLE IF NOT EXISTS project_products (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			project_id INTEGER NOT NULL,
			name TEXT NOT NULL,
			description TEXT,
			price REAL NOT NULL,
			shared_product_id INTEGER, -- どの共有商品からコピーされたか（NULL許容）
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
			FOREIGN KEY (shared_product_id) REFERENCES shared_products(id) ON DELETE SET NULL
		);`,
	}

	// 各SQL文を順番に実行
	for _, sql := range createTableSQLs {
		_, err = db.Exec(sql)
		if err != nil {
			log.Printf("Failed to execute DDL: %s", sql) // どのSQLで失敗したかログ出力
			return nil, fmt.Errorf("failed to create table: %w", err)
		}
	}

	log.Println("All necessary tables are ensured.")
	// 成功したら、データベース接続情報と、エラーなし(nil)を返す
	return db, nil
}