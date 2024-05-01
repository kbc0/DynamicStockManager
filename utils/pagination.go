package utils

import (
	"github.com/gofiber/fiber/v2"
	"strconv"
)

// ParsePagination reads pagination parameters from the request query
func ParsePagination(c *fiber.Ctx) (int64, int64) {
	// Default values
	limit := int64(10)
	offset := int64(0)

	// Read limit from query, if present
	if limitQuery := c.Query("limit"); limitQuery != "" {
		if newLimit, err := strconv.ParseInt(limitQuery, 10, 64); err == nil {
			limit = newLimit
		}
	}

	// Read offset from query, if present
	if offsetQuery := c.Query("offset"); offsetQuery != "" {
		if newOffset, err := strconv.ParseInt(offsetQuery, 10, 64); err == nil {
			offset = newOffset
		}
	}

	return limit, offset
}
