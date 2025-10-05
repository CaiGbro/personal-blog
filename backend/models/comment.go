// /backend/models/comment.go
package models

import "time"

type Comment struct {
    ID            int64     `json:"id"`
    Nickname      string    `json:"nickname"`
    Content       string    `json:"content"`
    CreatedAt     time.Time `json:"createdAt"`
    // 新增字段
    ParentID      *int64    `json:"parentId"` // 使用指针以接受 null 值
    AttachmentURL *string   `json:"attachmentUrl"`
    // reactions 会被处理成 map[string]int 返回
    Reactions     map[string]int `json:"reactions"` 
    // 用于构建评论树，不在数据库中
    Replies       []*Comment `json:"replies"` 
}