package main

import (
	"go-todo-app/backend/infrastructure"
	"go-todo-app/backend/interface/controller"
	"go-todo-app/backend/interface/repository"
	"go-todo-app/backend/usecase"
	"log"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	db, err := infrastructure.ConnectDB()
	if err != nil {
		log.Fatalf("DB接続失敗: %v", err)
	}
	defer db.Close()

	// --- 依存性の注入 (DI) ---
	todoRepo := repository.NewTodoRepository(db)
	userRepo := repository.NewUserRepository(db)
	sessionRepo := repository.NewSessionRepository(db)
	todoUsecase := usecase.NewTodoUsecase(todoRepo)
	userUsecase := usecase.NewUserUsecase(userRepo, sessionRepo)
	todoHandler := controller.NewTodoHandler(todoUsecase)
	userHandler := controller.NewUserHandler(userUsecase)
	
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.OPTIONS},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	// --- ルーティング ---
	// 認証が不要なルート (公開API)
	e.POST("/register", userHandler.Register)
	e.POST("/login", userHandler.Login)

	// 認証が必要なルートグループ
	api := e.Group("/api")
	
	api.Use(echojwt.WithConfig(echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(usecase.JwtCustomClaims)
		},
		SigningKey: []byte("my_super_secret_key"),
	}))
	
	// 同時ログインを禁止するカスタムミドルウェア
	api.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// echo-jwtミドルウェアを使うと、c.Get("user")の中身が *jwt.Token になります。
			user, ok := c.Get("user").(*jwt.Token)
			if !ok {
				// 型アサーションに失敗した場合
				log.Println("ERROR: c.Get(\"user\") is not a *jwt.Token")
				return echo.ErrInternalServerError
			}
			
			claims, ok := user.Claims.(*usecase.JwtCustomClaims)
			if !ok {
				// 型アサーションに失敗した場合
				log.Println("ERROR: user.Claims is not a *usecase.JwtCustomClaims")
				return echo.ErrInternalServerError
			}

			isValid, err := sessionRepo.Validate(claims.UserID, claims.ID)
			if err != nil || !isValid {
				return echo.ErrUnauthorized // セッションが無効なら401エラー
			}
            c.Set("userID", claims.UserID)
			return next(c)
		}
	})

	// 認証済みユーザーのみがアクセスできるAPI
	api.GET("/todos", todoHandler.GetTodos)
	api.POST("/todos", todoHandler.CreateTodo)
	api.DELETE("/todos/:id", todoHandler.DeleteTodo)

	e.Logger.Fatal(e.Start(":8080"))
}