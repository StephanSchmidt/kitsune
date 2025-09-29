package kitsune

import (
	"database/sql/driver"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/gofrs/uuid"
)

func NewUuid() uuid.UUID {
	return uuid.Must(uuid.NewV7())
}

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

func FromString(uuidStr string) (uuid.UUID, error) {
	u, err := uuid.FromString(uuidStr)
	if err != nil {
		return uuid.UUID{}, err
	}
	return u, nil
}

func FromStringMust(uuidStr string) uuid.UUID {
	u, err := uuid.FromString(uuidStr)
	if err != nil {
		panic(err)
	}
	return u
}

func FromByteArray(uuidBytes []byte) (uuid.UUID, error) {
	u, err := uuid.FromBytes(uuidBytes)
	if err != nil {
		return uuid.UUID{}, err
	}
	return u, nil
}

func ToString(u uuid.UUID) string {
	return hex.EncodeToString(u.Bytes()[:])
}

func Marshall(d any) (driver.Value, error) {
	return json.Marshal(d)
}

// Unmarshal is a generic function that unmarshals a value into a specified struct type.
func Unmarshal[T any](value interface{}, target *T) error {
	b, ok := value.(string)
	if ok {
		return json.Unmarshal([]byte(b), target)
	}
	data, ok := value.([]uint8)
	if ok {
		return json.Unmarshal(data, target)
	}
	return errors.New("type assertion to []byte and []string failed")
}
