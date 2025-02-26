package postgres

import (
	"context"

	"github.com/ZigaoWang/zebra-server/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WorkLogRepository struct {
	db *gorm.DB
}

func NewWorkLogRepository(db *gorm.DB) *WorkLogRepository {
	return &WorkLogRepository{db: db}
}

func (r *WorkLogRepository) Create(ctx context.Context, log *domain.WorkLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

func (r *WorkLogRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.WorkLog, error) {
	var workLog domain.WorkLog
	if err := r.db.WithContext(ctx).First(&workLog, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &workLog, nil
}

func (r *WorkLogRepository) GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]domain.WorkLog, error) {
	var workLogs []domain.WorkLog
	query := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at DESC")
	
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&workLogs).Error; err != nil {
		return nil, err
	}
	return workLogs, nil
}

func (r *WorkLogRepository) Update(ctx context.Context, log *domain.WorkLog) error {
	return r.db.WithContext(ctx).Save(log).Error
}

func (r *WorkLogRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&domain.WorkLog{}, "id = ?", id).Error
}
