package messages

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/brenfwd/gocraft/constants"
	"github.com/brenfwd/gocraft/data"
	"github.com/brenfwd/gocraft/network"
	"github.com/brenfwd/gocraft/shared"
)

type ServerboundInterface interface {
	Handle(*shared.ClientShared) error
}

type Serverbound struct{}

type Clientbound struct{}

// Serverbound message registry
type serverboundRegistryKey struct {
	State constants.ClientState
	Id    int
}

var serverboundRegistry = make(map[serverboundRegistryKey]reflect.Type)

func RegisterServerbound[T any](state constants.ClientState, id int) {
	t := reflect.TypeFor[T]()
	sbField, sbExists := t.FieldByName("Serverbound")
	if !sbExists || !sbField.Anonymous {
		panic(fmt.Sprint("registering non-serverbound or non-anonymously-typed serverbound type ", t))
	}
	serverboundRegistry[serverboundRegistryKey{State: state, Id: id}] = t
}

func LookupServerbound(state constants.ClientState, id int) (reflect.Type, bool) {
	t, found := serverboundRegistry[serverboundRegistryKey{State: state, Id: id}]
	return t, found
}

func DecodeServerbound(state constants.ClientState, packet *network.Packet) (ServerboundInterface, error) {
	// TODO: move to buffer like with WriteAny...
	// TODO: handle slice types
	t, found := LookupServerbound(state, packet.Id)
	if !found {
		return nil, fmt.Errorf("could not find handler for packet in state %v with ID 0x%02x (%d) -- did you forget to call RegisterServerbound?", state, packet.Id, packet.Id)
	}

	log.Printf("Decoding %v...\n", t)

	msg := reflect.New(t)

	buf := data.NewBufferFromBytes(packet.Body)

	for i := range t.NumField() {
		f := t.Field(i)
		if f.Anonymous {
			continue
		}

		target := msg.Elem().Field(i)

		if f.Type.Kind() == reflect.Slice {
			tag, ok := f.Tag.Lookup("message")
			if !ok {
				return nil, fmt.Errorf("packet %v field %v is a slice type but is missing a `message:\"length...\" tag", t, f)
			}
			const lengthPrefix string = "length:"
			suffix, ok := strings.CutPrefix(tag, lengthPrefix)
			if ok {
				value, err := buf.ReadReflectedSlice(f.Type.Elem(), data.BufferSliceLength(suffix))
				if err != nil {
					return nil, err
				}
				target.Set(value)
				continue
			} else {
				return nil, fmt.Errorf("packet %v field %v tag has unknown contents", t, f)
			}
		} else {
			value, err := buf.ReadReflected(f.Type)
			if err != nil {
				return nil, err
			}
			target.Set(value)
		}
	}

	return msg.Interface().(ServerboundInterface), nil
}

type clientboundRegistryValue struct {
	State constants.ClientState
	Id    int
}

var clientboundRegistry = make(map[reflect.Type]clientboundRegistryValue)

func RegisterClientbound[T any](state constants.ClientState, id int) {
	t := reflect.TypeFor[T]()
	sbField, sbExists := t.FieldByName("Clientbound")
	if !sbExists || !sbField.Anonymous {
		panic(fmt.Sprint("registering non-clientbound or non-anonymously-typed clientbound type ", t))
	}
	clientboundRegistry[t] = clientboundRegistryValue{State: state, Id: id}
}

func LookupClientbound(t reflect.Type) (clientboundRegistryValue, bool) {
	v, found := clientboundRegistry[t]
	return v, found
}

func Encode[T any](msg *T) (network.Packet, error) {
	t := reflect.TypeOf(msg).Elem()
	info, found := LookupClientbound(t)
	if !found {
		return network.Packet{}, fmt.Errorf("could not find ID for message %v -- did you forget to call RegisterClientbound?", t)
	}

	var wbuf data.Buffer

	for i := range t.NumField() {
		f := t.Field(i)
		if f.Anonymous {
			continue
		}

		value := reflect.ValueOf(msg).Elem().Field(i).Interface()

		// First we have to handle slice types
		if f.Type.Kind() == reflect.Slice {
			// For slices, a tag must be set: `message:"length,<ltype>"`
			// where <ltype> is one of { varint }
			tag, ok := f.Tag.Lookup("message")
			if !ok {
				return network.Packet{}, fmt.Errorf("packet %v field %v is a slice type but is missing a `message:\"length...\" tag", t, f)
			}
			const lengthPrefix string = "length:"
			suffix, ok := strings.CutPrefix(tag, lengthPrefix)
			if ok {
				err := wbuf.WriteSlice(value, data.BufferSliceLength(suffix))
				if err != nil {
					return network.Packet{}, err
				}
			} else {
				return network.Packet{}, fmt.Errorf("packet %v field %v tag has unknown contents", t, f)
			}
		} else {
			err := wbuf.WriteAny(value)
			if err != nil {
				return network.Packet{}, err
			}
		}
	}

	return network.Packet{Id: info.Id, Body: wbuf.Raw}, nil
}
