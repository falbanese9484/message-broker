package worker

import (
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/falbanese9484/message-broker/internal/queue"
	"github.com/falbanese9484/message-broker/internal/task"
)

type Worker struct {
	ID            int
	QueueID       int
	Queue         *queue.Queue
	Done          chan struct{}
	TCPConnection string
	Conn          net.Conn
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
	conn, err := net.Dial("tcp", w.TCPConnection)
	if err != nil {
		return
	}
	defer conn.Close()
	w.Conn = conn

	for task := range w.Queue.Tasks {
		w.RunTask(task)
	}
}

func (w *Worker) RunTask(t *task.Task) {
	w.Busy = true
	t.Lock.Lock()
	t.Status = task.TaskRunning
	t.Lock.Unlock()
	marshalled, err := json.Marshal(t.Data)
	if err != nil {
		t.Lock.Lock()
		t.Status = task.TaskFailed
		t.Error = err
		t.Lock.Unlock()
		return
	}
	_, err = w.Conn.Write(marshalled)
	if err != nil {
		t.Lock.Lock()
		t.Status = task.TaskFailed
		t.Error = err
		t.Lock.Unlock()
		return
	}
	time.Sleep(time.Duration(t.Timeout) * time.Second)
	t.Lock.Lock()
	t.Status = task.TaskCompleted
	t.Done = true
	t.Lock.Unlock()
	fmt.Printf("Task %d completed\n", t.Id)
	w.Busy = false
}
