package handler

type Tag struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type TagsPage struct {
	Items      []Tag  `json:"items"`
	NextCursor string `json:"next_cursor"`
}
