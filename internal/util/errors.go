package util

import "errors"

var (
	ErrUnknownArg = errors.New("unknown argument passed")

	ErrInvalidArg = errors.New("invalid argument passed")

	ErrDatabaseServiceUnavailable = errors.New("database service unavailable")
)

func FromString(s string) error {
	if s == "" {
		return nil
	}
	return errors.New(s)
}
