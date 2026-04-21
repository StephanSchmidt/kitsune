package kitsune

import (
	"gitlab.com/tozd/go/errors"
)

func WithDetails(message string, details ...interface{}) error {
	// Extract every second element from details to use as args for Wrapf
	var args []interface{}
	for i := 1; i < len(details); i += 2 {
		args = append(args, details[i])
	}
	// Count % placeholders in message (excluding %w)
	placeholderCount := 0
	for i := 0; i < len(message)-1; i++ {
		if message[i] == '%' && (message[i+1] == 'd' || message[i+1] == 's' || message[i+1] == 'v') {
			placeholderCount++
		}
	}

	// Cut args to match placeholder count
	if len(args) > placeholderCount {
		args = args[:placeholderCount]
	}

	return errors.WithDetails(
		errors.Errorf(message, args...),
		details...,
	)
}

func WrapWithDetails(err error, message string, details ...interface{}) error {
	if err == nil {
		return nil
	}

	// Extract every second element from details to use as args for message formatting
	var args []interface{}
	for i := 1; i < len(details); i += 2 {
		args = append(args, details[i])
	}

	// Count % placeholders in message (excluding %w)
	placeholderCount := 0
	for i := 0; i < len(message)-1; i++ {
		if message[i] == '%' && (message[i+1] == 'd' || message[i+1] == 's' || message[i+1] == 'v') {
			placeholderCount++
		}
	}

	// Cut args to match placeholder count
	if len(args) > placeholderCount {
		args = args[:placeholderCount]
	}

	// Single wrap with formatted message (no double wrapping)
	if len(args) > 0 {
		return errors.Wrapf(err, message, args...)
	}

	return errors.Wrap(err, message)
}

func AllDetails(err error) map[string]interface{} {
	return errors.AllDetails(err)
}
