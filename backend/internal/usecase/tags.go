package usecase

import (
	"context"
	"errors"
	"regexp"
	"strings"

	"github.com/Komura-Taichi/nipopo/backend/internal/entity"
	"github.com/Komura-Taichi/nipopo/backend/internal/repository"
)

var tagCursorRe = regexp.MustCompile(`^t_[1-9][0-9]*$`) // タグIDの形式

// TagsListerとTagCreatorの実体
type TagUsecase struct {
	repo repository.TagRepository
}

func NewTagUsecase(repo repository.TagRepository) *TagUsecase {
	return &TagUsecase{repo: repo}
}

func (t *TagUsecase) List(ctx context.Context, userID string, q string, limit int, cursor string) (entity.TagsPage, error) {
	if cursor != "" && !tagCursorRe.MatchString(cursor) {
		return entity.TagsPage{}, ErrInvalidCursor
	}

	// 不正なカーソルでなければ、repositoryのエラーもそのまま返す
	return t.repo.List(ctx, userID, q, limit, cursor)
}

func (t *TagUsecase) Create(ctx context.Context, userID, name string) (CreateTagResult, error) {
	// タグ名が空白または空文字列
	nameNoSpace := strings.TrimSpace(name)
	if nameNoSpace == "" {
		return CreateTagResult{}, ErrEmptyTagName
	}

	// 既存タグかどうか確認
	tag, found, err := t.repo.FindByName(ctx, userID, nameNoSpace)
	if err != nil {
		return CreateTagResult{}, err
	}
	if found {
		return CreateTagResult{
			Tag:     tag,
			Created: false,
		}, nil
	}

	// 既存でなければ作成
	createdTag, err := t.repo.Create(ctx, userID, nameNoSpace)
	if err == nil {
		return CreateTagResult{
			Tag:     createdTag,
			Created: true,
		}, nil
	}

	// 複数のリクエストが競合した場合
	if errors.Is(err, repository.ErrAlreadyTagExists) {
		// 再検索
		tag, found, refindErr := t.repo.FindByName(ctx, userID, nameNoSpace)
		if refindErr != nil {
			return CreateTagResult{}, refindErr
		}
		if !found {
			return CreateTagResult{}, ErrContradictoryRepoState
		}

		return CreateTagResult{
			Tag:     tag,
			Created: false,
		}, nil
	}

	// 未定義のCreate時エラー発生
	return CreateTagResult{}, err
}
