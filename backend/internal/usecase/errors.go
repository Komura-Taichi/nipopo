package usecase

import "errors"

var (
	ErrEmptyTagName  = errors.New("empty tag name")
	ErrInvalidCursor = errors.New("invalid cursor")
)
