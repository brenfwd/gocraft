package serverbound

import (
	"github.com/brenfwd/gocraft/constants"
	"github.com/brenfwd/gocraft/network/messages"
	"github.com/brenfwd/gocraft/network/messages/clientbound"
	"github.com/brenfwd/gocraft/shared"
)

func init() {
	messages.RegisterServerbound[StatusServerboundPing](constants.ClientStateStatus, 0x01)
}

type StatusServerboundPing struct {
	messages.Serverbound
	Challenge int64
}

func (p *StatusServerboundPing) Handle(c *shared.ClientShared) error {
	res := clientbound.StatusClientboundPing{
		Response: p.Challenge,
	}
	encoded, err := messages.Encode(&res)
	if err != nil {
		return err
	}
	c.SendPacket(&encoded)
	return nil
}
