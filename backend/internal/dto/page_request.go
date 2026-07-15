package dto

type PageListRequest struct {
	PaginationRequest
	PageType string `form:"page_type"`
	Title    string `form:"title"`
	Status   string `form:"status"`
}

// CreatePageRequest 页面创建请求
type CreatePageRequest struct {
	ProjectID       *uint64  `json:"project_id"`
	Title           string   `json:"title"`
	Slug            string   `json:"slug"`
	Content         string   `json:"content"`
	CoverImage      string   `json:"cover_image"`
	Tags            []string `json:"tags"`
	MetaTitle       string   `json:"meta_title"`
	MetaDescription string   `json:"meta_description"`
	Template        string   `json:"template"`
	PageType        string   `json:"page_type"`
	Status          string   `json:"status"`
	SortOrder       int      `json:"sort_order"`
	IsPinned        bool     `json:"is_pinned"`
}
