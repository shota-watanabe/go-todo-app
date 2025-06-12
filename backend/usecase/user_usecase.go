package usecase

import (
	"errors"
	"go-todo-app/backend/domain"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// JWTの秘密鍵（実際には環境変数などから読み込むべき）
var jwtSecret = []byte("my_super_secret_key")

// JWTのカスタムクレーム（トークンに含める情報）
type JwtCustomClaims struct {
	Username string `json:"username"`
	UserID   int64  `json:"user_id"`
	jwt.RegisteredClaims
}

// UserUsecase は、ユーザーに関するビジネスロジックを定義します。
type UserUsecase interface {
	Register(username, password string) error
	Login(username, password string) (string, error)
}

type userUsecase struct {
	userRepo    UserRepository
	sessionRepo SessionRepository
}

func NewUserUsecase(ur UserRepository, sr SessionRepository) UserUsecase {
	return &userUsecase{userRepo: ur, sessionRepo: sr}
}

// Register はユーザー登録のロジックです。
func (uu *userUsecase) Register(username, password string) error {
	// ユーザーが既に存在するかチェック
	if _, err := uu.userRepo.FindByUsername(username); err == nil {
		return errors.New("username already exists")
	}

	// パスワードをハッシュ化（平文で保存するのは絶対にNG）
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := domain.User{
		Username:     username,
		PasswordHash: string(hashedPassword),
	}
	_, err = uu.userRepo.Store(user)
	return err
}

func (uu *userUsecase) Login(username, password string) (string, error) {
	// ユーザーを検索
	user, err := uu.userRepo.FindByUsername(username)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	// パスワードハッシュを比較
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	// --- JWTとセッションの生成 ---
	tokenID := uuid.NewString() // トークンに一意のIDを付与
	expiresAt := time.Now().Add(time.Hour * 24) // 有効期限は24時間

	// JWTクレームを作成
	claims := &JwtCustomClaims{
		Username: user.Username,
		UserID:   user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        tokenID,
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	// 新しいセッションを保存（古いセッションは上書きされ、無効になる）
	if err := uu.sessionRepo.Save(user.ID, tokenID, expiresAt); err != nil {
		return "", err
	}

	return tokenString, nil
}