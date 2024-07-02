package clientbound

import (
	"github.com/brenfwd/gocraft/constants"
	"github.com/brenfwd/gocraft/network/messages"
)

func init() {
	messages.RegisterClientbound[LoginEncryptionRequest](constants.ClientStateLogin, 0x01)
}

type LoginEncryptionRequest struct {
	messages.Clientbound
	ServerID           string // Empty?
	PublicKey          []byte `message:"length:varint"`
	VerifyToken        []byte `message:"length:varint"`
	ShouldAuthenticate bool
}
