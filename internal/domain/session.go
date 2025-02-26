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

type Record struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	SessionID uuid.UUID `json:"session_id" gorm:"type:uuid"`
	Text      string    `json:"text"`
	GitLink   string    `json:"git_link,omitempty"`
	Files     []File    `json:"files" gorm:"foreignKey:RecordID"`
	AudioURL  string    `json:"audio_url,omitempty"`
	Timestamp time.Time `json:"timestamp"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type File struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	RecordID  uuid.UUID `json:"record_id" gorm:"type:uuid"`
	Name      string    `json:"name"`
	URL       string    `json:"url"`
	Type      string    `json:"type"`
	Size      int64     `json:"size"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
