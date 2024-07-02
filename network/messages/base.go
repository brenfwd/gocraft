package messages

import (
	"fmt"
	"log"
	"reflect"

	"github.com/brenfwd/gocraft/constants"
	"github.com/brenfwd/gocraft/data"
	"github.com/brenfwd/gocraft/ipc"
	"github.com/brenfwd/gocraft/network"
	"github.com/google/uuid"
)

type ServerboundInterface interface {
	Handle(*ipc.ClientIPC) error
}

type Serverbound struct{}

type Clientbound struct{}

// func (*Serverbound) Handle(*ipc.Client) {
// 	fmt.Println("Hello from Base.Handle()!")
// }

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
		panic(fmt.Sprint("registering non-serverbound or non-anonymously-typed serverbound type", t))
	}
	serverboundRegistry[serverboundRegistryKey{State: state, Id: id}] = t
}

func LookupServerbound(state constants.ClientState, id int) (reflect.Type, bool) {
	t, found := serverboundRegistry[serverboundRegistryKey{State: state, Id: id}]
	return t, found
}

func DecodeServerbound(state constants.ClientState, packet *network.Packet) (ServerboundInterface, error) {
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
		var value any
		var err error

		switch f.Type {
		case reflect.TypeFor[data.VarInt]():
			value, _, err = buf.ReadVarInt()
		case reflect.TypeFor[data.VarLong]():
			value, _, err = buf.ReadVarLong()
		case reflect.TypeFor[uuid.UUID]():
			value, err = buf.ReadUUID()
		default:
			switch f.Type.Kind() {
			// TODO: other types (float, bool...)
			case reflect.String:
				value, _, err = buf.ReadString()
			case reflect.Uint16:
				value, err = buf.ReadUShort()
			case reflect.Int:
			case reflect.Int32:
				value, err = buf.ReadInt()
			case reflect.Int64:
				value, err = buf.ReadLong()
			default:
				err = fmt.Errorf("unhandled type %v with kind %v", f.Type, f.Type.Kind())
			}
		}

		if err != nil {
			return nil, err
		}

		target.Set(reflect.ValueOf(value))
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
		panic(fmt.Sprint("registering non-clientbound or non-anonymously-typed clientbound type", t))
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

		switch f.Type {
		case reflect.TypeFor[data.VarInt]():
			_, err := wbuf.WriteVarInt(value.(data.VarInt))
			if err != nil {
				return network.Packet{}, err
			}
		case reflect.TypeFor[data.VarLong]():
			_, err := wbuf.WriteVarLong(value.(data.VarLong))
			if err != nil {
				return network.Packet{}, err
			}
		case reflect.TypeFor[data.Chat]():
			local := value.(data.Chat)
			str, err := local.String()
			if err != nil {
				return network.Packet{}, err
			}
			_, err = wbuf.WriteString(str)
			if err != nil {
				return network.Packet{}, err
			}
		case reflect.TypeFor[uuid.UUID]():
			err := wbuf.WriteUUID(value.(uuid.UUID))
			if err != nil {
				return network.Packet{}, err
			}
		default:
			switch f.Type.Kind() {
			case reflect.String:
				_, err := wbuf.WriteString(value.(string))
				if err != nil {
					return network.Packet{}, err
				}
			case reflect.Uint16:
				err := wbuf.WriteUShort(value.(uint16))
				if err != nil {
					return network.Packet{}, err
				}
			case reflect.Int:
			case reflect.Int32:
				err := wbuf.WriteInt(value.(int32))
				if err != nil {
					return network.Packet{}, err
				}
			case reflect.Int64:
				err := wbuf.WriteLong(value.(int64))
				if err != nil {
					return network.Packet{}, err
				}
			default:
				return network.Packet{}, fmt.Errorf("unhandled type %v with kind %v", f.Type, f.Type.Kind())
			}
		}
	}

	return network.Packet{Id: info.Id, Body: wbuf.Raw}, nil
}
