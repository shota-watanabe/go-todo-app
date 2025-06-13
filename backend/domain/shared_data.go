package domain

import "time"

// ViewGroupCompany は参照のみのグループ会社エンティティ
type ViewGroupCompany struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// SharedProduct はコピー可能な共有商品エンティティ
type SharedProduct struct {
	ID             uint      `json:"id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	Price          float64   `json:"price"`
	SKU            string    `json:"sku"` // Stock Keeping Unit
	OrganizationID uint      `json:"organization_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// ProjectProduct はプロジェクト固有の商品エンティティ
type ProjectProduct struct {
	ID                   uint      `json:"id"`
	ProjectID            uint      `json:"project_id"`
	Name                 string    `json:"name"`
	Description          string    `json:"description"`
	Price                float64   `json:"price"`
	SharedProductID *uint `json:"original_shared_product_id"` // どの共有商品からコピーされたか
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}