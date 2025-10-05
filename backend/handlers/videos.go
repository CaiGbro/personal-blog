// /backend/handlers/videos.go
package handlers

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

// GetVideos 处理获取视频列表的请求
func GetVideos(c *gin.Context) {
	videoDir := "./static/video"
	var videoURLs []string
	entries, err := os.ReadDir(videoDir)
	if err != nil {
		log.Printf("Error reading video directory: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not list videos"})
		return
	}
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(strings.ToLower(entry.Name()), ".mp4") {
			urlPath := filepath.ToSlash(filepath.Join("/static", "video", entry.Name()))
			videoURLs = append(videoURLs, urlPath)
		}
	}
	c.JSON(http.StatusOK, videoURLs)
}