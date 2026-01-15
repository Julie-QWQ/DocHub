package utils

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name    string
		password string
		wantErr bool
	}{
		{
			name:     "正常密码",
			password: "password123",
			wantErr:  false,
		},
		{
			name:     "短密码",
			password: "12345",
			wantErr:  true,
		},
		{
			name:     "空密码",
			password: "",
			wantErr:  true,
		},
		{
			name:     "长密码",
			password: "this_is_a_very_long_password_with_more_than_50_characters",
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hashed, err := HashPassword(tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("HashPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && hashed == "" {
				t.Error("HashPassword() 返回空字符串")
			}
			if !tt.wantErr && hashed == tt.password {
				t.Error("HashPassword() 密码未加密")
			}
		})
	}
}

func TestCheckPassword(t *testing.T) {
	password := "test123456"
	hashed, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword() 失败: %v", err)
	}

	tests := []struct {
		name           string
		password       string
		hashedPassword string
		want           bool
	}{
		{
			name:           "正确密码",
			password:       password,
			hashedPassword: hashed,
			want:           true,
		},
		{
			name:           "错误密码",
			password:       "wrongpassword",
			hashedPassword: hashed,
			want:           false,
		},
		{
			name:           "空密码",
			password:       "",
			hashedPassword: hashed,
			want:           false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckPassword(tt.password, tt.hashedPassword); got != tt.want {
				t.Errorf("CheckPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}
