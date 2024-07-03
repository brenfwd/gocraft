package clientbound

import (
	"github.com/brenfwd/gocraft/data"
	"github.com/brenfwd/gocraft/network/messages"
)

func init() {
}

type ConfigurationKnownPacks_Pack struct {
	Namespace string
	ID        string
	Version   string
}

func (p *ConfigurationKnownPacks_Pack) BufferWrite(buf *data.Buffer) (err error) {
	buf.WriteString(p.Namespace)
	buf.WriteString(p.ID)
	buf.WriteString(p.Version)
	return
}

// TODO: get a packet dump of this from vanilla

type ConfigurationKnownPacks struct {
	messages.Clientbound
	KnownPacks []ConfigurationKnownPacks_Pack `message:"length:varint"`
}
