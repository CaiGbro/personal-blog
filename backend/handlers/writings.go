// /backend/handlers/writings.go (最终健壮版)
package handlers

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort" // <-- 引入 sort 包
	"strconv"
	"strings"
	"github.com/gin-gonic/gin"
)

type WritingInfo struct {
	Filename     string `json:"filename"`
	URL          string `json:"url"`
	ThumbnailURL string `json:"thumbnailUrl"`
}

// GetWritings 处理获取文章列表的请求
func GetWritings(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "5"))
	if page < 1 { page = 1 }
	if pageSize < 1 { pageSize = 5 }
	searchQuery := c.DefaultQuery("search", "")

	writingsDir := "./static/writings"
	entries, err := os.ReadDir(writingsDir)
	if err != nil {
		log.Printf("Error reading writings directory: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not list writings"})
		return
	}

	supportedDocExtensions := map[string]bool{".pdf": true, ".md": true, ".docx": true}
	supportedImgExtensions := map[string]bool{".png": true, ".jpg": true, ".jpeg": true}

	// --- 核心修改：使用更清晰的逻辑进行去重 ---
	writingsMap := make(map[string]WritingInfo)

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		
		entryName := entry.Name()
		ext := strings.ToLower(filepath.Ext(entryName))

		if supportedDocExtensions[ext] {
			baseName := strings.TrimSuffix(entryName, ext)
			
			// 检查 map 中是否已存在该条目
			existingEntry, exists := writingsMap[baseName]

			// 决定是否要添加或覆盖的逻辑：
			// 1. 如果 map 中不存在此条目，直接添加。
			// 2. 如果 map 中已存在，但新文件是 .pdf 而旧文件不是，则用 .pdf 覆盖。
			if !exists || (ext == ".pdf" && filepath.Ext(existingEntry.Filename) != ".pdf") {
				var thumbnailURL string
				for imgExt := range supportedImgExtensions {
					thumbnailPath := filepath.Join(writingsDir, baseName+imgExt)
					if _, err := os.Stat(thumbnailPath); err == nil {
						thumbnailURL = filepath.ToSlash(filepath.Join("/static", "writings", baseName+imgExt))
						break
					}
				}
				
				writingsMap[baseName] = WritingInfo{
					Filename:     entryName,
					URL:          filepath.ToSlash(filepath.Join("/static", "writings", entryName)),
					ThumbnailURL: thumbnailURL,
				}
			}
		}
	}

	// 将 map 转换为切片
	var allWritings []WritingInfo
	for _, writing := range writingsMap {
		allWritings = append(allWritings, writing)
	}
	// --- 修改结束 ---

	// 新增：按文件名排序，确保每次请求的顺序一致
	sort.Slice(allWritings, func(i, j int) bool {
		return allWritings[i].Filename < allWritings[j].Filename
	})

	// 执行搜索过滤
	if searchQuery != "" {
		var filteredWritings []WritingInfo
		lowerSearchQuery := strings.ToLower(searchQuery)
		for _, writing := range allWritings {
			if strings.Contains(strings.ToLower(writing.Filename), lowerSearchQuery) {
				filteredWritings = append(filteredWritings, writing)
			}
		}
		allWritings = filteredWritings
	}

	// 对最终的列表进行分页
	start := (page - 1) * pageSize
	end := start + pageSize
	if start >= len(allWritings) {
		c.JSON(http.StatusOK, gin.H{"writings": []WritingInfo{}, "hasMore": false})
		return
	}
	if end > len(allWritings) {
		end = len(allWritings)
	}

	paginatedWritings := allWritings[start:end]
	hasMore := end < len(allWritings)

	c.JSON(http.StatusOK, gin.H{
		"writings": paginatedWritings,
		"hasMore":  hasMore,
	})
}