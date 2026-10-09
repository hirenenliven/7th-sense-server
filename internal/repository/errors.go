package repository

import "errors"

// ErrNotFound is returned when no record matches the lookup.
var ErrNotFound = errors.New("record not found")
