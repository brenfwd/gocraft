package data

import (
	"fmt"
)

//go:generate stringer -type=NBTTag
type NBTTag int

const (
	TAG_End NBTTag = iota
	TAG_Byte
	TAG_Short
	TAG_Int
	TAG_Long
	TAG_Float
	TAG_Double
	TAG_Byte_Array
	TAG_String
	TAG_List
	TAG_Compound
	TAG_Int_Array
	TAG_Long_Array
)

type NBTValue struct {
	Tag   NBTTag
	Name  *string
	Value any
}

func makeValue(t NBTTag, name *string, v any) *NBTValue {
	return &NBTValue{Tag: t, Name: name, Value: v}
}

func NBTCompoundValue(name *string, entries []*NBTValue) *NBTValue {
	return makeValue(TAG_Compound, name, entries)
}

func NBTListValue(name string, entries []*NBTValue) *NBTValue {
	t := TAG_End
	for _, entry := range entries {
		if t == TAG_End {
			t = entry.Tag
		} else if t != entry.Tag {
			panic(fmt.Errorf("mismatched tags in ListValue call, expected all to be of type %v but got %v", t, entry.Tag))
		}
	}
	return makeValue(TAG_List, &name, entries)
}

func NBTByteValue(name string, val byte) *NBTValue {
	return makeValue(TAG_Byte, &name, val)
}

func NBTShortValue(name string, val int16) *NBTValue {
	return makeValue(TAG_Byte, &name, val)
}

func NBTIntValue(name string, val int32) *NBTValue {
	return makeValue(TAG_Int, &name, val)
}

func NBTLongValue(name string, val int64) *NBTValue {
	return makeValue(TAG_Long, &name, val)
}

func NBTFloatValue(name string, val float32) *NBTValue {
	return makeValue(TAG_Float, &name, val)
}

func NBTDoubleValue(name string, val float64) *NBTValue {
	return makeValue(TAG_Double, &name, val)
}

func NBTByteArrayValue(name string, val []byte) *NBTValue {
	return makeValue(TAG_Byte_Array, &name, val)
}

func NBTStringValue(name string, val string) *NBTValue {
	return makeValue(TAG_String, &name, val)
}

func NBTIntArrayValue(name string, val []int32) *NBTValue {
	return makeValue(TAG_Int_Array, &name, val)
}

func NBTLongArrayValue(name string, val []int64) *NBTValue {
	return makeValue(TAG_Long_Array, &name, val)
}

type state int

const (
	state_Default state = iota
	state_InCompound
	state_InList
)

// Implement BufferWritable interface so we can call Buffer.WriteAny(...) with this

func (v *NBTValue) bufferWriteInternal(buf *Buffer, s state) error {
	endValue := makeValue(TAG_End, nil, nil)

	// If we are inside of a compound:
	// <type> <name length> <name> <payload>
	//
	// If we are inside of a list (in Payload mode):
	// <payload>
	//
	// Otherwise, if we are at the top level:
	// <type> <payload>

	// Write data before payload (if any)
	switch s {
	case state_Default:
		buf.Push(byte(v.Tag))
	case state_InCompound:
		// Compounds fields should be of type Value
		buf.Push(byte(v.Tag))
		if v.Name != nil {
			nameBytes := []byte(*v.Name)
			buf.WriteUShort(uint16(len(nameBytes)))
			buf.Write(nameBytes)
		}
	case state_InList:
		// Don't need to do anything in this case
	default:
		return fmt.Errorf("unknown state %v", s)
	}

	// Write payload itself
	switch v.Tag {
	case TAG_Compound:
		entries := v.Value.([]*NBTValue)
		for _, entry := range entries {
			entry.bufferWriteInternal(buf, state_InCompound)
		}
		endValue.bufferWriteInternal(buf, state_InCompound)
	case TAG_List:
		entries := v.Value.([]*NBTValue)
		var t NBTTag
		if len(entries) == 0 {
			t = TAG_End
		} else {
			t = entries[0].Tag
		}
		buf.Push(byte(t))
		buf.WriteInt(int32(len(entries)))
		for _, entry := range entries {
			if entry.Tag != t {
				return fmt.Errorf("inconsistent types in list: expected all to be of type %v but got an element of type %v", t, entry.Tag)
			}
			entry.bufferWriteInternal(buf, state_InList)
		}
	case TAG_Byte:
		buf.Push(v.Value.(byte))
	case TAG_Short:
		buf.WriteShort(v.Value.(int16))
	case TAG_Int:
		buf.WriteInt(v.Value.(int32))
	case TAG_Long:
		buf.WriteLong(v.Value.(int64))
	case TAG_Float:
		buf.WriteFloat(v.Value.(float32))
	case TAG_Double:
		buf.WriteDouble(v.Value.(float64))
	case TAG_Byte_Array:
		bytes := v.Value.([]byte)
		buf.WriteInt(int32(len(bytes)))
		buf.Write(bytes)
	case TAG_String:
		value := v.Value.(string)
		valueBytes := []byte(value)
		buf.WriteUShort(uint16(len(valueBytes)))
		buf.Write(valueBytes)
	case TAG_Int_Array:
		values := v.Value.([]int32)
		buf.WriteInt(int32(len(values)))
		for _, value := range values {
			buf.WriteInt(value)
		}
	case TAG_Long_Array:
		values := v.Value.([]int64)
		buf.WriteInt(int32(len(values)))
		for _, value := range values {
			buf.WriteLong(value)
		}
	default:
		return fmt.Errorf("unhandled nbt tag type %v", v.Tag)
	}

	return nil
}

func (v *NBTValue) BufferWrite(buf *Buffer) {
	err := v.bufferWriteInternal(buf, state_Default)
	if err != nil {
		panic(err)
	}
}

// TODO: buffer read...
