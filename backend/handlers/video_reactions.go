// /backend/handlers/video_reactions.go
package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"your-blog/backend/database"

	"github.com/gin-gonic/gin"
)

// GetVideoReactions 获取指定视频的所有表情回应
func GetVideoReactions(c *gin.Context) {
	// 使用 Param 方法来正确处理包含斜杠的文件路径
	videoFilename := c.Param("video_filename")
	if !strings.HasPrefix(videoFilename, "/") {
		videoFilename = "/" + videoFilename
	}

	var reactions json.RawMessage
	sql := `SELECT reactions FROM video_reactions WHERE video_filename = $1`

	// QueryRow 在找不到记录时会返回 err
	err := database.Pool.QueryRow(context.Background(), sql, videoFilename).Scan(&reactions)
	if err != nil {
		// 如果没有记录，则返回一个空的 JSON 对象，这是正常情况
		c.JSON(http.StatusOK, gin.H{})
		return
	}

	var reactionsMap map[string]int
	_ = json.Unmarshal(reactions, &reactionsMap)
	c.JSON(http.StatusOK, reactionsMap)
}

// AddVideoReaction 为指定视频添加一个表情回应
func AddVideoReaction(c *gin.Context) {
	videoFilename := c.Param("video_filename")
	if !strings.HasPrefix(videoFilename, "/") {
		videoFilename = "/" + videoFilename
	}

	var payload struct {
		Emoji string `json:"emoji"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil || payload.Emoji == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// 使用强大的 "UPSERT" (Update or Insert) SQL 语句
	sql := `
		INSERT INTO video_reactions (video_filename, reactions)
		VALUES ($1, jsonb_build_object($2::text, 1))
		ON CONFLICT (video_filename)
		DO UPDATE SET
			reactions = jsonb_set(
				video_reactions.reactions,
				ARRAY[$2::text], -- 关键修复：在这里显式指定类型为 text
				to_jsonb(COALESCE((video_reactions.reactions->>$2::text)::int, 0) + 1) -- 为保持健壮性，这里也加上
			)
		RETURNING reactions
	`

	var updatedReactions json.RawMessage
	err := database.Pool.QueryRow(context.Background(), sql, videoFilename, payload.Emoji).Scan(&updatedReactions)
	if err != nil {
		log.Printf("Database error on updating video reaction: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update video reaction"})
		return
	}

	var reactionsMap map[string]int
	_ = json.Unmarshal(updatedReactions, &reactionsMap)
	c.JSON(http.StatusOK, reactionsMap)
}