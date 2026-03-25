package handler

import (
    "github.com/Arthit3108/devsecops-platform/internal/models"
    "github.com/Arthit3108/devsecops-platform/internal/service"
    "github.com/gofiber/fiber/v2"
)

type ScanHandler struct {
    svc *services.ScanService
}

func NewScanHandler(svc *services.ScanService) *ScanHandler {
    return &ScanHandler{svc: svc}
}

func (h *ScanHandler) StartScan(c *fiber.Ctx) error {
    req := new(models.ScanRequest)
    if err := c.BodyParser(req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid request body",
        })
    }

    if req.RepoURL == "" {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Repository URL is required",
        })
    }

    var user *models.User
    if u := c.Locals("user"); u != nil {
        user = u.(*models.User)
    }

    runID, err := h.svc.InitScan(req.RepoURL, user)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": err.Error(),
        })
    }

    return c.Status(fiber.StatusAccepted).JSON(models.ScanResponse{
        RunID:   runID,
        Status:  "pending",
        RepoURL: req.RepoURL,
    })
}

func (h *ScanHandler) GetScan(c *fiber.Ctx) error {
    runID := c.Params("run_id")

    job, ok := h.svc.GetJob(runID)
    if !ok {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "error": "run_id not found",
        })
    }

    return c.JSON(job)
}

