package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/falbanese9484/message-broker/internal/workerclient"
)

func main() {
	worker := workerclient.NewWorkerClient("localhost:6901", "localhost:6900", 1, 4)
	if worker == nil {
		log.Fatal("Failed to create worker client")
	}
	log.Println("Worker client started successfully.")

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	<-sigchan
	log.Println("Shutting down gracefully...")

}
