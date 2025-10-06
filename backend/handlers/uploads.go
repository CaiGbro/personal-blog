// /backend/handlers/uploads.go (最终修复版)
package handlers

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// --- 辅助函数，用于获取基于可执行文件的绝对上传路径 ---
func getUploadPath(subDir string) (string, error) {
	// 获取当前运行的程序（即可执行文件）的路径
	exePath, err := os.Executable()
	if err != nil {
		return "", err
	}

	// 获取可执行文件所在的目录 (例如 /prodct/personal-blog/personal-blog_prduct/backend/)
	exeDir := filepath.Dir(exePath)

	// 构建到 static 文件夹下指定子目录的绝对路径
	return filepath.Join(exeDir, "static", subDir), nil
}


// --- 通用的验证和保存函数 ---
func handleFileUpload(c *gin.Context, destSubDir string) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file is received"})
		return
	}

	// 文件大小验证
	const MAX_UPLOAD_SIZE = 5 * 1024 * 1024 // 5 MB
	if file.Size > MAX_UPLOAD_SIZE {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is too large. Maximum size is 5MB."})
		return
	}

	// 文件类型验证
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
		return
	}
	defer src.Close()

	buffer := make([]byte, 512)
	_, err = src.Read(buffer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file for validation"})
		return
	}
	contentType := http.DetectContentType(buffer)

	allowedMIMETypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/gif":  true,
		"image/webp": true,
	}
	if !allowedMIMETypes[contentType] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type. Only JPG, PNG, GIF, WebP are allowed."})
		return
	}

	// --- 获取绝对路径并保存文件 ---
	destPath, err := getUploadPath(destSubDir)
	if err != nil {
		log.Printf("Error getting absolute upload path: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server configuration error"})
		return
	}

	// 确保目标目录存在
	if err := os.MkdirAll(destPath, 0755); err != nil {
		log.Printf("Error creating upload directory '%s': %v", destPath, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to create storage directory"})
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	newFileName := uuid.New().String() + ext
	uploadPath := filepath.Join(destPath, newFileName)

	if err := c.SaveUploadedFile(file, uploadPath); err != nil {
		log.Printf("Error saving uploaded file to '%s': %v", uploadPath, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save the file"})
		return
	}

	// 返回给前端的URL路径保持不变
	fileURL := filepath.ToSlash(filepath.Join("/static", destSubDir, newFileName))
	c.JSON(http.StatusOK, gin.H{"url": fileURL})
}

// --- 重构现有的 Handler ---
func UploadFile(c *gin.Context) {
	handleFileUpload(c, "uploads")
}

func UploadVideoCommentFile(c *gin.Context) {
	handleFileUpload(c, "video_comment_uploads")
}

func UploadWritingCommentFile(c *gin.Context) {
	handleFileUpload(c, "writing_comment_uploads")
}