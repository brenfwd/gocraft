package data

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/google/uuid"
)

type Buffer struct {
	Raw []byte
}

func NewBufferFromBytes(bytes []byte) Buffer {
	return Buffer{Raw: bytes}
}

func (buf *Buffer) Length() int {
	return len(buf.Raw)
}

func (buf *Buffer) Empty() bool {
	return len(buf.Raw) == 0
}

func (buf *Buffer) Write(data []byte) {
	buf.Raw = append(buf.Raw, data...)
}

func (buf *Buffer) Read(length int) ([]byte, error) {
	if len(buf.Raw) < length {
		return nil, errors.New("buffer is too short in Read")
	}
	data := buf.Raw[:length]
	buf.Raw = buf.Raw[length:]
	return data, nil
}

func (buf *Buffer) ReadByte() (byte, error) {
	if len(buf.Raw) == 0 {
		return 0, errors.New("bufer is empty in ReadByte")
	}
	first := buf.Raw[0]
	buf.Raw = buf.Raw[1:]
	return first, nil
}

func (buf *Buffer) Push(b byte) {
	buf.Raw = append(buf.Raw, b)
}

func (buf *Buffer) ReadBoolean() (bool, error) {
	v, err := buf.ReadByte()
	if err != nil {
		return false, err
	}
	switch v {
	case 0x00:
		return false, nil
	case 0x01:
		return true, nil
	default:
		return false, errors.New("invalid boolean value in ReadBoolean")
	}
}

func (buf *Buffer) WriteBoolean(v bool) {
	if v {
		buf.Push(0x01)
	} else {
		buf.Push(0x00)
	}
}

func (buf *Buffer) ReadUByte() (uint8, error) {
	v, err := buf.ReadByte()
	return uint8(v), err
}

func (buf *Buffer) WriteUByte(v uint8) {
	buf.Push(byte(v))
}

func (buf *Buffer) ReadShort() (int16, error) {
	if len(buf.Raw) < 2 {
		return 0, errors.New("buffer is too short in ReadShort")
	}
	v := int16(buf.Raw[0])<<8 | int16(buf.Raw[1])
	buf.Raw = buf.Raw[2:]
	return v, nil
}

func (buf *Buffer) WriteShort(v int16) {
	buf.Push(byte(v >> 8))
	buf.Push(byte(v))
}

func (buf *Buffer) ReadUShort() (uint16, error) {
	v, err := buf.ReadShort()
	return uint16(v), err
}

func (buf *Buffer) WriteUShort(v uint16) {
	buf.WriteShort(int16(v))
}

func (buf *Buffer) ReadInt() (int32, error) {
	if len(buf.Raw) < 4 {
		return 0, errors.New("buffer is too short in ReadInt")
	}
	v := int32(buf.Raw[0])<<24 | int32(buf.Raw[1])<<16 | int32(buf.Raw[2])<<8 | int32(buf.Raw[3])
	buf.Raw = buf.Raw[4:]
	return v, nil
}

func (buf *Buffer) WriteInt(v int32) {
	buf.Push(byte(v >> 24))
	buf.Push(byte(v >> 16))
	buf.Push(byte(v >> 8))
	buf.Push(byte(v))
}

func (buf *Buffer) ReadUInt() (uint32, error) {
	v, err := buf.ReadInt()
	return uint32(v), err
}

func (buf *Buffer) WriteUInt(v uint32) {
	buf.WriteInt(int32(v))
}

func (buf *Buffer) ReadLong() (int64, error) {
	if len(buf.Raw) < 8 {
		return 0, errors.New("buffer is too short in ReadLong")
	}
	v := int64(buf.Raw[0])<<56 | int64(buf.Raw[1])<<48 | int64(buf.Raw[2])<<40 | int64(buf.Raw[3])<<32 | int64(buf.Raw[4])<<24 | int64(buf.Raw[5])<<16 | int64(buf.Raw[6])<<8 | int64(buf.Raw[7])
	buf.Raw = buf.Raw[8:]
	return v, nil
}

func (buf *Buffer) WriteLong(v int64) {
	buf.Push(byte(v >> 56))
	buf.Push(byte(v >> 48))
	buf.Push(byte(v >> 40))
	buf.Push(byte(v >> 32))
	buf.Push(byte(v >> 24))
	buf.Push(byte(v >> 16))
	buf.Push(byte(v >> 8))
	buf.Push(byte(v))
}

func (buf *Buffer) ReadULong() (uint64, error) {
	v, err := buf.ReadLong()
	return uint64(v), err
}

func (buf *Buffer) WriteULong(v uint64) {
	buf.WriteLong(int64(v))
}

func (buf *Buffer) ReadFloat() (float32, error) {
	v, err := buf.ReadUInt()
	return float32(v), err
}

func (buf *Buffer) WriteFloat(v float32) {
	buf.WriteUInt(uint32(v))
}

func (buf *Buffer) ReadDouble() (float64, error) {
	v, err := buf.ReadULong()
	return float64(v), err
}

func (buf *Buffer) WriteDouble(v float64) {
	buf.WriteULong(uint64(v))
}

func (buf *Buffer) ReadVarInt() (value VarInt, bytes int, err error) {
	// var value int32 = 0
	var position int = 0
	err = nil

	for {
		currentByte, err := buf.ReadByte()
		if err != nil {
			return 0, 0, err
		}
		bytes++
		value |= VarInt((int32(currentByte) & 0x7F) << position)
		if int(currentByte)&0x80 == 0 {
			break
		}
		position += 7
		if position >= 32 {
			err = errors.New("varint is too big")
			return 0, 0, err
		}
	}

	return
}

func (buf *Buffer) WriteVarInt(v VarInt) (bytes int) {
	for {
		if v&^0x7F == 0 {
			buf.Push(byte(v))
			bytes++
			break
		}
		buf.Push(byte(v&0x7F | 0x80))
		bytes++
		v = VarInt(uint32(v) >> 7)
	}
	return
}

func (buf *Buffer) ReadVarLong() (value VarLong, bytes int, err error) {
	// var value int64 = 0
	var position int = 0
	err = nil

	for {
		currentByte, err := buf.ReadByte()
		if err != nil {
			return 0, 0, err
		}
		bytes++
		value |= VarLong((int64(currentByte) & 0x7F) << position)
		if int(currentByte)&0x80 == 0 {
			break
		}
		position += 7
		if position >= 64 {
			err = errors.New("varlong is too big")
			return 0, 0, err
		}
	}

	return
}

func (buf *Buffer) WriteVarLong(v VarLong) (bytes int) {
	for {
		if v&^0x7F == 0 {
			buf.Push(byte(v))
			bytes++
			break
		}
		buf.Push(byte(v&0x7F | 0x80))
		bytes++
		v = VarLong(uint64(v) >> 7)
	}
	return
}

func (buf *Buffer) ReadString() (value string, bytes int, err error) {
	var length VarInt
	length, bytes, err = buf.ReadVarInt()
	if err != nil {
		return
	}
	if length < 0 {
		err = errors.New("string length is negative")
		return
	}
	if length == 0 {
		return "", bytes, nil
	}
	if length > VarInt(len(buf.Raw)) {
		err = errors.New("string length is too big")
		return
	}
	value = string(buf.Raw[:length])
	buf.Raw = buf.Raw[length:]
	bytes += int(length)
	return
}

func (buf *Buffer) WriteString(str string) (bytes int) {
	lbytes := buf.WriteVarInt(VarInt(len(str)))
	bytes += lbytes
	strbytes := []byte(str)
	buf.Raw = append(buf.Raw, strbytes...)
	bytes += len(strbytes)
	return
}

func (buf *Buffer) ReadUUID() (uuid.UUID, error) {
	bytes, err := buf.Read(16)
	if err != nil {
		return uuid.UUID{}, err
	}
	return uuid.FromBytes(bytes)
}

func (buf *Buffer) WriteUUID(value uuid.UUID) {
	buf.Write(value[:])
}

type BufferWritable interface {
	BufferWrite(*Buffer)
}

type BufferReadable[T any] interface {
	BufferRead(*Buffer) (T, error)
}

type BufferSliceLength string

const (
	BufferSliceLengthVarInt BufferSliceLength = "varint"
)

func (buf *Buffer) ReadReflectedSlice(elemType reflect.Type, lengthType BufferSliceLength) (value reflect.Value, err error) {
	length, err := buf.readLength(lengthType)
	if err != nil {
		return reflect.Value{}, err
	}

	slice := reflect.MakeSlice(reflect.SliceOf(elemType), length, length)
	for i := 0; i < length; i++ {
		v, err := buf.ReadReflected(elemType)
		if err != nil {
			return reflect.Value{}, err
		}
		slice.Index(i).Set(v)
	}
	return slice, nil
}

func (buf *Buffer) ReadReflected(t reflect.Type) (value reflect.Value, err error) {
	// Check if value implements BufferReadable[T] for some T...
	if t.Implements(reflect.TypeFor[BufferReadable[any]]()) {
		readMethod := reflect.ValueOf(value).MethodByName("BufferRead")
		if !readMethod.IsValid() {
			return reflect.Value{}, fmt.Errorf("BufferReadable type %v does not have a BufferRead method", t)
		}
		args := []reflect.Value{reflect.ValueOf(buf)}
		results := readMethod.Call(args)
		if len(results) != 2 {
			return reflect.Value{}, fmt.Errorf("BufferReadable type %v BufferRead method did not return 2 values", t)
		}
		if !results[1].IsNil() {
			return reflect.Value{}, results[1].Interface().(error)
		}
		return results[0], nil
	}

	if t.Kind() == reflect.Slice {
		return buf.ReadReflectedSlice(t, BufferSliceLengthVarInt)
	}
	var v any
	switch t {
	case reflect.TypeFor[VarInt]():
		v, _, err = buf.ReadVarInt()
	case reflect.TypeFor[VarLong]():
		v, _, err = buf.ReadVarLong()
	case reflect.TypeFor[uuid.UUID]():
		v, err = buf.ReadUUID()
	case reflect.TypeFor[string]():
		v, _, err = buf.ReadString()
	case reflect.TypeFor[bool]():
		v, err = buf.ReadBoolean()
	case reflect.TypeFor[byte]():
		v, err = buf.ReadByte()
	case reflect.TypeFor[uint16]():
		v, err = buf.ReadUShort()
	case reflect.TypeFor[int32]():
		v, err = buf.ReadInt()
	case reflect.TypeFor[int]():
		v, err = buf.ReadInt()
	case reflect.TypeFor[int64]():
		v, err = buf.ReadLong()
	default:
		err = fmt.Errorf("unhandled type %v with kind %v", t, t.Kind())
	}
	if err != nil {
		return reflect.Value{}, err
	}
	return reflect.ValueOf(v), nil
}

func (buf *Buffer) readLength(lengthType BufferSliceLength) (int, error) {
	switch lengthType {
	case BufferSliceLengthVarInt:
		v, _, err := buf.ReadVarInt()
		return int(v), err
	default:
		return 0, fmt.Errorf("unhandled length type for ReadSlice: %v", lengthType)
	}
}

func (buf *Buffer) WriteSlice(value any, lengthType BufferSliceLength) (err error) {
	t := reflect.TypeOf(value)
	if t.Kind() != reflect.Slice {
		return fmt.Errorf("value passed to WriteSlice is not a slice: %v", t.Kind())
	}

	reflected := reflect.ValueOf(value)
	if err := buf.writeLength(lengthType, reflected.Len()); err != nil {
		return err
	}

	for i := 0; i < reflected.Len(); i++ {
		v := reflected.Index(i).Interface()
		if err := buf.WriteAny(v); err != nil {
			return err
		}
	}
	return
}

func (buf *Buffer) writeLength(lengthType BufferSliceLength, length int) error {
	switch lengthType {
	case BufferSliceLengthVarInt:
		buf.WriteVarInt(VarInt(length))
		return nil
	default:
		return fmt.Errorf("unhandled length type for WriteSlice: %v", lengthType)
	}
}

func (buf *Buffer) WriteAny(value any) error {
	if writable, ok := value.(BufferWritable); ok {
		writable.BufferWrite(buf)
		return nil
	}

	switch v := value.(type) {
	case VarInt:
		buf.WriteVarInt(v)
		return nil
	case VarLong:
		buf.WriteVarLong(v)
		return nil
	case Chat:
		s, err := v.String()
		if err != nil {
			return err
		}
		buf.WriteString(s)
		return nil
	case uuid.UUID:
		buf.WriteUUID(v)
		return nil
	case string:
		buf.WriteString(v)
		return nil
	case bool:
		buf.WriteBoolean(v)
		return nil
	case byte:
		buf.Push(v)
		return nil
	case uint16:
		buf.WriteUShort(v)
		return nil
	case int32:
		buf.WriteInt(v)
		return nil
	case int:
		buf.WriteInt(int32(v))
		return nil
	case int64:
		buf.WriteLong(v)
		return nil
	default:
		return fmt.Errorf("unhandled type for WriteAny: %T", value)
	}
}
