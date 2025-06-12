package controller

import (
	"go-todo-app/backend/usecase"
	"net/http"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userUsecase usecase.UserUsecase
}

func NewUserHandler(uu usecase.UserUsecase) *UserHandler {
	return &UserHandler{userUsecase: uu}
}

// ユーザー関連のリクエストを処理するためのハンドラーの構造体
type userRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (uh *UserHandler) Register(c echo.Context) error {
	var req userRequest
	// c.Bind(&req)で、リクエストのJSONデータをreq変数にバインド
	if err := c.Bind(&req); err != nil {
		// JSONの形式が不正など、バインドに失敗した場合はエラー
		return c.JSON(http.StatusBadRequest, "Invalid request")
	}
	// バインドしたデータ（ユーザー名とパスワード）を使って、usecase層のRegisterメソッドを呼び出す
	if err := uh.userUsecase.Register(req.Username, req.Password); err != nil {
		return c.JSON(http.StatusConflict, err.Error())
	}
	return c.JSON(http.StatusCreated, "User created successfully")	
}

func (uh *UserHandler) Login(c echo.Context) error {
	var req userRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid request")
	}
	token, err := uh.userUsecase.Login(req.Username, req.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, err.Error())
	}
	return c.JSON(http.StatusOK, echo.Map{"token": token})
}