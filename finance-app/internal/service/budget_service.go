package service

import (
	"context"

	"kr/finance-app/internal/domain"

	"github.com/google/uuid"
)

type BudgetRepository interface {
	Create(ctx context.Context, b *domain.Budget) error

	GetByID(ctx context.Context, id uuid.UUID) (*domain.Budget, error)
	GetAll(ctx context.Context) ([]domain.Budget, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Budget, error)

	Update(ctx context.Context, b *domain.Budget) error
	Delete(ctx context.Context, id uuid.UUID) error
}
