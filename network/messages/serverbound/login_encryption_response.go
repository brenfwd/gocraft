package messages

import (
	"fmt"

	"github.com/brenfwd/gocraft/constants"
	"github.com/brenfwd/gocraft/network/messages"
	"github.com/brenfwd/gocraft/shared"
)

func init() {
	messages.RegisterServerbound[LoginEncryptionResponse](constants.ClientStateLogin, 0x01)
}

type LoginEncryptionResponse struct {
	messages.Serverbound
	SharedSecret []byte `message:"length:varint"`
	VerifyToken  []byte `message:"length:varint"`
}

func (p *LoginEncryptionResponse) Handle(c *shared.ClientShared) error {
	fmt.Println(p)
	return nil
}
