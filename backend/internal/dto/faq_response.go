package dto

// FAQProjectInfo is a lightweight project reference used in FAQ project filters.
type FAQProjectInfo struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// FAQResponse is the API response struct for FAQ entries, including project info.
type FAQResponse struct {
	ID          uint64  `json:"id"`
	Question    string  `json:"question"`
	Answer      string  `json:"answer"`
	ProjectID   *uint64 `json:"project_id"`
	ProjectName string  `json:"project_name"`
	ProjectSlug string  `json:"project_slug"`
	IsGlobal    bool    `json:"is_global"`
	SortOrder   int     `json:"sort_order"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}
