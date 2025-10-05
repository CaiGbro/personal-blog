// /backend/handlers/ads.go
package handlers

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

// GetAds 处理获取广告列表的请求
func GetAds(c *gin.Context) {
	adsDir := "./static/ads"
	var adURLs []string

	entries, err := os.ReadDir(adsDir)
	if err != nil {
		// 如果目录不存在，返回空列表是合理的，而不是错误
		if os.IsNotExist(err) {
			c.JSON(http.StatusOK, []string{})
			return
		}
		log.Printf("Error reading ads directory: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not list ads"})
		return
	}

	supportedExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			ext := strings.ToLower(filepath.Ext(entry.Name()))
			if supportedExtensions[ext] {
				// 构造可从前端访问的 URL 路径
				urlPath := filepath.ToSlash(filepath.Join("/static", "ads", entry.Name()))
				adURLs = append(adURLs, urlPath)
			}
		}
	}

	c.JSON(http.StatusOK, adURLs)
}