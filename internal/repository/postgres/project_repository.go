package postgres

import (
	"context"

	"github.com/ZigaoWang/zebra-server/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProjectRepository struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

func (r *ProjectRepository) Create(ctx context.Context, project *domain.Project) error {
	return r.db.WithContext(ctx).Create(project).Error
}

func (r *ProjectRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Project, error) {
	var project domain.Project
	if err := r.db.WithContext(ctx).Preload("Sessions").Preload("Sessions.Records").Preload("Sessions.Records.Files").First(&project, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &project, nil
}

func (r *ProjectRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Project, error) {
	var projects []domain.Project
	if err := r.db.WithContext(ctx).Preload("Sessions").Preload("Sessions.Records").Preload("Sessions.Records.Files").Where("user_id = ?", userID).Find(&projects).Error; err != nil {
		return nil, err
	}
	return projects, nil
}

func (r *ProjectRepository) Update(ctx context.Context, project *domain.Project) error {
	return r.db.WithContext(ctx).Save(project).Error
}

func (r *ProjectRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&domain.Project{}, "id = ?", id).Error
}
