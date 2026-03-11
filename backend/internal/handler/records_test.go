package handler_test

import (
	"bytes"
	"context"
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
			"body": "内容",
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
}
