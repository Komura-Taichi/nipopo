package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Komura-Taichi/nipopo/backend/internal/handler"
	"github.com/Komura-Taichi/nipopo/backend/internal/usecase"
)

type mockTagsLister struct {
	listCalled bool
	listQ      string
	listLimit  int
	listCursor string

	listResponse handler.TagsPage
	listErr      error
}

type mockTagCreator struct {
	createCalled bool
	createName   string

	createResponse usecase.CreateTagResult
	createErr      error
}

func (f *mockTagsLister) List(ctx context.Context, q string, limit int, cursor string) (handler.TagsPage, error) {
	f.listCalled = true
	f.listQ, f.listLimit, f.listCursor = q, limit, cursor
	return f.listResponse, f.listErr
}

func (f *mockTagCreator) Create(ctx context.Context, name string) (usecase.CreateTagResult, error) {
	f.createCalled = true
	f.createName = name
	return f.createResponse, f.createErr
}

func TestListTags(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		m := &mockTagsLister{
			listResponse: handler.TagsPage{
				Items: []handler.Tag{
					{ID: "t1", Name: "タグ1"},
					{ID: "t2", Name: "タグ2"},
				},
				NextCursor: "next123",
			},
		}

		h := handler.ListTags(m)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/v1/tags?q=タグ&limit=5&cursor=cur123", nil)
		h.ServeHTTP(rec, req)

		// ステータスコードとヘッダの確認
		assertStatus(t, rec, http.StatusOK)
		assertContentType(t, rec)

		// 呼び出し確認
		if !m.listCalled {
			t.Fatalf("List was not called")
		}
		if m.listQ != "タグ" || m.listLimit != 5 || m.listCursor != "cur123" {
			t.Fatalf("List args mismatch: q=%q limit=%d cursor=%q", m.listQ, m.listLimit, m.listCursor)
		}

		// レスポンスの確認
		var got handler.TagsPage
		unmarshalJSON(t, rec.Body.Bytes(), &got)
		if got.NextCursor != "next123" {
			t.Fatalf("NextCursor mismatch: got=%q", got.NextCursor)
		}
		if len(got.Items) != 2 ||
			got.Items[0].ID != "t1" || got.Items[0].Name != "タグ1" ||
			got.Items[1].ID != "t2" || got.Items[1].Name != "タグ2" {
			t.Fatalf("Items mismatch: got=%+v", got.Items)
		}
	})

	t.Run("BadRequest_limit_not_int", func(t *testing.T) {
		m := &mockTagsLister{}
		h := handler.ListTags(m)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/v1/tags?limit=abc", nil)
		h.ServeHTTP(rec, req)

		// ステータスコードの確認 (400)
		assertStatus(t, rec, http.StatusBadRequest)

		// 表示数が不適切なのにListは呼び出してほしくない
		if m.listCalled {
			t.Fatalf("List should not be called on bad request")
		}
	})

	t.Run("InternalServerError_usecase_error", func(t *testing.T) {
		m := &mockTagsLister{listErr: errors.New("intentional error")}
		h := handler.ListTags(m)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/v1/tags?limit=5", nil)
		h.ServeHTTP(rec, req)

		// ステータスコードの確認 (500)
		assertStatus(t, rec, http.StatusInternalServerError)

		// そもそもListが呼び出されてないなら、おかしい
		if !m.listCalled {
			t.Fatalf("List was not called")
		}
	})
}

func TestCreateTag(t *testing.T) {
	t.Run("OK_new", func(t *testing.T) {
		m := &mockTagCreator{
			createResponse: usecase.CreateTagResult{
				Tag:     usecase.Tag{ID: "t10", Name: "タグ10"},
				Created: true,
			},
		}
		h := handler.CreateTag(m)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/v1/tags", bytes.NewReader([]byte(`{"name": "タグ10"}`)))
		req.Header.Set("Content-Type", "application/json")
		h.ServeHTTP(rec, req)

		// ステータスコード (201) とヘッダの確認
		assertStatus(t, rec, http.StatusCreated)
		assertContentType(t, rec)

		// 呼び出し確認
		if !m.createCalled {
			t.Fatalf("Create was not called")
		}
		if m.createName != "タグ10" {
			t.Fatalf("Create args mismatch: name=%q", m.createName)
		}

		// レスポンスはhandler.Tag型であるはず
		var got handler.Tag
		unmarshalJSON(t, rec.Body.Bytes(), &got)
		if got.ID != "t10" || got.Name != "タグ10" {
			t.Fatalf("Tag mismatch: got=%+v", got)
		}
	})

	t.Run("OK_existing", func(t *testing.T) {
		m := &mockTagCreator{
			createResponse: usecase.CreateTagResult{
				Tag:     usecase.Tag{ID: "t10", Name: "タグ10"},
				Created: false,
			},
		}
		h := handler.CreateTag(m)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/v1/tags", bytes.NewReader([]byte(`{"name": "タグ10"}`)))
		req.Header.Set("Content-Type", "application/json")
		h.ServeHTTP(rec, req)

		// ステータスコードの確認 (200)
		assertStatus(t, rec, http.StatusOK)
	})

	t.Run("BadRequest_invalid_json", func(t *testing.T) {
		m := &mockTagCreator{}
		h := handler.CreateTag(m)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/v1/tags", bytes.NewReader([]byte(`{`)))
		req.Header.Set("Content-Type", "application/json")
		h.ServeHTTP(rec, req)

		// ステータスコードの確認 (400)
		assertStatus(t, rec, http.StatusBadRequest)

		// 不適切な形式のJSONについて、Createメソッドを呼び出してほしくない
		if m.createCalled {
			t.Fatal("Create should not be called on bad request")
		}
	})

	t.Run("BadRequest_empty_name", func(t *testing.T) {
		m := &mockTagCreator{}
		h := handler.CreateTag(m)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/v1/tags", bytes.NewReader([]byte(`{"name": ""}`)))
		req.Header.Set("Content-Type", "application/json")
		h.ServeHTTP(rec, req)

		// ステータスコードの確認 (400)
		assertStatus(t, rec, http.StatusBadRequest)

		// タグ名が空なのに、Createメソッドを呼び出してほしくない
		if m.createCalled {
			t.Fatal("Create should not be called on bad request")
		}
	})

	t.Run("InternalServerError_usecase_error", func(t *testing.T) {
		m := &mockTagCreator{createErr: errors.New("intentional error")}
		h := handler.CreateTag(m)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/v1/tags", bytes.NewReader([]byte(`{"name": "タグ10"}`)))
		req.Header.Set("Content-Type", "application/json")
		h.ServeHTTP(rec, req)

		// ステータスコードの確認 (500)
		assertStatus(t, rec, http.StatusInternalServerError)

		// そもそもCreateが呼び出されてないなら、おかしい
		if !m.createCalled {
			t.Fatalf("Create was not called")
		}
	})
}

func assertStatus(t *testing.T, rec *httptest.ResponseRecorder, want int) {
	t.Helper()

	if rec.Code != want {
		t.Fatalf("status=%d want=%d", rec.Code, want)
	}
}

func assertContentType(t *testing.T, rec *httptest.ResponseRecorder) {
	t.Helper()

	if ct := rec.Header().Get("Content-Type"); ct != "application/json; charset=utf-8" {
		t.Fatalf("content-type=%q", ct)
	}
}

func unmarshalJSON(t *testing.T, data []byte, v any) {
	t.Helper()
	if err := json.Unmarshal(data, v); err != nil {
		t.Fatalf("invalid json: %v body=%s", err, string(data))
	}
}
