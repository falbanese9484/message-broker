package queue

import (
	"message-broker/task"
	"sync"
)

type Queue struct {
	Id    int
	Tasks []*task.Task
	Lock  sync.RWMutex
}
