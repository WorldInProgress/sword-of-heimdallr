package utils

import (
	zmq "github.com/pebbe/zmq4"
)

// CreateSocket 创建并配置 ZMQ socket
func CreateSocket(socketType zmq.Type, address string, bind bool) *zmq.Socket {
	socket, err := zmq.NewSocket(socketType)
	if err != nil {
		panic(err)
    }

    if bind {
        err = socket.Bind(address)
    } else {
        err = socket.Connect(address)
    }
    if err != nil {
        panic(err)
    }

	return socket
}

// SendMessage 发送消息
func SendMessage(socket *zmq.Socket, message string) {
    _, err := socket.Send(message, 0)
    if err != nil {
        panic(err)
    }
}

// ReceiveMessage 接收消息
func ReceiveMessage(socket *zmq.Socket) string {
	message, err := socket.Recv(0)
	if err != nil {
        panic(err)
    }
    return message
}