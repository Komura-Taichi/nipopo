package usecase

import "context"

type Tag struct {
	ID   string
	Name string
}

type TagsPage struct {
	Items      []Tag
	NextCursor string
}

type CreateTagResult struct {
	Tag     Tag
	Created bool
}

type TagsLister interface {
	List(ctx context.Context, q string, limit int, cursor string) (TagsPage, error)
}

type TagCreator interface {
	Create(ctx context.Context, name string) (CreateTagResult, error)
}
