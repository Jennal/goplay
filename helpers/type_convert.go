package helpers

import (
	"bytes"
	"encoding/binary"
	"errors"
)

func GetBytes(i interface{}) ([]byte, error) {
	var buffer bytes.Buffer
	err := binary.Write(&buffer, binary.BigEndian, i)

	return buffer.Bytes(), err
}

func ToUInt32(buffer []byte) (uint32, error) {
	if len(buffer) < 4 {
		return 0, errors.New("length of buffer < 4")
	}

	var i uint32 = 0
	r := bytes.NewReader(buffer[:4])
	err := binary.Read(r, binary.BigEndian, &i)

	return i, err
}

func ToUInt16(buffer []byte) (uint16, error) {
	if len(buffer) < 2 {
		return 0, errors.New("length of buffer < 2")
	}

	var i uint16 = 0
	r := bytes.NewReader(buffer[:2])
	err := binary.Read(r, binary.BigEndian, &i)

	return i, err
}
