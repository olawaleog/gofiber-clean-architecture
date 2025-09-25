package utils

import (
	"github.com/RizkiMufrizal/gofiber-clean-architecture/model"
	"github.com/gofiber/fiber/v2"
)

// ResponseBuilder provides a consistent way to build API responses
type ResponseBuilder struct{}

// NewResponseBuilder creates a new ResponseBuilder instance
func NewResponseBuilder() *ResponseBuilder {
	return &ResponseBuilder{}
}

// Success creates a successful response with standard structure
func (rb *ResponseBuilder) Success(data interface{}, message string) model.GeneralResponse {
	return model.GeneralResponse{
		Code:    fiber.StatusOK,
		Message: message,
		Data:    data,
		Success: true,
	}
}

// Error creates an error response with standard structure
func (rb *ResponseBuilder) Error(statusCode int, message string) model.GeneralResponse {
	return model.GeneralResponse{
		Code:    statusCode,
		Message: message,
		Data:    nil,
		Success: false,
	}
}

// Pagination creates a paginated response
func (rb *ResponseBuilder) Pagination(data interface{}, page, limit int, totalCount int64) model.GeneralResponse {
	paginatedResponse := model.NewPaginationResponse(data, page, limit, totalCount)

	return model.GeneralResponse{
		Code:    fiber.StatusOK,
		Message: "Successful",
		Data:    paginatedResponse,
		Success: true,
	}
}
