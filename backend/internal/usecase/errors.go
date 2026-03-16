package usecase

import "errors"

var (
	ErrEmptyTagName           = errors.New("empty tag name")
	ErrContradictoryRepoState = errors.New("contradictory repo state")
	ErrInvalidCursor          = errors.New("invalid cursor")
)

type RecordAlreadyExistsError struct {
	ExistingID string
}

func (e *RecordAlreadyExistsError) Error() string {
	return "record of that day already exists"
}
