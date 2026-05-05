package dto

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type PaginatedResponse struct {
	Code       int         `json:"code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
	Pagination *Pagination `json:"pagination"`
}

type Pagination struct {
	Page    int   `json:"page"`
	PerPage int   `json:"per_page"`
	Total   int64 `json:"total"`
}

type ErrorResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Errors  []FieldError `json:"errors,omitempty"`
}

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func Success(data interface{}) Response {
	return Response{Code: 200, Message: "ok", Data: data}
}

func SuccessPaginated(data interface{}, page, perPage int, total int64) PaginatedResponse {
	return PaginatedResponse{
		Code:    200,
		Message: "ok",
		Data:    data,
		Pagination: &Pagination{
			Page:    page,
			PerPage: perPage,
			Total:   total,
		},
	}
}

func Error(code int, message string) Response {
	return Response{Code: code, Message: message}
}
