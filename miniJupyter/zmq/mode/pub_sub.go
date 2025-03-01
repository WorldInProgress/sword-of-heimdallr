package mode

import (
	"fmt"
	"time"

	"zmq/utils"

	zmq "github.com/pebbe/zmq4"
)

// Publisher
func RunPublisher(address string) {
	pub := utils.CreateSocket(zmq.PUB, address, true)
	defer pub.Close()

	for i := 0; i < 10; i++ {
		utils.SendMessage(pub, fmt.Sprintf("Topic1 Message %d", i))
		time.Sleep(time.Second)
	}
}

// Subscriber
func RunSubscriber(address string) {
	sub := utils.CreateSocket(zmq.SUB, address, false)
	defer sub.Close()
	sub.SetSubscribe("Topic1") // 订阅 "Topic1"

	for {
		msg := utils.ReceiveMessage(sub)
		fmt.Println("[Subscriber] Received:", msg)
	}
}
