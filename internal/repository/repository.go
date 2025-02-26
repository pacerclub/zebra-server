package repository

import (
	"context"

	"github.com/ZigaoWang/zebra-server/internal/domain"
	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type WorkLogRepository interface {
	Create(ctx context.Context, log *domain.WorkLog) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.WorkLog, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]domain.WorkLog, error)
	Update(ctx context.Context, log *domain.WorkLog) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type LogEntryRepository interface {
	Create(ctx context.Context, entry *domain.LogEntry) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.LogEntry, error)
	GetByWorkLogID(ctx context.Context, workLogID uuid.UUID) ([]domain.LogEntry, error)
	Update(ctx context.Context, entry *domain.LogEntry) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type ProjectRepository interface {
	Create(ctx context.Context, project *domain.Project) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Project, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Project, error)
	Update(ctx context.Context, project *domain.Project) error
	Delete(ctx context.Context, id uuid.UUID) error
}
