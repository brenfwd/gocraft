package core

import (
	"log"
	"sync"

	"github.com/brenfwd/gocraft/constants"
	"github.com/brenfwd/gocraft/network"
	"github.com/brenfwd/gocraft/network/messages"
	"github.com/brenfwd/gocraft/shared"
)

type Client struct {
	Shared     *shared.ClientShared
	State      constants.ClientState
	connection network.Connection
}

func NewClient(connection network.Connection) Client {
	return Client{
		Shared:     shared.NewClientShared(connection.Keypair),
		State:      constants.ClientStateHandshaking,
		connection: connection,
	}
}

func (c *Client) processPacket(packet *network.Packet) error {
	log.Printf("Packet from %s (%v): %+v", c.connection.RemoteAddr(), c.State, *packet)

	decoded, err := messages.DecodeServerbound(c.State, packet)
	if err != nil {
		return err
	}

	if err := decoded.Handle(c.Shared); err != nil {
		return err
	}

	return nil
}

func (c *Client) handleSharedMessage(msg *shared.ClientMessage) error {
	switch inner := (*msg).(type) {
	case shared.ClientChangeState:
		log.Printf("Changing state to %v", inner.NewState)
		c.State = inner.NewState
	case shared.ClientSend:
		log.Printf("Sending packet with ID 0x%02x (%d)", inner.Packet.Id, inner.Packet.Id)
		log.Println(inner.Packet.Body)
		if err := c.connection.WritePacket(inner.Packet); err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) Handle() {
	var wg sync.WaitGroup
	defer wg.Wait()

	wg.Add(1)
	go func() {
		defer wg.Done()
		c.connection.Receive()
	}()

	for {
		// Process pending IPC messages first
		for more_ipc := true; more_ipc; {
			select {
			case msg := <-c.Shared.C:
				if err := c.handleSharedMessage(msg); err != nil {
					log.Println("Error handling shared message:", err)
					goto end
				}
			default:
				more_ipc = false
			}
		}

		// Then process network packets & events and later IPC messages
		select {
		case <-c.connection.Eof:
			log.Println("EOF", c.connection.RemoteAddr())
			goto end
		case packet := <-c.connection.Packets:
			err := c.processPacket(&packet)
			if err != nil {
				log.Println("Error processing packet:", err)
				goto end
			}
		case msg := <-c.Shared.C:
			// in this case, we should handle this message immediately
			// but then continue to the next iteration of the outer loop
			// to handle any further IPC messages
			if err := c.handleSharedMessage(msg); err != nil {
				log.Println("Error handling shared message:", err)
				goto end
			}
		}
	}

end:
	c.connection.Close()
}
