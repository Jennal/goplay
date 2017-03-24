package session

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var idgen = NewIDGen()

func TestIdGen(t *testing.T) {
	assert.Equal(t, 0, idgen.NextID())
	assert.Equal(t, 1, idgen.NextID())
}
