package utils

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
)

const (
	// DefaultCost bcrypt 加密成本
	DefaultCost = bcrypt.DefaultCost
)

var (
	// ErrPasswordTooShort 密码过短错误
	ErrPasswordTooShort = errors.New("密码长度至少为6位")
)

// HashPassword 对密码进行哈希加密
func HashPassword(password string) (string, error) {
	if len(password) < 6 {
		return "", ErrPasswordTooShort
	}

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedBytes), nil
}

// CheckPassword 验证密码是否正确
func CheckPassword(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
