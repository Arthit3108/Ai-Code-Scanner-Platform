package models

import (
	"time"

	"github.com/Arthit3108/devsecops-platform/internal/scanner"
)

type ScanStatus string

const (
	StatusPending  ScanStatus = "pending"
	StatusCloning  ScanStatus = "cloning"
	StatusScanning ScanStatus = "scanning"
	StatusDone     ScanStatus = "done"
	StatusFailed   ScanStatus = "failed"
)

// ScanJob represents a single security scanning execution for a repository.
type ScanJob struct {
	RunID        string               `gorm:"primaryKey" json:"run_id"`
	UserID       *uint                `gorm:"index" json:"user_id,omitempty"` 
	RepositoryID *uint                `gorm:"index" json:"repository_id,omitempty"` 
	RepoURL      string               `json:"repo_url"`
	Status       ScanStatus           `json:"status"`
	Error        string               `json:"error,omitempty"`
	// Results are stored as a JSON blob in the database for flexibility (GORM Serializer).
	Results      []scanner.ScanResult `gorm:"serializer:json" json:"results,omitempty"`
	CreatedAt    time.Time            `json:"created_at"`
	UpdatedAt    time.Time            `json:"updated_at"`
}
