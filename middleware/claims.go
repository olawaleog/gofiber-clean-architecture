package middleware

import (
	"github.com/RizkiMufrizal/gofiber-clean-architecture/exception"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/service"
	"github.com/gofiber/fiber/v2"
)

// ExtractClaims middleware extracts JWT claims and adds them to the local context
func ExtractClaims(userService service.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Get("Authorization")
		if token != "" {
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
		}

		return c.Next()
	}
}

// GetClaims is a helper function to retrieve claims from context
func GetClaims(c *fiber.Ctx) (map[string]interface{}, bool) {
	claims, ok := c.Locals("claims").(map[string]interface{})
	return claims, ok
}

// RequireClaims middleware ensures that valid claims exist
func RequireClaims() fiber.Handler {
	return func(c *fiber.Ctx) error {
		_, ok := GetClaims(c)
		if !ok {
			return exception.UnauthorizedError{
				Message: "Authentication required",
			}
		}
		return c.Next()
	}
}
