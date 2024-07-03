package clientbound

import (
	"github.com/brenfwd/gocraft/constants"
	"github.com/brenfwd/gocraft/data"
	"github.com/brenfwd/gocraft/network/messages"
)

func init() {
	messages.RegisterClientbound[ConfigurationDisconnect](constants.ClientStateConfiguration, 0x02)
}

type ConfigurationDisconnect struct {
	messages.Clientbound
	Reason *data.NBTValue
}
