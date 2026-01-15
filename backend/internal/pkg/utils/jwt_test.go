package utils

import (
	"testing"
	"time"
)

func TestJWTManager_GenerateTokenPair(t *testing.T) {
	secret := "test-secret-key"
	jwtManager := NewJWTManager(secret, time.Hour, 24*time.Hour, "test")

	userID := uint(1)
	role := "student"

	accessToken, refreshToken, err := jwtManager.GenerateTokenPair(userID, role)
	if err != nil {
		t.Fatalf("GenerateTokenPair() 失败: %v", err)
	}

	if accessToken == "" {
		t.Error("AccessToken 为空")
	}
	if refreshToken == "" {
		t.Error("RefreshToken 为空")
	}
}

func TestJWTManager_ValidateAccessToken(t *testing.T) {
	secret := "test-secret-key"
	jwtManager := NewJWTManager(secret, time.Hour, 24*time.Hour, "test")

	userID := uint(123)
	role := "admin"

	// 生成访问 Token
	accessToken, err := jwtManager.GenerateAccessToken(userID, role)
	if err != nil {
		t.Fatalf("GenerateAccessToken() 失败: %v", err)
	}

	// 验证访问 Token
	claims, err := jwtManager.ValidateAccessToken(accessToken)
	if err != nil {
		t.Fatalf("ValidateAccessToken() 失败: %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("UserID = %v, want %v", claims.UserID, userID)
	}
	if claims.Role != role {
		t.Errorf("Role = %v, want %v", claims.Role, role)
	}
	if claims.Type != "access" {
		t.Errorf("Type = %v, want access", claims.Type)
	}
}

func TestJWTManager_ValidateRefreshToken(t *testing.T) {
	secret := "test-secret-key"
	jwtManager := NewJWTManager(secret, time.Hour, 24*time.Hour, "test")

	userID := uint(456)
	role := "committee"

	// 生成刷新 Token
	refreshToken, err := jwtManager.GenerateRefreshToken(userID, role)
	if err != nil {
		t.Fatalf("GenerateRefreshToken() 失败: %v", err)
	}

	// 验证刷新 Token
	claims, err := jwtManager.ValidateRefreshToken(refreshToken)
	if err != nil {
		t.Fatalf("ValidateRefreshToken() 失败: %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("UserID = %v, want %v", claims.UserID, userID)
	}
	if claims.Type != "refresh" {
		t.Errorf("Type = %v, want refresh", claims.Type)
	}
}

func TestJWTManager_InvalidToken(t *testing.T) {
	secret := "test-secret-key"
	jwtManager := NewJWTManager(secret, time.Hour, 24*time.Hour, "test")

	tests := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{
			name:    "无效 Token",
			token:   "invalid.token.string",
			wantErr: true,
		},
		{
			name:    "空 Token",
			token:   "",
			wantErr: true,
		},
		{
			name:    "错误的签名",
			token:   "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJyb2xlIjoic3R1ZGVudCJ9.invalid",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := jwtManager.ValidateAccessToken(tt.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateAccessToken() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestJWTManager_WrongTokenType(t *testing.T) {
	secret := "test-secret-key"
	jwtManager := NewJWTManager(secret, time.Hour, 24*time.Hour, "test")

	userID := uint(789)
	role := "student"

	// 生成访问 Token
	accessToken, err := jwtManager.GenerateAccessToken(userID, role)
	if err != nil {
		t.Fatalf("GenerateAccessToken() 失败: %v", err)
	}

	// 用验证刷新 Token 的方法验证访问 Token
	_, err = jwtManager.ValidateRefreshToken(accessToken)
	if err != ErrTokenTypeWrong {
		t.Errorf("期望返回 ErrTokenTypeWrong, 实际: %v", err)
	}
}
