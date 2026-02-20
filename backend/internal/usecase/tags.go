package usecase

import (
	"context"

	"github.com/Komura-Taichi/nipopo/backend/internal/entity"
	"github.com/Komura-Taichi/nipopo/backend/internal/repository"
)

// TagsListerとTagCreatorの実体
type TagUsecase struct {
	repo repository.TagRepository
}

func NewTagUsecase(repo repository.TagRepository) *TagUsecase {
	return &TagUsecase{repo: repo}
}

func (t *TagUsecase) List(ctx context.Context, userID string, q string, limit int, cursor string) (entity.TagsPage, error) {
	// TODO: 実装
	return entity.TagsPage{}, nil
}

func (t *TagUsecase) Create(ctx context.Context, userID, name string) (CreateTagResult, error) {
	// TODO: 実装
	return CreateTagResult{}, nil
}
