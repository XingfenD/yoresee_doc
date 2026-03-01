package service

type Pagination struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

type SortArgs struct {
	Field string `json:"field"`
	Desc  bool   `json:"desc"`
}
