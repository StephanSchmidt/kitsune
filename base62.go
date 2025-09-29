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

	var id uuid.UUID
	copy(id[:], i.Bytes())
	return id, nil
}
