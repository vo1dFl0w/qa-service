package handlers

import "errors"

var (
	ErrMethodNotAllowed = errors.New("method not allowed")
	ErrNotFound = errors.New("not found")
	ErrBadRequest = errors.New("bad request")
)