package entity

type Tag struct {
	ID   string
	Name string
}

type TagsPage struct {
	Items      []Tag
	NextCursor string
}
