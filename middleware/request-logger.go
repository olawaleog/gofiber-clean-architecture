package middleware

import (
	"github.com/RizkiMufrizal/gofiber-clean-architecture/logger"
	"github.com/gofiber/fiber/v2"
	"strings"
)

func maskPassword(body string) string {
	// Split the body into key-value pairs
	pairs := strings.Split(body, "&")
	for i, pair := range pairs {
		// Split each pair into key and value
		kv := strings.Split(pair, "=")
		if len(kv) == 2 && strings.ToLower(kv[0]) == "password" {
			// Mask the password value
			kv[1] = "****"
			pairs[i] = strings.Join(kv, "=")
		}
	}
	return strings.Join(pairs, "&")
}

func RequestLogger(c *fiber.Ctx) error {
	// Log request details
	var body string
	if c.Body() != nil {
		body = string(c.Body())
		body = maskPassword(body)
	} else {
		body = ""
	}
	logger.Logger.Infof("Request: %s %s %s %s", c.Method(), c.Path(), c.IP(), body)

	// Call the next handler
	err := c.Next()

	// Log response details
	logger.Logger.Infof("Response: %d %s", c.Response().StatusCode(), c.Response().Body())

	return err
}
