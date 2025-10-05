// /backend/handlers/comments.go
package handlers

import (
	"context"
	"encoding/json"
	// "fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"your-blog/backend/database"
	"your-blog/backend/models"
	"github.com/microcosm-cc/bluemonday"
	"github.com/gin-gonic/gin"
)

// CreateComment 处理创建新评论的请求
func CreateComment(c *gin.Context) {
	var newComment models.Comment
	if err := c.ShouldBindJSON(&newComment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	if newComment.Nickname == "" || newComment.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nickname and content cannot be empty"})
		return
	}
	// --- XSS 防护 ---
    p := bluemonday.UGCPolicy() // 创建一个用于用户生成内容的策略
    newComment.Content = p.Sanitize(newComment.Content) // 清理 HTML
    // --- 防护结束 ---
	
	sql := `INSERT INTO comments (nickname, content, parent_id, attachment_url) VALUES ($1, $2, $3, $4) RETURNING id, created_at`
	err := database.Pool.QueryRow(context.Background(), sql, newComment.Nickname, newComment.Content, newComment.ParentID, newComment.AttachmentURL).Scan(&newComment.ID, &newComment.CreatedAt)
	if err != nil {
		log.Printf("Database error on creating comment: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		return
	}
	c.JSON(http.StatusCreated, newComment)
}

// GetComments 处理获取评论列表（包含树状结构）的请求
func GetComments(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "5"))
	if page < 1 { page = 1 }
	if pageSize < 1 { pageSize = 5 }
	offset := (page - 1) * pageSize

	var totalCount int
	countSql := `SELECT COUNT(*) FROM comments WHERE parent_id IS NULL`
	err := database.Pool.QueryRow(context.Background(), countSql).Scan(&totalCount)
	if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get comments count"}); return }

	rowsSql := `SELECT id FROM comments WHERE parent_id IS NULL ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	rows, err := database.Pool.Query(context.Background(), rowsSql, pageSize, offset)
	if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get comments"}); return }
	defer rows.Close()

	var topLevelIDs []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil { continue }
		topLevelIDs = append(topLevelIDs, id)
	}
	if len(topLevelIDs) == 0 {
		c.JSON(http.StatusOK, gin.H{"totalCount": totalCount, "comments": []models.Comment{}}); return
	}

	recursiveSql := `
		WITH RECURSIVE comment_tree AS (
			SELECT id, nickname, content, created_at, parent_id, attachment_url, reactions FROM comments WHERE id = ANY($1)
			UNION ALL
			SELECT c.id, c.nickname, c.content, c.created_at, c.parent_id, c.attachment_url, c.reactions FROM comments c JOIN comment_tree ct ON c.parent_id = ct.id
		)
		SELECT id, nickname, content, created_at, parent_id, attachment_url, reactions FROM comment_tree`
	allRows, err := database.Pool.Query(context.Background(), recursiveSql, topLevelIDs)
	if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get comment tree"}); return }
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

	var commentsTree []*models.Comment
	for _, comment := range allComments {
		if comment.ParentID != nil {
			if parent, ok := commentMap[*comment.ParentID]; ok { parent.Replies = append(parent.Replies, comment) }
		} else {
			commentsTree = append(commentsTree, comment)
		}
	}
	
	finalSortedComments := make([]*models.Comment, 0, len(topLevelIDs))
	for _, id := range topLevelIDs {
		if comment, ok := commentMap[id]; ok { finalSortedComments = append(finalSortedComments, comment) }
	}

	c.JSON(http.StatusOK, gin.H{"totalCount": totalCount, "comments": finalSortedComments})
}


// AddReaction 处理为评论添加表情回应的请求
func AddReaction(c *gin.Context) {
	commentID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"}); return }
	
	var payload struct { Emoji string `json:"emoji"` }
	if err := c.ShouldBindJSON(&payload); err != nil || payload.Emoji == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"}); return
	}

	sql := `UPDATE comments SET reactions = jsonb_set(COALESCE(reactions, '{}'::jsonb), ARRAY[$1], to_jsonb(COALESCE((reactions->>$1)::int, 0) + 1)) WHERE id = $2 RETURNING reactions`
	var updatedReactions json.RawMessage
	err = database.Pool.QueryRow(context.Background(), sql, payload.Emoji, commentID).Scan(&updatedReactions)
	if err != nil {
		log.Printf("Database error on updating reaction: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update reaction"}); return
	}

	var reactionsMap map[string]int
	if err := json.Unmarshal(updatedReactions, &reactionsMap); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse reactions"}); return
	}
	c.JSON(http.StatusOK, reactionsMap)
}