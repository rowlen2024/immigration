package dto

type OptionPaginationRequest struct {
	Page    int `form:"page" binding:"omitempty,min=1"`
	PerPage int `form:"per_page" binding:"omitempty,min=1,max=500"`
}

type ProjectOptionRequest struct {
	OptionPaginationRequest
	Name string `form:"name"`
}

type CaseOptionRequest struct {
	OptionPaginationRequest
	ProjectID *uint64 `form:"project_id"`
	Name      string  `form:"name"`
}

type TestimonialOptionRequest struct {
	OptionPaginationRequest
	Nickname string `form:"nickname"`
}

type PageOptionRequest struct {
	OptionPaginationRequest
	PageType string `form:"page_type"`
	Title    string `form:"title"`
	Status   string `form:"status"`
}

type ProjectOption struct {
	ID   uint64 `json:"id"`
	Slug string `json:"slug"`
	Name string `json:"name"`
}

type CaseOption struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

type TestimonialOption struct {
	ID       uint64 `json:"id"`
	Nickname string `json:"nickname"`
}

type PageOption struct {
	ID    uint64 `json:"id"`
	Slug  string `json:"slug"`
	Title string `json:"title"`
}
