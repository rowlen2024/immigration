package dto

// TestimonialListRequest 客户评价列表查询参数
type TestimonialListRequest struct {
	PaginationRequest
	ProjectID *uint64 `form:"project_id"`
	Nickname  string  `form:"nickname"`
	Rating    *uint8  `form:"rating"`
}
