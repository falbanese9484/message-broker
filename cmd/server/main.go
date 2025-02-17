package main

import (
	"context"
	"log"
	"net"

	pb "github.com/falbanese9484/message-broker/proto/message-broker/proto"
	"github.com/falbanese9484/message-broker/queue"
	"github.com/falbanese9484/message-broker/system"
	"github.com/falbanese9484/message-broker/task"
	"github.com/falbanese9484/message-broker/worker"
	"google.golang.org/grpc"
)

type Broker struct {
	pb.UnimplementedBrokerServer
	System *system.System
}

func (b *Broker) RegisterWorker(ctx context.Context, in *pb.WorkerRequest) (*pb.WorkerResponse, error) {
	worker := worker.NewWorker(int(in.Queue), in.TcpConnection)
	b.System.WorkerRegister <- worker
	return &pb.WorkerResponse{Status: 200}, nil
}

func (b *Broker) RegisterQueue(ctx context.Context, in *pb.QueueRequest) (*pb.QueueResponse, error) {
	queue := queue.NewQueue(in.Name, int(in.BufferLength))
	b.System.Register <- queue
	return &pb.QueueResponse{Status: 200}, nil
}

func (b *Broker) SendTask(ctx context.Context, in *pb.TaskRequest) (*pb.TaskResponse, error) {
	task := task.NewTask(int(in.Queue), in.Data, nil)
	b.System.Router <- task
	return &pb.TaskResponse{Status: 200}, nil
}

func main() {
	s := system.NewSystem()
	go s.Run()
	lis, err := net.Listen("tcp", ":6900")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	broker := &Broker{System: s}
	pb.RegisterBrokerServer(server, broker)
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
