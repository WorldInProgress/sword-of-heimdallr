package mode

import (
	"fmt"

	"zmq/utils"

	zmq "github.com/pebbe/zmq4"
)

// Router 服务器
func RunRouter(address string) {
	router := utils.CreateSocket(zmq.ROUTER, address, true)
	defer router.Close()

	for {
		identity, _ := router.Recv(0) // 先接收 Dealer 标识符
		msg := utils.ReceiveMessage(router) // 再接收实际数据
		fmt.Printf("[Router] Received from %s: %s\n", identity, msg)

		// 发送回复
		router.Send(identity, zmq.SNDMORE)
		utils.SendMessage(router, "ACK:"+msg)
	}
}

// Dealer 客户端
func RunDealer(address string) {
	dealer := utils.CreateSocket(zmq.DEALER, address, false)
	defer dealer.Close()

	for i := 0; i < 5; i++ {
		utils.SendMessage(dealer, fmt.Sprintf("Hello %d", i))
		reply := utils.ReceiveMessage(dealer)
		fmt.Println("[Dealer] Received:", reply)
	}
}
