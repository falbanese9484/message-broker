package broker

import (
	"context"

	"github.com/falbanese9484/message-broker/internal/queue"
	"github.com/falbanese9484/message-broker/internal/system"
	"github.com/falbanese9484/message-broker/internal/task"
	"github.com/falbanese9484/message-broker/internal/worker"
	pb "github.com/falbanese9484/message-broker/proto/message-broker/proto"
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
