package main

import (
	"message-broker/queue"
	"message-broker/system"
	"message-broker/task"
	"message-broker/worker"
	"sync"
	"time"
)

func main() {
	s := system.NewSystem()
	go s.Run()
	q := &queue.Queue{
		Id:    1,
		Tasks: make([]*task.Task, 0),
		Lock:  sync.RWMutex{},
	}
	s.Register <- q
	w := worker.NewWorker(1, q)
	w2 := worker.NewWorker(2, q)
	go w2.ReadPump()
	go w.ReadPump()
	s.WorkerRegister <- w
	go TaskCreator(s)
	select {}
}

func TaskCreator(s *system.System) {
	id := 1
	for {
		time.Sleep(200 * time.Millisecond)
		t := task.NewTask(id, 1, "key", []byte("value"), nil)
		s.Router <- t
		id++
	}
}

// func main() {
// 	s := server.Server{
// 		Port: 8080,
// 	}

// 	server.RunServer(s)
// }
