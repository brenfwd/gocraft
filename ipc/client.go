package ipc

import (
	"github.com/brenfwd/gocraft/constants"
	"github.com/brenfwd/gocraft/network"
)

type ClientMessage interface{}
type ClientIPC struct {
	C chan *ClientMessage
}

type ClientChangeState struct {
	NewState constants.ClientState
}

func (i *ClientIPC) ChangeState(newState constants.ClientState) {
	cm := ClientMessage(ClientChangeState{NewState: newState})
	i.C <- &cm
}

type ClientSend struct {
	Packet *network.Packet
}

func (i *ClientIPC) SendPacket(packet *network.Packet) {
	cm := ClientMessage(ClientSend{Packet: packet})
	i.C <- &cm
}

const maxClientMessages = 1024

func NewClient() ClientIPC {
	// All channels have to be buffered because the channel is sent data during a select statement
	// so it must be buffered to prevent blocking since nothing will read from it until the select
	// statement is re-run.
	return ClientIPC{
		C: make(chan *ClientMessage, maxClientMessages),
	}
}
