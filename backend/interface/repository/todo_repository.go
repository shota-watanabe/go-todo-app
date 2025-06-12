// データベース操作の具体的な実装
package repository

import (
	"database/sql"
	"go-todo-app/backend/domain"
	"go-todo-app/backend/usecase"
)

type todoRepository struct {
	db *sql.DB
}

func NewTodoRepository(db *sql.DB) usecase.TodoRepository {
	return &todoRepository{db: db}
}

func (tr *todoRepository) FindAll() ([]domain.Todo, error) {
	rows, err := tr.db.Query("SELECT * FROM todos ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []domain.Todo
	// rows.Next()で結果を1行ずつループ処理
	for rows.Next() {
		// 1行分のデータを受け取るためのTodo構造体を準備
		var t domain.Todo
		// 現在の行のカラムを、tの各フィールドにスキャン（マッピング）
		// SELECTで指定した順序と、Scanの引数の順序を一致させる必要がある
		if err := rows.Scan(&t.ID, &t.Content, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		todos = append(todos, t)
	}
	return todos, nil
}

func (tr *todoRepository) Store(todo domain.Todo) (domain.Todo, error) {
	stmt, err := tr.db.Prepare("INSERT INTO todos(content) VALUES(?)")
	if err != nil {
		return domain.Todo{}, err
	}
	defer stmt.Close()
	// SQLの ? の部分に、実際の値(todo.Content)を埋め込んでSQLを実行
	res, err := stmt.Exec(todo.Content)
	if err != nil {
		return domain.Todo{}, err
	}
	// 最後に挿入された行のID（自動採番されたID）を取得
	id, err := res.LastInsertId()
	if err != nil {
		return domain.Todo{}, err
	}
	todo.ID = id
	// DBが自動で設定したcreated_atとupdated_atを取得して、todoオブジェクトにセットし直す
	err = tr.db.QueryRow("SELECT created_at, updated_at FROM todos WHERE id = ?", id).Scan(&todo.CreatedAt, &todo.UpdatedAt)
	if err != nil {
		return domain.Todo{}, err
	}
	
	return todo, nil
}

func (tr *todoRepository) Delete(id int64) error {
	stmt, err := tr.db.Prepare("DELETE FROM todos WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// SQLを実行します。
	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}

	// エラーがなければnilを返します。
	return nil
}


// FindByID, Update, Deleteも同様に実装します（今回は省略）
func (tr *todoRepository) FindByID(id int64) (domain.Todo, error) { /* ... */ return domain.Todo{}, nil }
func (tr *todoRepository) Update(todo domain.Todo) error { /* ... */ return nil }