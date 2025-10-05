// backend/utils/hash.go (新增文件)
package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

// HashString 使用 SHA-256 对字符串进行哈希处理
func HashString(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}