// DB接続などのインフラ層のコード
package infrastructure

import (
	"database/sql"
	// アンダースコア(_)でインポートするのは、このドライバの関数を直接は使わないが、
	// 内部的に "database/sql" に "sqlite3" ドライバを登録させるために必要だから
	_ "github.com/mattn/go-sqlite3"
)

// データベースに接続する関数
func ConnectDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./todo.db?_foreign_keys=on")
	// "sqlite3"ドライバを使って、カレントディレクトリの "todo.db" ファイルを開く（なければ新規作成）
	if err != nil {
		return nil, err
	}

	createTableSQL := `
    CREATE TABLE IF NOT EXISTS todos (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
        content TEXT NOT NULL,
        created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    );`
	// 定義したSQL文をデータベースで実行
	_, err = db.Exec(createTableSQL)
	if err != nil {
		// SQLの実行に失敗したら、エラーを返す
		return nil, err
	}

	// Userテーブル
	userTableSQL := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE, -- ユーザー名は重複不可
		password_hash TEXT NOT NULL
	);`
	if _, err = db.Exec(userTableSQL); err != nil {
		return nil, err
	}

	// UserSessionテーブル (同時ログイン禁止のキー)
	sessionTableSQL := `
	CREATE TABLE IF NOT EXISTS user_sessions (
		user_id INTEGER PRIMARY KEY, -- ユーザーIDを主キーにすることで、1ユーザー1セッションを保証
		token_id TEXT NOT NULL,
		expires_at TIMESTAMP NOT NULL
	);`
	if _, err = db.Exec(sessionTableSQL); err != nil {
		return nil, err
	}

	// 成功したら、データベース接続情報と、エラーなし(nil)を返す
	return db, nil
}