package kitsune

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

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
