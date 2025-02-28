package protocol

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

// NewMessage 创建新消息
func NewMessage(msgType string, content interface{}, transport Transport) *Message {
    return &Message{
        Header: Header{
            MsgId:       GenerateUUID(),
            SessionId:   GenerateUUID(),
            UserId:      "",  // 需要外部设置
            Timestamp:   time.Now(),
            MsgType:     msgType,
            Compression: CompressNone,
            Encoding:    EncodeJSON,
            Transport:   transport,
            Version:     "0.3",
        },
        Meta: Metadata{
            Priority: PriorityNormal,
            Tags:     []string{},
        },
        Content:  content,
        Security: SecurityConfig{},
    }
}

// GenerateUUID 生成UUID
func GenerateUUID() string {
    bytes := make([]byte, 16)
    rand.Read(bytes)
    return hex.EncodeToString(bytes)
}