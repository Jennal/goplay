package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNextID(t *testing.T) {
	for i := 0; i < 512; i++ {
		id := NextID()
		assert.Equal(t, PackageIDType(i%256), id)
	}
}
