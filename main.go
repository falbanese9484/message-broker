package main

import (
	"fmt"
	"sync"
	"time"
)

type Task struct {
	Id        int
	QueueID   int
	Lock      sync.RWMutex
	KV        map[string][]byte
	Done      bool
	Priority  int
	Timeout   int
	Status    TaskStatus
	CreatedAt time.Time
	Error     error
}

func newTask(id int, queueID int, key string, value []byte, opts *TaskOptions) *Task {
	var priority int
	var timeout int

	if opts == nil {
		priority = 0
		timeout = 10
	} else {
		priority = opts.Priority
		timeout = opts.Timeout
	}

	t := &Task{
		Id:        id,
		QueueID:   queueID,
		Done:      false,
		Priority:  priority,
		Timeout:   timeout,
		Status:    TaskPending,
		CreatedAt: time.Now(),
		Error:     nil,
		KV:        make(map[string][]byte),
	}

	t.KV[key] = value
	return t
}

type TaskOptions struct {
	Priority int
	Timeout  int
}

type TaskStatus int

const (
	TaskPending TaskStatus = iota
	TaskRunning
	TaskCompleted
	TaskFailed
)

type Queue struct {
	Id    int
	Tasks []*Task
	Lock  sync.RWMutex
}

type System struct {
	Queues           []*Queue
	Lock             sync.RWMutex
	Register         chan *Queue
	Unregister       chan *Queue
	WorkerRegister   chan *Worker
	WorkerUnregister chan *Worker
	Router           chan *Task
	Done             chan struct{}
}

func newSystem() *System {
	return &System{
		Queues:           make([]*Queue, 0),
		Lock:             sync.RWMutex{},
		Register:         make(chan *Queue),
		Unregister:       make(chan *Queue),
		WorkerRegister:   make(chan *Worker),
		WorkerUnregister: make(chan *Worker),
		Router:           make(chan *Task),
		Done:             make(chan struct{}),
	}
}

func (s *System) run() {
	var workers = make([]*Worker, 0)
	for {
		select {
		case queue := <-s.Register:
			s.Lock.Lock()
			s.Queues = append(s.Queues, queue)
			s.Lock.Unlock()
		case queue := <-s.Unregister:
			s.Lock.Lock()
			for i, q := range s.Queues {
				if q == queue {
					s.Queues = append(s.Queues[:i], s.Queues[i+1:]...)
					break
				}
			}
			s.Lock.Unlock()
		case task := <-s.Router:
			assigned := false
			for _, w := range workers {
				if w.Queue.Id == task.QueueID {
					if len(w.TaskChan) < cap(w.TaskChan) {
						w.TaskChan <- task
						assigned = true
						break
					}
				}
			}
			if !assigned {
				fmt.Printf("No Worker Available for Task %d\n", task.Id)
				// Route these to a different queue
			}
		case worker := <-s.WorkerRegister:
			s.Lock.Lock()
			workers = append(workers, worker)
			s.Lock.Unlock()
		case worker := <-s.WorkerUnregister:
			s.Lock.Lock()
			for i, w := range workers {
				if w == worker {
					workers = append(workers[:i], workers[i+1:]...)
					s.Lock.Unlock()
					break
				}
			}
		case <-s.Done:
			return
		}
	}
}

type Worker struct {
	ID       int
	Queue    *Queue
	Done     chan struct{}
	Busy     bool
	TaskChan chan *Task
}

func newWorker(id int, queue *Queue) *Worker {

	return &Worker{
		ID:       id,
		Queue:    queue,
		Done:     make(chan struct{}),
		TaskChan: make(chan *Task, 10),
	}
}

func (w *Worker) readPump() {
	for {
		select {
		case task := <-w.TaskChan:
			w.Busy = true
			task.Lock.Lock()
			task.Status = TaskRunning
			task.Lock.Unlock()
			fmt.Printf("Worker %d is Processing Task %d\n", w.ID, task.Id)
			time.Sleep(10 * time.Second)
			task.Lock.Lock()
			task.Status = TaskCompleted
			task.Done = true
			task.Lock.Unlock()
			fmt.Printf("Worker %d completed Task %d\n", w.ID, task.Id)
			w.Busy = false
		case <-w.Done:
			fmt.Printf("Worker %d Shutting Down\n", w.ID)
		}
	}
}

func main() {
	s := newSystem()
	go s.run()
	q := &Queue{
		Id:    1,
		Tasks: make([]*Task, 0),
		Lock:  sync.RWMutex{},
	}
	s.Register <- q
	w := newWorker(1, q)
	go w.readPump()
	s.WorkerRegister <- w
	go TaskCreator(s)
	select {}
}

func TaskCreator(s *System) {
	id := 1
	for {
		time.Sleep(2 * time.Second)
		t := newTask(id, 1, "key", []byte("value"), nil)
		s.Router <- t
		id++
	}
}
