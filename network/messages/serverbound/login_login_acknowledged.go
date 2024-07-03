package serverbound

import (
	"github.com/brenfwd/gocraft/constants"
	"github.com/brenfwd/gocraft/data"
	"github.com/brenfwd/gocraft/network/messages"
	"github.com/brenfwd/gocraft/network/messages/clientbound"
	"github.com/brenfwd/gocraft/shared"
)

func init() {
	messages.RegisterServerbound[LoginAcknowledged](constants.ClientStateLogin, 0x03)
}

type LoginAcknowledged struct {
	messages.Serverbound
}

func (p *LoginAcknowledged) Handle(c *shared.ClientShared) error {
	// Switch to configuration state
	c.ChangeState(constants.ClientStateConfiguration)

	// TODO: Send registry data...

	res := clientbound.ConfigurationDisconnect{
		// Reason: data.NBTCompoundValue(nil, []*data.NBTValue{
		// 	data.NBTStringValue("color", "red"),
		// 	data.NBTStringValue("text", "Hello world"),
		// }),
		Reason: data.MakeChat().
			SetText("Hello world!!!").
			SetColor(data.ChatColorAqua).
			AddExtra(
				data.MakeChat().
					SetText(" Test")).
			ToNBT(nil),
	}
	encoded, err := messages.Encode(&res)
	if err != nil {
		return err
	}
	c.SendPacket(&encoded)

	return nil
}
