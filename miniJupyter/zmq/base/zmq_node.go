package base

import (
	"os"

	zmq "github.com/pebbe/zmq4"
	"gopkg.in/yaml.v2"
)

// 读取配置文件
type Config struct {
    Zmq struct {
        RouterAddress     string `yaml:"router_address"`
        DealerAddress     string `yaml:"dealer_address"`
        PubAddress        string `yaml:"pub_address"`
        SubAddress        string `yaml:"sub_address"`
        HeartbeatInterval int    `yaml:"heartbeat_interval"`
        HeartbeatTimeout  int    `yaml:"heartbeat_timeout"`
    } `yaml:"zmq"`
}

// 解析 YAML 配置
func LoadConfig(filename string) (*Config, error) {
    data, err := os.ReadFile(filename)
    if err != nil {
        return nil, err
    }
    var config Config
    err = yaml.Unmarshal(data, &config)
    if err != nil {
        return nil, err
    }
    return &config, nil
}

// ZmqNode 结构体，封装 ZMQ 逻辑
type ZmqNode struct {
    socket *zmq.Socket
}

// 创建 ZMQ 端点
func NewZmqNode(socketType zmq.Type, address string, bind bool) (*ZmqNode, error) {
    socket, err := zmq.NewSocket(socketType)
    if err != nil {
        return nil, err
    }
    if bind {
        socket.Bind(address)
    } else {
        socket.Connect(address)
    }
    return &ZmqNode{socket: socket}, nil
}

// 发送消息，可以发送字符串和字节数组
func (z *ZmqNode) Send(msg ...string) error {
    _, err := z.socket.SendMessage(msg)
    return err
}

// 接收消息
func (z *ZmqNode) Receive() ([]string, error) {
    return z.socket.RecvMessage(0)
}

// 关闭 ZMQ 连接
func (z *ZmqNode) Close() {
    z.socket.Close()
}

// SetSubscribe 设置订阅主题
func (z *ZmqNode) SetSubscribe(topic string) error {
    return z.socket.SetSubscribe(topic)
}
