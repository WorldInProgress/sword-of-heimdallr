package protocol

import (
	"encoding/json"
	"time"
)

// MessageTrace 定义消息追踪结构
type MessageTrace struct {
    TraceId    string       `json:"trace_id"`     // 追踪ID
    StartTime  time.Time    `json:"start_time"`   // 消息创建时间
    Hops       []MessageHop `json:"hops"`         // 消息经过的服务节点
    TotalTime  Duration     `json:"total_time"`   // 总处理时间
}

// MessageHop 定义消息经过的每个服务节点信息
type MessageHop struct {
    ServiceId   string   `json:"service_id"`    // 服务ID
    ServiceName string   `json:"service_name"`  // 服务名称
    HostName    string   `json:"host_name"`     // 主机名
    EntryTime   time.Time `json:"entry_time"`   // 进入服务时间
    ExitTime    time.Time `json:"exit_time"`    // 离开服务时间
    Duration    Duration  `json:"duration"`      // 处理耗时
    Status      string    `json:"status"`       // 处理状态
    Error       string    `json:"error,omitempty"` // 错误信息（如果有）
}

// Duration 自定义时间类型，支持更友好的JSON序列化
type Duration time.Duration

// MarshalJSON 实现Duration的JSON序列化
func (d Duration) MarshalJSON() ([]byte, error) {
    return json.Marshal(time.Duration(d).String())
}

// UnmarshalJSON 实现Duration的JSON反序列化
func (d *Duration) UnmarshalJSON(b []byte) error {
    var v string
    if err := json.Unmarshal(b, &v); err != nil {
        return err
    }
    parsed, err := time.ParseDuration(v)
    if err != nil {
        return err
    }
    *d = Duration(parsed)
    return nil
}

// NewMessageTrace 创建新的消息追踪
func NewMessageTrace() *MessageTrace {
    return &MessageTrace{
        TraceId:   GenerateUUID(),
        StartTime: time.Now(),
        Hops:      make([]MessageHop, 0),
    }
}

// AddHop 添加一个服务节点
func (mt *MessageTrace) AddHop(serviceId, serviceName, hostName string) *MessageHop {
    hop := MessageHop{
        ServiceId:   serviceId,
        ServiceName: serviceName,
        HostName:    hostName,
        EntryTime:   time.Now(),
    }
    mt.Hops = append(mt.Hops, hop)
    return &mt.Hops[len(mt.Hops)-1]
}

// CompleteHop 完成当前服务节点的处理
func (h *MessageHop) Complete(status string, err error) {
    h.ExitTime = time.Now()
    h.Duration = Duration(h.ExitTime.Sub(h.EntryTime))
    h.Status = status
    if err != nil {
        h.Error = err.Error()
    }
}

// CalculateTotalTime 计算消息总处理时间
func (mt *MessageTrace) CalculateTotalTime() {
    if len(mt.Hops) == 0 {
        mt.TotalTime = 0
        return
    }
    
    firstHop := mt.Hops[0]
    lastHop := mt.Hops[len(mt.Hops)-1]
    mt.TotalTime = Duration(lastHop.ExitTime.Sub(firstHop.EntryTime))
}

// GetHopByService 获取指定服务的处理信息
func (mt *MessageTrace) GetHopByService(serviceName string) *MessageHop {
    for i := range mt.Hops {
        if mt.Hops[i].ServiceName == serviceName {
            return &mt.Hops[i]
        }
    }
    return nil
}

// String 实现消息追踪的字符串表示
func (mt *MessageTrace) String() string {
    data, _ := json.MarshalIndent(mt, "", "  ")
    return string(data)
}