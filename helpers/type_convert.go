package helpers

import (
	"bytes"
	"encoding/binary"
	"errors"
)

type UInt32 uint32
type Bytes []byte

func (i UInt32) GetBytes() ([]byte, error) {
	var buffer bytes.Buffer
	err := binary.Write(&buffer, binary.BigEndian, i)

	return buffer.Bytes(), err
}

func (buffer Bytes) ToInt() (uint32, error) {
	if len(buffer) < 4 {
		return 0, errors.New("length of buffer < 4")
	}

	var i uint32 = 0
	r := bytes.NewReader(buffer[:4])
	err := binary.Read(r, binary.BigEndian, &i)

	return i, err
}
