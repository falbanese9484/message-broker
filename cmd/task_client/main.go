// test_client.go
// Used for testing the Worker service, as well as the client service that will issue tasks and queues to the system.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	pb "github.com/falbanese9484/message-broker/proto/message-broker/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	grpcAddr       = "localhost:6900"
	workerTCPAddr  = "localhost:9000"
	taskInterval   = 2 * time.Second
	requestTimeout = 5 * time.Second
)

func startWorkerListener() {
	listener, err := net.Listen("tcp", workerTCPAddr)
	if err != nil {
		log.Fatalf("Worker listener error: %v", err)
	}
	defer listener.Close()
	log.Printf("Worker listener is running on %s...", workerTCPAddr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		go handleWorkerConnection(conn)
	}
}

func handleWorkerConnection(conn net.Conn) {
	defer conn.Close()
	log.Printf("New worker connection from %s", conn.RemoteAddr())
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		task := buf[:n]
		json.Unmarshal(task, &task)
		if err != nil {
			if err != io.EOF {
				log.Printf("Error reading from worker connection: %v", err)
			}
			return
		}
		if n > 0 {
			// Print the task received by the worker.
			fmt.Printf("Worker received task: %s\n", task)
		}
	}
}

func main() {

	// Give the worker listener a moment to start.
	time.Sleep(1 * time.Second)

	// Connect to the gRPC broker.
	conn, err := grpc.NewClient(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to broker: %v", err)
	}
	defer conn.Close()

	client := pb.NewBrokerClient(conn)

	// Use a common context for our initial requests.
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	// Register a queue.
	queueReq := &pb.QueueRequest{
		Name:         "test-queue",
		BufferLength: 100,
	}
	qResp, err := client.RegisterQueue(ctx, queueReq)
	if err != nil {
		log.Fatalf("failed to register queue: %v", err)
	}
	log.Printf("Queue registered, status: %d", qResp.Status)

	// Continuously send tasks so you can see the worker process them.
	ticker := time.NewTicker(taskInterval)
	defer ticker.Stop()
	taskCount := 1

	log.Println("Starting task submission...")
	for range ticker.C {
		taskMsg := fmt.Sprintf(`{"message": "Task number %d"}`, taskCount)
		taskData := []byte(taskMsg)
		taskReq := &pb.TaskRequest{
			Queue: 1, // send tasks to queue ID 1
			Data:  taskData,
		}

		// Create a new context for each task submission.
		ctxTask, cancelTask := context.WithTimeout(context.Background(), requestTimeout)
		tResp, err := client.SendTask(ctxTask, taskReq)
		cancelTask()

		if err != nil {
			log.Printf("failed to send task %d: %v", taskCount, err)
		} else {
			log.Printf("Task %d sent, status: %d", taskCount, tResp.Status)
		}
		taskCount++
	}
}
