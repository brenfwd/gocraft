package network

import (
	"errors"
	"fmt"
	"log"
	"net"
)

type Listener struct {
	inner         net.Listener
	incoming_send chan<- Connection
	Incoming      <-chan Connection
}

func NewListener(host string, port uint16) (Listener, error) {
	netListener, err := net.Listen("tcp4", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return Listener{}, err
	}

	c := make(chan Connection, 256) // backlog

	return Listener{
		inner:         netListener,
		incoming_send: c,
		Incoming:      c,
	}, nil
}

func (l *Listener) Close() error {
	return l.inner.Close()
}

func (l *Listener) Listen() {
	for {
		netConn, err := l.inner.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return
			}
			log.Println("Error during Listener.Listen Accept call:", err)
		}
		conn := NewConnection(netConn)
		l.incoming_send <- conn
	}
}
