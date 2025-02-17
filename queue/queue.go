package queue

import (
	"sync"

	"github.com/falbanese9484/message-broker/task"
)

type Queue struct {
	Id          int             `json:"id"`
	Namespace   string          `json:"namespace"`
	Tasks       chan *task.Task `json:"-"`
	TaskCounter int             `json:"task_counter"`
	Lock        sync.RWMutex    `json:"-"`
}

func NewQueue(namespace string, bufferLength int) *Queue {
	return &Queue{
		Namespace:   namespace,
		Tasks:       make(chan *task.Task, bufferLength),
		TaskCounter: 0,
	}
}
