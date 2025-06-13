package usecase

import "go-todo-app/backend/domain"

// ProjectRepository は Project の永続化を担当するインターフェース
type ProjectRepository interface {
	Create(project *domain.Project) (*domain.Project, error)
	FindByID(id uint) (*domain.Project, error)
	FindByOrganizationID(orgID uint) ([]domain.Project, error)
	Update(project *domain.Project) error
	Delete(id uint) error // 論理削除
	// TODO: 物理削除のバッチ処理関連も考慮
}