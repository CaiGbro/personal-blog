// /backend/handlers/uploads.go
package handlers

import (
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)
// --- 新增：一个通用的验证和保存函数 ---
func handleFileUpload(c *gin.Context, destPath string, destSubDir string) {
    file, err := c.FormFile("file")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "No file is received"})
        return
    }

    // --- 1. 后端文件大小验证 ---
    // r.MaxMultipartMemory 已经限制了总体大小，这里可以再对单个文件做更精细的控制
    const MAX_UPLOAD_SIZE = 5 * 1024 * 1024 // 5 MB
    if file.Size > MAX_UPLOAD_SIZE {
        c.JSON(http.StatusBadRequest, gin.H{"error": "File is too large. Maximum size is 5MB."})
        return
    }

    // --- 2. 后端文件类型验证 ---
    src, err := file.Open()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
        return
    }
    defer src.Close()

    // 读取文件头 512 字节来判断真实类型
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
        "image/webp": true, // 现代浏览器支持
    }
    if !allowedMIMETypes[contentType] {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type. Only JPG, PNG, GIF, WebP are allowed."})
        return
    }

    // --- 3. 保存文件 ---
    ext := filepath.Ext(file.Filename)
    newFileName := uuid.New().String() + ext
    uploadPath := filepath.Join(destPath, newFileName)

    if err := c.SaveUploadedFile(file, uploadPath); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save the file"})
        return
    }

    fileURL := filepath.ToSlash(filepath.Join(destSubDir, newFileName))
    c.JSON(http.StatusOK, gin.H{"url": fileURL})
}

// --- 重构现有的 Handler ---
func UploadFile(c *gin.Context) {
	handleFileUpload(c, "./static/uploads", "/static/uploads")
}

func UploadVideoCommentFile(c *gin.Context) {
	handleFileUpload(c, "./static/video_comment_uploads", "/static/video_comment_uploads")
}

func UploadWritingCommentFile(c *gin.Context) {
	handleFileUpload(c, "./static/writing_comment_uploads", "/static/writing_comment_uploads")
}