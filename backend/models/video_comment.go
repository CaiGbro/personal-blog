// /backend/models/video_comment.go
package models

import "time"

// VideoComment 定义了与视频关联的评论结构
type VideoComment struct {
	ID            int64          `json:"id"`
	VideoFilename string         `json:"videoFilename"` // 关联的视频文件名
	Nickname      string         `json:"nickname"`
	Content       string         `json:"content"`
	CreatedAt     time.Time      `json:"createdAt"`
	ParentID      *int64         `json:"parentId"`
	AttachmentURL *string        `json:"attachmentUrl"`
	Reactions     map[string]int `json:"reactions"`
	Replies       []*Comment     `json:"replies"` // 这里可以复用 Comment 结构来构建树
}