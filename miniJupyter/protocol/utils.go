package protocol

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateMessageId 生成唯一消息ID
func GenerateMessageId() string {
    bytes := make([]byte, 8)
    rand.Read(bytes)
    return hex.EncodeToString(bytes)
}

// GenerateSessionId 生成会话ID
func GenerateSessionId() string {
    bytes := make([]byte, 16)
    rand.Read(bytes)
    return hex.EncodeToString(bytes)
}