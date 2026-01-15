package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	// ErrTokenInvalid Token 无效错误
	ErrTokenInvalid = errors.New("Token 无效")
	// ErrTokenExpired Token 已过期错误
	ErrTokenExpired = errors.New("Token 已过期")
	// ErrTokenTypeWrong Token 类型错误
	ErrTokenTypeWrong = errors.New("Token 类型错误")
)

// JWTManager JWT 管理器
type JWTManager struct {
	secret           string
	accessTTL        time.Duration
	refreshTTL       time.Duration
	issuer           string
}

// TokenClaims JWT Token 声明
type TokenClaims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	Type   string `json:"type"` // access 或 refresh
	jwt.RegisteredClaims
}

// NewJWTManager 创建 JWT 管理器
func NewJWTManager(secret string, accessTTL, refreshTTL time.Duration, issuer string) *JWTManager {
	return &JWTManager{
		secret:     secret,
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
		issuer:     issuer,
	}
}

// GenerateAccessToken 生成访问 Token
func (j *JWTManager) GenerateAccessToken(userID uint, role string) (string, error) {
	claims := TokenClaims{
		UserID: userID,
		Role:   role,
		Type:   "access",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.issuer,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.accessTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	return j.generateToken(claims)
}

// GenerateRefreshToken 生成刷新 Token
func (j *JWTManager) GenerateRefreshToken(userID uint, role string) (string, error) {
	claims := TokenClaims{
		UserID: userID,
		Role:   role,
		Type:   "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.issuer,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.refreshTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	return j.generateToken(claims)
}

// GenerateTokenPair 生成 Token 对
func (j *JWTManager) GenerateTokenPair(userID uint, role string) (accessToken, refreshToken string, err error) {
	accessToken, err = j.GenerateAccessToken(userID, role)
	if err != nil {
		return "", "", err
	}

	refreshToken, err = j.GenerateRefreshToken(userID, role)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// generateToken 生成 Token
func (j *JWTManager) generateToken(claims TokenClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secret))
}

// ParseToken 解析并验证 Token
func (j *JWTManager) ParseToken(tokenString string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrTokenInvalid
		}
		return []byte(j.secret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, ErrTokenInvalid
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok || !token.Valid {
		return nil, ErrTokenInvalid
	}

	return claims, nil
}

// ValidateAccessToken 验证访问 Token
func (j *JWTManager) ValidateAccessToken(tokenString string) (*TokenClaims, error) {
	claims, err := j.ParseToken(tokenString)
	if err != nil {
		return nil, err
	}

	if claims.Type != "access" {
		return nil, ErrTokenTypeWrong
	}

	return claims, nil
}

// ValidateRefreshToken 验证刷新 Token
func (j *JWTManager) ValidateRefreshToken(tokenString string) (*TokenClaims, error) {
	claims, err := j.ParseToken(tokenString)
	if err != nil {
		return nil, err
	}

	if claims.Type != "refresh" {
		return nil, ErrTokenTypeWrong
	}

	return claims, nil
}

// GetAccessTTL 获取访问 Token 过期时间（秒）
func (j *JWTManager) GetAccessTTL() int64 {
	return int64(j.accessTTL.Seconds())
}
