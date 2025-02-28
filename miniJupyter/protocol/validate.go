package protocol

import (
	"errors"
	"fmt"
)

// Validator 接口定义消息验证方法
type Validator interface {
    Validate() error
}

// ValidateMessage 验证整个消息结构
func ValidateMessage(msg *Message) error {
    // 验证Header
    if err := validateHeader(&msg.Header); err != nil {
        return fmt.Errorf("invalid header: %w", err)
    }

    // 验证Content
    if validator, ok := msg.Content.(Validator); ok {
        if err := validator.Validate(); err != nil {
            return fmt.Errorf("invalid content: %w", err)
        }
    }

    return nil
}

// validateHeader 验证消息头
func validateHeader(h *Header) error {
    if h.MsgId == "" {
        return errors.New("msg_id is required")
    }
    if h.SessionId == "" {
        return errors.New("session_id is required")
    }
    if h.UserId == "" {
        return errors.New("user_id is required")
    }
    if !IsValidMessageType(h.MsgType) {
        return fmt.Errorf("invalid message type: %s", h.MsgType)
    }
    if h.Version != ProtocolVersion {
        return fmt.Errorf("unsupported version: %s", h.Version)
    }
    return nil
}

// ExecuteRequestContent 验证
func (c *ExecuteRequestContent) Validate() error {
    if c.CommandId == "" {
        return errors.New("command_id is required")
    }
    if c.Service == "" {
        return errors.New("service is required")
    }
    if c.Method == "" {
        return errors.New("method is required")
    }
    if c.Timeout < 0 {
        return errors.New("timeout cannot be negative")
    }
    if c.Retry.MaxAttempts < 0 {
        return errors.New("retry max_attempts cannot be negative")
    }
    return nil
}

// ExecuteReplyContent 验证
func (c *ExecuteReplyContent) Validate() error {
    switch c.Status {
    case StatusError, StatusStarting, StatusWaiting:
        return nil
    default:
        return fmt.Errorf("invalid status: %s", c.Status)
    }
}

// CoreInfoContent 验证
func (c *CoreInfoContent) Validate() error {
    if c.CoreVersion == "" {
        return errors.New("core_version is required")
    }
    if c.ActiveConnections < 0 {
        return errors.New("active_connections cannot be negative")
    }
    if c.RunningTasks < 0 {
        return errors.New("running_tasks cannot be negative")
    }
    if c.TaskQueueSize < 0 {
        return errors.New("task_queue_size cannot be negative")
    }
    return nil
}

// ExecuteResultContent 验证
func (c *ExecuteResultContent) Validate() error {
    switch c.Status {
    case StatusSuccess, StatusError:
        return nil
    default:
        return fmt.Errorf("invalid status: %s", c.Status)
    }
}

// StreamContent 验证
func (c *StreamContent) Validate() error {
    switch c.Type {
    case StreamStdout, StreamStderr:
        return nil
    default:
        return fmt.Errorf("invalid stream type: %s", c.Type)
    }
}

// CommOpenContent 验证
func (c *CommOpenContent) Validate() error {
    if c.CommId == "" {
        return errors.New("comm_id is required")
    }
    if c.TargetName == "" {
        return errors.New("target_name is required")
    }
    return nil
}

// CommMsgContent 验证
func (c *CommMsgContent) Validate() error {
    if c.CommId == "" {
        return errors.New("comm_id is required")
    }
    return nil
}