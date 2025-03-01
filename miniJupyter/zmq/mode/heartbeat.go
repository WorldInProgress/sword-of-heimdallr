package mode

import (
	"fmt"
	"time"

	"zmq/utils"

	zmq "github.com/pebbe/zmq4"
)

// 心跳检测服务器
func RunHeartbeatServer(address string) {
	router := utils.CreateSocket(zmq.ROUTER, address, true)
	defer router.Close()

	for {
		identity, _ := router.Recv(0)
		msg := utils.ReceiveMessage(router)
		if msg == "PING" {
			fmt.Printf("[Heartbeat Server] Received PING from %s\n", identity)
			router.Send(identity, zmq.SNDMORE)
			utils.SendMessage(router, "PONG")
		}
	}
}

// 心跳检测客户端
func RunHeartbeatClient(address string, interval int) {
	dealer := utils.CreateSocket(zmq.DEALER, address, false)
	defer dealer.Close()

	for {
		utils.SendMessage(dealer, "PING")
		reply := utils.ReceiveMessage(dealer)
		if reply == "PONG" {
			fmt.Println("[Heartbeat Client] Received PONG")
		}
		time.Sleep(time.Duration(interval) * time.Second)
	}
}
