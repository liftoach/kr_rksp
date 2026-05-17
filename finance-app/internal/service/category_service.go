package service

import (
	"context"
	"kr/finance-app/internal/domain"

	"github.com/google/uuid"
)

type CategoryRepository interface {
	Create(ctx context.Context, c *domain.Category) error

	GetByID(ctx context.Context, id uuid.UUID) (*domain.Category, error)
	GetAll(ctx context.Context) ([]domain.Category, error)

	Update(ctx context.Context, c *domain.Category) error
	Delete(ctx context.Context, id uuid.UUID) error
}
