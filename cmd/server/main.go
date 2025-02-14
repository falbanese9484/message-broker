package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	pb "github.com/falbanese9484/message-broker/proto"

	"google.golang.org/grpc"
)

type brokerServer struct {
	pb.UnimplementedBrokerServer
}

func (s *brokerServer) RegisterWorker(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	log.Printf("Registering %d worker(s) for queue %d", req.WorkerCount, req.QueueId)

	return &pb.RegisterResponse{
		Success: true,
		Message: fmt.Sprintf("Registered %d workers for queue %d", req.WorkerCount, req.QueueId),
	}, nil
}

func (s *brokerServer) GetTask(ctx context.Context, req *pb.TaskRequest) (*pb.TaskRequest, error) {
	log.Printf("Received task request for queue %d", req.QueueId)

	task := &pb.Task{
		Id:          time.Now().Unix(),
		Description: "Example Task",
	}
	return &pb.TaskResponse{Task: []*pb.Task{task}}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterBrokerServer(grpcServer, &brokerServer{})

	log.Println("gRPC server listening on :9000")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
