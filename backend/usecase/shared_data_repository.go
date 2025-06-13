package usecase

import "go-todo-app/backend/domain"

// ViewGroupCompanyRepository は ViewGroupCompany の読み取りを担当するインターフェース
type ViewGroupCompanyRepository interface {
	GetAll() ([]domain.ViewGroupCompany, error)
}

// SharedProductRepository は SharedProduct の永続化を担当するインターフェース
type SharedProductRepository interface {
	FindByID(id uint) (*domain.SharedProduct, error)
	FindByOrganizationID(orgID uint) ([]domain.SharedProduct, error)
	Create(product *domain.SharedProduct) (*domain.SharedProduct, error)
	// TODO: 親オーガナイゼーションによるCRUD操作メソッド
}

// ProjectProductRepository は ProjectProduct の永続化を担当するインターフェース
type ProjectProductRepository interface {
	Create(product *domain.ProjectProduct) (*domain.ProjectProduct, error)
	FindByProjectID(projectID uint) ([]domain.ProjectProduct, error)
	FindByID(id uint) (*domain.ProjectProduct, error)
	Update(product *domain.ProjectProduct) error
	Delete(id uint) error // 物理削除
}