package kitsune

import (
	"fmt"
	"math/big"

	"github.com/gofrs/uuid"
)

func ToBase62(uuid uuid.UUID) string {
	var i big.Int
	i.SetBytes(uuid[:])
	return i.Text(62)
}

func FromBase62(s string) (uuid.UUID, error) {
	var i big.Int
	_, ok := i.SetString(s, 62)
	if !ok {
		return uuid.UUID{}, fmt.Errorf("cannot parse base62: %q", s)
	}

	if i.Sign() < 0 {
		return uuid.UUID{}, fmt.Errorf("cannot parse base62: %q is negative", s)
	}
	b := i.Bytes()
	if len(b) > len(uuid.UUID{}) {
		return uuid.UUID{}, fmt.Errorf("cannot parse base62: %q is longer than a UUID", s)
	}
	// Bytes drops leading zero bytes; right-align so UUIDs starting with
	// zero bytes decode to themselves.
	var id uuid.UUID
	copy(id[len(id)-len(b):], b)
	return id, nil
}
