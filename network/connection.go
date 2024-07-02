package network

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/brenfwd/gocraft/network/encryption"
)

type Connection struct {
	inner        net.Conn
	eofSend      chan<- bool
	Eof          <-chan bool
	packetsSend  chan<- Packet
	Packets      <-chan Packet
	unmarshaller PacketUnmarshaller
	Keypair      *encryption.KeypairBytes
}

func MakeConnection(inner net.Conn, keypair *encryption.KeypairBytes) Connection {
	eof := make(chan bool, 10)
	packets := make(chan Packet, 10)
	return Connection{inner: inner,
		eofSend:     eof,
		Eof:         eof,
		packetsSend: packets,
		Packets:     packets,
		Keypair:     keypair,
	}
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.inner.RemoteAddr()
}

func (c *Connection) Close() error {
	err := c.inner.Close()
	c.eofSend <- true
	return err
}

func (c *Connection) WriteBytes(bytes []byte) error {
	_, err := c.inner.Write(bytes)
	return err
}

func (c *Connection) WritePacket(packet *Packet) error {
	bytes, err := packet.Marshal()
	if err != nil {
		return err
	}
	return c.WriteBytes(bytes)
}

func (c *Connection) Receive() {
	reader := bufio.NewReader(c.inner)
	buf := make([]byte, 4096)

	for {
		n, err := reader.Read(buf)
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				log.Println("Connection closed by remote host", c.inner.RemoteAddr())
				c.Close()
			}
			if !errors.Is(err, io.EOF) {
				log.Println("reader.Read error:", err, c.inner.RemoteAddr())
			}
			// c.eofSend <- true
			c.Close()
			return
		}
		if n != 0 {
			// c.packetsSend <- Packet{Data: buf[0:n]}
			packets, err := c.unmarshaller.Unmarshal(buf[0:n])
			if err != nil {
				fmt.Println("Error during unmarshal:", err)
				// c.eofSend <- true
				c.Close()
				return
			}
			for _, packet := range packets {
				c.packetsSend <- packet
			}
		}
	}
}
