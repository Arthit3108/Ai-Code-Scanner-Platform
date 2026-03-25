package handler

import (
	"github.com/Arthit3108/devsecops-platform/internal/db"
	"github.com/Arthit3108/devsecops-platform/internal/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetUserRepos(c *fiber.Ctx) error {
	u := c.Locals("user")
	if u == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}
	user := u.(*models.User)

	var repos []models.Repository
	// Preload all ScanJobs for each repository, sort by updated_at descending
	if err := db.DB.Preload("ScanJobs", func(d *gorm.DB) *gorm.DB {
		return d.Order("updated_at DESC")
	}).Where("user_id = ?", user.ID).Order("updated_at DESC").Find(&repos).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch repositories"})
	}

	return c.JSON(repos)
}
