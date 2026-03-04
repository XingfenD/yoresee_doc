package service

type Pagination struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

func (p *Pagination) Validate() bool {
	if p == nil {
		return false
	}
	if p.Page <= 0 || p.PageSize <= 0 {
		return false
	}
	return true
}

type SortArgs struct {
	Field string `json:"field"`
	Desc  bool   `json:"desc"`
}
