package protocol

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

// MessageBuilder 使用构建器模式创建消息
type MessageBuilder struct {
    message *Message
}

func NewMessageBuilder() *MessageBuilder {
    return &MessageBuilder{
        message: &Message{
            Header: Header{
                MsgId:       GenerateUUID(),
                SessionId:   "",  // 需要显式设置
                UserId:      "",  // 需要显式设置
                Timestamp:   time.Now(),
                MsgType:     "",  // 需要显式设置
                Compression: CompressNone,
                Encoding:    EncodeJSON,
                Transport:   "",  // 需要显式设置
                Version:     ProtocolVersion,
            },
            ParentHeader: Header{},  // 保留空初始化，可通过WithParentMessage/WithParentHeader设置
            Meta:        Metadata{}, // 保留空初始化，可通过WithPriority/WithTags等方法设置
            Trace:       NewMessageTrace(), // 保留初始化，确保追踪功能可用
            Security:    SecurityConfig{},  // 保留空初始化，可通过WithSecurity等方法设置
        },
    }
}

// 必需的设置方法
func (b *MessageBuilder) WithType(msgType string) *MessageBuilder {
    b.message.Header.MsgType = msgType
    return b
}

func (b *MessageBuilder) WithSession(sessionId string) *MessageBuilder {
    b.message.Header.SessionId = sessionId
    return b
}

func (b *MessageBuilder) WithUser(userId string) *MessageBuilder {
    b.message.Header.UserId = userId
    return b
}

func (b *MessageBuilder) WithTransport(transport Transport) *MessageBuilder {
    b.message.Header.Transport = transport
    return b
}

// 可选的设置方法
func (b *MessageBuilder) WithCompression(compression Compression) *MessageBuilder {
    b.message.Header.Compression = compression
    return b
}

func (b *MessageBuilder) WithEncoding(encoding Encoding) *MessageBuilder {
    b.message.Header.Encoding = encoding
    return b
}

func (b *MessageBuilder) WithContent(content interface{}) *MessageBuilder {
    b.message.Content = content
    return b
}

// Meta 相关方法
func (b *MessageBuilder) WithPriority(priority Priority) *MessageBuilder {
    b.message.Meta.Priority = priority
    return b
}

func (b *MessageBuilder) WithTags(tags []string) *MessageBuilder {
    b.message.Meta.Tags = tags
    return b
}

func (b *MessageBuilder) AddTag(tag string) *MessageBuilder {
    b.message.Meta.Tags = append(b.message.Meta.Tags, tag)
    return b
}

// Security 相关方法
func (b *MessageBuilder) WithToken(token string) *MessageBuilder {
    b.message.Security.Token = token
    return b
}

func (b *MessageBuilder) WithEncryption(encryption string) *MessageBuilder {
    b.message.Security.Encryption = encryption
    return b
}

// 便捷方法：同时设置 token 和加密方式
func (b *MessageBuilder) WithSecurity(token, encryption string) *MessageBuilder {
    b.message.Security.Token = token
    b.message.Security.Encryption = encryption
    return b
}

// ParentHeader 相关方法
func (b *MessageBuilder) WithParentMessage(parent *Message) *MessageBuilder {
    if parent != nil {
        b.message.ParentHeader = parent.Header
    }
    return b
}

// 也可以直接设置 ParentHeader
func (b *MessageBuilder) WithParentHeader(header Header) *MessageBuilder {
    b.message.ParentHeader = header
    return b
}

// Trace 相关方法
func (b *MessageBuilder) WithNewTrace() *MessageBuilder {
    b.message.Trace = NewMessageTrace()
    return b
}

func (b *MessageBuilder) WithTrace(trace *MessageTrace) *MessageBuilder {
    b.message.Trace = trace
    return b
}

// 便捷方法：同时添加 trace 并记录第一个 hop
func (b *MessageBuilder) WithTraceHop(serviceId, serviceName, hostName string) *MessageBuilder {
    if b.message.Trace == nil {
        b.message.Trace = NewMessageTrace()
    }
    b.message.Trace.AddHop(serviceId, serviceName, hostName)
    return b
}

// GenerateUUID 生成UUID
func GenerateUUID() string {
    bytes := make([]byte, 16)
    rand.Read(bytes)
    return hex.EncodeToString(bytes)
}

// Build 方法包含必要的验证
func (b *MessageBuilder) Build() (*Message, error) {
    // 验证必需字段
    if b.message.Header.MsgType == "" {
        return nil, ErrInvalidMessageType
    }
    if b.message.Header.SessionId == "" {
        return nil, NewProtocolError(ErrCodeInvalidMessage, "session_id is required", nil)
    }
    if b.message.Header.UserId == "" {
        return nil, NewProtocolError(ErrCodeInvalidMessage, "user_id is required", nil)
    }
    if b.message.Header.Transport == "" {
        return nil, NewProtocolError(ErrCodeInvalidMessage, "transport is required", nil)
    }

    return b.message, nil
}