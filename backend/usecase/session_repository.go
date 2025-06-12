package usecase

import "time"

// SessionRepositoryは、同時ログイン制御のためのセッション情報を管理します。
type SessionRepository interface {
	Save(userID int64, tokenID string, expiresAt time.Time) error
	Validate(userID int64, tokenID string) (bool, error)
}