package domain

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	ProjectID uuid.UUID `json:"project_id" gorm:"type:uuid"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time,omitempty"`
	Duration  int64     `json:"duration"` // Duration in milliseconds
	Records   []Record  `json:"records,omitempty" gorm:"foreignKey:SessionID"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Record represents a record of a session
type Record struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	SessionID uuid.UUID `json:"session_id" gorm:"type:uuid;not null"`
	Text      string    `json:"text"`
	GitLink   string    `json:"git_link"`
	AudioURL  string    `json:"audio_url"`
	AudioData []byte    `json:"audio_data" gorm:"type:bytea"` // Actual audio data
	Timestamp time.Time `json:"timestamp" gorm:"not null;default:now()"`
	Files     []File    `json:"files" gorm:"foreignKey:RecordID"`
	CreatedAt time.Time `json:"created_at" gorm:"not null;default:now()"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null;default:now()"`
}

// File represents a file uploaded by a user
type File struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	RecordID  uuid.UUID `json:"record_id" gorm:"type:uuid;not null"`
	Name      string    `json:"name" gorm:"not null"`
	URL       string    `json:"url"`
	Type      string    `json:"type"`
	Size      int64     `json:"size"`
	Data      []byte    `json:"data" gorm:"type:bytea"` // Actual file data
	CreatedAt time.Time `json:"created_at" gorm:"not null;default:now()"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null;default:now()"`
}
