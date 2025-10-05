// /backend/middleware/auth.go
package middleware

import (
	"context"
	"net/http"
	"strings"
	"your-blog/backend/auth"
	"your-blog/backend/database"

	"github.com/gin-gonic/gin"
)

const ADMIN_EMAIL = "1612632430@qq.com"

// CheckAccess 中间件，检查游客模式或JWT
func CheckAccess() gin.HandlerFunc {
	return func(c *gin.Context) {
		var value string
		sql := `SELECT value FROM system_settings WHERE key = 'visitor_mode_enabled'`
		err := database.Pool.QueryRow(context.Background(), sql).Scan(&value)

		// 如果查询出错或值为 true，则认为是游客模式
		if err != nil || value == "true" {
			c.Next()
			return
		}

		// 游客模式关闭，需要验证 JWT
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "需要认证"})
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := auth.ValidateJWT(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "无效的 Token"})
			return
		}

		if claims.Email != ADMIN_EMAIL {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "无权访问"})
			return
		}

		// 将用户信息存入上下文，以便后续 handler 使用
		c.Set("user_email", claims.Email)
		c.Next()
	}
}

// RequireAdmin 中间件，强制要求管理员权限 (用于保护 /admin 接口)
func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "需要认证"})
			return
		}
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := auth.ValidateJWT(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "无效的 Token"})
			return
		}
		if claims.Email != ADMIN_EMAIL {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "需要管理员权限"})
			return
		}
		c.Set("user_email", claims.Email)
		c.Next()
	}
}