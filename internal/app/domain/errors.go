package domain

import "errors"

var (
	ErrNotFound = errors.New("not found")
	ErrNoRowDeleted = errors.New("no row deleted")
	ErrEmptyText = errors.New("empty text")
)