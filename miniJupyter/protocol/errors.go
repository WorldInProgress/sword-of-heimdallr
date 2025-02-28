package protocol

import "fmt"

// ProtocolError 定义协议错误类型
type ProtocolError struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Details interface{} `json:"details,omitempty"`
}

// Error 实现 error 接口
func (e *ProtocolError) Error() string {
    if e.Details != nil {
        return fmt.Sprintf("[%d] %s: %v", e.Code, e.Message, e.Details)
    }
    return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// NewProtocolError 创建新的协议错误
func NewProtocolError(code int, message string, details interface{}) *ProtocolError {
    return &ProtocolError{
        Code:    code,
        Message: message,
        Details: details,
    }
}

// 预定义错误码
const (
    // 1000-1099: Protocol level errors
    ErrCodeInvalidMessage     = 1000
    ErrCodeInvalidMessageType = 1001
    ErrCodeInvalidVersion     = 1002
    ErrCodeInvalidFormat      = 1003
    ErrCodeValidationFailed   = 1004
    ErrCodeSerializeFailed    = 1005
    ErrCodeDeserializeFailed  = 1006

    // 1100-1199: Authentication/Authorization errors
    ErrCodeUnauthorized       = 1100
    ErrCodeInvalidToken       = 1101
    ErrCodeInsufficientPerms  = 1102
    ErrCodeSessionExpired     = 1103

    // 1200-1299: Execution errors
    ErrCodeExecutionFailed    = 1200
    ErrCodeTimeout           = 1201
    ErrCodeDependencyFailed  = 1202
    ErrCodeServiceNotFound   = 1203
    ErrCodeMethodNotFound    = 1204
    ErrCodeInvalidParams     = 1205

    // 1300-1399: Communication errors
    ErrCodeConnectionFailed  = 1300
    ErrCodeHeartbeatTimeout  = 1301
    ErrCodeSubscribeFailed   = 1302
    ErrCodePublishFailed     = 1303
    ErrCodeCommFailed        = 1304
)

// 预定义错误实例
var (
    // Protocol errors
    ErrInvalidMessage     = NewProtocolError(ErrCodeInvalidMessage, "Invalid message format", nil)
    ErrInvalidMessageType = NewProtocolError(ErrCodeInvalidMessageType, "Invalid message type", nil)
    ErrInvalidVersion     = NewProtocolError(ErrCodeInvalidVersion, "Invalid protocol version", nil)
    ErrInvalidFormat      = NewProtocolError(ErrCodeInvalidFormat, "Invalid message format", nil)
    ErrValidationFailed   = NewProtocolError(ErrCodeValidationFailed, "Message validation failed", nil)
    ErrSerializeFailed    = NewProtocolError(ErrCodeSerializeFailed, "Message serialization failed", nil)
    ErrDeserializeFailed  = NewProtocolError(ErrCodeDeserializeFailed, "Message deserialization failed", nil)

    // Auth errors
    ErrUnauthorized       = NewProtocolError(ErrCodeUnauthorized, "Unauthorized access", nil)
    ErrInvalidToken       = NewProtocolError(ErrCodeInvalidToken, "Invalid token", nil)
    ErrInsufficientPerms  = NewProtocolError(ErrCodeInsufficientPerms, "Insufficient permissions", nil)
    ErrSessionExpired     = NewProtocolError(ErrCodeSessionExpired, "Session expired", nil)

    // Execution errors
    ErrExecutionFailed    = NewProtocolError(ErrCodeExecutionFailed, "Execution failed", nil)
    ErrTimeout           = NewProtocolError(ErrCodeTimeout, "Operation timeout", nil)
    ErrDependencyFailed  = NewProtocolError(ErrCodeDependencyFailed, "Dependency execution failed", nil)
    ErrServiceNotFound   = NewProtocolError(ErrCodeServiceNotFound, "Service not found", nil)
    ErrMethodNotFound    = NewProtocolError(ErrCodeMethodNotFound, "Method not found", nil)
    ErrInvalidParams     = NewProtocolError(ErrCodeInvalidParams, "Invalid parameters", nil)

    // Communication errors
    ErrConnectionFailed  = NewProtocolError(ErrCodeConnectionFailed, "Connection failed", nil)
    ErrHeartbeatTimeout  = NewProtocolError(ErrCodeHeartbeatTimeout, "Heartbeat timeout", nil)
    ErrSubscribeFailed   = NewProtocolError(ErrCodeSubscribeFailed, "Subscribe failed", nil)
    ErrPublishFailed     = NewProtocolError(ErrCodePublishFailed, "Publish failed", nil)
    ErrCommFailed        = NewProtocolError(ErrCodeCommFailed, "Comm operation failed", nil)
)

// WithDetails 添加错误详情
func (e *ProtocolError) WithDetails(details interface{}) *ProtocolError {
    return &ProtocolError{
        Code:    e.Code,
        Message: e.Message,
        Details: details,
    }
}

// IsProtocolError 检查是否为协议错误
func IsProtocolError(err error) bool {
    _, ok := err.(*ProtocolError)
    return ok
}

// GetErrorCode 获取错误码，如果不是协议错误则返回0
func GetErrorCode(err error) int {
    if pe, ok := err.(*ProtocolError); ok {
        return pe.Code
    }
    return 0
}