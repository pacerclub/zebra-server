package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/ZigaoWang/zebra-server/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

func (r *SessionRepository) Create(ctx context.Context, projectID uuid.UUID, session *domain.Session) error {
	// Start a transaction
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Set session ID if not set
	if session.ID == uuid.Nil {
		session.ID = uuid.New()
	}
	
	// Set project ID
	session.ProjectID = projectID

	// Set timestamps
	now := time.Now()
	if session.CreatedAt.IsZero() {
		session.CreatedAt = now
	}
	session.UpdatedAt = now

	// Set end time if not set
	if session.EndTime.IsZero() {
		session.EndTime = now
	}

	// Calculate duration in seconds
	if session.Duration == 0 {
		// Use UTC times to avoid timezone issues
		startTimeUTC := session.StartTime.UTC()
		endTimeUTC := session.EndTime.UTC()
		
		// Print times for debugging
		fmt.Printf("StartTime: %v, EndTime: %v\n", startTimeUTC, endTimeUTC)
		
		// Calculate duration in seconds
		duration := endTimeUTC.Sub(startTimeUTC).Seconds()
		session.Duration = int64(duration)
		
		fmt.Printf("Calculated duration: %d seconds\n", session.Duration)
	}

	// Create the session first without the records
	sessionCopy := domain.Session{
		ID:        session.ID,
		ProjectID: session.ProjectID,
		StartTime: session.StartTime,
		EndTime:   session.EndTime,
		Duration:  session.Duration,
		CreatedAt: session.CreatedAt,
		UpdatedAt: session.UpdatedAt,
	}
	
	if err := tx.Create(&sessionCopy).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Create records one by one
	for i := range session.Records {
		record := &session.Records[i]

		// Always generate a new ID for records
		record.ID = uuid.New()
		record.SessionID = session.ID
		
		// Set timestamps
		if record.CreatedAt.IsZero() {
			record.CreatedAt = now
		}
		record.UpdatedAt = now

		// Set timestamp if not set
		if record.Timestamp.IsZero() {
			record.Timestamp = now
		}

		// Create a copy of the record without the files
		recordCopy := domain.Record{
			ID:        record.ID,
			SessionID: record.SessionID,
			Text:      record.Text,
			GitLink:   record.GitLink,
			AudioURL:  record.AudioURL,
			AudioData: record.AudioData,
			Timestamp: record.Timestamp,
			CreatedAt: record.CreatedAt,
			UpdatedAt: record.UpdatedAt,
		}

		if err := tx.Create(&recordCopy).Error; err != nil {
			tx.Rollback()
			return err
		}

		// Create files for this record
		for j := range record.Files {
			file := &record.Files[j]

			// Always generate a new ID for files
			file.ID = uuid.New()
			file.RecordID = record.ID
			
			// Set timestamps
			if file.CreatedAt.IsZero() {
				file.CreatedAt = now
			}
			file.UpdatedAt = now

			// Create a copy of the file
			fileCopy := domain.File{
				ID:        file.ID,
				RecordID:  file.RecordID,
				Name:      file.Name,
				URL:       file.URL,
				Type:      file.Type,
				Size:      file.Size,
				Data:      file.Data,
				CreatedAt: file.CreatedAt,
				UpdatedAt: file.UpdatedAt,
			}

			if err := tx.Create(&fileCopy).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	// Commit the transaction
	return tx.Commit().Error
}

func (r *SessionRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Session, error) {
	var session domain.Session
	if err := r.db.WithContext(ctx).Preload("Records.Files").First(&session, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *SessionRepository) GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]domain.Session, error) {
	var sessions []domain.Session
	if err := r.db.WithContext(ctx).Preload("Records.Files").Where("project_id = ?", projectID).Order("start_time desc").Find(&sessions).Error; err != nil {
		return nil, err
	}
	return sessions, nil
}

func (r *SessionRepository) Update(ctx context.Context, session *domain.Session) error {
	// Start a transaction
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Update the session
	if err := tx.Save(session).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Update or create associated records
	for i := range session.Records {
		record := &session.Records[i]
		if record.ID == uuid.Nil {
			record.ID = uuid.New()
		}
		record.SessionID = session.ID
		if err := tx.Save(record).Error; err != nil {
			tx.Rollback()
			return err
		}

		// Update or create associated files
		for j := range record.Files {
			file := &record.Files[j]
			if file.ID == uuid.Nil {
				file.ID = uuid.New()
			}
			file.RecordID = record.ID
			if err := tx.Save(file).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	// Commit the transaction
	return tx.Commit().Error
}

func (r *SessionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	// Start a transaction
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Get the session with its relationships
	var session domain.Session
	if err := tx.Preload("Records.Files").First(&session, "id = ?", id).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Delete files first
	for _, record := range session.Records {
		if err := tx.Where("record_id = ?", record.ID).Delete(&domain.File{}).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// Delete records
	if err := tx.Where("session_id = ?", id).Delete(&domain.Record{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Delete session
	if err := tx.Delete(&session).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction
	return tx.Commit().Error
}
