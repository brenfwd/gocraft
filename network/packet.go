package network

import (
	"errors"

	"github.com/brenfwd/gocraft/data"
)

type Packet struct {
	Id   int
	Body []byte
}

type PacketUnmarshaller struct {
	buffer data.Buffer
}

func (pu *PacketUnmarshaller) Unmarshal(newData []byte) ([]Packet, error) {
	if err := pu.buffer.Write(newData); err != nil {
		return []Packet{}, err
	}
	if pu.buffer.Length() > 4*1024*1024 {
		return []Packet{}, errors.New("buffer grew too big")
	}

	var packets []Packet
	for !pu.buffer.Empty() {
		save := pu.buffer.Raw

		length, _, err := pu.buffer.ReadVarInt()
		if err != nil {
			pu.buffer.Raw = save
			break
		}

		if length < 1 {
			return []Packet{}, errors.New("invalid packet length")
		}

		if int(length) > pu.buffer.Length() {
			pu.buffer.Raw = save
			break
		}

		packetId, packetIdBytes, err := pu.buffer.ReadVarInt()
		if err != nil {
			return []Packet{}, err
		}

		body, err := pu.buffer.Read(int(length) - packetIdBytes)
		if err != nil {
			return []Packet{}, err
		}

		packets = append(packets, Packet{Id: int(packetId), Body: body})
	}

	return packets, nil
}

func (p *Packet) Marshal() ([]byte, error) {
	var headWBuf data.Buffer
	var bodyBuf data.Buffer
	bodyBuf.WriteVarInt(data.VarInt(p.Id))
	bodyBuf.Write(p.Body)
	headWBuf.WriteVarInt(data.VarInt(bodyBuf.Length()))
	return append(headWBuf.Raw, bodyBuf.Raw...), nil
}
