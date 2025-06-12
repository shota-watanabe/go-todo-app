// データベース操作のルールブック（インターフェース） を定義
package usecase

import "go-todo-app/backend/domain"

// Todoの永続化を抽象化するインターフェース
type TodoRepository interface {
	FindAll() ([]domain.Todo, error)
	FindByID(id int64) (domain.Todo, error)
	Store(todo domain.Todo) (domain.Todo, error)
	Update(todo domain.Todo) error
	Delete(id int64) error
}