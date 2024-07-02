package messages

import (
	_ "embed"

	"github.com/brenfwd/gocraft/constants"
	"github.com/brenfwd/gocraft/data"
	"github.com/brenfwd/gocraft/ipc"
	"github.com/brenfwd/gocraft/network"
	"github.com/brenfwd/gocraft/network/messages"
)

//go:embed tempresponse.json
var tempresponse string

func init() {
	messages.RegisterServerbound[StatusServerboundStatusRequest](constants.ClientStateStatus, 0x00)
}

type StatusServerboundStatusRequest struct {
	messages.Serverbound
}

func (p *StatusServerboundStatusRequest) Handle(i *ipc.ClientIPC) error {
	var wbuf data.Buffer
	_, err := wbuf.WriteString(tempresponse)
	if err != nil {
		return err
	}
	outPacket := network.Packet{Id: 0, Body: wbuf.Raw}
	i.SendPacket(&outPacket)
	return nil
}
