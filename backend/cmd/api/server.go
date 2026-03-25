package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/gofiber/fiber/v2"
	"github.com/Arthit3108/devsecops-platform/internal/config"
	"github.com/Arthit3108/devsecops-platform/internal/db"
	"github.com/Arthit3108/devsecops-platform/internal/handler"
	"github.com/Arthit3108/devsecops-platform/internal/middleware"
	"github.com/Arthit3108/devsecops-platform/internal/routes"
	"github.com/Arthit3108/devsecops-platform/internal/service"
)

func main() {
	// Try loading .env from two levels up (backend root)
	if err := godotenv.Load("../../.env"); err != nil {
		log.Println("Notice: could not load .env file, continuing with system env vars")
	}

	app := fiber.New(fiber.Config{
		AppName: "Ai DevSecOps Platform",
	})

	// Init DB
	db.InitDB()
	config.InitOAuth()

	// Initialize Services
	scanSvc := services.NewScanService()

	// Initialize Handlers
	scanHandler := handler.NewScanHandler(scanSvc)

	middleware.SetupMiddleware(app)
	routes.SetupRoutes(app, scanHandler)

	log.Fatal(app.Listen(":3000"))
}