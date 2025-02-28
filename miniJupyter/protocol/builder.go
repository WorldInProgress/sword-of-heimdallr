package protocol

import "time"

// MessageBuilder 使用构建器模式创建消息
type MessageBuilder struct {
    message *Message
}

func NewMessageBuilder() *MessageBuilder {
    return &MessageBuilder{
        message: &Message{
            Header:   Header{Version: "0.3"},
            Meta:     Metadata{},
            Security: SecurityConfig{},
        },
    }
}

func (b *MessageBuilder) WithType(msgType string) *MessageBuilder {
    b.message.Header.MsgType = msgType
    return b
}

func (b *MessageBuilder) WithContent(content interface{}) *MessageBuilder {
    b.message.Content = content
    return b
}

func (b *MessageBuilder) WithPriority(priority Priority) *MessageBuilder {
    b.message.Meta.Priority = priority
    return b
}

func (b *MessageBuilder) Build() (*Message, error) {
    // 验证并返回消息
    if !IsValidMessageType(b.message.Header.MsgType) {
        return nil, ErrInvalidMessageType
    }
    
    b.message.Header.MsgId = GenerateUUID()
    b.message.Header.Timestamp = time.Now()
    
    return b.message, nil
}