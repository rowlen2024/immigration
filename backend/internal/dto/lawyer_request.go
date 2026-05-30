package dto

// LawyerListRequest 律师列表查询请求
type LawyerListRequest struct {
	PaginationRequest
	Name string `form:"name"`
}

// CreateLawyerInput 律师创建/更新请求
type CreateLawyerInput struct {
	PhotoURL  string   `json:"photo_url"`
	Name      string   `json:"name"`
	Title     string   `json:"title"`
	Tags      []string `json:"tags"`
	SortOrder int      `json:"sort_order"`
}

