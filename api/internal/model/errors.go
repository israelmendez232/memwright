package model

import "errors"

var (
	ErrNotFound     = errors.New("resource not found")
	ErrDuplicateKey = errors.New("duplicate key")
	ErrInvalidInput = errors.New("invalid input")
)
