package middleware

import (
	"github.com/study-upc/backend/internal/model"
	"github.com/study-upc/backend/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

// RequireRole 要求特定角色的中间件
func RequireRole(roles ...model.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从上下文获取用户角色（由 JWT 中间件设置）
		userRole, exists := c.Get("user_role")
		if !exists {
			response.Error(c, response.ErrUnauthorized, "未认证")
			c.Abort()
			return
		}

		// 检查角色是否在允许列表中
		roleStr := userRole.(string)
		allowed := false
		for _, r := range roles {
			if string(r) == roleStr {
				allowed = true
				break
			}
		}

		if !allowed {
			response.Error(c, response.ErrForbidden, "权限不足")
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAdmin 要求管理员角色的中间件
func RequireAdmin() gin.HandlerFunc {
	return RequireRole(model.RoleAdmin)
}

// RequireCommittee 要求学委或管理员角色的中间件
func RequireCommittee() gin.HandlerFunc {
	return RequireRole(model.RoleCommittee, model.RoleAdmin)
}

// RequireStudent 要求学生角色的中间件（实际上所有用户都至少是学生）
func RequireStudent() gin.HandlerFunc {
	return RequireRole(model.RoleStudent, model.RoleCommittee, model.RoleAdmin)
}

// IsAdmin 判断当前用户是否是管理员
func IsAdmin(c *gin.Context) bool {
	userRole, exists := c.Get("user_role")
	if !exists {
		return false
	}
	return userRole.(string) == string(model.RoleAdmin)
}

// IsCommittee 判断当前用户是否是学委或管理员
func IsCommittee(c *gin.Context) bool {
	userRole, exists := c.Get("user_role")
	if !exists {
		return false
	}
	roleStr := userRole.(string)
	return roleStr == string(model.RoleCommittee) || roleStr == string(model.RoleAdmin)
}

// GetUserID 从上下文获取用户 ID
func GetUserID(c *gin.Context) (uint, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}
	return userID.(uint), true
}

// GetUserRole 从上下文获取用户角色
func GetUserRole(c *gin.Context) (string, bool) {
	userRole, exists := c.Get("user_role")
	if !exists {
		return "", false
	}
	return userRole.(string), true
}
