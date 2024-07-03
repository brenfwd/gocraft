package serverbound

import (
	"fmt"
	"log"

	"github.com/brenfwd/gocraft/constants"
	"github.com/brenfwd/gocraft/network/messages"
	"github.com/brenfwd/gocraft/network/messages/clientbound"
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
	fmt.Printf("%+v\n", p)
	// SharedSecret and VerifyToken are encrypted using the server public key
	decSharedSecret := p.SharedSecret
	c.ListenerKeypair.DecryptWithPrivateKey(&decSharedSecret)
	decVerifyToken := p.VerifyToken
	c.ListenerKeypair.DecryptWithPrivateKey(&decVerifyToken)

	// Store the shared secret and check the verify token
	c.SharedSecret = decSharedSecret

	// Compare array with slice
	if len(decVerifyToken) != len(c.EncryptionVerifyToken) {
		return fmt.Errorf("verify token length mismatch")
	}
	for i := range decVerifyToken {
		if decVerifyToken[i] != c.EncryptionVerifyToken[i] {
			return fmt.Errorf("verify token mismatch")
		}
	}

	log.Printf("Encryption response: using shared secret %x", c.SharedSecret)

	// Enable encryption
	c.EnableEncryption()

	// Send login success
	res := clientbound.LoginSuccess{
		UUID:       c.AllegedUUID,
		Username:   c.AllegedUsername,
		Properties: []clientbound.LoginSuccess_Property{},
	}
	encoded, err := messages.Encode(&res)
	if err != nil {
		return err
	}
	c.SendPacket(&encoded)

	return nil
}
