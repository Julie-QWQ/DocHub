package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/study-upc/backend/internal/middleware"
	"github.com/study-upc/backend/internal/model"
	"github.com/study-upc/backend/internal/pkg/utils"
	"github.com/study-upc/backend/internal/repository"
	"github.com/study-upc/backend/internal/service"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupAuthTest(t *testing.T) (*gin.Engine, *gorm.DB) {
	// 设置 Gin 为测试模式
	gin.SetMode(gin.TestMode)

	// 创建内存数据库
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	require.NoError(t, err)

	// 自动迁移
	err = db.AutoMigrate(&model.User{})
	require.NoError(t, err)

	// 初始化依赖
	userRepo := repository.NewUserRepository(db)
	jwtManager := utils.NewJWTManager("test-secret", time.Hour, 24*time.Hour, "test")
	authService := service.NewAuthService(userRepo, jwtManager, nil)
	authHandler := NewAuthHandler(authService)

	// 设置路由
	router := gin.New()
	authGroup := router.Group("/api/v1/auth")
	{
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/login", authHandler.Login)
		authGroup.POST("/refresh", authHandler.RefreshToken)

		// 需要认证的路由
		authenticated := authGroup.Use(middleware.JWTAuth(jwtManager, nil))
		authenticated.POST("/logout", authHandler.Logout)
		authenticated.GET("/me", authHandler.GetUserInfo)
	}

	return router, db
}

func TestAuthHandler_Register(t *testing.T) {
	router, _ := setupAuthTest(t)

	body := map[string]string{
		"username":  "testuser",
		"email":     "test@example.com",
		"password":  "password123",
		"real_name": "Test User",
		"major":     "计算机科学",
		"class":     "2101",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, float64(0), resp["code"])

	data := resp["data"].(map[string]interface{})
	assert.Equal(t, "testuser", data["username"])
	assert.Equal(t, "test@example.com", data["email"])
}

func TestAuthHandler_Login(t *testing.T) {
	router, db := setupAuthTest(t)

	// 先创建用户
	hashedPassword, _ := utils.HashPassword("password123")
	user := &model.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: hashedPassword,
		Role:     model.RoleStudent,
		Status:   model.StatusActive,
	}
	db.Create(user)

	// 测试登录
	body := map[string]string{
		"username": "testuser",
		"password": "password123",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, float64(0), resp["code"])

	data := resp["data"].(map[string]interface{})
	assert.NotEmpty(t, data["access_token"])
	assert.NotEmpty(t, data["refresh_token"])
	assert.NotEmpty(t, data["expires_in"])

	userData := data["user"].(map[string]interface{})
	assert.Equal(t, "testuser", userData["username"])
}

func TestAuthHandler_Login_InvalidPassword(t *testing.T) {
	router, db := setupAuthTest(t)

	// 先创建用户
	hashedPassword, _ := utils.HashPassword("password123")
	user := &model.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: hashedPassword,
		Role:     model.RoleStudent,
		Status:   model.StatusActive,
	}
	db.Create(user)

	// 测试错误密码登录
	body := map[string]string{
		"username": "testuser",
		"password": "wrongpassword",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.NotEqual(t, float64(0), resp["code"])
}
