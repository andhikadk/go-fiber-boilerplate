package models

// APIResponse is the standard API response wrapper
type APIResponse struct {
	Status  int         `json:"status" example:"200"`
	Message string      `json:"message" example:"Success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty" example:""`
}

// PaginatedResponse is the response wrapper for paginated data
type PaginatedResponse struct {
	Status  int         `json:"status" example:"200"`
	Message string      `json:"message" example:"Success"`
	Data    interface{} `json:"data"`
	Page    int         `json:"page" example:"1"`
	Limit   int         `json:"limit" example:"10"`
	Total   int64       `json:"total" example:"100"`
}
