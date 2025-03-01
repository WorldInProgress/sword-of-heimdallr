package base

import (
	"os"

	"github.com/pebbe/zmq4"
	"gopkg.in/yaml.v2"
)

// 读取配置文件
type Config struct {
    Zmq struct {
        RouterAddress     string `yaml:"router_address"`
        DealerAddress     string `yaml:"dealer_address"`
        HeartbeatInterval int    `yaml:"heartbeat_interval"`
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
    socket *zmq4.Socket
}

// 创建 ZMQ 端点
func NewZmqNode(socketType zmq4.Type, address string, bind bool) (*ZmqNode, error) {
    socket, err := zmq4.NewSocket(socketType)
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

// 发送消息
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
