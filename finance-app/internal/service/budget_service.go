package service

import (
	"context"
	"kr/internal/domain"

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

func (s *Service) CreateBudget(ctx context.Context, b *domain.Budget) error {
	return s.BudgetRepository.Create(ctx, b)
}

func (s *Service) GetBudgets(ctx context.Context, userID uuid.UUID) ([]domain.Budget, error) {
	return s.BudgetRepository.GetByUserID(ctx, userID)
}

func (s *Service) DeleteBudget(ctx context.Context, id uuid.UUID) error {
	return s.BudgetRepository.Delete(ctx, id)
}
