package model

// PaginationResponse represents a standard paginated response
type PaginationResponse struct {
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalItems int64       `json:"totalItems"`
	TotalPages int         `json:"totalPages"`
}

// NewPaginationResponse creates a new pagination response with calculated total pages
func NewPaginationResponse(data interface{}, page, limit int, totalItems int64) PaginationResponse {
	totalPages := int(totalItems) / limit
	if int(totalItems)%limit > 0 {
		totalPages++
	}

	return PaginationResponse{
		Data:       data,
		Page:       page,
		Limit:      limit,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}
}
