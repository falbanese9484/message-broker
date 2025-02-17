package task

import (
	"sync"
	"time"
)

type Task struct {
	Id        int
	QueueID   int
	Lock      sync.RWMutex
	Data      []byte
	Done      bool
	Priority  int
	Timeout   int
	Status    TaskStatus
	CreatedAt time.Time
	Error     error
}

func NewTask(queueID int, data []byte, opts *TaskOptions) *Task {
	var priority, timeout int

	if opts == nil {
		priority = 0
		timeout = 10
	} else {
		priority = opts.Priority
		timeout = opts.Timeout
	}

	return &Task{
		QueueID:   queueID,
		Data:      data,
		Done:      false,
		Priority:  priority,
		Timeout:   timeout,
		Status:    TaskPending,
		CreatedAt: time.Now(),
		Error:     nil,
	}
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
