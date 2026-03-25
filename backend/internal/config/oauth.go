package config

import (
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
)

var (
	GoogleOAuthConfig *oauth2.Config
	GithubOAuthConfig *oauth2.Config
	JwtSecret         []byte
)

func InitOAuth() {
	backendUrl := os.Getenv("BACKEND_URL")
	if backendUrl == "" {
		backendUrl = "http://localhost:3000"
	}

	GoogleOAuthConfig = &oauth2.Config{
		RedirectURL:  backendUrl + "/api/v1/auth/google/callback",
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}

	GithubOAuthConfig = &oauth2.Config{
		RedirectURL:  backendUrl + "/api/v1/auth/github/callback",
		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		Scopes:       []string{"user:email"},
		Endpoint:     github.Endpoint,
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "super-secret-local-key"
	}
	JwtSecret = []byte(secret)
}
