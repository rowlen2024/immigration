package dto

// CaseListRequest 案例列表查询参数
type CaseListRequest struct {
	PaginationRequest
	ProjectID   *uint64 `form:"project_id"`
	CountryFrom string  `form:"country_from"`
	Name        string  `form:"name"`
}
