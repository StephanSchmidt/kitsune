package kitsune

import (
	"encoding/hex"

	"github.com/gofrs/uuid"
)

func NewUuid() uuid.UUID {
	return uuid.Must(uuid.NewV7())
}

func FromString(uuidStr string) (uuid.UUID, error) {
	u, err := uuid.FromString(uuidStr)
	if err != nil {
		return uuid.UUID{}, WrapWithDetails(err, "failed to parse UUID string", "uuid", uuidStr)
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
		return uuid.UUID{}, WrapWithDetails(err, "failed to parse UUID from bytes", "uuid", uuidBytes)
	}
	return u, nil
}

func ToString(u uuid.UUID) string {
	return hex.EncodeToString(u.Bytes()[:])
}
