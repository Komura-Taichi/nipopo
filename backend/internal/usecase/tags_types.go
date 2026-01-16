package usecase

import (
	"context"

	"github.com/Komura-Taichi/nipopo/backend/internal/entity"
)

type CreateTagResult struct {
	Tag     entity.Tag
	Created bool
}

type TagsLister interface {
	List(ctx context.Context, userID string, q string, limit int, cursor string) (entity.TagsPage, error)
}

type TagCreator interface {
	Create(ctx context.Context, userID string, name string) (CreateTagResult, error)
}
