// /backend/middleware/ratelimit.go
package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// IPRateLimiter 结构体用于存储每个 IP 的请求记录
type IPRateLimiter struct {
	ips    map[string][]time.Time // 键是 IP 地址，值是请求时间戳的列表
	mu     sync.Mutex             // 互斥锁，用于保护 ips map 的并发访问安全
	limit  int                    // 在时间窗口内允许的最大请求数
	period time.Duration          // 时间窗口的长度
}

// NewIPRateLimiter 是 IPRateLimiter 的构造函数
func NewIPRateLimiter(limit int, period time.Duration) *IPRateLimiter {
	return &IPRateLimiter{
		ips:    make(map[string][]time.Time),
		limit:  limit,
		period: period,
	}
}

// RateLimitMiddleware 创建一个 Gin 中间件函数
func (limiter *IPRateLimiter) RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		limiter.mu.Lock()         // 加锁以安全地访问 map
		defer limiter.mu.Unlock() // 确保函数结束时解锁

		ip := c.ClientIP()
		now := time.Now()

		// 1. 清理该 IP 的过期时间戳
		// 创建一个新的切片，只保留在时间窗口内的记录
		var recentTimestamps []time.Time
		if timestamps, found := limiter.ips[ip]; found {
			for _, ts := range timestamps {
				if now.Sub(ts) < limiter.period {
					recentTimestamps = append(recentTimestamps, ts)
				}
			}
		}

		// 2. 检查请求次数是否超限
		if len(recentTimestamps) >= limiter.limit {
			// 如果超限，则拒绝请求
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests, please try again later."})
			return
		}

		// 3. 记录本次请求的时间戳
		limiter.ips[ip] = append(recentTimestamps, now)

		// 4. 继续处理请求
		c.Next()
	}
}