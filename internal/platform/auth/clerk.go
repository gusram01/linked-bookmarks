package auth

import (
	"strings"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/jwt"
	"github.com/clerk/clerk-sdk-go/v2/user"
	"github.com/gofiber/fiber/v2"
	"github.com/gusram01/linked-bookmarks/internal"
	"github.com/gusram01/linked-bookmarks/internal/platform/logger"
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

		c.SetUserContext(
			clerk.ContextWithSessionClaims(c.UserContext(), claims),
		)

		return c.Next()
	}
}

func WithSessionClaims(c *fiber.Ctx) (clerk.SessionClaims, error) {
	claims, ok := clerk.SessionClaimsFromContext(c.UserContext())

	if !ok {
		logger.GetLogger().Error("Failed to get session claims")

		return clerk.SessionClaims{}, internal.NewErrorf(internal.ErrorCodeInvalidClaims, "Unauthorized")
	}

	if claims == nil {
		logger.GetLogger().Error("Session claims are nil")
		return clerk.SessionClaims{}, internal.NewErrorf(internal.ErrorCodeInvalidClaims, "Unauthorized")
	}

	return *claims, nil
}
