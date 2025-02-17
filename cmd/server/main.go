package main

import (
	"log"
	"net"

	"github.com/falbanese9484/message-broker/internal/broker"
	"github.com/falbanese9484/message-broker/internal/system"
	pb "github.com/falbanese9484/message-broker/proto/message-broker/proto"
	"google.golang.org/grpc"
)

func main() {
	s := system.NewSystem()
	go s.Run()
	lis, err := net.Listen("tcp", ":6900")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	broker := &broker.Broker{System: s}
	pb.RegisterBrokerServer(server, broker)
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
