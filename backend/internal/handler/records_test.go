package handler_test

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Komura-Taichi/nipopo/backend/internal/entity"
	"github.com/Komura-Taichi/nipopo/backend/internal/handler"
	"github.com/Komura-Taichi/nipopo/backend/internal/handler/middleware"
	"github.com/Komura-Taichi/nipopo/backend/internal/usecase"
)

type mockRecordsLister struct {
	listCalled bool
	listUserID string
	listInput  usecase.ListRecordsInput

	listResponse entity.RecordsPage
	listErr      error
}

type mockRecordRetriever struct {
	retrieveByIDCalled   bool
	retrieveByIDUserID   string
	retrieveByIDRecordID string

	retrieveByIDResponse entity.Record
	retrieveByIDErr      error
}

type mockRecordCreator struct {
	createCalled bool
	createUserID string
	createInput  usecase.CreateRecordInput

	createResponse entity.Record
	createErr      error
}

type mockRecordUpdater struct {
	updateCalled   bool
	updateUserID   string
	updateRecordID string
	updateInput    usecase.UpdateRecordInput

	updateResponse entity.Record
	updateErr      error
}

type mockRecordDeleter struct {
	deleteCalled   bool
	deleteUserID   string
	deleteRecordID string

	deleteErr error
}

func (m *mockRecordsLister) List(ctx context.Context, userID string, input usecase.ListRecordsInput) (entity.RecordsPage, error) {
	_ = ctx
	m.listCalled = true
	m.listUserID, m.listInput = userID, input
	return m.listResponse, m.listErr
}

func (m *mockRecordRetriever) RetrieveByID(ctx context.Context, userID, recordID string) (entity.Record, error) {
	_ = ctx
	m.retrieveByIDCalled = true
	m.retrieveByIDUserID, m.retrieveByIDRecordID = userID, recordID
	return m.retrieveByIDResponse, m.retrieveByIDErr
}

func (m *mockRecordCreator) Create(ctx context.Context, userID string, input usecase.CreateRecordInput) (entity.Record, error) {
	_ = ctx
	m.createCalled = true
	m.createUserID, m.createInput = userID, input
	return m.createResponse, m.createErr
}

func (m *mockRecordUpdater) Update(ctx context.Context, userID, recordID string, input usecase.UpdateRecordInput) (entity.Record, error) {
	_ = ctx
	m.updateCalled = true
	m.updateUserID, m.updateRecordID, m.updateInput = userID, recordID, input
	return m.updateResponse, m.updateErr
}

func (m *mockRecordDeleter) Delete(ctx context.Context, userID, recordID string) error {
	_ = ctx
	m.deleteCalled = true
	m.deleteUserID, m.deleteRecordID = userID, recordID
	return m.deleteErr
}

func TestCreateRecord(t *testing.T) {
	const dateFormat = "2006-01-02"
	t.Run("OK_new", func(t *testing.T) {
		m := &mockRecordCreator{
			createResponse: entity.Record{
				UserID:    "u_1",
				ID:        "r_1",
				Date:      time.Date(2026, 3, 11, 0, 0, 0, 0, time.UTC),
				CreatedAt: time.Date(2026, 3, 11, 12, 34, 56, 0, time.UTC),
				UpdatedAt: time.Date(2026, 3, 11, 12, 34, 56, 0, time.UTC),
				Effort:    3,
				TagIDs:    []string{"t_1", "t_2"},
				Body:      "内容",
			},
		}

		h := middleware.AuthStub("u_1")(handler.CreateRecord(m))

		reqBody := `{
			"date": "2026-03-11",
			"effort": 3,
			"tag_ids": ["t_1", "t_2"],
			"body": "内容"
		}`

		rec := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/v1/records", bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		h.ServeHTTP(rec, req)

		// ステータスコード（201）とヘッダの確認
		assertStatus(t, rec, http.StatusCreated)
		assertContentType(t, rec)

		// usecaseの呼び出し確認
		if !m.createCalled {
			t.Fatalf("Create was not called")
		}
		if m.createUserID != "u_1" {
			t.Fatalf("Create userID mismatch: got=%q want=%q", m.createUserID, "u_1")
		}

		// usecaseの呼び出し時引数ごとの確認
		if got := m.createInput.Date.Format(dateFormat); got != "2026-03-11" {
			t.Fatalf("Create date mismatch: got=%q want=%q", got, "2026-03-11")
		}
		if m.createInput.Effort != 3 {
			t.Fatalf("Create effort mismatch: got=%d want=%d", m.createInput.Effort, 3)
		}
		if len(m.createInput.TagIDs) != 2 || m.createInput.TagIDs[0] != "t_1" || m.createInput.TagIDs[1] != "t_2" {
			t.Fatalf("Create tags mismatch: got=%v want=%v", m.createInput.TagIDs, []string{"t_1", "t_2"})
		}
		if m.createInput.Body != "内容" {
			t.Fatalf("Create body mismatch: got=%q want=%q", m.createInput.Body, "内容")
		}

		// レスポンスの確認
		var got handler.Record
		unmarshalJSON(t, rec.Body.Bytes(), &got)
		if got.ID != "r_1" {
			t.Fatalf("id mismatch: got=%q want=%q", got.ID, "r_1")
		}
		if got.Date != "2026-03-11" {
			t.Fatalf("date mismatch: got=%q want=%q", got.Date, "2026-03-11")
		}
		if got.Effort != 3 {
			t.Fatalf("effort mismatch: got=%d want=%d", got.Effort, 3)
		}
		if got.Body != "内容" {
			t.Fatalf("body mismatch: got=%q want=%q", got.Body, "内容")
		}

		// created_atとupdated_atはRFC3339形式にパースできなければならない
		if _, err := time.Parse(time.RFC3339, got.CreatedAt); err != nil {
			t.Fatalf("created_at is not RFC3339: got=%q err=%v", got.CreatedAt, err)
		}
		if _, err := time.Parse(time.RFC3339, got.UpdatedAt); err != nil {
			t.Fatalf("updated_at is not RFC3339: got=%q err=%v", got.UpdatedAt, err)
		}

		// レスポンスのタグが合ってるか確認
		if len(got.Tags) != 2 || got.Tags[0].ID != "t_1" || got.Tags[1].ID != "t_2" {
			t.Fatalf("response tags mismatch: got=%+v", got.Tags)
		}
	})

	t.Run("BadRequest_invalid_json", func(t *testing.T) {
		m := &mockRecordCreator{}

		h := middleware.AuthStub("u_1")(handler.CreateRecord(m))

		// 壊れたJSON（末尾カンマ）
		reqBody := `{
			"date": "2026-03-11",
			"effort": 3,
			"tag_ids": ["t_1"],
			"body": "内容",
		}`

		rec := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/v1/records", bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		h.ServeHTTP(rec, req)

		// ステータスコードの確認（400）
		assertStatus(t, rec, http.StatusBadRequest)

		// 壊れたJSONについて、Createメソッドを呼び出してほしくない
		if m.createCalled {
			t.Fatal("Create should not be called on bad request")
		}
	})

	t.Run("BadRequest_effort_greater_than_5", func(t *testing.T) {
		m := &mockRecordCreator{}

		h := middleware.AuthStub("u_1")(handler.CreateRecord(m))

		// 頑張り度が[0...5]の範囲外
		reqBody := `{
			"date": "2026-03-11",
			"effort": 6,
			"tag_ids": ["t_1"],
			"body": "内容"
		}`

		rec := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/v1/records", bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		h.ServeHTTP(rec, req)

		// ステータスコードの確認（400）
		assertStatus(t, rec, http.StatusBadRequest)

		// > 5 の頑張り度が渡された場合について、Createメソッドを呼び出してほしくない
		if m.createCalled {
			t.Fatal("Create should not be called on bad request")
		}
	})

	t.Run("BadRequest_missing_tag_ids", func(t *testing.T) {
		m := &mockRecordCreator{}

		h := middleware.AuthStub("u_1")(handler.CreateRecord(m))

		// tag_idsが抜けてる
		reqBody := `{
			"date": "2026-03-11",
			"effort": 3,
			"body": "内容"
		}`

		rec := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/v1/records", bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		h.ServeHTTP(rec, req)

		// ステータスコードの確認（400）
		assertStatus(t, rec, http.StatusBadRequest)

		// タグが指定されなかった場合について、Createメソッドを呼び出してほしくない
		if m.createCalled {
			t.Fatal("Create should not be called on bad request")
		}
	})

	t.Run("empty_tag_id_str", func(t *testing.T) {
		m := &mockRecordCreator{}

		h := middleware.AuthStub("u_1")(handler.CreateRecord(m))

		// tag_idsに含まれるIDが空文字列
		reqBody := `{
			"date": "2026-03-11",
			"effort": 3,
			"tag_ids": [""],
			"body": "内容"
		}`

		rec := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/v1/records", bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		h.ServeHTTP(rec, req)

		// ステータスコードの確認（400）
		assertStatus(t, rec, http.StatusBadRequest)

		// タグに空文字列が含まれる場合について、Createメソッドを呼び出してほしくない
		if m.createCalled {
			t.Fatal("Create should not be called on bad request")
		}
	})

	t.Run("BadRequest_duplicated_tag_ids", func(t *testing.T) {
		m := &mockRecordCreator{}

		h := middleware.AuthStub("u_1")(handler.CreateRecord(m))

		// タグが重複してる
		reqBody := `{
			"date": "2026-03-11",
			"effort": 3,
			"tag_ids": ["t_1", "t_1"],
			"body": "内容"
		}`

		rec := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/v1/records", bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		h.ServeHTTP(rec, req)

		// ステータスコードの確認（400）
		assertStatus(t, rec, http.StatusBadRequest)

		// タグが重複している場合について、Createメソッドを呼び出してほしくない
		if m.createCalled {
			t.Fatal("Create should not be called on bad request")
		}
	})

	t.Run("BadRequest_tag_ids_empty", func(t *testing.T) {
		m := &mockRecordCreator{}

		h := middleware.AuthStub("u_1")(handler.CreateRecord(m))

		// 紐づけられたタグがない
		reqBody := `{
			"date": "2026-03-11",
			"effort": 3,
			"tag_ids": [],
			"body": "内容"
		}`

		rec := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/v1/records", bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		h.ServeHTTP(rec, req)

		// ステータスコードの確認（400）
		assertStatus(t, rec, http.StatusBadRequest)

		// タグが紐づけられなかった場合について、Createメソッドを呼び出してほしくない
		if m.createCalled {
			t.Fatal("Create should not be called on bad request")
		}
	})

	t.Run("BadRequest_invalid_date", func(t *testing.T) {
		m := &mockRecordCreator{}

		h := middleware.AuthStub("u_1")(handler.CreateRecord(m))

		// dateのフォーマットがYYYY-mm-ddじゃない
		reqBody := `{
			"date": "2026/03/11",
			"effort": 3,
			"tag_ids": ["t_1"],
			"body": "内容"
		}`

		rec := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/v1/records", bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		h.ServeHTTP(rec, req)

		// ステータスコードの確認（400）
		assertStatus(t, rec, http.StatusBadRequest)

		// dateのフォーマットが不適切な場合について、Createメソッドを呼び出してほしくない
		if m.createCalled {
			t.Fatal("Create should not be called on bad request")
		}
	})

	t.Run("BadRequest_missing_body", func(t *testing.T) {
		m := &mockRecordCreator{}

		h := middleware.AuthStub("u_1")(handler.CreateRecord(m))

		// bodyが抜けてる
		reqBody := `{
			"date": "2026/03/11",
			"effort": 3,
			"tag_ids": ["t_1"],
		}`

		rec := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/v1/records", bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		h.ServeHTTP(rec, req)

		// ステータスコードの確認（400）
		assertStatus(t, rec, http.StatusBadRequest)

		// bodyが抜けてる場合について、Createメソッドを呼び出してほしくない
		if m.createCalled {
			t.Fatal("Create should not be called on bad request")
		}
	})

	t.Run("BadRequest_body_blank", func(t *testing.T) {
		m := &mockRecordCreator{}

		h := middleware.AuthStub("u_1")(handler.CreateRecord(m))

		// bodyが空白のみ
		reqBody := `{
			"date": "2026/03/11",
			"effort": 3,
			"tag_ids": ["t_1"],
			"body": " "
		}`

		rec := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/v1/records", bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		h.ServeHTTP(rec, req)

		// ステータスコードの確認（400）
		assertStatus(t, rec, http.StatusBadRequest)

		// bodyが空白のみの場合について、Createメソッドを呼び出してほしくない
		if m.createCalled {
			t.Fatal("Create should not be called on bad request")
		}
	})

	t.Run("Conflict_existing_date", func(t *testing.T) {
		m := &mockRecordCreator{
			createErr: &usecase.RecordAlreadyExistsError{ExistingID: "r_1"},
		}

		h := middleware.AuthStub("u_1")(handler.CreateRecord(m))

		// 指定した日に既に記録が存在している
		reqBody := `{
			"date": "2026-03-11",
			"effort": 3,
			"tag_ids": ["t_1"],
			"body": "内容"
		}`

		rec := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/v1/records", bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		h.ServeHTTP(rec, req)

		// ステータスコードの確認（409）とヘッダの確認
		assertStatus(t, rec, http.StatusConflict)
		assertContentType(t, rec)

		// リクエスト自体は正常なので、Createメソッドが呼ばれているはず
		if !m.createCalled {
			t.Fatal("Create was not called")
		}

		var got handler.ErrorResponse
		unmarshalJSON(t, rec.Body.Bytes(), &got)
		if got.Error.Code != http.StatusConflict {
			t.Fatalf("status=%d want=%d", got.Error.Code, http.StatusConflict)
		}
		if got.Error.Details == nil {
			t.Fatal("details should not be nil")
		}
		v, ok := got.Error.Details["existing_id"]
		if !ok {
			t.Fatal("existing_id not found")
		}
		existingID, ok := v.(string)
		if !ok {
			t.Fatalf("existing_id should be string: got=%T", v)
		}
		if existingID != "r_1" {
			t.Fatalf("existing_id mismatch: got=%q want=%q", existingID, "r_1")
		}
	})

	t.Run("InternalServerError_usecase_error", func(t *testing.T) {
		m := &mockRecordCreator{
			createErr: errors.New("intentional error"),
		}

		h := middleware.AuthStub("u_1")(handler.CreateRecord(m))

		// 紐づけられたタグがない
		reqBody := `{
			"date": "2026-03-11",
			"effort": 3,
			"tag_ids": [],
			"body": "内容"
		}`

		rec := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/v1/records", bytes.NewReader([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")

		h.ServeHTTP(rec, req)

		// ステータスコードの確認（500）
		assertStatus(t, rec, http.StatusInternalServerError)

		// そもそもCreateメソッドが呼び出されてないなら、おかしい
		if !m.createCalled {
			t.Fatal("Create was not called")
		}
	})

}
