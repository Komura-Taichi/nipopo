package handler_test

import (
	"context"

	"github.com/Komura-Taichi/nipopo/backend/internal/entity"
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
