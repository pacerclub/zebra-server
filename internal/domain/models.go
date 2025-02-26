package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	Email     string    `json:"email" gorm:"unique;not null"`
	Password  string    `json:"-" gorm:"not null"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type WorkLog struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	UserID      uuid.UUID `json:"user_id" gorm:"type:uuid"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time,omitempty"`
	Duration    int64     `json:"duration"` // in seconds
	Status      string    `json:"status"`   // active, paused, completed
	Tags        []string  `json:"tags" gorm:"type:text[]"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type LogEntry struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	WorkLogID uuid.UUID `json:"work_log_id" gorm:"type:uuid"`
	Type      string    `json:"type"` // commit, note, voice, media
	Content   string    `json:"content"`
	Metadata  JSON      `json:"metadata" gorm:"type:jsonb"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Project struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	UserID      uuid.UUID `json:"user_id" gorm:"type:uuid"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	GitHubRepo  string    `json:"github_repo,omitempty"`
	Sessions    []Session `json:"sessions" gorm:"foreignKey:ProjectID"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// JSON is a wrapper for handling JSONB in PostgreSQL
type JSON map[string]interface{}
