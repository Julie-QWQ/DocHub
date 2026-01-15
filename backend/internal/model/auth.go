package model

import "github.com/golang-jwt/jwt/v5"

// RegisterRequest 用户注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50,alphanum"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=50"`
	RealName string `json:"real_name" binding:"required,min=2,max=50"`
	Major    string `json:"major" binding:"required"`
	Class    string `json:"class" binding:"required"`
}

// LoginRequest 用户登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required,min=6,max=50"`
	NewPassword string `json:"new_password" binding:"required,min=6,max=50"`
}

// RefreshTokenRequest 刷新Token请求
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	User         UserInfo `json:"user"`
}

// UserInfo 用户信息（不包含敏感信息）
type UserInfo struct {
	ID        uint       `json:"id"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	RealName  string     `json:"real_name"`
	Role      UserRole   `json:"role"`
	Status    UserStatus `json:"status"`
	Avatar    string     `json:"avatar"`
	Phone     string     `json:"phone"`
	Major     string     `json:"major"`
	Class     string     `json:"class"`
	CreatedAt string     `json:"created_at"`
}

// TokenClaims JWT Token 声明
type TokenClaims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	Type   string `json:"type"` // access 或 refresh
	jwt.RegisteredClaims
}

// ToUserInfo 将 User 转换为 UserInfo
func (u *User) ToUserInfo() UserInfo {
	return UserInfo{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		RealName:  u.RealName,
		Role:      u.Role,
		Status:    u.Status,
		Avatar:    u.Avatar,
		Phone:     u.Phone,
		Major:     u.Major,
		Class:     u.Class,
		CreatedAt: u.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}
