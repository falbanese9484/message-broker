package queue

import (
	"sync"

	"github.com/falbanese9484/message-broker/task"
)

type Queue struct {
	Id    int
	Tasks []*task.Task
	Lock  sync.RWMutex
}
