// /backend/handlers/video_comments.go
package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"strconv"
	"your-blog/backend/database"
	"your-blog/backend/models"

	"github.com/gin-gonic/gin"
	"github.com/microcosm-cc/bluemonday"
)

// CreateVideoComment 处理创建新视频评论的请求
func CreateVideoComment(c *gin.Context) {
	var newComment models.VideoComment
	if err := c.ShouldBindJSON(&newComment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	if newComment.Nickname == "" || newComment.Content == "" || newComment.VideoFilename == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nickname, content, and videoFilename cannot be empty"})
		return
	}

	// --- 2. XSS 防护：净化用户输入的内容 ---
	p := bluemonday.UGCPolicy() // 创建一个适用于用户生成内容的策略
	newComment.Content = p.Sanitize(newComment.Content)
	// --- 防护结束 ---

	sql := `INSERT INTO video_comments (video_filename, nickname, content, parent_id, attachment_url) VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at`
	err := database.Pool.QueryRow(context.Background(), sql, newComment.VideoFilename, newComment.Nickname, newComment.Content, newComment.ParentID, newComment.AttachmentURL).Scan(&newComment.ID, &newComment.CreatedAt)
	if err != nil {
		log.Printf("Database error on creating video comment: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create video comment"})
		return
	}
	c.JSON(http.StatusCreated, newComment)
}

// GetVideoComments 处理根据视频文件名获取评论的请求
func GetVideoComments(c *gin.Context) {
	videoFilename := c.Query("video_filename")
	if videoFilename == "" { c.JSON(http.StatusBadRequest, gin.H{"error": "video_filename parameter is required"}); return }

	var totalCount int
	countSql := `SELECT COUNT(*) FROM video_comments WHERE video_filename = $1`
	if err := database.Pool.QueryRow(context.Background(), countSql, videoFilename).Scan(&totalCount); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get video comments count"}); return
	}

	rowsSql := `SELECT id FROM video_comments WHERE parent_id IS NULL AND video_filename = $1 ORDER BY created_at DESC`
	rows, err := database.Pool.Query(context.Background(), rowsSql, videoFilename)
	if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get video comments"}); return }
	defer rows.Close()

	var topLevelIDs []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil { continue }
		topLevelIDs = append(topLevelIDs, id)
	}

	if len(topLevelIDs) == 0 { c.JSON(http.StatusOK, gin.H{"totalCount": totalCount, "comments": []models.Comment{}}); return }

	recursiveSql := `
        WITH RECURSIVE comment_tree AS (
            SELECT id, nickname, content, created_at, parent_id, attachment_url, reactions FROM video_comments WHERE id = ANY($1)
            UNION ALL
            SELECT c.id, c.nickname, c.content, c.created_at, c.parent_id, c.attachment_url, c.reactions FROM video_comments c JOIN comment_tree ct ON c.parent_id = ct.id
        )
        SELECT id, nickname, content, created_at, parent_id, attachment_url, reactions FROM comment_tree`
	allRows, err := database.Pool.Query(context.Background(), recursiveSql, topLevelIDs)
	if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get video comment tree"}); return }
	defer allRows.Close()
	
	commentMap := make(map[int64]*models.Comment)
	var allComments []*models.Comment
	for allRows.Next() {
		var comment models.Comment
		var reactions json.RawMessage
		if err := allRows.Scan(&comment.ID, &comment.Nickname, &comment.Content, &comment.CreatedAt, &comment.ParentID, &comment.AttachmentURL, &reactions); err != nil { continue }
		if reactions != nil { _ = json.Unmarshal(reactions, &comment.Reactions) }
		comment.Replies = []*models.Comment{}
		commentMap[comment.ID] = &comment
		allComments = append(allComments, &comment)
	}

	sort.Slice(allComments, func(i, j int) bool { return allComments[i].CreatedAt.Before(allComments[j].CreatedAt) })

	for _, comment := range allComments {
		if comment.ParentID != nil {
			if parent, ok := commentMap[*comment.ParentID]; ok { parent.Replies = append(parent.Replies, comment) }
		}
	}
	
	finalSortedComments := make([]*models.Comment, 0, len(topLevelIDs))
	for _, id := range topLevelIDs {
		if comment, ok := commentMap[id]; ok { finalSortedComments = append(finalSortedComments, comment) }
	}

	c.JSON(http.StatusOK, gin.H{"totalCount": totalCount, "comments": finalSortedComments})
}


// AddVideoCommentReaction 处理为视频评论添加表情回应的请求
func AddVideoCommentReaction(c *gin.Context) {
	commentID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"}); return }

	var payload struct { Emoji string `json:"emoji"` }
	if err := c.ShouldBindJSON(&payload); err != nil || payload.Emoji == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"}); return
	}

	sql := `UPDATE video_comments SET reactions = jsonb_set(COALESCE(reactions, '{}'::jsonb), ARRAY[$1], to_jsonb(COALESCE((reactions->>$1)::int, 0) + 1)) WHERE id = $2 RETURNING reactions`
	var updatedReactions json.RawMessage
	err = database.Pool.QueryRow(context.Background(), sql, payload.Emoji, commentID).Scan(&updatedReactions)
	if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update reaction"}); return }

	var reactionsMap map[string]int
	_ = json.Unmarshal(updatedReactions, &reactionsMap)
	c.JSON(http.StatusOK, reactionsMap)
}