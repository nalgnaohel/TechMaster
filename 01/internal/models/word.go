package models

type Word struct {
	ID        int64  `json:"id"`
	Lang      string `json:"lang"`
	Content   string `json:"content"`
	Translate string `json:"translate"`
}

type WordList struct {
	TotalCount int     `json:"total_count"`
	TotalPages int     `json:"total_pages"`
	Page       int     `json:"page"`
	Size       int     `json:"size"`
	HasMore    bool    `json:"has_more"`
	Words      []*Word `json:"aircrafts"`
}
