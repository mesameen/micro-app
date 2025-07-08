package repository

import "errors"

// ErrNotFound is returned when the requested data not found
var ErrNotFound = errors.New("not found")
