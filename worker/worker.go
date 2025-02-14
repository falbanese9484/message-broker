package worker

import (
	"fmt"
	"sort"
	"time"

	"github.com/falbanese9484/message-broker/queue"
	"github.com/falbanese9484/message-broker/task"
)

type Worker struct {
	ID    int
	Queue *queue.Queue
	Done  chan struct{}
	Busy  bool
}

func NewWorker(id int, queue *queue.Queue) *Worker {

	return &Worker{
		ID:    id,
		Queue: queue,
		Done:  make(chan struct{}),
	}
}

func (w *Worker) ReadPump() {
	for {
		task := w.CheckQueue()
		if task != nil {
			w.RunTask(task)
		}
	}
}

func (w *Worker) CheckQueue() *task.Task {
	w.Queue.Lock.Lock()
	defer w.Queue.Lock.Unlock()
	if len(w.Queue.Tasks) == 0 {
		return nil
	}
	if len(w.Queue.Tasks) == 1 {
		t := w.Queue.Tasks[0]
		w.Queue.Tasks = w.Queue.Tasks[:0]
		return t
	}
	sort.Slice(w.Queue.Tasks, func(i, j int) bool {
		return w.Queue.Tasks[i].Id < w.Queue.Tasks[j].Id
	})
	t := w.Queue.Tasks[0]
	w.Queue.Tasks = w.Queue.Tasks[1:]
	return t
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
