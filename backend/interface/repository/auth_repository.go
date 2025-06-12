package repository

import (
	"database/sql"
	"go-todo-app/backend/domain"
	"go-todo-app/backend/usecase"
	"time"
)

// --- UserRepositoryの実装 ---
type userRepository struct{ db *sql.DB }

func NewUserRepository(db *sql.DB) usecase.UserRepository {
	return &userRepository{db: db}
}

func (ur *userRepository) Store(user domain.User) (int64, error) {
	res, err := ur.db.Exec("INSERT INTO users (username, password_hash) VALUES (?, ?)", user.Username, user.PasswordHash)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (ur *userRepository) FindByUsername(username string) (*domain.User, error) {
	var user domain.User
	err := ur.db.QueryRow("SELECT id, username, password_hash FROM users WHERE username = ?", username).Scan(&user.ID, &user.Username, &user.PasswordHash)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
func (ur *userRepository) FindByID(id int64) (*domain.User, error) { /* ... */ return nil, nil }

// --- SessionRepositoryの実装 ---
type sessionRepository struct{ db *sql.DB }

func NewSessionRepository(db *sql.DB) usecase.SessionRepository {
	return &sessionRepository{db: db}
}

func (sr *sessionRepository) Save(userID int64, tokenID string, expiresAt time.Time) error {
	// INSERT OR REPLACE: 既に同じuser_idのレコードがあれば上書き、なければ新規作成する。
	// これが同時ログイン禁止の核となる処理。
	_, err := sr.db.Exec(
		"INSERT OR REPLACE INTO user_sessions (user_id, token_id, expires_at) VALUES (?, ?, ?)",
		userID, tokenID, expiresAt,
	)
	return err
}

func (sr *sessionRepository) Validate(userID int64, tokenID string) (bool, error) {
	var storedTokenID string
	err := sr.db.QueryRow("SELECT token_id FROM user_sessions WHERE user_id = ?", userID).Scan(&storedTokenID)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil // セッションが存在しない
		}
		return false, err
	}
	// DBに保存されているトークンIDと、クライアントから送られてきたトークンIDが一致するか検証
	return storedTokenID == tokenID, nil
}