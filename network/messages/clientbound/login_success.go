package clientbound

import (
	"github.com/brenfwd/gocraft/constants"
	"github.com/brenfwd/gocraft/data"
	"github.com/brenfwd/gocraft/network/messages"
	"github.com/google/uuid"
)

func init() {
	messages.RegisterClientbound[LoginSuccess](constants.ClientStateLogin, 0x02)
}

type LoginSuccess_Property struct {
	Name      string
	Value     string
	Signature *string
}

func (p *LoginSuccess_Property) BufferWrite(buf *data.Buffer) (err error) {
	buf.WriteString(p.Name)
	buf.WriteString(p.Value)
	buf.WriteBoolean(p.Signature != nil)
	if p.Signature != nil {
		buf.WriteString(*p.Signature)
	}
	return
}

type LoginSuccess struct {
	messages.Clientbound
	UUID                uuid.UUID
	Username            string
	Properties          []LoginSuccess_Property `message:"length:varint"`
	StrictErrorHandling bool
}
