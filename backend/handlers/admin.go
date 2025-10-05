// /backend/handlers/admin.go (已修复)
package handlers

import (
	"context"
	"log"
	"net/http"
	"strconv" // <-- 1. 导入 strconv 包
	"your-blog/backend/database"

	"github.com/gin-gonic/gin"
)

// GetSystemSettings 函数保持不变
func GetSystemSettings(c *gin.Context) {
	var key, value string
	sql := `SELECT key, value FROM system_settings WHERE key = 'visitor_mode_enabled'`
	err := database.Pool.QueryRow(context.Background(), sql).Scan(&key, &value)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"visitor_mode_enabled": true})
		return
	}
	c.JSON(http.StatusOK, gin.H{"visitor_mode_enabled": value == "true"})
}

// UpdateSystemSettings 函数被修改
func UpdateSystemSettings(c *gin.Context) {
	var payload struct {
		VisitorModeEnabled bool `json:"visitor_mode_enabled"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求体"})
		return
	}

	sql := `
        INSERT INTO system_settings (key, value)
        VALUES ('visitor_mode_enabled', $1)
        ON CONFLICT (key) DO UPDATE
        SET value = $1
    `

	// 2. 将布尔值转换为字符串 "true" 或 "false"
	valueStr := strconv.FormatBool(payload.VisitorModeEnabled)

	// 3. 将转换后的字符串传递给数据库
	_, err := database.Pool.Exec(context.Background(), sql, valueStr)
	if err != nil {
		// 如果仍然出错，打印详细日志以帮助调试
		c.Error(err) 
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法更新设置"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "设置已更新", "visitor_mode_enabled": payload.VisitorModeEnabled})
}

// GetVisitorStats 获取访客统计数据
func GetVisitorStats(c *gin.Context) {
	// 默认获取过去30天的数据
	sql := `
        SELECT
            to_char(visit_date, 'YYYY-MM-DD') as date,
            COUNT(id) as count
        FROM visitor_logs
        WHERE visit_date >= CURRENT_DATE - INTERVAL '30 days'
        GROUP BY visit_date
        ORDER BY visit_date ASC
    `
	rows, err := database.Pool.Query(context.Background(), sql)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法获取访客统计数据"})
		return
	}
	defer rows.Close()

	// 使用匿名结构体以保持 handler 的简洁性
	var stats []struct {
		Date  string `json:"date"`
		Count int    `json:"count"`
	}

	for rows.Next() {
		var stat struct {
			Date  string
			Count int
		}
		if err := rows.Scan(&stat.Date, &stat.Count); err != nil {
			// 仅记录错误，继续处理下一行
			log.Printf("Error scanning visitor stat row: %v", err)
			continue
		}
		stats = append(stats, struct {
			Date  string `json:"date"`
			Count int    `json:"count"`
		}{Date: stat.Date, Count: stat.Count})
	}

	c.JSON(http.StatusOK, stats)
}