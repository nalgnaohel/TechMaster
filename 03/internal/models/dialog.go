package models

type Dialog struct {
	ID      int64  `json:"id"`
	Lang    string `json:"lang"`
	Content string `json:"content"`
}

type DialogList struct {
	TotalCount int       `json:"total_count"`
	TotalPages int       `json:"total_pages"`
	Page       int       `json:"page"`
	Size       int       `json:"size"`
	HasMore    bool      `json:"has_more"`
	Dialogs    []*Dialog `json:"dialogs"`
}
