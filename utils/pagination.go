package utils

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// PaginationParams holds pagination parameters
type PaginationParams struct {
	Page  int
	Limit int
}

// DefaultPagination provides default pagination values
var DefaultPagination = PaginationParams{
	Page:  1,
	Limit: 10,
}

// ExtractPaginationParams extracts pagination parameters from request
func ExtractPaginationParams(c *fiber.Ctx) PaginationParams {
	// Get pagination parameters from the query string
	pageStr := c.Query("page", "1")    // default to page 1 if not provided
	limitStr := c.Query("limit", "10") // default to 10 items per page if not provided

	// Convert string parameters to integers
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = DefaultPagination.Page
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = DefaultPagination.Limit
	}

	return PaginationParams{
		Page:  page,
		Limit: limit,
	}
}
