// /backend/handlers/auth.go
package handlers

import (
	"context"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
	"your-blog/backend/auth"
	"your-blog/backend/database"
	"your-blog/backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
)

const ADMIN_EMAIL = "1612632430@qq.com"

// SendVerificationCode 发送验证码
func SendVerificationCode(c *gin.Context) {
	var payload struct {
		Email string `json:"email"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil || payload.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "需要提供有效的邮箱"})
		return
	}

	// 只有管理员邮箱才能收到验证码
	if payload.Email != ADMIN_EMAIL {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权登录，此非开发者邮箱"})
		return
	}

	code := strconv.Itoa(100000 + rand.Intn(900000)) // 生成6位随机数
	expiresAt := time.Now().Add(5 * time.Minute)

	sql := `
        INSERT INTO verification_codes (email, code, expires_at)
        VALUES ($1, $2, $3)
        ON CONFLICT (email) DO UPDATE
        SET code = $2, expires_at = $3
    `
	_, err := database.Pool.Exec(context.Background(), sql, payload.Email, code, expiresAt)
	if err != nil {
		log.Printf("数据库错误，无法保存验证码: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法保存验证码"})
		return
	}

	// 发送邮件 (在后台协程中执行以避免阻塞)
	go func() {
		err := utils.SendVerificationCode(payload.Email, code)
		if err != nil {
			log.Printf("发送邮件失败: %v", err)
		}
	}()

	c.JSON(http.StatusOK, gin.H{"message": "验证码已发送，请查收邮件"})
}

// Login 使用邮箱和验证码登录
func Login(c *gin.Context) {
	var payload struct {
		Email string `json:"email"`
		Code  string `json:"code"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil || payload.Email == "" || payload.Code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "邮箱和验证码不能为空"})
		return
	}

	if payload.Email != ADMIN_EMAIL {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权登录，此非开发者邮箱"})
		return
	}

	var storedCode string
	var expiresAt time.Time
	sql := `SELECT code, expires_at FROM verification_codes WHERE email = $1`
	err := database.Pool.QueryRow(context.Background(), sql, payload.Email).Scan(&storedCode, &expiresAt)
	if err == pgx.ErrNoRows {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请先获取验证码"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误"})
		return
	}

	if time.Now().After(expiresAt) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "验证码已过期"})
		return
	}

	if storedCode != payload.Code {
		c.JSON(http.StatusBadRequest, gin.H{"error": "验证码不正确"})
		return
	}

	// 验证成功，生成 JWT
	token, err := auth.GenerateJWT(payload.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法生成 token"})
		return
	}

	// 登录成功后删除验证码
	deleteSql := `DELETE FROM verification_codes WHERE email = $1`
	database.Pool.Exec(context.Background(), deleteSql, payload.Email)

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"email": payload.Email,
		},
	})
}