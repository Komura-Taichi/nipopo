package usecase_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/Komura-Taichi/nipopo/backend/internal/entity"
	"github.com/Komura-Taichi/nipopo/backend/internal/repository"
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

// --- インタフェースを満たすためのダミーメソッド ---
func (m *mockTagRepoForCreate) List(ctx context.Context, userID, q string, limit int, cursor string) (entity.TagsPage, error) {
	return entity.TagsPage{}, nil
}

// --- ここまで ---

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

// --- ここまで ---

func TestTagCreator_Create(t *testing.T) {
	const (
		userID = "u1"
		name   = "タグ1"
	)
	// --- 正常系のテスト ---
	t.Run("OK_new_tag", func(t *testing.T) {
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
		assertUnexpectedError(t, err)
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

	t.Run("OK_existing_tag", func(t *testing.T) {
		existingTag := entity.Tag{ID: "t_1", UserID: userID, Name: name}

		m := &mockTagRepoForCreate{
			// FindByNameで見つかる
			findTags:   []entity.Tag{existingTag},
			findFounds: []bool{true},
			findErrs:   []error{nil},
		}

		tagUsecase := usecase.NewTagUsecase(m)

		got, err := tagUsecase.Create(context.Background(), userID, name)
		assertUnexpectedError(t, err)
		if got.Created {
			t.Fatalf("Created should be false")
		}
		if got.Tag.ID != existingTag.ID || got.Tag.Name != existingTag.Name || got.Tag.UserID != existingTag.UserID {
			t.Fatalf("tag mismatch: got=%+v want=%+v", got.Tag, existingTag)
		}

		// リポジトリ層のメソッド呼び出しに関する確認
		if m.findCalled != 1 {
			t.Fatalf("FindByName called %d times, want 1", m.findCalled)
		}
		if m.findUserID != existingTag.UserID || m.findName != existingTag.Name {
			t.Fatalf("FindByName args mismatch: userID=%q name=%q", m.findUserID, m.findName)
		}

		// 既存ならCreateは呼ばれないはず
		if m.createCalled {
			t.Fatalf("Create should not be called for existing tag")
		}
	})

	t.Run("OK_conflict_request", func(t *testing.T) {
		existingTag := entity.Tag{ID: "t_1", UserID: userID, Name: name}

		m := &mockTagRepoForCreate{
			// 1回目のFindByNameで見つからず、
			// 2回目のFindByNameで見つかる
			findTags:   []entity.Tag{{}, existingTag},
			findFounds: []bool{false, true},
			findErrs:   []error{nil, nil},

			// Create時に既に存在
			createErr: repository.ErrAlreadyTagExists,
		}

		tagUsecase := usecase.NewTagUsecase(m)

		got, err := tagUsecase.Create(context.Background(), existingTag.UserID, existingTag.Name)
		assertUnexpectedError(t, err)

		if got.Created {
			t.Fatalf("Created should be false")
		}
		// Createの挙動としては、1回目FindByNameで見つからない -> CreateでErr -> もう一度FindByNameで見つかる -> 見つかったTagを返す が理想
		if got.Tag.ID != existingTag.ID || got.Tag.Name != existingTag.Name || got.Tag.UserID != existingTag.UserID {
			t.Fatalf("tag mismatch: got=%+v want=%+v", got.Tag, existingTag)
		}

		// リポジトリ層のメソッド呼び出しに関する確認
		if m.findCalled != 2 {
			t.Fatalf("FindByName called %d times, want 2", m.findCalled)
		}
		if !m.createCalled {
			t.Fatalf("Create should be called")
		}
	})

	// --- 異常系のテスト ---
	t.Run("Err_empty_tag_name", func(t *testing.T) {
		m := &mockTagRepoForCreate{}

		tagUsecase := usecase.NewTagUsecase(m)

		_, err := tagUsecase.Create(context.Background(), userID, "")
		if !errors.Is(err, usecase.ErrEmptyTagName) {
			t.Fatalf("unexpected error: %v", err)
		}

		// repositoryのFindByNameもCreateも呼ばれないはず
		if m.findCalled != 0 {
			t.Fatalf("FindByName should not be called, but called %d times", m.findCalled)
		}
		if m.createCalled {
			t.Fatal("Create should not be called")
		}
	})

	t.Run("Err_find_error", func(t *testing.T) {
		intentionalFindErr := errors.New("FindByName error")

		m := &mockTagRepoForCreate{
			// 1回目のFindByNameでエラーが出る
			findTags:   []entity.Tag{{}},
			findFounds: []bool{false},
			findErrs:   []error{intentionalFindErr},
		}

		tagUsecase := usecase.NewTagUsecase(m)

		_, err := tagUsecase.Create(context.Background(), userID, name)
		// 未知のエラーであれば、そのままエラーを返すはず
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, intentionalFindErr) {
			t.Fatalf("unexpected error: got=%v want=%v", err, intentionalFindErr)
		}

		// repositoryのFindByNameは1度だけ呼ばれるはず
		if m.findCalled != 1 {
			t.Fatalf("FindByName should be called once, but called %d times", m.findCalled)
		}
		// repositoryのCreateは呼ばれないはず
		if m.createCalled {
			t.Fatal("Create should not be called when FindByName returns error")
		}
	})

	t.Run("Err_create_error", func(t *testing.T) {
		intentionalCreateErr := errors.New("Create error")

		m := &mockTagRepoForCreate{
			// 1回目のFindByNameで見つからない
			// 2回目のFindByNameは実行されない
			findTags:   []entity.Tag{{}},
			findFounds: []bool{false},
			findErrs:   []error{nil},

			// Create時に未知のエラーが発生
			createErr: intentionalCreateErr,
		}

		tagUsecase := usecase.NewTagUsecase(m)

		_, err := tagUsecase.Create(context.Background(), userID, name)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, intentionalCreateErr) {
			t.Fatalf("unexpected error: got=%v want=%v", err, intentionalCreateErr)
		}

		// repositoryのFindByNameは1度だけ呼ばれるはず
		if m.findCalled != 1 {
			t.Fatalf("FindByName should be called once, but called %d times", m.findCalled)
		}
		// repositoryのCreateは呼ばれるはず
		if !m.createCalled {
			t.Fatalf("Create should be called")
		}
	})

	t.Run("Err_conflict_refind_error", func(t *testing.T) {
		intentionalRefindError := errors.New("refind error")

		m := &mockTagRepoForCreate{
			// 1回目のFindByNameで見つからない
			// 2回目のFindByNameでerror発生
			findTags:   []entity.Tag{{}, {}},
			findFounds: []bool{false, false},
			findErrs:   []error{nil, intentionalRefindError},

			// Create時に既に存在
			createErr: repository.ErrAlreadyTagExists,
		}

		tagUsecase := usecase.NewTagUsecase(m)

		_, err := tagUsecase.Create(context.Background(), userID, name)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, intentionalRefindError) {
			t.Fatalf("unexpected error: got=%v want=%v", err, intentionalRefindError)
		}

		// repositoryのFindByNameは2回呼ばれるはず
		if m.findCalled != 2 {
			t.Fatalf("FindByName should be called 2 times, but called %d times", m.findCalled)
		}
		// repositoryのCreateは呼ばれるはず
		if !m.createCalled {
			t.Fatalf("Create should be called")
		}
	})

	t.Run("Err_conflict_refind_notfound_error", func(t *testing.T) {
		m := &mockTagRepoForCreate{
			// 1回目のFindByNameで見つからない
			// 2回目のFindByNameでも見つからない (矛盾)
			findTags:   []entity.Tag{{}, {}},
			findFounds: []bool{false, false},
			findErrs:   []error{nil, nil},

			// Create時に既に存在
			createErr: repository.ErrAlreadyTagExists,
		}

		tagUsecase := usecase.NewTagUsecase(m)

		_, err := tagUsecase.Create(context.Background(), userID, name)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, usecase.ErrContradictoryRepoState) {
			t.Fatalf("unexpected error: got=%v want=%v", err, usecase.ErrContradictoryRepoState)
		}

		// repositoryのFindByNameは2回呼ばれるはず
		if m.findCalled != 2 {
			t.Fatalf("FindByName should be called 2 times, but called %d times", m.findCalled)
		}
		// repositoryのCreateは呼ばれるはず
		if !m.createCalled {
			t.Fatalf("Create should be called")
		}
	})
}

func TestTagsLister_List(t *testing.T) {
	const (
		userID = "u1"
		limit  = 20
	)
	t.Run("OK_first_page", func(t *testing.T) {
		const (
			q      = ""
			cursor = ""
		)

		want := entity.TagsPage{
			Items: []entity.Tag{
				{ID: "t_1", UserID: userID, Name: "bar"},
				{ID: "t_2", UserID: userID, Name: "foo"},
			},
			NextCursor: "t_2",
		}

		m := &mockTagRepoForList{
			listPage: want,
			listErr:  nil,
		}

		tagUsecase := usecase.NewTagUsecase(m)

		got, err := tagUsecase.List(context.Background(), userID, q, limit, cursor)
		assertUnexpectedError(t, err)

		// タグ一覧やNextCursorが完全に一致するか
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("page mismatch: got=%+v want=%+v", got, want)
		}

		// repositoryのListは呼ばれるはず
		if !m.listCalled {
			t.Fatal("List should be called")
		}
		// repositoryのListに渡される引数は正しいか？
		if m.listUserID != userID || m.listQ != q || m.listLimit != limit || m.listCursor != cursor {
			t.Fatalf("List args mismatch: userID=%q q=%q limit=%d cursor=%q",
				m.listUserID, m.listQ, m.listLimit, m.listCursor)
		}
	})

	t.Run("OK_with_query_and_cursor", func(t *testing.T) {
		const (
			q      = "foo"
			cursor = "t_2"
		)

		want := entity.TagsPage{
			Items: []entity.Tag{
				{ID: "t_2", UserID: userID, Name: "foo"},
				{ID: "t_3", UserID: userID, Name: "foobar"},
			},
			NextCursor: "t_3",
		}

		m := &mockTagRepoForList{
			listPage: want,
			listErr:  nil,
		}

		tagUsecase := usecase.NewTagUsecase(m)

		got, err := tagUsecase.List(context.Background(), userID, q, limit, cursor)
		assertUnexpectedError(t, err)

		// タグ一覧やNextCursorが完全に一致するか
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("page mismatch: got=%+v want=%+v", got, want)
		}

		// repositoryのListは呼ばれるはず
		if !m.listCalled {
			t.Fatal("List should be called")
		}
		// repositoryのListに渡される引数は正しいか？
		if m.listUserID != userID || m.listQ != q || m.listLimit != limit || m.listCursor != cursor {
			t.Fatalf("List args mismatch: userID=%q q=%q limit=%d cursor=%q",
				m.listUserID, m.listQ, m.listLimit, m.listCursor)
		}
	})

	t.Run("Err_list_error", func(t *testing.T) {
		const (
			q      = "foo"
			cursor = "t_2"
		)

		intentionalErr := errors.New("List error")

		m := &mockTagRepoForList{
			listPage: entity.TagsPage{},
			listErr:  intentionalErr,
		}

		tagUsecase := usecase.NewTagUsecase(m)

		_, err := tagUsecase.List(context.Background(), userID, q, limit, cursor)
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		// 発生したエラーがそのまま返ってくるか
		if !errors.Is(err, intentionalErr) {
			t.Fatalf("unexpected error: got=%v want=%v", err, intentionalErr)
		}

		// repositoryのListは呼ばれるはず
		if !m.listCalled {
			t.Fatal("List should be called")
		}
	})

	t.Run("Err_invalid_cursor_error", func(t *testing.T) {
		const (
			q      = "foo"
			cursor = "foo" // 不正なcursor (空文字列でもt_<連番>でもない)
		)

		m := &mockTagRepoForList{}

		tagUsecase := usecase.NewTagUsecase(m)

		_, err := tagUsecase.List(context.Background(), userID, q, limit, cursor)
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		// ErrInvalidCursorエラーが返るはず
		if !errors.Is(err, usecase.ErrInvalidCursor) {
			t.Fatalf("unexpected error: got=%v want=%v", err, usecase.ErrInvalidCursor)
		}

		// repositoryのListは呼ばれないはず
		if m.listCalled {
			t.Fatal("List should not be called when cursor invalid")
		}
	})
}

func assertUnexpectedError(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
