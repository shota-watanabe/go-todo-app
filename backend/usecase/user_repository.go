package usecase

import "go-todo-app/backend/domain"

// UserRepositoryは、Userの永続化を抽象化するインターフェースです。
type UserRepository interface {
	Store(user domain.User) (int64, error)
	FindByUsername(username string) (*domain.User, error)
	FindByID(id int64) (*domain.User, error)
}