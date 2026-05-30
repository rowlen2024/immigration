package dto

// LeadListRequest 咨询列表查询参数
type LeadListRequest struct {
	PaginationRequest
	Status            string `form:"status"`
	Name              string `form:"name"`
	Email             string `form:"email"`
	InterestedProject string `form:"interested_project"`
}
