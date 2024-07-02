package messages

import (
	"fmt"

	"github.com/brenfwd/gocraft/constants"
	"github.com/brenfwd/gocraft/data"
	"github.com/brenfwd/gocraft/ipc"
	"github.com/brenfwd/gocraft/network/messages"
)

func init() {
	messages.RegisterServerbound[HandshakingServerboundHandshake](constants.ClientStateHandshaking, 0x00)
}

type HandshakingServerboundHandshake struct {
	messages.Serverbound
	ProtocolVersion data.VarInt
	ServerAddress   string
	ServerPort      uint16
	NextState       data.VarInt
}

func (p *HandshakingServerboundHandshake) Handle(i *ipc.ClientIPC) error {
	nextStateDecode, validState := constants.ClientStateFromInt(int(p.NextState))
	if !validState {
		return fmt.Errorf("invalid next state %v", p.NextState)
	}
	i.ChangeState(nextStateDecode)
	return nil
}
