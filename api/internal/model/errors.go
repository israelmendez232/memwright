package model

import "errors"

var (
	ErrNotFound     = errors.New("resource not found")
	ErrDuplicateKey = errors.New("duplicate key")
	ErrDuplicateEmail = errors.New("email already exists")
	ErrDuplicateName  = errors.New("name already exists")
	ErrInvalidInput   = errors.New("invalid input")
	ErrUnauthorized   = errors.New("unauthorized")
	ErrForbidden      = errors.New("forbidden")
)
