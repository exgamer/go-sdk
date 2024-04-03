package pagination

type Paging struct {
	Page    int      `json:"page"`
	OrderBy []string `json:"order_by"`
	Limit   int      `json:"limit"`
	ShowSQL bool
}
