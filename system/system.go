package system

import (
	"fmt"
	"sync"

	"github.com/falbanese9484/message-broker/queue"
	"github.com/falbanese9484/message-broker/task"
	"github.com/falbanese9484/message-broker/worker"
)

type System struct {
	Queues           []*queue.Queue
	Lock             sync.RWMutex
	Register         chan *queue.Queue
	Unregister       chan *queue.Queue
	WorkerRegister   chan *worker.Worker
	WorkerUnregister chan *worker.Worker
	Router           chan *task.Task
	Done             chan struct{}
}

func NewSystem() *System {
	return &System{
		Queues:           make([]*queue.Queue, 0),
		Lock:             sync.RWMutex{},
		Register:         make(chan *queue.Queue),
		Unregister:       make(chan *queue.Queue),
		WorkerRegister:   make(chan *worker.Worker),
		WorkerUnregister: make(chan *worker.Worker),
		Router:           make(chan *task.Task),
		Done:             make(chan struct{}),
	}
}

func (s *System) Run() {
	var workers = make([]*worker.Worker, 0)
	for {
		select {
		case queue := <-s.Register:
			s.Lock.Lock()
			queue.Id = len(s.Queues) + 1
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
				if task.QueueID == q.Id {
					q.Lock.Lock()
					task.Id = q.TaskCounter + 1
					q.TaskCounter++
					q.Tasks <- task
					q.Lock.Unlock()
					fmt.Printf("Added %d Task to Queue %d\n", task.Id, q.Id)
				}
			}
		case worker := <-s.WorkerRegister:
			s.Lock.Lock()
			for _, q := range s.Queues {
				if worker.QueueID == q.Id {
					worker.Queue = q
					break
				}
			}
			workers = append(workers, worker)
			go worker.Run()
			s.Lock.Unlock()
		case worker := <-s.WorkerUnregister:
			s.Lock.Lock()
			worker.ID = workers[len(workers)-1].ID + 1
			workers = append(workers, worker)
			s.Lock.Unlock()
		case <-s.Done:
			return
		}
	}
}
