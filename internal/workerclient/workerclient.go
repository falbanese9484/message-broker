package workerclient

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/falbanese9484/message-broker/internal/task"
	pb "github.com/falbanese9484/message-broker/proto/message-broker/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type WorkerBroker struct {
	WorkerClient *WorkerClient
	pb.UnimplementedBrokerServer
}

type WorkerClient struct {
	TCPAddress string
	Task       chan []byte
	Connection net.Conn
}

func (w *WorkerBroker) SendTask(ctx context.Context, in *pb.TaskRequest) (*pb.TaskResponse, error) {
	task := task.NewTask(int(in.Queue), in.Data, nil)
	w.WorkerClient.Task <- task.Data
	return &pb.TaskResponse{Status: 200}, nil
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
	ctx := context.Background()

	workerRequest := &pb.WorkerRequest{
		Queue:         int32(queueId),
		TcpConnection: tcpAddress,
	}
	for i := 0; i < workers; i++ {
		wResp, err := brokerClient.RegisterWorker(ctx, workerRequest)
		if err != nil {
			log.Fatalf("Failed to register worker: %v", err)
			return nil
		}
		log.Printf("Worker registered with status: %d", wResp.Status)
	}
	// wResp, err := brokerClient.RegisterWorker(ctx, workerRequest)
	// if err != nil {
	// 	log.Fatalf("Failed to register worker: %v", err)
	// 	return nil
	// }
	go client.startWorkerListener()

	for i := 0; i < workers; i++ {
		go client.run()
	}
	return client
}

func (w *WorkerClient) run() {
	for task := range w.Task {
		fmt.Println("Received Task: ", string(task))
	}
}

func (w *WorkerClient) startWorkerListener() {
	go w.run()
	listener, err := net.Listen("tcp", w.TCPAddress)
	if err != nil {
		log.Fatalf("Worker listener error: %v", err)
	}
	defer listener.Close()
	log.Printf("Worker listener is running on %s...", w.TCPAddress)

	server := grpc.NewServer()
	workerBroker := &WorkerBroker{WorkerClient: w}
	pb.RegisterBrokerServer(server, workerBroker)
	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

//Example Implementation
//func ReceiveTask(tcpAddress, BrokerAddress string, queueId int) {
//	worker := workerClient.NewWorkerClient(tcpAddress, BrokerAddress, queueId)
//	for task := range worker.Tasks() {
//		fmt.Println("Received Task: ", string(task))
//	}
//}
