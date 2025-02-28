package protocol

import "time"

// 基础消息结构
type Message struct {
    Header        Header                 `json:"header"`
    ParentHeader  Header                 `json:"parent_header"`
    Meta          Metadata               `json:"meta"`
    Content       interface{}            `json:"content"`
    Security      SecurityConfig         `json:"security"`
    Trace         *MessageTrace          `json:"trace"`
}

// Header 定义
type Header struct {
    MsgId       string      `json:"msg_id"`
    SessionId   string      `json:"session_id"`
    UserId      string      `json:"user_id"`
    Timestamp   time.Time   `json:"timestamp"`
    MsgType     string      `json:"msg_type"`
    Compression Compression `json:"compression"`
    Encoding    Encoding    `json:"encoding"`
    Transport   Transport   `json:"transport"`
    Version     string      `json:"version"`
}

// Metadata 定义
type Metadata struct {
    Priority Priority `json:"priority"`
    Tags     []string `json:"tags"`
}

// Security 定义
type SecurityConfig struct {
    Token      string `json:"token"`
    Encryption string `json:"encryption"`
}

///////////////////////////////////////////////////////////////////////////////////////

// Execute Request Content
type ExecuteRequestContent struct {
    CommandId    string                 `json:"command_id"`
    Service      string                 `json:"service"`
    Method       string                 `json:"method"`
    Params       map[string]interface{} `json:"params"`
    Condition    map[string]interface{} `json:"condition"`
    Dependency   []string               `json:"dependency"`
    Timeout      int                    `json:"timeout"`
    Retry        RetryConfig           `json:"retry"`
    StopOnError  bool                  `json:"stop_on_error"`
    AllowedUsers []string              `json:"allowed_users"`
}

type RetryConfig struct {
    MaxAttempts int           `json:"max_attempts"`
    Strategy    RetryStrategy `json:"strategy"`
}

// Execute Reply Content
type ExecuteReplyContent struct {
    Status Status `json:"status"`
}

// Core Info Reply Content
type CoreInfoContent struct {
    Status            Status `json:"status"`
    CoreStatus        string `json:"core_status"`
    CoreVersion       string `json:"core_version"`
    CPUUsage         string `json:"cpu_usage"`
    MemoryUsage      string `json:"memory_usage"`
    DiskUsage        string `json:"disk_usage"`
    NetworkIO        string `json:"network_io"`
    ActiveConnections int    `json:"active_connections"`
    RunningTasks      int    `json:"running_tasks"`
    TaskQueueSize     int    `json:"task_queue_size"`
}

// Execute Result Content
type ExecuteResultContent struct {
    Status Status      `json:"status"`
    Result interface{} `json:"result"`
}

// Stream Content
type StreamContent struct {
    Type StreamType `json:"type"`
    Text string     `json:"text"`
}

// Comm Messages
type CommOpenContent struct {
    CommId     string      `json:"comm_id"`
    TargetName string      `json:"target_name"`
    Data       interface{} `json:"data"`
}

type CommMsgContent struct {
    CommId string      `json:"comm_id"`
    Data   interface{} `json:"data"`
}

///////////////////////////////////////////////////////////////////////////////////////

// Message 的追踪相关方法
func (m *Message) AddTrace(serviceId, serviceName, hostName string) *MessageHop {
    // 如果还没有追踪信息，创建一个
    if m.Trace == nil {
        m.Trace = NewMessageTrace()
    }
    
    // 添加服务节点
    return m.Trace.AddHop(serviceId, serviceName, hostName)
}