package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTypeConvert(t *testing.T) {
	var u32 uint32 = 1024
	buf, err := GetBytes(u32)
	assert.Nil(t, err)
	t.Log(buf)
	u32New, err := ToUInt32(buf)
	assert.Equal(t, u32, u32New)

	var u16 uint16 = 1024
	buf, err = GetBytes(u16)
	assert.Nil(t, err)
	t.Log(buf)
	u16New, err := ToUInt16(buf)
	assert.Equal(t, u16, u16New)
}
