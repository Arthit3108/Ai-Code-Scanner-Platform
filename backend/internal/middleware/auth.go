package middleware

import (
	"strings"

	"github.com/Arthit3108/devsecops-platform/internal/config"
	"github.com/Arthit3108/devsecops-platform/internal/db"
	"github.com/Arthit3108/devsecops-platform/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// RequireAuth is a middleware that enforces authentication.
// It checks for a valid JWT token in cookies or the Authorization header.
// If the token is missing or invalid, it returns a 401 Unauthorized response.
func RequireAuth(c *fiber.Ctx) error {
	// 1. Extract token from Cookie "jwt"
	tokenString := c.Cookies("jwt")

	// 2. Fallback: Extract token from Authorization header (Bearer token)
	if tokenString == "" {
		authHeader := c.Get("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = strings.TrimPrefix(authHeader, "Bearer ")
		}
	}

	// If no token found, deny access
	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	// 3. Parse and validate the JWT token using the secret key
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return config.JwtSecret, nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
	}

	// 4. Extract claims and retrieve user_id
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token claims"})
	}

	userID := claims["user_id"]
	var user models.User

	// 5. Fetch user from database to ensure they still exist
	if err := db.DB.First(&user, userID).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User not found"})
	}

	// 6. Store user in context locals for use in subsequent handlers
	c.Locals("user", &user)
	return c.Next()
}

// OptionalAuth is a middleware that attempts to authenticate the user but does not block the request on failure.
// If a valid token is provided, it populates c.Locals("user").
func OptionalAuth(c *fiber.Ctx) error {
	// Attempt to extract token from cookies or header
	tokenString := c.Cookies("jwt")

	if tokenString == "" {
		authHeader := c.Get("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = strings.TrimPrefix(authHeader, "Bearer ")
		}
	}

	// If no token is provided, just proceed to the next handler
	if tokenString == "" {
		return c.Next()
	}

	// Attempt to parse token; if invalid, just proceed without setting user context
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return config.JwtSecret, nil
	})

	if err != nil || !token.Valid {
		return c.Next()
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.Next()
	}

	userID := claims["user_id"]
	var user models.User
	
	// If user exists in DB, store in locals
	if err := db.DB.First(&user, userID).Error; err == nil {
		c.Locals("user", &user)
	}

	return c.Next()
}
