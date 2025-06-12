// HTTPリクエストを受け取り、適切なUsecaseを呼び出し、結果をレスポンスとして返す受付係
package controller

import (
	"net/http"
	"strconv"
	"go-todo-app/backend/usecase"
	"github.com/labstack/echo/v4"
)

type TodoHandler struct {
	// 内部にビジネスロジック担当のusecaseを保持
	todoUsecase usecase.TodoUsecase
}

// 引数でusecaseを受け取り、それを内部に保持したTodoHandlerを返す
func NewTodoHandler(tu usecase.TodoUsecase) *TodoHandler {
	return &TodoHandler{todoUsecase: tu}
}

func (th *TodoHandler) GetTodos(c echo.Context) error {
	todos, err := th.todoUsecase.GetAllTodos()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, todos)
}

func (th *TodoHandler) CreateTodo(c echo.Context) error {
	type Request struct {
		Content string `json:"content"`
	}
	var req Request
	// リクエストボディのJSONを、req構造体にマッピング（バインド）
	if err := c.Bind(&req); err != nil {
		// マッピングに失敗したら(JSON形式がおかしい等)、400 Bad Request
		return c.JSON(http.StatusBadRequest, "Invalid request")
	}
	// Contentが空文字でないかチェック
	if req.Content == "" {
		return c.JSON(http.StatusBadRequest, "Content is required")
	}

	todo, err := th.todoUsecase.CreateTodo(req.Content)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	// 成功したら、作成されたTodoオブジェクトをステータスコード201 Createdと共にJSON形式で返す
	return c.JSON(http.StatusCreated, todo)
}

func (th *TodoHandler) DeleteTodo(c echo.Context) error {
	// URLのパスパラメータから "id" を取得します (例: /api/todos/123 -> "123")。
	idStr := c.Param("id")
	// 取得したID（文字列）を64ビット整数(int64)に変換します。
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		// 変換に失敗した場合（IDが数値でないなど）、400 Bad Requestエラーを返します。
		return c.JSON(http.StatusBadRequest, "Invalid ID format")
	}

	// Usecaseに「このIDのTodoを削除して」と依頼します。
	err = th.todoUsecase.DeleteTodo(id)
	if err != nil {
		// Usecaseでエラーが発生したら、500 Internal Server Errorを返します。
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	// 成功した場合、ボディなしの 204 No Content ステータスを返すのが一般的です。
	return c.NoContent(http.StatusNoContent)
}