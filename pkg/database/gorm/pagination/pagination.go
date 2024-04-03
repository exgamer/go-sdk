package pagination

type Pagination struct {
	TotalRecords int64 `json:"total_records"`
	TotalPage    int   `json:"total_page"`
	Offset       int   `json:"offset"`
	Limit        int   `json:"limit"`
	Page         int   `json:"page"`
	PrevPage     int   `json:"prev_page"`
	NextPage     int   `json:"next_page"`
}
