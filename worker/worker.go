package worker

import (
	"fmt"
	"message-broker/queue"
	"message-broker/task"
	"time"
)

type Worker struct {
	ID         int
	Queue      *queue.Queue
	Done       chan struct{}
	Busy       bool
	BufferFull bool
	TaskChan   chan *task.Task
}

func NewWorker(id int, queue *queue.Queue) *Worker {

	return &Worker{
		ID:       id,
		Queue:    queue,
		Done:     make(chan struct{}),
		TaskChan: make(chan *task.Task, 10),
	}
}

func (w *Worker) ReadPump() {
	for {
		task := w.CheckQueue()
		if task != nil {
			w.RunTask(task)
		}
		select {
		case task := <-w.TaskChan:
			w.BufferFull = false
			w.RunTask(task)
		case <-w.Done:
			fmt.Printf("Worker %d Shutting Down\n", w.ID)
		}
	}
}

func (w *Worker) CheckQueue() *task.Task {
	if len(w.Queue.Tasks) > 1 {
		t := w.Queue.Tasks[0]
		w.Queue.Tasks = w.Queue.Tasks[1:]
		return t
	}
	return nil
}

func (w *Worker) RunTask(t *task.Task) {
	w.Busy = true
	t.Lock.Lock()
	t.Status = task.TaskRunning
	t.Lock.Unlock()
	fmt.Printf("Worker %d is Processing Task %d\n", w.ID, t.Id)
	time.Sleep(10 * time.Second)
	t.Lock.Lock()
	t.Status = task.TaskCompleted
	t.Done = true
	t.Lock.Unlock()
	fmt.Printf("Worker %d completed Task %d\n", w.ID, t.Id)
	w.Busy = false
}
