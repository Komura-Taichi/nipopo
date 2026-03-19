package repository

import (
	"context"

	"github.com/Komura-Taichi/nipopo/backend/internal/entity"
)

type TagRepository interface {
	FindByIDs(ctx context.Context, userID string, ids []string) ([]entity.Tag, error)
	FindByName(ctx context.Context, userID string, name string) (entity.Tag, bool, error)
	Create(ctx context.Context, userID string, name string) (entity.Tag, error)
	List(ctx context.Context, userID string, q string, limit int, cursor string) (entity.TagsPage, error)
}
