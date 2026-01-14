package repository

import "errors"

var (
	ErrAlreadyTagExists = errors.New("tag already exists")
	ErrCursorNotFound   = errors.New("cursor not found")
	ErrEmptyTagName     = errors.New("empty tag name")
)
