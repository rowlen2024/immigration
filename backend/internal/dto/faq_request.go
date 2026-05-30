package dto

type FAQListRequest struct {
	PaginationRequest
	ProjectID *uint64 `form:"project_id"`
	IsGlobal  *bool   `form:"is_global"`
	Search    string  `form:"search"`
}
