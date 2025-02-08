package main

import (
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
	Queues     []*Queue
	Lock       sync.RWMutex
	Register   chan *Queue
	Unregister chan *Queue
	Router     chan *Task
	Done       chan struct{}
}

func newSystem() *System {
	return &System{
		Queues:     make([]*Queue, 0),
		Lock:       sync.RWMutex{},
		Register:   make(chan *Queue),
		Unregister: make(chan *Queue),
		Router:     make(chan *Task),
		Done:       make(chan struct{}),
	}
}

func (s *System) run() {
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
			for _, q := range s.Queues {
				if task.Id == q.Id {
					q.Lock.Lock()
					q.Tasks = append(q.Tasks, task)
					q.Lock.Unlock()
					break
				}
			}
		default:
			return
		}
	}
}

type Worker struct {
	ID       int
	Queue    *Queue
	Done     chan struct{}
	TaskChan chan Task
}
