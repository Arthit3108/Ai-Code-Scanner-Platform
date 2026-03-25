package models

import "time"

type User struct {
	ID         uint         `gorm:"primaryKey" json:"id"`
	Email      string       `gorm:"uniqueIndex;not null" json:"email"`
	Name       string       `json:"name"`
	AvatarURL  string       `json:"avatar_url"`
	Provider   string       `json:"provider"`
	ProviderID string       `json:"provider_id"`
	CreatedAt  time.Time    `json:"created_at"`
	UpdatedAt  time.Time    `json:"updated_at"`
}

type Repository struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index" json:"user_id"`
	URL       string    `json:"url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ScanJobs  []ScanJob `json:"scan_jobs,omitempty"`
}
