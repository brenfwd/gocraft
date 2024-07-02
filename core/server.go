package core

import (
	"log"
	"sync"

	"github.com/brenfwd/gocraft/network"
)

type Server struct {
	listener network.Listener
	clients  []*Client
}

func NewServer() (Server, error) {
	listener, err := network.NewListener("0.0.0.0", 25565)
	if err != nil {
		return Server{}, err
	}

	return Server{listener, make([]*Client, 0)}, nil
}

func (s *Server) Close() error {
	log.Println("gocraft server is shutting down...")
	if err := s.listener.Close(); err != nil {
		return err
	}
	// s.wg.Wait()
	return nil
}

func (s *Server) Run() {
	var wg sync.WaitGroup
	defer wg.Wait()

	log.Println("gocraft server is starting...")

	// Start listener
	wg.Add(1)
	go func() {
		defer wg.Done()
		s.listener.Listen()
	}()

	log.Println("Server is ready")

	for conn := range s.listener.Incoming {
		log.Println("Got connection:", conn.RemoteAddr())

		client := NewClient(conn)
		s.clients = append(s.clients, &client)

		wg.Add(1)
		go func() {
			defer wg.Done()
			client.Handle()
		}()
	}
}
