package domain

import "time"

// Project は利用スコープを表すエンティティ
type Project struct {
	ID            uint       `json:"id"`
	Name          string     `json:"name"`
	OrganizationID uint      `json:"organization_id"` // 所属オーガナイゼーションID
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at"`      // 論理削除用
	// TODO: プロジェクトのブラックリストなど、その他の属性を追加
}