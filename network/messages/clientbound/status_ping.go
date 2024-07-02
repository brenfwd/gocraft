package clientbound

import (
	"github.com/brenfwd/gocraft/constants"
	"github.com/brenfwd/gocraft/network/messages"
)

func init() {
	messages.RegisterClientbound[StatusClientboundPing](constants.ClientStateStatus, 0x01)
}

type StatusClientboundPing struct {
	messages.Clientbound
	Response int64
}
