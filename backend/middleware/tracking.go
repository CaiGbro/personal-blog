// backend/middleware/tracking.go (新增文件)
package middleware

import (
	"context"
	"log"
	"your-blog/backend/database"
	"your-blog/backend/utils"

	"github.com/gin-gonic/gin"
)

// TrackVisitor 中间件用于追踪独立访客
func TrackVisitor() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 将追踪逻辑放入 goroutine，使其完全不会阻塞主请求流程
		go func(clientIP string) {
			// 如果无法获取IP，则直接返回
			if clientIP == "" {
				return
			}
			
			// 使用盐值增加哈希的安全性，防止彩虹表攻击
			// 盐值应该是一个固定的、保密的字符串，可以从环境变量中读取
			// const visitorSalt = os.Getenv("VISITOR_IP_SALT")
			const visitorSalt = "a-fixed-secret-salt-for-visitors" // 示例盐值
			
			hashedIP := utils.HashString(clientIP + visitorSalt)

			// 使用 ON CONFLICT DO NOTHING，如果记录已存在则忽略，性能极高
			sql := `INSERT INTO visitor_logs (hashed_ip) VALUES ($1) ON CONFLICT (hashed_ip, visit_date) DO NOTHING`
			
			_, err := database.Pool.Exec(context.Background(), sql, hashedIP)
			if err != nil {
				// 只在后台记录错误，不影响用户
				log.Printf("Error tracking visitor: %v", err)
			}
		}(c.ClientIP())

		// 立即继续处理请求
		c.Next()
	}
}