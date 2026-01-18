package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/study-upc/backend/internal/pkg/utils"
	"github.com/study-upc/backend/internal/service"
)

// EmailVerificationHandler 邮箱验证处理器
type EmailVerificationHandler struct {
	emailService service.EmailVerificationService
	jwtManager   *utils.JWTManager
}

// NewEmailVerificationHandler 创建邮箱验证处理器
func NewEmailVerificationHandler(
	emailService service.EmailVerificationService,
	jwtManager *utils.JWTManager,
) *EmailVerificationHandler {
	return &EmailVerificationHandler{
		emailService: emailService,
		jwtManager:   jwtManager,
	}
}

// SendCodeRequest 发送验证码请求
type SendCodeRequest struct {
	Email   string `json:"email" binding:"required,email"`
	Purpose string `json:"purpose" binding:"required,oneof=register login reset_password"`
}

// SendCodeResponse 发送验证码响应
type SendCodeResponse struct {
	Message string `json:"message"`
}

// SendVerificationCode 发送验证码
// @Summary 发送邮箱验证码
// @Description 发送邮箱验证码用于注册、登录或重置密码
// @Tags 邮箱验证
// @Accept json
// @Produce json
// @Param request body SendCodeRequest true "发送验证码请求"
// @Success 200 {object} SendCodeResponse
// @Router /api/v1/verification/send [post]
func (h *EmailVerificationHandler) SendVerificationCode(c *gin.Context) {
	var req SendCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    20001,
			"message": "参数错误: " + err.Error(),
			"data":    nil,
		})
		return
	}

	// 发送验证码
	if err := h.emailService.SendVerificationCode(c.Request.Context(), req.Email, req.Purpose); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    50001,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "验证码已发送到您的邮箱",
		"data": SendCodeResponse{
			Message: "验证码已发送,请查收邮件",
		},
	})
}

// VerifyCodeRequest 验证验证码请求
type VerifyCodeRequest struct {
	Email   string `json:"email" binding:"required,email"`
	Code    string `json:"code" binding:"required,len=6"`
	Purpose string `json:"purpose" binding:"required,oneof=register login reset_password"`
}

// VerifyCode 验证验证码
// @Summary 验证邮箱验证码
// @Description 验证邮箱验证码是否正确
// @Tags 邮箱验证
// @Accept json
// @Produce json
// @Param request body VerifyCodeRequest true "验证验证码请求"
// @Success 200 {object} SendCodeResponse
// @Router /api/v1/verification/verify [post]
func (h *EmailVerificationHandler) VerifyCode(c *gin.Context) {
	var req VerifyCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    20001,
			"message": "参数错误: " + err.Error(),
			"data":    nil,
		})
		return
	}

	// 验证验证码
	if err := h.emailService.VerifyCode(c.Request.Context(), req.Email, req.Code, req.Purpose); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    40001,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "验证码验证成功",
		"data": nil,
	})
}

// RegisterWithCodeRequest 使用验证码注册请求
type RegisterWithCodeRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=50"`
	Code     string `json:"code" binding:"required,len=6"`
}

// RegisterWithCode 使用验证码注册
// @Summary 使用邮箱验证码注册
// @Description 使用邮箱验证码完成用户注册
// @Tags 用户认证
// @Accept json
// @Produce json
// @Param request body RegisterWithCodeRequest true "注册请求"
// @Success 200 {object} SendCodeResponse
// @Router /api/v1/auth/register [post]
func (h *EmailVerificationHandler) RegisterWithCode(c *gin.Context) {
	var req RegisterWithCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    20001,
			"message": "参数错误: " + err.Error(),
			"data":    nil,
		})
		return
	}

	// 注册用户
	if err := h.emailService.RegisterWithEmailCode(
		c.Request.Context(),
		req.Username,
		req.Email,
		req.Password,
		req.Code,
	); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    40002,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "注册成功",
		"data": nil,
	})
}

// LoginWithCodeRequest 使用验证码登录请求
type LoginWithCodeRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"omitempty"`
	Code     string `json:"code" binding:"required,len=6"`
}

// LoginWithCode 使用验证码登录
// @Summary 使用邮箱验证码登录
// @Description 使用邮箱验证码和密码登录
// @Tags 用户认证
// @Accept json
// @Produce json
// @Param request body LoginWithCodeRequest true "登录请求"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/auth/login [post]
func (h *EmailVerificationHandler) LoginWithCode(c *gin.Context) {
	var req LoginWithCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    20001,
			"message": "参数错误: " + err.Error(),
			"data":    nil,
		})
		return
	}

	// 登录验证
	_, user, err := h.emailService.LoginWithEmailCode(
		c.Request.Context(),
		req.Email,
		req.Password,
		req.Code,
	)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    40003,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	// 生成JWT token
	token, err := h.jwtManager.GenerateAccessToken(user.ID, string(user.Role))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    50001,
			"message": "生成token失败",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "登录成功",
		"data": gin.H{
			"token": token,
			"user": gin.H{
				"id":       user.ID,
				"username": user.Username,
				"email":    user.Email,
				"role":     user.Role,
				"status":   user.Status,
			},
		},
	})
}
