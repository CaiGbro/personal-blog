// /backend/router/router.go
package router

import (
	"log"
	"os"
	"path/filepath" // <--- 1. 导入 path/filepath 包
	"time"
	"your-blog/backend/handlers"
	"your-blog/backend/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRouter 配置并返回一个 Gin 引擎
func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.MaxMultipartMemory = 8 << 20

	// --- 核心修复：使用绝对路径来设置静态文件服务 ---
	// 2. 获取可执行文件的路径
	exePath, err := os.Executable()
	if err != nil {
		log.Fatalf("Failed to get executable path: %v", err)
	}
	// 3. 获取可执行文件所在的目录
	exeDir := filepath.Dir(exePath)
	// 4. 构建到 static 文件夹的绝对路径
	staticPath := filepath.Join(exeDir, "static")

	// 5. 使用这个绝对路径来创建目录和设置静态文件服务
	if err := os.MkdirAll(filepath.Join(staticPath, "uploads"), 0755); err != nil {
		log.Fatalf("Failed to create upload directory: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(staticPath, "video_comment_uploads"), 0755); err != nil {
		log.Fatalf("Failed to create video comment upload directory: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(staticPath, "writing_comment_uploads"), 0755); err != nil {
		log.Fatalf("Failed to create writing comment upload directory: %v", err)
	}

	// 6. 使用绝对路径设置静态文件服务
	r.Static("/static", staticPath)
	// --- 修复结束 ---


	// API 路由组
	api := r.Group("/api")

	// 这个中间件应该放在最前面，以确保所有 API 请求都被追踪
	api.Use(middleware.TrackVisitor())

	{
		// --- 公开的系统状态路由 ---
		api.GET("/settings", handlers.GetSystemSettings)
		// --- 新增：公开的广告路由 ---
		api.GET("/ads", handlers.GetAds)

		// --- 认证路由组 (已添加速率限制) ---
		// 1. 实例化我们自己的速率限制器
		//    规则：来自同一个 IP 地址，每分钟最多允许 5 次请求
		authLimiter := middleware.NewIPRateLimiter(5, time.Minute)

		auth := api.Group("/auth")
		// 2. 将速率限制中间件应用到整个 /auth 组
		auth.Use(authLimiter.RateLimitMiddleware())
		{
			auth.POST("/send-code", handlers.SendVerificationCode)
			auth.POST("/login", handlers.Login)
		}

		// --- 受访问控制保护的业务路由组 ---
		authorized := api.Group("/")
		authorized.Use(middleware.CheckAccess()) // 应用核心的访问控制中间件

		

		{
			// 留言板评论
			authorized.POST("/comments", handlers.CreateComment)
			authorized.GET("/comments", handlers.GetComments)
			authorized.POST("/comments/:id/react", handlers.AddReaction)

			// 视频评论
			authorized.POST("/video_comments", handlers.CreateVideoComment)
			authorized.GET("/video_comments", handlers.GetVideoComments)
			authorized.POST("/video_comments/:id/react", handlers.AddVideoCommentReaction)

			// 文章与文章评论
			authorized.GET("/writings", handlers.GetWritings)
			authorized.POST("/writing_comments", handlers.CreateWritingComment)
			authorized.GET("/writing_comments", handlers.GetWritingComments)
			authorized.POST("/writing_comments/:id/react", handlers.AddWritingCommentReaction)
			authorized.GET("/writing_comments/count", handlers.GetWritingCommentsCount)

			// 视频与回应
			authorized.GET("/videos", handlers.GetVideos)
			authorized.GET("/videos/reactions/*video_filename", handlers.GetVideoReactions)
			authorized.POST("/videos/react/*video_filename", handlers.AddVideoReaction)
			
			// 文章回应
			authorized.GET("/writings/reactions/*writing_filename", handlers.GetWritingReactions)
			authorized.POST("/writings/react/*writing_filename", handlers.AddWritingReaction)

			// 文件上传
			authorized.POST("/upload", handlers.UploadFile)
			authorized.POST("/upload/video_comment", handlers.UploadVideoCommentFile)
			authorized.POST("/upload/writing_comment", handlers.UploadWritingCommentFile)
		}

		// --- 管理员专属路由组 ---
		// 强制要求管理员登录，无论游客模式是否开启
		admin := api.Group("/admin")
		admin.Use(middleware.RequireAdmin()) // 应用强制管理员中间件
		{
			// GET /api/admin/settings 和 POST /api/admin/settings
			// 注意：GetSystemSettings 同时存在于公开路由和管理路由中，
			// 这是可以的，因为前端在不同场景下会携带不同的凭证访问。
			admin.GET("/settings", handlers.GetSystemSettings)
			admin.POST("/settings", handlers.UpdateSystemSettings)

			// --- 新增：注册访客统计路由 ---
			admin.GET("/stats/visitors", handlers.GetVisitorStats)
		}
	}

	return r
}