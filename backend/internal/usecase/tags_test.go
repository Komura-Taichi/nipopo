package usecase_test

import (
	"context"

	"github.com/Komura-Taichi/nipopo/backend/internal/entity"
)

type mockTagRepo struct {
	findCalled int
	findName   string

	createCalled bool
	createName   string

	listCalled bool
	listQ      string
	listLimit  int
	listCursor string

	// FindByNameは多重で追加しようとした場合に2回呼ばれる
	findTags   []entity.Tag
	findFounds []bool
	findErrs   []error

	createTag entity.Tag
	createErr error

	listPage entity.TagsPage
	listErr  error
}

func (m *mockTagRepo) FindByName(ctx context.Context, name string) (entity.Tag, bool, error) {
	m.findCalled++
	m.findName = name

	i := m.findCalled - 1
	// 入れるデータが足りない場合
	if i >= len(m.findTags) {
		return entity.Tag{}, false, nil
	}

	return m.findTags[i], m.findFounds[i], m.findErrs[i]
}

func (m *mockTagRepo) Create(ctx context.Context, name string) (entity.Tag, error) {
	m.createCalled = true
	m.createName = name
	return m.createTag, m.createErr
}

func (m *mockTagRepo) List(ctx context.Context, q string, limit int, cursor string) (entity.TagsPage, error) {
	m.listCalled = true
	m.listQ, m.listLimit, m.listCursor = q, limit, cursor
	return m.listPage, m.listErr
}
