package kitsune

import (
	"encoding/hex"
	"encoding/json"
	"strings"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewUuid(t *testing.T) {
	t.Run("creates valid UUID v7", func(t *testing.T) {
		id := NewUuid()
		assert.Equal(t, uint8(1), id.Variant())
		assert.Equal(t, uint8(7), id.Version())
		assert.False(t, id.IsNil())
	})

	t.Run("creates unique UUIDs", func(t *testing.T) {
		id1 := NewUuid()
		id2 := NewUuid()
		assert.NotEqual(t, id1, id2)
	})

	t.Run("creates monotonic UUIDs", func(t *testing.T) {
		// UUIDv7 should be monotonic when created in sequence
		id1 := NewUuid()
		id2 := NewUuid()
		assert.True(t, strings.Compare(id1.String(), id2.String()) < 0,
			"Expected %s < %s", id1.String(), id2.String())
	})
}

func TestToBase62(t *testing.T) {
	tests := []struct {
		name     string
		uuid     uuid.UUID
		checkLen bool
	}{
		{
			name:     "valid UUID to base62",
			uuid:     NewUuid(),
			checkLen: true,
		},
		{
			name:     "nil UUID to base62",
			uuid:     uuid.Nil,
			checkLen: false,
		},
		{
			name:     "custom max UUID to base62",
			uuid:     uuid.FromStringOrNil("ffffffff-ffff-ffff-ffff-ffffffffffff"),
			checkLen: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToBase62(tt.uuid)
			assert.NotEmpty(t, result)

			// Base62 string should only contain valid base62 characters
			for _, c := range result {
				assert.True(t, (c >= '0' && c <= '9') || (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z'),
					"Invalid base62 character: %c", c)
			}

			if tt.checkLen && tt.uuid != uuid.Nil {
				// Base62 should be shorter than standard UUID representation
				assert.Less(t, len(result), len(tt.uuid.String()))
			}
		})
	}
}

func TestFromBase62(t *testing.T) {
	t.Run("valid base62 to UUID", func(t *testing.T) {
		original := NewUuid()
		base62 := ToBase62(original)

		recovered, err := FromBase62(base62)
		assert.NoError(t, err)
		assert.Equal(t, original, recovered)
	})

	t.Run("nil UUID roundtrip", func(t *testing.T) {
		base62 := ToBase62(uuid.Nil)
		recovered, err := FromBase62(base62)
		assert.NoError(t, err)
		assert.Equal(t, uuid.Nil, recovered)
	})

	t.Run("invalid base62 string", func(t *testing.T) {
		invalidInputs := []string{
			"!@#$%",       // Invalid characters
			"",            // Empty string
			"hello world", // Space
			"test+test",   // Plus sign
			"test-test",   // Minus sign
		}

		for _, input := range invalidInputs {
			_, err := FromBase62(input)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "cannot parse base62")
		}
	})
}

func TestFromString(t *testing.T) {
	t.Run("valid UUID string", func(t *testing.T) {
		id := NewUuid()
		idFromString, err := FromString(id.String())
		assert.NoError(t, err)
		assert.Equal(t, id, idFromString)
	})

	t.Run("valid UUID formats", func(t *testing.T) {
		validFormats := []string{
			"550e8400-e29b-41d4-a716-446655440000",
			"550E8400-E29B-41D4-A716-446655440000",          // Uppercase
			"{550e8400-e29b-41d4-a716-446655440000}",        // With braces
			"urn:uuid:550e8400-e29b-41d4-a716-446655440000", // URN format
		}

		for _, format := range validFormats {
			u, err := FromString(format)
			assert.NoError(t, err)
			assert.NotEqual(t, uuid.Nil, u)
		}
	})

	t.Run("invalid UUID string", func(t *testing.T) {
		invalidInputs := []string{
			"not-a-uuid",
			"550e8400-e29b-41d4-a716", // Too short
			"550e8400-e29b-41d4-a716-446655440000-extra", // Too long
			// Note: UUID without hyphens is actually valid in gofrs/uuid
		}

		for _, input := range invalidInputs {
			_, err := FromString(input)
			assert.Error(t, err, "Expected error for input: %s", input)
		}
	})

	t.Run("UUID without hyphens", func(t *testing.T) {
		// gofrs/uuid accepts UUIDs without hyphens
		_, err := FromString("550e8400e29b41d4a716446655440000")
		assert.NoError(t, err)
	})

	t.Run("empty string", func(t *testing.T) {
		// Empty string returns an error
		_, err := FromString("")
		assert.Error(t, err)
	})

	t.Run("nil UUID string", func(t *testing.T) {
		u, err := FromString(uuid.Nil.String())
		assert.NoError(t, err)
		assert.Equal(t, uuid.Nil, u)
	})
}

func TestFromStringMust(t *testing.T) {
	t.Run("valid UUID string", func(t *testing.T) {
		id := NewUuid()
		idFromString := FromStringMust(id.String())
		assert.Equal(t, id, idFromString)
	})

	t.Run("panic on invalid UUID", func(t *testing.T) {
		assert.Panics(t, func() {
			FromStringMust("not-a-uuid")
		})
	})

	t.Run("nil UUID string", func(t *testing.T) {
		u := FromStringMust(uuid.Nil.String())
		assert.Equal(t, uuid.Nil, u)
	})
}

func TestFromByteArray(t *testing.T) {
	t.Run("valid byte array", func(t *testing.T) {
		original := NewUuid()
		bytes := original.Bytes()

		recovered, err := FromByteArray(bytes[:])
		assert.NoError(t, err)
		assert.Equal(t, original, recovered)
	})

	t.Run("nil UUID bytes", func(t *testing.T) {
		bytes := uuid.Nil.Bytes()
		u, err := FromByteArray(bytes[:])
		assert.NoError(t, err)
		assert.Equal(t, uuid.Nil, u)
	})

	t.Run("invalid byte array length", func(t *testing.T) {
		invalidLengths := [][]byte{
			make([]byte, 0),  // Empty
			make([]byte, 8),  // Too short
			make([]byte, 32), // Too long
			make([]byte, 15), // Almost right
		}

		for _, bytes := range invalidLengths {
			_, err := FromByteArray(bytes)
			assert.Error(t, err)
		}
	})

	t.Run("valid 16-byte array", func(t *testing.T) {
		bytes := make([]byte, 16)
		for i := range bytes {
			bytes[i] = byte(i)
		}

		u, err := FromByteArray(bytes)
		assert.NoError(t, err)
		assert.Equal(t, bytes, u.Bytes()[:])
	})
}

func TestToString(t *testing.T) {
	t.Run("UUID to hex string", func(t *testing.T) {
		u := NewUuid()
		hexStr := ToString(u)

		// Should be valid hex
		_, err := hex.DecodeString(hexStr)
		assert.NoError(t, err)

		// Should be 32 characters (16 bytes * 2 hex chars per byte)
		assert.Equal(t, 32, len(hexStr))

		// Should be lowercase
		assert.Equal(t, strings.ToLower(hexStr), hexStr)
	})

	t.Run("nil UUID to hex", func(t *testing.T) {
		hexStr := ToString(uuid.Nil)
		assert.Equal(t, "00000000000000000000000000000000", hexStr)
	})

	t.Run("max UUID to hex", func(t *testing.T) {
		maxUUID := uuid.FromStringOrNil("ffffffff-ffff-ffff-ffff-ffffffffffff")
		hexStr := ToString(maxUUID)
		assert.Equal(t, "ffffffffffffffffffffffffffffffff", hexStr)
	})

	t.Run("roundtrip with FromByteArray", func(t *testing.T) {
		original := NewUuid()
		hexStr := ToString(original)

		bytes, err := hex.DecodeString(hexStr)
		require.NoError(t, err)

		recovered, err := FromByteArray(bytes)
		assert.NoError(t, err)
		assert.Equal(t, original, recovered)
	})
}

func TestMarshall(t *testing.T) {
	tests := []struct {
		name  string
		input any
		check func(t *testing.T, result []byte)
	}{
		{
			name:  "simple string",
			input: "hello world",
			check: func(t *testing.T, result []byte) {
				var s string
				err := json.Unmarshal(result, &s)
				assert.NoError(t, err)
				assert.Equal(t, "hello world", s)
			},
		},
		{
			name:  "struct",
			input: struct{ Name string }{Name: "test"},
			check: func(t *testing.T, result []byte) {
				var s struct{ Name string }
				err := json.Unmarshal(result, &s)
				assert.NoError(t, err)
				assert.Equal(t, "test", s.Name)
			},
		},
		{
			name:  "map",
			input: map[string]int{"a": 1, "b": 2},
			check: func(t *testing.T, result []byte) {
				var m map[string]int
				err := json.Unmarshal(result, &m)
				assert.NoError(t, err)
				assert.Equal(t, 1, m["a"])
				assert.Equal(t, 2, m["b"])
			},
		},
		{
			name:  "slice",
			input: []int{1, 2, 3},
			check: func(t *testing.T, result []byte) {
				var s []int
				err := json.Unmarshal(result, &s)
				assert.NoError(t, err)
				assert.Equal(t, []int{1, 2, 3}, s)
			},
		},
		{
			name:  "nil value",
			input: nil,
			check: func(t *testing.T, result []byte) {
				assert.Equal(t, []byte("null"), result)
			},
		},
		{
			name:  "UUID",
			input: uuid.Nil,
			check: func(t *testing.T, result []byte) {
				var u uuid.UUID
				err := json.Unmarshal(result, &u)
				assert.NoError(t, err)
				assert.Equal(t, uuid.Nil, u)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, err := Marshall(tt.input)
			assert.NoError(t, err)

			bytes, ok := value.([]byte)
			assert.True(t, ok, "Expected []byte from Marshall")

			tt.check(t, bytes)
		})
	}
}

func TestUnmarshal(t *testing.T) {
	type TestStruct struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	t.Run("unmarshal from string", func(t *testing.T) {
		jsonStr := `{"name":"test","value":42}`
		var result TestStruct

		err := Unmarshal(jsonStr, &result)
		assert.NoError(t, err)
		assert.Equal(t, "test", result.Name)
		assert.Equal(t, 42, result.Value)
	})

	t.Run("unmarshal from []uint8", func(t *testing.T) {
		jsonData := []uint8(`{"name":"test2","value":100}`)
		var result TestStruct

		err := Unmarshal(jsonData, &result)
		assert.NoError(t, err)
		assert.Equal(t, "test2", result.Name)
		assert.Equal(t, 100, result.Value)
	})

	t.Run("unmarshal array", func(t *testing.T) {
		jsonStr := `[1,2,3,4,5]`
		var result []int

		err := Unmarshal(jsonStr, &result)
		assert.NoError(t, err)
		assert.Equal(t, []int{1, 2, 3, 4, 5}, result)
	})

	t.Run("unmarshal map", func(t *testing.T) {
		jsonStr := `{"key1":"value1","key2":"value2"}`
		var result map[string]string

		err := Unmarshal(jsonStr, &result)
		assert.NoError(t, err)
		assert.Equal(t, "value1", result["key1"])
		assert.Equal(t, "value2", result["key2"])
	})

	t.Run("invalid JSON", func(t *testing.T) {
		invalidJSON := "not valid json"
		var result TestStruct

		err := Unmarshal(invalidJSON, &result)
		assert.Error(t, err)
	})

	t.Run("type mismatch", func(t *testing.T) {
		jsonStr := `{"name":"test","value":"not a number"}`
		var result TestStruct

		err := Unmarshal(jsonStr, &result)
		assert.Error(t, err)
	})

	t.Run("unsupported type", func(t *testing.T) {
		var result TestStruct

		err := Unmarshal(123, &result) // int is not supported
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "type assertion")

		err = Unmarshal([]int{1, 2, 3}, &result) // []int is not supported
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "type assertion")
	})

	t.Run("null value", func(t *testing.T) {
		jsonStr := "null"
		var result *TestStruct

		err := Unmarshal(jsonStr, &result)
		assert.NoError(t, err)
		assert.Nil(t, result)
	})

	t.Run("empty string", func(t *testing.T) {
		var result TestStruct
		err := Unmarshal("", &result)
		assert.Error(t, err)
	})

	t.Run("complex nested struct", func(t *testing.T) {
		type Nested struct {
			Items []string          `json:"items"`
			Meta  map[string]string `json:"meta"`
		}
		type Complex struct {
			ID     uuid.UUID `json:"id"`
			Nested Nested    `json:"nested"`
		}

		id := NewUuid()
		original := Complex{
			ID: id,
			Nested: Nested{
				Items: []string{"a", "b", "c"},
				Meta:  map[string]string{"key": "value"},
			},
		}

		// Marshall it first
		marshalled, err := Marshall(original)
		require.NoError(t, err)

		// Now unmarshal
		var result Complex
		err = Unmarshal(string(marshalled.([]byte)), &result)
		assert.NoError(t, err)
		assert.Equal(t, original.ID, result.ID)
		assert.Equal(t, original.Nested.Items, result.Nested.Items)
		assert.Equal(t, original.Nested.Meta, result.Nested.Meta)
	})
}

// Benchmarks
func BenchmarkNewUuid(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NewUuid()
	}
}

func BenchmarkToBase62(b *testing.B) {
	u := NewUuid()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ToBase62(u)
	}
}

func BenchmarkFromBase62(b *testing.B) {
	u := NewUuid()
	base62 := ToBase62(u)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = FromBase62(base62)
	}
}

func BenchmarkMarshall(b *testing.B) {
	data := map[string]string{"key": "value", "foo": "bar"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Marshall(data)
	}
}

func BenchmarkUnmarshal(b *testing.B) {
	jsonStr := `{"key":"value","foo":"bar"}`
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var result map[string]string
		_ = Unmarshal(jsonStr, &result)
	}
}
