package data

import "testing"
import "github.com/stretchr/testify/assert"

func TestMap(t *testing.T) {
	m := NewMap()
	m.Set("a", 1)
	assert.Equal(t, m.Int("a"), 1)
}
