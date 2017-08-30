package helpers

import (
	"bytes"
)

type Buffer struct {
	bytes.Buffer
}

func NewBuffer(data []byte) *Buffer {
	return &Buffer{
		Buffer: *bytes.NewBuffer(data),
	}
}

func (b *Buffer) WriteUInt16(i interface{}) (int, error) {
	buf, err := GetBytes(i)
	if err != nil {
		return 0, err
	}
	return b.Buffer.Write(buf)
}

func (b *Buffer) WriteUInt32(i interface{}) (int, error) {
	buf, err := GetBytes(i)
	if err != nil {
		return 0, err
	}
	return b.Buffer.Write(buf)
}

func (b *Buffer) ReadUInt16() (uint16, error) {
	data := []byte{0, 0}
	for i := 0; i < 2; i++ {
		v, err := b.ReadByte()
		if err != nil {
			return 0, err
		}
		data[i] = v
	}

	result, err := ToUInt16(data)
	if err != nil {
		return 0, err
	}

	return result, nil
}

func (b *Buffer) ReadUInt32() (uint32, error) {
	data := []byte{0, 0, 0, 0}
	for i := 0; i < 4; i++ {
		v, err := b.ReadByte()
		if err != nil {
			return 0, err
		}
		data[i] = v
	}

	result, err := ToUInt32(data)
	if err != nil {
		return 0, err
	}

	return result, nil
}
