package main

import (
	"message-broker/queue"
	"message-broker/task"
	"message-broker/worker"
	"sync"
	"time"
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

func newSystem() *System {
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

func (s *System) run() {
	var workers = make([]*worker.Worker, 0)
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
			for _, w := range workers {
				if w.Queue.Id == task.QueueID {
					if len(w.TaskChan) < cap(w.TaskChan) {
						if len(w.TaskChan) == 9 {
							w.BufferFull = true
						}
						w.TaskChan <- task
						break
					}
					for _, q := range s.Queues {
						if w.Queue == q {
							q.Tasks = append(q.Tasks, task)
						}
					}
				}
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

func main() {
	s := newSystem()
	go s.run()
	q := &queue.Queue{
		Id:    1,
		Tasks: make([]*task.Task, 0),
		Lock:  sync.RWMutex{},
	}
	s.Register <- q
	w := worker.NewWorker(1, q)
	go w.ReadPump()
	s.WorkerRegister <- w
	go TaskCreator(s)
	select {}
}

func TaskCreator(s *System) {
	id := 1
	for {
		time.Sleep(2 * time.Second)
		t := task.NewTask(id, 1, "key", []byte("value"), nil)
		s.Router <- t
		id++
	}
}
