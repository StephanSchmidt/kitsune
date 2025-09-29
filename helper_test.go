package kitsune

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUuid(t *testing.T) {
	id := NewUuid()
	assert.Equal(t, uint8(1), id.Variant())
	assert.Equal(t, uint8(7), id.Version())
	assert.False(t, id.IsNil())
}

func TestFromString(t *testing.T) {
	id := NewUuid()
	idFromString, _ := FromString(id.String())
	assert.Equal(t, id, idFromString)
}
