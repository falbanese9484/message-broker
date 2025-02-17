package main

import (
	"fmt"

	"github.com/falbanese9484/message-broker/internal/workerclient"
)

func startLoopingWorker(tcpAddress, grpcAddress string, queueId int, concurrency int) {
	w := workerclient.NewWorkerClient(tcpAddress, grpcAddress, queueId, concurrency)
	if w == nil {
		return
	}
	for i := 1; i < concurrency; i++ {
		go func() {
			for task := range w.Tasks() {
				fmt.Println("Received Task: ", string(task))
			}
		}()
	}
}

func main() {
	go startLoopingWorker("localhost:6901", "localhost:6900", 1, 4)
	select {}
}
