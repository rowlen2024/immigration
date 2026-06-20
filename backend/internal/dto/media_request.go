package dto

// MediaListRequest holds media list query params.
type MediaListRequest struct {
	PaginationRequest
	Search string `form:"search"`
}
