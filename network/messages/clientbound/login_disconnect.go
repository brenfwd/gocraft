package clientbound

import (
	"github.com/brenfwd/gocraft/constants"
	"github.com/brenfwd/gocraft/data"
	"github.com/brenfwd/gocraft/network/messages"
)

func init() {
	messages.RegisterClientbound[LoginClientboundDisconnect](constants.ClientStateLogin, 0x00)
}

type LoginClientboundDisconnect struct {
	messages.Clientbound
	Reason data.Chat
}
