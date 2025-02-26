package postgres

import (
	"context"
	"errors"

	"github.com/ZigaoWang/zebra-server/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// GetFileByID retrieves a file by its ID
func (r *SessionRepository) GetFileByID(ctx context.Context, id uuid.UUID) (*domain.File, error) {
	var file domain.File
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&file).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &file, nil
}

// GetRecordByID retrieves a record by its ID
func (r *SessionRepository) GetRecordByID(ctx context.Context, id uuid.UUID) (*domain.Record, error) {
	var record domain.Record
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &record, nil
}
