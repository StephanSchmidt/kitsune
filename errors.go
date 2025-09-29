package kitsune

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gitlab.com/tozd/go/errors"
)

func LogError(err error) *zerolog.Event {
	return log.Error().Err(err).Fields(errors.AllDetails(err)).Err(err)
}

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
		errors.Wrapf(err, message, args...),
		details...,
	)
}

func AllDetails(err error) map[string]interface{} {
	return errors.AllDetails(err)
}
