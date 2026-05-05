package dto

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LeadRequest struct {
	Name              string `json:"name" binding:"required"`
	Phone             string `json:"phone" binding:"required"`
	Email             string `json:"email"`
	InterestedProject string `json:"interested_project"`
	Message           string `json:"message"`
}

type PaginationRequest struct {
	Page    int    `form:"page" binding:"omitempty,min=1"`
	PerPage int    `form:"per_page" binding:"omitempty,min=1,max=100"`
	Status  string `form:"status"`
	Q       string `form:"q"`
}
