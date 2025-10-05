// /backend/models/writing_comment.go
package models

import "time"

// WritingComment 定义了与文章关联的评论结构
type WritingComment struct {
	ID            int64          `json:"id"`
	WritingFilename string       `json:"writingFilename"` // 关联的文章文件名
	Nickname      string         `json:"nickname"`
	Content       string         `json:"content"`
	CreatedAt     time.Time      `json:"createdAt"`
	ParentID      *int64         `json:"parentId"`
	AttachmentURL *string        `json:"attachmentUrl"`
	Reactions     map[string]int `json:"reactions"`
	Replies       []*Comment     `json:"replies"` // 复用 Comment 结构来构建树
}