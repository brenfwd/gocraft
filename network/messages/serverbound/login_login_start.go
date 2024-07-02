package messages

import (
	"log"

	"github.com/brenfwd/gocraft/constants"
	"github.com/brenfwd/gocraft/network/messages"
	"github.com/brenfwd/gocraft/network/messages/clientbound"
	"github.com/brenfwd/gocraft/shared"
	"github.com/google/uuid"
)

func init() {
	messages.RegisterServerbound[LoginServerboundLoginStart](constants.ClientStateLogin, 0x00)
}

type LoginServerboundLoginStart struct {
	messages.Serverbound
	Name       string
	PlayerUUID uuid.UUID
}

func (p *LoginServerboundLoginStart) Handle(c *shared.ClientShared) error {
	log.Println(p)

	res := clientbound.LoginEncryptionRequest{
		ServerID:           "",
		PublicKey:          c.ListenerKeypair.PublicKey,
		VerifyToken:        c.EncryptionVerifyToken[:],
		ShouldAuthenticate: true,
	}
	encoded, err := messages.Encode(&res)
	if err != nil {
		return err
	}
	c.SendPacket(&encoded)

	return nil
}
