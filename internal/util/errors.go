package util

import "errors"

var (
	ErrUnknownArg = errors.New("unknown argument passed")

	ErrInvalidArg = errors.New("invalid argument passed")
)
