package constants

import "errors"

var (
	// ErrNotFound is returned when a  requested resource not found
	ErrNotFound = errors.New("not found")
)
