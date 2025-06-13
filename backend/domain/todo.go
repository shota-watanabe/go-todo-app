package domain

import "time"

// 「Todo」というデータ構造（エンティティ）を定義
type Todo struct {
	ID        int64     `json:"id"`
	UserID    uint      `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}