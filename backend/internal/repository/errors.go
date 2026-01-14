package repository

import "errors"

var ErrAlreadyTagExists = errors.New("tag already exists")
var ErrCursorNotFound = errors.New("cursor not found")
var ErrEmptyTagName = errors.New("empty tag name")
