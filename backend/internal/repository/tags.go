package repository

import (
	"context"

	"github.com/Komura-Taichi/nipopo/backend/internal/entity"
)

type TagRepository interface {
	FindByName(ctx context.Context, name string) (entity.Tag, bool, error)
	Create(ctx context.Context, name string) (entity.Tag, error)
	List(ctx context.Context, q string, limit int, cursor string) (entity.TagsPage, error)
}
