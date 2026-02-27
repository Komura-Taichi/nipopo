package entity

import "time"

type Record struct {
	ID        string
	UserID    string
	Date      time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	Effort    int
	TagIDs    []string
	Body      string
}

type RecordsPage struct {
	Items      []Record
	NextCursor string
}
