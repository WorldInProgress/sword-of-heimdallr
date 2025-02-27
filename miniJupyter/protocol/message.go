package protocol

import (
	"time"
)

// Header 定义消息头
type Header struct {
    MsgId    string    `json:"msg_id"`
    MsgType  string    `json:"msg_type"`
    Username string    `json:"username"`
    Session  string    `json:"session"`
    Date     time.Time `json:"date"`
    Version  string    `json:"version"`
}

// Metadata 定义元数据
type Metadata map[string]interface{}

// Message 定义基础消息结构
type Message struct {
    Header   Header                 `json:"header"`
    Parent   Header                 `json:"parent_header"`
    Metadata Metadata              `json:"metadata"`
    Content  interface{}           `json:"content"`
}

// 定义不同类型消息的Content结构
type ExecuteRequest struct {
    Code         string `json:"code"`
    Silent       bool   `json:"silent"`
    StoreHistory bool   `json:"store_history"`
}

type ExecuteReply struct {
    Status        string `json:"status"`
    ExecutionCount int    `json:"execution_count"`
    Payload       []interface{} `json:"payload"`
}

type HeartbeatRequest struct {}

type HeartbeatReply struct {
    Status string `json:"status"`
}

// 创建新消息的辅助函数
func NewMessage(msgType string, content interface{}) *Message {
    return &Message{
        Header: Header{
            MsgId:    GenerateMessageId(),
            MsgType:  msgType,
            Username: "anonymous",
            Session:  GenerateSessionId(),
            Date:     time.Now(),
            Version:  "5.3",
        },
        Metadata: make(Metadata),
        Content:  content,
    }
}

// 根据消息类型获取对应的Content类型
func GetContentType(msgType string) interface{} {
    switch msgType {
    case "execute_request":
        return &ExecuteRequest{}
    case "execute_reply":
        return &ExecuteReply{}
    case "heartbeat_request":
        return &HeartbeatRequest{}
    case "heartbeat_reply":
        return &HeartbeatReply{}
    default:
        return nil
    }
}