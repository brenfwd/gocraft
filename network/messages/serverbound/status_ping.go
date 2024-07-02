package messages

import (
	"github.com/brenfwd/gocraft/constants"
	"github.com/brenfwd/gocraft/ipc"
	"github.com/brenfwd/gocraft/network/messages"
	"github.com/brenfwd/gocraft/network/messages/clientbound"
)

func init() {
	messages.RegisterServerbound[StatusServerboundPing](constants.ClientStateStatus, 0x01)
}

type StatusServerboundPing struct {
	messages.Serverbound
	Challenge int64
}

func (p *StatusServerboundPing) Handle(i *ipc.ClientIPC) error {
	res := clientbound.StatusClientboundPing{
		Response: p.Challenge,
	}
	encoded, err := messages.Encode(&res)
	if err != nil {
		return err
	}
	i.SendPacket(&encoded)
	return nil
}
