package domain

import "time"

// Organization は法人を表すエンティティ
type Organization struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	ParentID  *uint     `json:"parent_id"` // 親オーガナイゼーションのID (NULL許容)
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// TODO: ブラックリストなど、その他の属性を追加
}