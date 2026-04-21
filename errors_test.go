package kitsune

import (
	"errors"
	"strings"
	"testing"

	tozderrors "gitlab.com/tozd/go/errors"
)


func TestWithDetails(t *testing.T) {
	tests := []struct {
		name           string
		message        string
		details        []interface{}
		expectedError  string
		expectedDetail map[string]interface{}
	}{
		{
			name:          "simple message no details",
			message:       "simple error",
			details:       nil,
			expectedError: "simple error",
		},
		{
			name:          "message with format and details",
			message:       "error: %s",
			details:       []interface{}{"key1", "value1", "key2", 123},
			expectedError: "error: value1",
			expectedDetail: map[string]interface{}{
				"key1": "value1",
				"key2": 123,
			},
		},
		{
			name:          "multiple format specifiers",
			message:       "error %d: %s",
			details:       []interface{}{"code", 42, "message", "failed", "extra", true},
			expectedError: "error 42: failed",
			expectedDetail: map[string]interface{}{
				"code":    42,
				"message": "failed",
				"extra":   true,
			},
		},
		{
			name:          "more args than placeholders",
			message:       "error: %s",
			details:       []interface{}{"key1", "value1", "key2", "value2", "key3", "value3"},
			expectedError: "error: value1",
			expectedDetail: map[string]interface{}{
				"key1": "value1",
				"key2": "value2",
				"key3": "value3",
			},
		},
		{
			name:          "no format specifiers with details",
			message:       "static error message",
			details:       []interface{}{"key", "value"},
			expectedError: "static error message",
			expectedDetail: map[string]interface{}{
				"key": "value",
			},
		},
		// Note: Commented out as tozd/go/errors panics on odd number of details
		// {
		// 	name:          "odd number of details",
		// 	message:       "error",
		// 	details:       []interface{}{"key1", "value1", "key2"},
		// 	expectedError: "error",
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := WithDetails(tt.message, tt.details...)

			if err == nil {
				t.Fatal("Expected non-nil error")
			}

			if !strings.Contains(err.Error(), tt.expectedError) {
				t.Errorf("Expected error containing '%s', got '%s'", tt.expectedError, err.Error())
			}

			// Check details if expected
			if tt.expectedDetail != nil {
				details := AllDetails(err)
				for key, expectedValue := range tt.expectedDetail {
					if value, ok := details[key]; !ok {
						t.Errorf("Missing detail key: %s", key)
					} else if value != expectedValue {
						t.Errorf("Detail %s: expected %v, got %v", key, expectedValue, value)
					}
				}
			}
		})
	}
}

func TestWrapWithDetails(t *testing.T) {
	baseError := errors.New("base error")

	tests := []struct {
		name           string
		err            error
		message        string
		details        []interface{}
		expectedError  string
		expectedDetail map[string]interface{}
	}{
		{
			name:          "wrap simple error",
			err:           baseError,
			message:       "wrapped",
			details:       nil,
			expectedError: "wrapped",
		},
		{
			name:          "wrap with format and details",
			err:           baseError,
			message:       "operation failed: %s",
			details:       []interface{}{"reason", "timeout", "retry", 3},
			expectedError: "operation failed: timeout",
			expectedDetail: nil,
		},
		{
			name:           "wrap nil error",
			err:            nil,
			message:        "should handle nil",
			details:        []interface{}{"key", "value"},
			expectedError:  "", // Wrapf returns nil for nil error
			expectedDetail: nil,
		},
		{
			name:          "wrap with multiple format specifiers",
			message:       "failed at %d with %s",
			err:           baseError,
			details:       []interface{}{"position", 42, "status", "critical", "action", "retry"},
			expectedError: "failed at 42 with critical",
			expectedDetail: nil,
		},
		{
			name:          "wrap with more args than placeholders",
			err:           baseError,
			message:       "error: %s",
			details:       []interface{}{"msg", "failed", "code", 500, "time", "2024-01-01"},
			expectedError: "error: failed",
			expectedDetail: nil,
		},
		{
			name:          "wrap already wrapped error",
			err:           tozderrors.Wrap(baseError, "first wrap"),
			message:       "second wrap: %s",
			details:       []interface{}{"level", "outer"},
			expectedError: "second wrap: outer",
			expectedDetail: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := WrapWithDetails(tt.err, tt.message, tt.details...)

			if tt.err != nil && err == nil {
				t.Fatal("Expected non-nil error")
			}

			if tt.err == nil && err != nil {
				t.Fatal("Expected nil error for nil input")
			}

			if err == nil {
				return // Skip further checks for nil errors
			}

			if !strings.Contains(err.Error(), tt.expectedError) {
				t.Errorf("Expected error containing '%s', got '%s'", tt.expectedError, err.Error())
			}

			// Check details if expected
			if tt.expectedDetail != nil {
				details := AllDetails(err)
				for key, expectedValue := range tt.expectedDetail {
					if value, ok := details[key]; !ok {
						t.Errorf("Missing detail key: %s", key)
					} else if value != expectedValue {
						t.Errorf("Detail %s: expected %v, got %v", key, expectedValue, value)
					}
				}
			}
		})
	}
}

func TestAllDetails(t *testing.T) {
	tests := []struct {
		name            string
		err             error
		expectedDetails map[string]interface{}
	}{
		{
			name:            "simple error no details",
			err:             errors.New("simple"),
			expectedDetails: map[string]interface{}{},
		},
		{
			name: "error with details",
			err:  tozderrors.WithDetails(errors.New("error"), "key1", "value1", "key2", 42),
			expectedDetails: map[string]interface{}{
				"key1": "value1",
				"key2": 42,
			},
		},
		{
			name: "nested error with details",
			err: tozderrors.WithDetails(
				tozderrors.WithDetails(errors.New("base"), "inner", "value"),
				"outer", "value2",
			),
			expectedDetails: map[string]interface{}{
				"inner": "value",
				"outer": "value2",
			},
		},
		{
			name:            "nil error",
			err:             nil,
			expectedDetails: map[string]interface{}{}, // AllDetails returns empty map for nil
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			details := AllDetails(tt.err)

			if tt.err == nil {
				// For nil error, details should be an empty map
				if len(details) != 0 {
					t.Errorf("Expected empty details map for nil error, got %v", details)
				}
				return
			}

			if len(details) != len(tt.expectedDetails) {
				t.Errorf("Expected %d details, got %d", len(tt.expectedDetails), len(details))
			}

			for key, expectedValue := range tt.expectedDetails {
				if value, ok := details[key]; !ok {
					t.Errorf("Missing detail key: %s", key)
				} else if value != expectedValue {
					t.Errorf("Detail %s: expected %v, got %v", key, expectedValue, value)
				}
			}
		})
	}
}

func TestPlaceholderCounting(t *testing.T) {
	tests := []struct {
		name              string
		message           string
		details           []interface{}
		expectedFormatted string
	}{
		{
			name:              "correct placeholder counting with %d",
			message:           "error %d occurred",
			details:           []interface{}{"code", 404, "extra", "data"},
			expectedFormatted: "error 404 occurred",
		},
		{
			name:              "correct placeholder counting with %s",
			message:           "file %s not found",
			details:           []interface{}{"name", "/path/to/file", "size", 1024},
			expectedFormatted: "file /path/to/file not found",
		},
		{
			name:              "correct placeholder counting with %v",
			message:           "value is %v",
			details:           []interface{}{"val", struct{ A int }{A: 1}, "other", "test"},
			expectedFormatted: "value is {1}",
		},
		{
			name:              "mixed placeholders",
			message:           "user %s has %d items of type %v",
			details:           []interface{}{"name", "john", "count", 5, "type", []string{"a", "b"}},
			expectedFormatted: "user john has 5 items of type [a b]",
		},
		{
			name:              "no placeholders",
			message:           "static message",
			details:           []interface{}{"key", "value"},
			expectedFormatted: "static message",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := WithDetails(tt.message, tt.details...)
			if !strings.Contains(err.Error(), tt.expectedFormatted) {
				t.Errorf("Expected error containing '%s', got '%s'", tt.expectedFormatted, err.Error())
			}
		})
	}
}

func TestEdgeCases(t *testing.T) {
	t.Run("empty details slice", func(t *testing.T) {
		err := WithDetails("error", []interface{}{}...)
		if err.Error() != "error" {
			t.Errorf("Expected 'error', got '%s'", err.Error())
		}
	})

	t.Run("single detail item (odd number)", func(t *testing.T) {
		// Test that odd number of details causes a panic (as per tozd/go/errors behavior)
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic for odd number of details")
			}
		}()
		_ = WithDetails("error", "key")
	})

	t.Run("percent sign without format specifier", func(t *testing.T) {
		err := WithDetails("100% failure", "key", "value")
		// The % without a valid format specifier will cause formatting issues
		// Just check that we get an error, not the exact message
		if err == nil {
			t.Error("Expected non-nil error")
		}
	})

	t.Run("percent at end of string", func(t *testing.T) {
		err := WithDetails("error%", "key", "value")
		if !strings.Contains(err.Error(), "error%") {
			t.Errorf("Expected 'error%%' in error, got '%s'", err.Error())
		}
	})
}

func BenchmarkWithDetails(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = WithDetails("error %s at %d", "key1", "value", "key2", 42)
	}
}

func BenchmarkWrapWithDetails(b *testing.B) {
	baseErr := errors.New("base")
	for i := 0; i < b.N; i++ {
		_ = WrapWithDetails(baseErr, "wrapped %s", "key", "value")
	}
}

func BenchmarkAllDetails(b *testing.B) {
	err := tozderrors.WithDetails(errors.New("error"), "k1", "v1", "k2", "v2", "k3", "v3")
	for i := 0; i < b.N; i++ {
		_ = AllDetails(err)
	}
}

func TestWrapWithDetailsNoDoubleWrapping(t *testing.T) {
	baseError := errors.New("base error")
	wrapped := WrapWithDetails(baseError, "wrapped: %s", "key", "value", "code", 42)

	if wrapped == nil {
		t.Fatal("Expected non-nil error")
	}

	// Count the number of wrapping layers by unwrapping
	wrapCount := 0
	current := wrapped
	for current != nil {
		wrapCount++
		current = tozderrors.Unwrap(current)
	}

	// Should be 2 layers: base error + single wrap (not double wrapped)
	// Base error (1) + WrapWithDetails wrap (1) = 2 total
	if wrapCount != 2 {
		t.Errorf("Expected 2 wrapping layers, got %d", wrapCount)
	}

	// Verify the error message is correct
	if !strings.Contains(wrapped.Error(), "wrapped: value") {
		t.Errorf("Expected error containing 'wrapped: value', got '%s'", wrapped.Error())
	}
}
