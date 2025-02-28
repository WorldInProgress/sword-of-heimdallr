package protocol

import (
	"encoding/json"
	"fmt"
)

// GetContentType 根据消息类型返回对应的 Content 结构体
func GetContentType(msgType string) interface{} {
	switch msgType {
	// ROUTER/DEALER 消息
	case MsgTypeExecuteRequest:
		return &ExecuteRequestContent{}
	case MsgTypeExecuteReply:
		return &ExecuteReplyContent{}
	case MsgTypeCoreInfoRequest:
		return &struct{}{} // 空结构体，因为该请求没有content
	case MsgTypeCoreInfoReply:
		return &CoreInfoContent{}

	// PUB/SUB 消息
	case MsgTypeExecuteResult:
		return &ExecuteResultContent{}
	case MsgTypeStream:
		return &StreamContent{}

	// Comm 消息
	case MsgTypeCommOpen:
		return &CommOpenContent{}
	case MsgTypeCommMsg:
		return &CommMsgContent{}
	case MsgTypeCommClose:
		return &CommMsgContent{} // CommClose 使用相同的结构

	default:
		return nil
	}
}

// ParseMessage 智能解析消息
func ParseMessage(data []byte) (*Message, error) {
	// 1. 先解析基础消息结构
	var msg Message
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, fmt.Errorf("failed to parse message: %w", err)
	}
	
	// 2. 获取正确的 Content 类型
	contentType := GetContentType(msg.Header.MsgType)
	if contentType == nil {
		return nil, fmt.Errorf("unknown message type: %s", msg.Header.MsgType)
	}
	
	// 3. 重新解析 Content 到正确的类型
	contentData, err := json.Marshal(msg.Content)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal content: %w", err)
	}
	
	if err := json.Unmarshal(contentData, contentType); err != nil {
		return nil, fmt.Errorf("failed to parse content: %w", err)
	}
	
	msg.Content = contentType
	return &msg, nil
}