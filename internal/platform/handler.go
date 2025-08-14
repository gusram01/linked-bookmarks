package platform

import (
	"strings"

	"github.com/clerk/clerk-sdk-go/v2/jwt"
	"github.com/clerk/clerk-sdk-go/v2/user"
	"github.com/gofiber/fiber/v2"
)

func JwtClerkMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		sessionToken := strings.TrimPrefix(c.Get("Authorization", ""), "Bearer")

		if sessionToken == "" {
			return c.Status(401).JSON(fiber.Map{
				"success": false,
				"error":   "Authentication required",
				"data":    nil,
			})
		}

		claims, err := jwt.Verify(c.UserContext(), &jwt.VerifyParams{
			Token: sessionToken,
		})

		if err != nil {
			return c.Status(403).JSON(fiber.Map{
				"success": false,
				"error":   "You do not have permission to access this resource",
				"data":    nil,
			})
		}

		usr, err := user.Get(c.UserContext(), claims.Subject)
		if err != nil {
			return c.Status(403).JSON(fiber.Map{
				"success": false,
				"error":   "You do not have permission to access this resource",
				"data":    nil,
			})
		}

		if usr.Banned {
			return c.Status(403).JSON(fiber.Map{
				"success": false,
				"error":   "You do not have permission to access this resource",
				"data":    nil,
			})
		}

		return c.Next()
	}
}
