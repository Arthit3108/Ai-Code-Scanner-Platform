package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Arthit3108/devsecops-platform/internal/config"
	"github.com/Arthit3108/devsecops-platform/internal/db"
	"github.com/Arthit3108/devsecops-platform/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func GoogleLogin(c *fiber.Ctx) error {
	url := config.GoogleOAuthConfig.AuthCodeURL("state")
	return c.Redirect(url)
}

func GoogleCallback(c *fiber.Ctx) error {
	code := c.Query("code")
	token, err := config.GoogleOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString("Failed to exchange token")
	}

	client := config.GoogleOAuthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil || resp.StatusCode != http.StatusOK {
		return c.Status(fiber.StatusUnauthorized).SendString("Failed to get user info")
	}
	defer resp.Body.Close()

	var userInfo struct {
		ID      string `json:"id"`
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to parse user info")
	}

	user, err := findOrCreateUser("google", userInfo.ID, userInfo.Email, userInfo.Name, userInfo.Picture)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Database error")
	}

	return setJWTAndRedirect(c, user)
}

func GithubLogin(c *fiber.Ctx) error {
	url := config.GithubOAuthConfig.AuthCodeURL("state")
	return c.Redirect(url)
}

func GithubCallback(c *fiber.Ctx) error {
	code := c.Query("code")
	token, err := config.GithubOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString("Failed to exchange token")
	}

	client := config.GithubOAuthConfig.Client(context.Background(), token)

	// Get User Profile
	resp, err := client.Get("https://api.github.com/user")
	if err != nil || resp.StatusCode != http.StatusOK {
		return c.Status(fiber.StatusUnauthorized).SendString("Failed to get user info")
	}
	defer resp.Body.Close()

	var userInfo struct {
		ID        float64 `json:"id"` // Github sends numeric ID
		Login     string  `json:"login"`
		Name      string  `json:"name"`
		AvatarURL string  `json:"avatar_url"`
		Email     string  `json:"email"` // Might be empty if private
	}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to parse user info")
	}

	// Wait, if email is empty we need to fetch from /user/emails
	email := userInfo.Email
	if email == "" {
		emailResp, err := client.Get("https://api.github.com/user/emails")
		if err == nil && emailResp.StatusCode == http.StatusOK {
			defer emailResp.Body.Close()
			var emails []struct {
				Email   string `json:"email"`
				Primary bool   `json:"primary"`
			}
			json.NewDecoder(emailResp.Body).Decode(&emails)
			for _, e := range emails {
				if e.Primary {
					email = e.Email
					break
				}
			}
		}
	}

	// Use numeric ID as string
	providerID := ""
	if userInfo.ID > 0 {
		// Float64 to string
		providerID = string(rune(userInfo.ID)) // naive, use fmt.Sprint
		providerID = fmt.Sprintf("%.0f", userInfo.ID)
	}

	user, err := findOrCreateUser("github", providerID, email, userInfo.Name, userInfo.AvatarURL)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Database error: " + err.Error())
	}

	return setJWTAndRedirect(c, user)
}

func findOrCreateUser(provider, providerID, email, name, avatarURL string) (*models.User, error) {
	var user models.User
	result := db.DB.Where("provider = ? AND provider_id = ?", provider, providerID).First(&user)
	if result.Error != nil {
		// Create
		user = models.User{
			Provider:   provider,
			ProviderID: providerID,
			Email:      email,
			Name:       name,
			AvatarURL:  avatarURL,
		}
		if err := db.DB.Create(&user).Error; err != nil {
			return nil, err
		}
	}
	return &user, nil
}

func setJWTAndRedirect(c *fiber.Ctx, user *models.User) error {
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := jwtToken.SignedString(config.JwtSecret)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to generate token")
	}

	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    tokenString,
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
		SameSite: "Lax",
		Path:     "/",
	})

	frontendUrl := os.Getenv("FRONTEND_URL")
	if frontendUrl == "" {
		frontendUrl = "http://localhost:5174" // fallback to 5174 based on current terminal
	}
	return c.Redirect(frontendUrl)
}

func GetMe(c *fiber.Ctx) error {
	user := c.Locals("user")
	if user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}
	return c.JSON(user)
}

func Logout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		HTTPOnly: true,
		SameSite: "Lax",
		Path:     "/",
	})
	return c.JSON(fiber.Map{"message": "Logged out"})
}
