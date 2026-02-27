package handler

type SimpleTag struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Record struct {
	ID        string      `json:"id"`
	Date      string      `json:"date"`
	CreatedAt string      `json:"created_at"`
	UpdatedAt string      `json:"updated_at"`
	Effort    int         `json:"effort"`
	Tags      []SimpleTag `json:"tags"`
	Body      string      `json:"body"`
}

type SimpleRecord struct {
	ID          string      `json:"id"`
	Date        string      `json:"date"`
	Effort      int         `json:"effort"`
	Tags        []SimpleTag `json:"tags"`
	BodySnippet string      `json:"body"`
}

type SimpleRecordsPage struct {
	Items      []SimpleRecord `json:"items"`
	NextCursor string         `json:"next_cursor"`
}
