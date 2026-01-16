package entity

type Tag struct {
	ID     string
	UserID string
	Name   string
}

type TagsPage struct {
	Items      []Tag
	NextCursor string
}
