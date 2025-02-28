package protocol

// 协议版本
const (
    ProtocolVersion = "0.4"
)

// 枚举类型定义
type (
    Compression string
    Encoding    string
    Transport   string
    Priority    string
    Status      string
    StreamType  string
    RetryStrategy string
)

const (
    // Compression
    CompressNone   Compression = "none"
    CompressGzip   Compression = "gzip"
    CompressSnappy Compression = "snappy"

    // Encoding
    EncodeJSON     Encoding = "json"
    EncodeProtobuf Encoding = "protobuf"
    EncodeCustom   Encoding = "custom"

    // Transport
    TransportZMQ   Transport = "zmq"
    TransportGRPC  Transport = "grpc"

    // Priority
    PriorityHigh   Priority = "HIGH"
    PriorityNormal Priority = "NORMAL"
    PriorityLow    Priority = "LOW"

    // Status
    StatusOK       Status = "ok"
    StatusError    Status = "error"
    StatusStarting Status = "starting"
    StatusWaiting  Status = "waiting"
    StatusSuccess  Status = "success"

    // StreamType
    StreamStdout StreamType = "stdout"
    StreamStderr StreamType = "stderr"

    // RetryStrategy
    RetryExponentialBackoff RetryStrategy = "exponential_backoff"

    // 定义消息类型常量
    MsgTypeExecuteRequest  = "execute_request"
    MsgTypeExecuteReply    = "execute_reply"
    MsgTypeExecuteResult   = "execute_result"
    MsgTypeCoreInfoRequest = "core_info_request"
    MsgTypeCoreInfoReply   = "core_info_reply"
    MsgTypeStream         = "stream"
    MsgTypeCommOpen       = "comm_open"
    MsgTypeCommMsg        = "comm_msg"
    MsgTypeCommClose      = "comm_close"
)

// 添加消息类型检查
func IsValidMessageType(msgType string) bool {
    switch msgType {
    case MsgTypeExecuteRequest, MsgTypeExecuteReply, MsgTypeExecuteResult,
         MsgTypeCoreInfoRequest, MsgTypeCoreInfoReply, MsgTypeStream,
         MsgTypeCommOpen, MsgTypeCommMsg, MsgTypeCommClose:
        return true
    }
    return false
}