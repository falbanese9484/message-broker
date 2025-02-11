package task

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

func NewTask(id int, queueID int, key string, value []byte, opts *TaskOptions) *Task {
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
