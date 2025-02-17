package workerclient

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net"
	"time"

	pb "github.com/falbanese9484/message-broker/proto/message-broker/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type WorkerClient struct {
	TCPAddress string
	Task       chan []byte
	Connection net.Conn
}

func NewWorkerClient(tcpAddress string, brokerAddress string, queueId int, workers int) *WorkerClient {
	client := &WorkerClient{
		TCPAddress: tcpAddress,
		Task:       make(chan []byte),
	}
	brokerConn, err := grpc.NewClient(brokerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to broker: %v", err)
		return nil
	}
	brokerClient := pb.NewBrokerClient(brokerConn)
	requestTimeout := 5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	workerRequest := &pb.WorkerRequest{
		Queue:         int32(queueId),
		TcpConnection: tcpAddress,
	}
	// Attempt at concurrent listener registration <-- Did not work
	for i := 1; i < workers; i++ {
		wResp, err := brokerClient.RegisterWorker(ctx, workerRequest)
		if err != nil {
			log.Fatalf("Failed to register worker: %v", err)
			return nil
		}
		log.Printf("Worker registered with status: %d", wResp.Status)
	}
	wResp, err := brokerClient.RegisterWorker(ctx, workerRequest)
	if err != nil {
		log.Fatalf("Failed to register worker: %v", err)
		return nil
	}
	log.Printf("Worker registered with status: %d", wResp.Status)
	go client.startWorkerListener()
	return client
}

func (w *WorkerClient) startWorkerListener() {
	listener, err := net.Listen("tcp", w.TCPAddress)
	if err != nil {
		log.Fatalf("Worker listener error: %v", err)
	}

	defer listener.Close()

	log.Printf("Worker listener is running on %s...", w.TCPAddress)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}

		go w.handleWorkerConnection(conn)
	}
}

func (w *WorkerClient) handleWorkerConnection(conn net.Conn) {
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
			// fmt.Printf("Worker received task: %s\n", task)
			w.Task <- task
		}
	}
}

func (w *WorkerClient) Tasks() <-chan []byte {
	return w.Task
}

//Example Implementation
//func ReceiveTask(tcpAddress, BrokerAddress string, queueId int) {
//	worker := workerClient.NewWorkerClient(tcpAddress, BrokerAddress, queueId)
//	for task := range worker.Tasks() {
//		fmt.Println("Received Task: ", string(task))
//	}
//}
