package usecase_test

import (
	"context"
	"testing"

	"github.com/Komura-Taichi/nipopo/backend/internal/entity"
	"github.com/Komura-Taichi/nipopo/backend/internal/usecase"
)

// usecaseのCreateが呼び出しうるメソッドのみ定義
type mockTagRepoForCreate struct {
	findCalled int
	findUserID string
	findName   string

	createCalled bool
	createUserID string
	createName   string

	// FindByNameは多重で追加しようとした場合に2回呼ばれる
	findTags   []entity.Tag
	findFounds []bool
	findErrs   []error

	createTag entity.Tag
	createErr error
}

// usecaseのListが呼び出しうるメソッドのみ定義
type mockTagRepoForList struct {
	listCalled bool
	listUserID string
	listQ      string
	listLimit  int
	listCursor string

	listPage entity.TagsPage
	listErr  error
}

func (m *mockTagRepoForCreate) FindByName(ctx context.Context, userID, name string) (entity.Tag, bool, error) {
	_ = ctx
	m.findCalled++
	m.findUserID, m.findName = userID, name

	i := m.findCalled - 1
	// 入れるデータが足りない場合
	if i >= len(m.findTags) {
		return entity.Tag{}, false, nil
	}

	return m.findTags[i], m.findFounds[i], m.findErrs[i]
}

func (m *mockTagRepoForCreate) Create(ctx context.Context, userID, name string) (entity.Tag, error) {
	_ = ctx
	m.createCalled = true
	m.createUserID, m.createName = userID, name
	return m.createTag, m.createErr
}

// インタフェースを満たすためのダミーメソッド
func (m *mockTagRepoForCreate) List(ctx context.Context, userID, q string, limit int, cursor string) (entity.TagsPage, error) {
	return entity.TagsPage{}, nil
}

func (m *mockTagRepoForList) List(ctx context.Context, userID, q string, limit int, cursor string) (entity.TagsPage, error) {
	_ = ctx
	m.listCalled = true
	m.listUserID, m.listQ, m.listLimit, m.listCursor = userID, q, limit, cursor
	return m.listPage, m.listErr
}

// --- インタフェースを満たすためのダミーメソッド ---
func (m *mockTagRepoForList) FindByName(ctx context.Context, userID, name string) (entity.Tag, bool, error) {
	return entity.Tag{}, false, nil
}

func (m *mockTagRepoForList) Create(ctx context.Context, userID, name string) (entity.Tag, error) {
	return entity.Tag{}, nil
}

// -- ここまで ---

func TestTagCreator_Create(t *testing.T) {
	const (
		userID = "u1"
		name   = "タグ1"
	)
	t.Run("OK_new", func(t *testing.T) {
		m := &mockTagRepoForCreate{
			// FindByNameで見つからない
			findTags:   []entity.Tag{{}},
			findFounds: []bool{false},
			findErrs:   []error{nil},

			// Createが成功
			createTag: entity.Tag{ID: "t1", UserID: userID, Name: name},
			createErr: nil,
		}

		tagUsecase := usecase.NewTagUsecase(m)

		got, err := tagUsecase.Create(context.Background(), userID, name)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !got.Created {
			t.Fatalf("Created should be true")
		}
		if got.Tag.ID != "t1" || got.Tag.Name != name || got.Tag.UserID != userID {
			t.Fatalf("tag mismatch: got=%+v", got.Tag)
		}

		// リポジトリ層のメソッド呼び出しに関する確認
		if m.findCalled != 1 {
			t.Fatalf("FindByName called %d times, want 1", m.findCalled)
		}
		if m.findUserID != userID || m.findName != name {
			t.Fatalf("FindByName args mismatch: userID=%q name=%q", m.findUserID, m.findName)
		}

		if !m.createCalled {
			t.Fatalf("Create should be called")
		}
		if m.createUserID != userID || m.createName != name {
			t.Fatalf("Create args mismatch: userID=%q name=%q", m.createUserID, m.createName)
		}
	})
}
