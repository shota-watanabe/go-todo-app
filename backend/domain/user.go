package domain

// Userエンティティを定義
type User struct {
	ID           int64
	Username     string
	PasswordHash string // パスワードのハッシュ値。クライアントには絶対に返さない。
}