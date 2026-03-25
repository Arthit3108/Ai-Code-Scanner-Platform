package  routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/Arthit3108/devsecops-platform/internal/handler"
	"github.com/Arthit3108/devsecops-platform/internal/middleware"
)

func SetupRoutes(app *fiber.App, h *handler.ScanHandler) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Welcome to the AI DevSecOps Platform",
		})
	})
 
	api := app.Group("/api/v1")

	auth := api.Group("/auth")
	auth.Get("/google/login", handler.GoogleLogin)
	auth.Get("/google/callback", handler.GoogleCallback)
	auth.Get("/github/login", handler.GithubLogin)
	auth.Get("/github/callback", handler.GithubCallback)
	auth.Get("/me", middleware.RequireAuth, handler.GetMe)
	auth.Post("/logout", handler.Logout)

	api.Post("/scan", middleware.RequireAuth, h.StartScan)
	api.Get("/scan/:run_id", middleware.RequireAuth, h.GetScan)
	api.Get("/repos", middleware.RequireAuth, handler.GetUserRepos)
}