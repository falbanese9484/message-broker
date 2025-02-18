package worker

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/falbanese9484/message-broker/internal/queue"
	"github.com/falbanese9484/message-broker/internal/task"
	pb "github.com/falbanese9484/message-broker/proto/message-broker/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Worker struct {
	ID            int
	QueueID       int
	Queue         *queue.Queue
	Done          chan struct{}
	TCPConnection string
	Conn          net.Conn
	Ctx           context.Context
	Cancel        context.CancelFunc
	BrokerConn    pb.BrokerClient
	Busy          bool
}

func NewWorker(queueID int, tcp string) *Worker {

	return &Worker{
		QueueID:       queueID,
		Done:          make(chan struct{}),
		TCPConnection: tcp,
	}
}

func (w *Worker) Run() {
	brokerConn, err := grpc.NewClient(w.TCPConnection, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to broker: %v", err)
		return
	}
	brokerClient := pb.NewBrokerClient(brokerConn)
	ctx := context.Background()

	w.BrokerConn = brokerClient
	w.Ctx = ctx

	log.Println("Worker started, waiting for tasks...")
	for t := range w.Queue.Tasks {
		w.RunTask(t)
	}
	log.Println("Worker task channel closed, exiting Run()")
}

func (w *Worker) RunTask(t *task.Task) {
	w.Busy = true
	t.Lock.Lock()
	t.Status = task.TaskRunning
	t.Lock.Unlock()
	t.Lock.Lock()
	newTask := &pb.TaskRequest{
		Queue: int32(t.QueueID),
		Data:  []byte(t.Data),
	}
	t.Lock.Unlock()
	resp, err := w.BrokerConn.SendTask(w.Ctx, newTask)
	if err != nil {
		t.Lock.Lock()
		t.Status = task.TaskFailed
		t.Error = err
		t.Lock.Unlock()
		return
	}
	if resp == nil {
		t.Lock.Lock()
		t.Status = task.TaskFailed
		t.Lock.Unlock()
		log.Fatalf("Received nil response from broker")
		w.Queue.Tasks <- t
	}
	// time.Sleep(time.Duration(t.Timeout) * time.Second)
	t.Lock.Lock()
	t.Status = task.TaskCompleted
	t.Done = true
	t.Lock.Unlock()
	fmt.Printf("Task %d completed\n", t.Id)
	w.Busy = false
}
