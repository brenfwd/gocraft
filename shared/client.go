package shared

import (
	"crypto/rand"
	"sync"

	"github.com/brenfwd/gocraft/constants"
	"github.com/brenfwd/gocraft/network"
	"github.com/brenfwd/gocraft/network/encryption"
	"github.com/google/uuid"
)

type ClientMessage interface{}
type ClientShared struct {
	Mutex                 sync.Mutex
	C                     chan *ClientMessage
	ListenerKeypair       *encryption.KeypairBytes
	EncryptionVerifyToken [4]byte
	AllegedUsername       string
	AllegedUUID           uuid.UUID
	SharedSecret          []byte
}

type ClientChangeState struct {
	NewState constants.ClientState
}

func (i *ClientShared) ChangeState(newState constants.ClientState) {
	cm := ClientMessage(ClientChangeState{NewState: newState})
	i.C <- &cm
}

type ClientSend struct {
	Packet *network.Packet
}

func (i *ClientShared) SendPacket(packet *network.Packet) {
	cm := ClientMessage(ClientSend{Packet: packet})
	i.C <- &cm
}

type ClientEnableEncryption struct{}

func (i *ClientShared) EnableEncryption() {
	cm := ClientMessage(ClientEnableEncryption{})
	i.C <- &cm
}

const maxClientMessages = 1024

func NewClientShared(keypair *encryption.KeypairBytes) *ClientShared {
	// All channels have to be buffered because the channel is sent data during a select statement
	// so it must be buffered to prevent blocking since nothing will read from it until the select
	// statement is re-run.

	cs := ClientShared{
		C:               make(chan *ClientMessage, maxClientMessages),
		ListenerKeypair: keypair,
	}
	rand.Read(cs.EncryptionVerifyToken[:])

	return &cs
}
