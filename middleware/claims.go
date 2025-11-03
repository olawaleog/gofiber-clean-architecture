package middleware

import (
	"strings"

	"github.com/RizkiMufrizal/gofiber-clean-architecture/exception"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/service"
	"github.com/gofiber/fiber/v2"
)

// ExtractClaims middleware extracts JWT claims and adds them to the local context
func ExtractClaims(userService service.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Skip token check for unprotected routes
		if isPublicRoute(c.Path()) {
			return c.Next()
		}

		token := c.Get("Authorization")
		if token == "" {
			return exception.UnauthorizedError{
				Message: "Authentication required",
			}
		}

		// Check for Bearer prefix and remove it if present
		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}

		claims, err := userService.GetClaimsFromToken(c.Context(), token)
		if err != nil {
			return exception.UnauthorizedError{
				Message: "Invalid or expired token",
			}
		}

		// Store claims in locals for easy access in controller methods
		c.Locals("claims", claims)
		return c.Next()
	}
}

// isPublicRoute checks if a route should skip authentication
func isPublicRoute(path string) bool {
	publicPaths := []string{
		"/v1/api/authentication",
		"/v1/api/register",
		"/v1/api/register-customer",
		"/v1/api/reset-password",
		"/v1/api/verify-otp",
		"/v1/api/post-new-password",
		"/v1/api/verify-phone",
		"/v1/api/update-fcm-token",
		"/v1/api/user/reset-password",
		"/swagger",
	}

	for _, publicPath := range publicPaths {
		if strings.HasPrefix(path, publicPath) {
			return true
		}
	}
	return false
}

// GetClaims is a helper function to retrieve claims from context
func GetClaims(c *fiber.Ctx) (map[string]interface{}, bool) {
	claims, ok := c.Locals("claims").(map[string]interface{})
	return claims, ok
}

// RequireClaims middleware ensures that valid claims exist
func RequireClaims() fiber.Handler {
	return func(c *fiber.Ctx) error {

		if isPublicRoute(c.Path()) {
			return c.Next()
		}
		_, ok := GetClaims(c)
		if !ok {
			return exception.UnauthorizedError{
				Message: "Authentication required",
			}
		}
		return c.Next()
	}
}
