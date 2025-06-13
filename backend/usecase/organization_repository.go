package usecase

import "go-todo-app/backend/domain"

// OrganizationRepository は Organization の永続化を担当するインターフェース
type OrganizationRepository interface {
	FindByID(id uint) (*domain.Organization, error)
	// TODO: 親子関係やブラックリストに関するメソッドを追加
}