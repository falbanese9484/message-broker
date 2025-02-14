package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

type Server struct {
	Port int
}

/*
This is placeholder server that will eventually host the protocol layer
of the system
*/

func RunServer(s Server) {
	addr := fmt.Sprintf(":%d", s.Port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Error starting TCP server: %v", err)
		return
	}
	defer ln.Close()
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v\n", err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	log.Printf("Client connected: %s\n", conn.RemoteAddr().String())

	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Connection closed by client: %s\n", conn.RemoteAddr().String())
			return
		}

		if message == "herro" {
			fmt.Println("You win!")
		}

		fmt.Printf("Connection closed by client: %s\n", conn.RemoteAddr().String())
		if message == "herro" {
			_, writeErr := conn.Write([]byte("Echo: " + message))
			if writeErr != nil {
				log.Printf("Error writing to client: %v\n", writeErr)
				return
			}
		}

	}
}
