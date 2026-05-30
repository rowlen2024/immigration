package dto

type PageListRequest struct {
	PaginationRequest
	PageType string `form:"page_type"`
	Title    string `form:"title"`
	Status   string `form:"status"`
}
