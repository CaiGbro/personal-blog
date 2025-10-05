// /backend/handlers/writing_reactions.go
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

// GetWritingReactions 获取指定文章的所有表情回应
func GetWritingReactions(c *gin.Context) {
	writingFilename := c.Param("writing_filename")
	if !strings.HasPrefix(writingFilename, "/") {
		writingFilename = "/" + writingFilename
	}

	var reactions json.RawMessage
	sql := `SELECT reactions FROM writing_reactions WHERE writing_filename = $1`

	err := database.Pool.QueryRow(context.Background(), sql, writingFilename).Scan(&reactions)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{})
		return
	}

	var reactionsMap map[string]int
	_ = json.Unmarshal(reactions, &reactionsMap)
	c.JSON(http.StatusOK, reactionsMap)
}

// AddWritingReaction 为指定文章添加一个表情回应
func AddWritingReaction(c *gin.Context) {
	writingFilename := c.Param("writing_filename")
	if !strings.HasPrefix(writingFilename, "/") {
		writingFilename = "/" + writingFilename
	}

	var payload struct {
		Emoji string `json:"emoji"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil || payload.Emoji == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	sql := `
		INSERT INTO writing_reactions (writing_filename, reactions)
		VALUES ($1, jsonb_build_object($2::text, 1))
		ON CONFLICT (writing_filename)
		DO UPDATE SET
			reactions = jsonb_set(
				writing_reactions.reactions,
				ARRAY[$2::text],
				to_jsonb(COALESCE((writing_reactions.reactions->>$2::text)::int, 0) + 1)
			)
		RETURNING reactions
	`

	var updatedReactions json.RawMessage
	err := database.Pool.QueryRow(context.Background(), sql, writingFilename, payload.Emoji).Scan(&updatedReactions)
	if err != nil {
		log.Printf("Database error on updating writing reaction: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update writing reaction"})
		return
	}

	var reactionsMap map[string]int
	_ = json.Unmarshal(updatedReactions, &reactionsMap)
	c.JSON(http.StatusOK, reactionsMap)
}