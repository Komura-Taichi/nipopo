package handler

type createRecordRequest struct {
	Date   string   `json:"date"`
	Effort int      `json:"effort"`
	TagIDs []string `json:"tag_ids"`
	Body   string   `json:"body"`
}
