package usecase

import (
	"context"
	"time"

	"github.com/Komura-Taichi/nipopo/backend/internal/entity"
)

type ListRecordsInput struct {
	Q       string
	Limit   int
	Cursor  string
	OrderBy string
	Order   string

	DateFrom *time.Time
	DateTo   *time.Time

	TagIDs []string

	EffortFrom int
	EffortTo   int
}

type CreateRecordInput struct {
	Date   time.Time
	TagIDs []string
	Effort int
	Body   string
}

type UpdateRecordInput struct {
	Date   *time.Time
	TagIDs *[]string
	Effort *int
	Body   *string
}

type RecordsLister interface {
	List(ctx context.Context, userID string, input ListRecordsInput) (entity.RecordsPage, error)
}

type RecordRetriever interface {
	RetrieveByID(ctx context.Context, userID, recordID string) (entity.Record, error)
}

type RecordCreator interface {
	Create(ctx context.Context, userID string, input CreateRecordInput) (entity.Record, error)
}

type RecordUpdater interface {
	Update(ctx context.Context, userID, recordID string, input UpdateRecordInput) (entity.Record, error)
}

type RecordDeleter interface {
	Delete(ctx context.Context, userID, recordID string) error
}
