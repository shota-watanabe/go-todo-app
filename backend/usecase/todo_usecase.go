// アプリケーション固有のビジネスロジックを実装
package usecase

import "go-todo-app/backend/domain"

type TodoUsecase interface {
	GetAllTodos() ([]domain.Todo, error)
	CreateTodo(content string) (domain.Todo, error)
	DeleteTodo(id int64) error
}

type todoUsecase struct {
	todoRepo TodoRepository
}

// NewTodoUsecase は、新しいTodoUsecaseを生成
func NewTodoUsecase(tr TodoRepository) TodoUsecase {
	return &todoUsecase{todoRepo: tr}
}

// 「全てのTodoを取得する」というビジネスロジックを実装
func (tu *todoUsecase) GetAllTodos() ([]domain.Todo, error) {
	return tu.todoRepo.FindAll()
}

// 「新しいTodoを作成する」というビジネスロジックを実装
func (tu *todoUsecase) CreateTodo(content string) (domain.Todo, error) {
	todo := domain.Todo{Content: content}
	return tu.todoRepo.Store(todo)
}

func (tu *todoUsecase) DeleteTodo(id int64) error {
	return tu.todoRepo.Delete(id)
}