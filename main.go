package main

import (
	"github.com/brenfwd/gocraft/core"
)

func unwrap(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	server, err := core.NewServer()
	unwrap(err)

	defer func() {
		unwrap(server.Close())
	}()

	server.Run()
}
