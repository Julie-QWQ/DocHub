package handler

import (
	"errors"
	"strings"

	"github.com/study-upc/backend/internal/model"
	"github.com/study-upc/backend/internal/pkg/response"
	"github.com/study-upc/backend/internal/pkg/utils"
	"github.com/study-upc/backend/internal/repository"
	"github.com/study-upc/backend/internal/service"

	"github.com/gin-gonic/gin"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	authService              service.AuthService
	emailVerificationService service.EmailVerificationService
	statisticsService        service.StatisticsService
	jwtManager               *utils.JWTManager
}

// NewAuthHandler 创建认证处理器实例
func NewAuthHandler(
	authService service.AuthService,
	statisticsService service.StatisticsService,
) *AuthHandler {
	return &AuthHandler{
		authService:       authService,
		statisticsService: statisticsService,
	}
}

// SetEmailVerificationService 设置邮箱验证服务（解决循环依赖）
func (h *AuthHandler) SetEmailVerificationService(emailVerificationService service.EmailVerificationService) {
	h.emailVerificationService = emailVerificationService
}

// SetJWTManager 设置 JWT 管理器（解决循环依赖）
func (h *AuthHandler) SetJWTManager(jwtManager *utils.JWTManager) {
	h.jwtManager = jwtManager
}

// Register 用户注册
// @Summary 用户注册
// @Description 注册新用户
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body model.RegisterRequest true "注册信息"
// @Success 200 {object} response.Response{data=model.UserInfo}
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.ErrInvalidParams, err.Error())
		return
	}

	userInfo, err := h.authService.Register(c.Request.Context(), &req)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrUserAlreadyExists):
			response.Error(c, response.ErrUserExists, err.Error())
		default:
			response.Error(c, response.ErrInternal, err.Error())
		}
		return
	}

	response.Success(c, userInfo)
}

// Login 用户登录
// @Summary 用户登录
// @Description 用户登录获取 Token（支持用户名登录或邮箱验证码登录）
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body model.LoginRequest true "登录信息"
// @Success 200 {object} response.Response{data=model.LoginResponse}
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	// 尝试解析为带验证码的登录请求
	var loginWithCodeReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Code     string `json:"code"`
		Username string `json:"username"`
	}

	if err := c.ShouldBindJSON(&loginWithCodeReq); err != nil {
		response.Error(c, response.ErrInvalidParams, err.Error())
		return
	}

	var loginResp *model.LoginResponse
	var err error

	// 如果有验证码，使用邮箱验证码登录
	if loginWithCodeReq.Code != "" {
		// 验证邮箱格式
		if loginWithCodeReq.Email == "" {
			response.Error(c, response.ErrInvalidParams, "使用验证码登录时必须提供邮箱")
			return
		}

		// 使用邮箱验证码登录
		_, user, loginErr := h.emailVerificationService.LoginWithEmailCode(
			c.Request.Context(),
			loginWithCodeReq.Email,
			loginWithCodeReq.Password,
			loginWithCodeReq.Code,
		)
		if loginErr != nil {
			// 记录登录失败的日志
			_ = h.statisticsService.RecordLoginLog(0, c.ClientIP(), c.Request.UserAgent(), false)

			if errors.Is(loginErr, service.ErrUserDisabled) || errors.Is(loginErr, service.ErrUserInactive) {
				response.Error(c, response.ErrUserDisabled, loginErr.Error())
				return
			}
			response.Error(c, response.ErrInvalidCredentials, loginErr.Error())
			return
		}

		// 生成 Token
		accessToken, refreshToken, err := h.jwtManager.GenerateTokenPair(user.ID, string(user.Role))
		if err != nil {
			response.Error(c, response.ErrInternal, "生成token失败")
			return
		}

		loginResp = &model.LoginResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			ExpiresIn:    h.jwtManager.GetAccessTTL(),
			User:         user.ToUserInfo(),
		}
	} else {
		// 普通登录
		req := model.LoginRequest{
			Username: loginWithCodeReq.Username,
			Password: loginWithCodeReq.Password,
		}

		loginResp, err = h.authService.Login(c.Request.Context(), &req)
		if err != nil {
			// 记录登录失败的日志
			_ = h.statisticsService.RecordLoginLog(0, c.ClientIP(), c.Request.UserAgent(), false)

			if errors.Is(err, service.ErrInvalidCredentials) {
				response.Error(c, response.ErrInvalidCredentials, err.Error())
				return
			}
			if errors.Is(err, service.ErrUserDisabled) || errors.Is(err, service.ErrUserInactive) {
				response.Error(c, response.ErrUserDisabled, err.Error())
				return
			}
			response.Error(c, response.ErrInternal, err.Error())
			return
		}
	}

	// 记录登录成功的日志（异步，不影响登录流程）
	go func() {
		_ = h.statisticsService.RecordLoginLog(loginResp.User.ID, c.ClientIP(), c.Request.UserAgent(), true)
	}()

	response.Success(c, loginResp)
}

// Logout 用户登出
// @Summary 用户登出
// @Description 用户登出，将 Token 加入黑名单
// @Tags 认证
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response
// @Router /api/v1/auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	// 从 Authorization Header 获取 Token
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		response.Error(c, response.ErrUnauthorized, "未提供认证 Token")
		return
	}

	// Bearer Token 格式
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		response.Error(c, response.ErrUnauthorized, "Token 格式错误")
		return
	}

	accessToken := parts[1]

	if err := h.authService.Logout(c.Request.Context(), accessToken); err != nil {
		response.Error(c, response.ErrInternal, err.Error())
		return
	}

	response.Success(c, nil)
}

// RefreshToken 刷新 Token
// @Summary 刷新 Token
// @Description 使用刷新 Token 获取新的访问 Token
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body model.RefreshTokenRequest true "刷新 Token"
// @Success 200 {object} response.Response{data=model.LoginResponse}
// @Router /api/v1/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req model.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.ErrInvalidParams, err.Error())
		return
	}

	loginResp, err := h.authService.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		if errors.Is(err, service.ErrTokenExpired) || errors.Is(err, service.ErrTokenInvalid) || errors.Is(err, service.ErrTokenInBlacklist) {
			response.Error(c, response.ErrInvalidToken, err.Error())
			return
		}
		if errors.Is(err, service.ErrUserDisabled) || errors.Is(err, service.ErrUserInactive) {
			response.Error(c, response.ErrUserDisabled, err.Error())
			return
		}
		response.Error(c, response.ErrInternal, err.Error())
		return
	}

	response.Success(c, loginResp)
}

// ChangePassword 修改密码
// @Summary 修改密码
// @Description 修改当前用户密码
// @Tags 认证
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body model.ChangePasswordRequest true "密码信息"
// @Success 200 {object} response.Response
// @Router /api/v1/auth/change-password [post]
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var req model.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.ErrInvalidParams, err.Error())
		return
	}

	// 从上下文获取用户 ID（由 JWT 中间件设置）
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, response.ErrUnauthorized, "未认证")
		return
	}

	if err := h.authService.ChangePassword(c.Request.Context(), userID.(uint), &req); err != nil {
		switch err {
		case service.ErrWrongPassword:
			response.Error(c, response.ErrWrongPassword, err.Error())
		default:
			response.Error(c, response.ErrInternal, err.Error())
		}
		return
	}

	response.Success(c, nil)
}

// GetUserInfo 获取当前用户信息
// @Summary 获取当前用户信息
// @Description 获取当前登录用户的详细信息
// @Tags 认证
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=model.UserInfo}
// @Router /api/v1/auth/me [get]
func (h *AuthHandler) GetUserInfo(c *gin.Context) {
	// 从上下文获取用户 ID（由 JWT 中间件设置）
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, response.ErrUnauthorized, "未认证")
		return
	}

	userInfo, err := h.authService.GetUserInfo(c.Request.Context(), userID.(uint))
	if err != nil {
		response.Error(c, response.ErrInternal, err.Error())
		return
	}

	response.Success(c, userInfo)
}
