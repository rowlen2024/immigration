package dto

// CreateLawyerInput 律师创建/更新请求
type CreateLawyerInput struct {
	PhotoURL  string   `json:"photo_url"`
	Name      string   `json:"name"`
	Title     string   `json:"title"`
	Tags      []string `json:"tags"`
	SortOrder int      `json:"sort_order"`
}

// CreateUserRequest 用户创建请求
type CreateUserRequest struct {
	Username    string `json:"username" binding:"required"`
	Password    string `json:"password" binding:"required"`
	DisplayName string `json:"display_name"`
	Role        string `json:"role"`
}
