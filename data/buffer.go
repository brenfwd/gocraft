package data

import (
	"errors"

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

func (buf *Buffer) Write(data []byte) error {
	buf.Raw = append(buf.Raw, data...)
	return nil
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

func (buf *Buffer) WriteByte(b byte) error {
	buf.Raw = append(buf.Raw, b)
	return nil
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

func (buf *Buffer) WriteBoolean(v bool) error {
	if v {
		buf.WriteByte(0x01)
	} else {
		buf.WriteByte(0x00)
	}
	return nil
}

func (buf *Buffer) ReadUByte() (uint8, error) {
	v, err := buf.ReadByte()
	return uint8(v), err
}

func (buf *Buffer) WriteUByte(v uint8) error {
	return buf.WriteByte(byte(v))
}

func (buf *Buffer) ReadShort() (int16, error) {
	if len(buf.Raw) < 2 {
		return 0, errors.New("buffer is too short in ReadShort")
	}
	v := int16(buf.Raw[0])<<8 | int16(buf.Raw[1])
	buf.Raw = buf.Raw[2:]
	return v, nil
}

func (buf *Buffer) WriteShort(v int16) error {
	buf.WriteByte(byte(v >> 8))
	buf.WriteByte(byte(v))
	return nil
}

func (buf *Buffer) ReadUShort() (uint16, error) {
	v, err := buf.ReadShort()
	return uint16(v), err
}

func (buf *Buffer) WriteUShort(v uint16) error {
	return buf.WriteShort(int16(v))
}

func (buf *Buffer) ReadInt() (int32, error) {
	if len(buf.Raw) < 4 {
		return 0, errors.New("buffer is too short in ReadInt")
	}
	v := int32(buf.Raw[0])<<24 | int32(buf.Raw[1])<<16 | int32(buf.Raw[2])<<8 | int32(buf.Raw[3])
	buf.Raw = buf.Raw[4:]
	return v, nil
}

func (buf *Buffer) WriteInt(v int32) error {
	var errs []error
	errs = append(errs, buf.WriteByte(byte(v>>24)))
	errs = append(errs, buf.WriteByte(byte(v>>16)))
	errs = append(errs, buf.WriteByte(byte(v>>8)))
	errs = append(errs, buf.WriteByte(byte(v)))

	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}

func (buf *Buffer) ReadUInt() (uint32, error) {
	v, err := buf.ReadInt()
	return uint32(v), err
}

func (buf *Buffer) WriteUInt(v uint32) error {
	return buf.WriteInt(int32(v))
}

func (buf *Buffer) ReadLong() (int64, error) {
	if len(buf.Raw) < 8 {
		return 0, errors.New("buffer is too short in ReadLong")
	}
	v := int64(buf.Raw[0])<<56 | int64(buf.Raw[1])<<48 | int64(buf.Raw[2])<<40 | int64(buf.Raw[3])<<32 | int64(buf.Raw[4])<<24 | int64(buf.Raw[5])<<16 | int64(buf.Raw[6])<<8 | int64(buf.Raw[7])
	buf.Raw = buf.Raw[8:]
	return v, nil
}

func (buf *Buffer) WriteLong(v int64) error {
	var errs []error
	errs = append(errs, buf.WriteByte(byte(v>>56)))
	errs = append(errs, buf.WriteByte(byte(v>>48)))
	errs = append(errs, buf.WriteByte(byte(v>>40)))
	errs = append(errs, buf.WriteByte(byte(v>>32)))
	errs = append(errs, buf.WriteByte(byte(v>>24)))
	errs = append(errs, buf.WriteByte(byte(v>>16)))
	errs = append(errs, buf.WriteByte(byte(v>>8)))
	errs = append(errs, buf.WriteByte(byte(v)))

	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}

func (buf *Buffer) ReadULong() (uint64, error) {
	v, err := buf.ReadLong()
	return uint64(v), err
}

func (buf *Buffer) WriteULong(v uint64) error {
	return buf.WriteLong(int64(v))
}

func (buf *Buffer) ReadFloat() (float32, error) {
	v, err := buf.ReadUInt()
	return float32(v), err
}

func (buf *Buffer) WriteFloat(v float32) error {
	return buf.WriteUInt(uint32(v))
}

func (buf *Buffer) ReadDouble() (float64, error) {
	v, err := buf.ReadULong()
	return float64(v), err
}

func (buf *Buffer) WriteDouble(v float64) error {
	return buf.WriteULong(uint64(v))
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

func (buf *Buffer) WriteVarInt(v VarInt) (bytes int, err error) {
	for {
		if v&^0x7F == 0 {
			if err = buf.WriteByte(byte(v)); err != nil {
				return
			}
			bytes++
			break
		}
		if err = buf.WriteByte(byte(v&0x7F | 0x80)); err != nil {
			return
		}
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

func (buf *Buffer) WriteVarLong(v VarLong) (bytes int, err error) {
	for {
		if v&^0x7F == 0 {
			if err = buf.WriteByte(byte(v)); err != nil {
				return
			}
			bytes++
			break
		}
		if err = buf.WriteByte(byte(v&0x7F | 0x80)); err != nil {
			return
		}
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

func (buf *Buffer) WriteString(str string) (bytes int, err error) {
	var lbytes int
	lbytes, err = buf.WriteVarInt(VarInt(len(str)))
	if err != nil {
		return
	}
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

func (buf *Buffer) WriteUUID(value uuid.UUID) error {
	return buf.Write(value[:])
}