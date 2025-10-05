// /backend/main.go
package main

import (
	"context"
	"log"
	"net/http" // <-- 导入 net/http
	"os"
	"os/signal" // <-- 导入 os/signal
	"syscall"   // <-- 导入 syscall
	"time"      // <-- 导入 time
	"your-blog/backend/database"
	"your-blog/backend/router"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, loading from system environment variables")
	}

	if err := database.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Pool.Close()

	r := router.SetupRouter()

	// --- 创建并配置一个自定义的 HTTP 服务器 ---
	srv := &http.Server{
		Addr:    ":8081",
		Handler: r,
		// ReadTimeout: 读取整个请求（包括头和体）的最长时间
		ReadTimeout: 10 * time.Second,
		// WriteTimeout: 写入响应的最长时间
		WriteTimeout: 15 * time.Second,
		// IdleTimeout: 连接在 Keep-Alive 状态下的最长空闲时间
		IdleTimeout: 60 * time.Second,
	}

	// --- 使用 goroutine 启动服务器，这样不会阻塞主线程 ---
	go func() {
		log.Println("Server is running on port 8081")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// --- 实现优雅停机 (Graceful Shutdown) ---
	// 创建一个通道，用于接收操作系统信号
	quit := make(chan os.Signal, 1)
	// 我们只关心 SIGINT (Ctrl+C) 和 SIGTERM (kill 命令)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 阻塞主线程，直到接收到一个信号
	<-quit
	log.Println("Shutting down server...")

	// 创建一个有 5 秒超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 调用 Shutdown()，它会尝试优雅地关闭服务器：
	// 1. 停止接受新的连接
	// 2. 等待现有请求处理完成（在上下文超时之前）
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}