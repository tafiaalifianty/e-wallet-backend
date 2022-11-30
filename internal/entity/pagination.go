package entity

type Pagination struct {
	Limit      int    `json:"limit"`
	Page       int    `json:"page"`
	Search     string `json:"search,omitempty"`
	Sort       string `json:"sort"`
	SortBy     string `json:"sort_by"`
	TotalRows  int    `json:"total_rows"`
	TotalPages int    `json:"total_pages"`
}
