package service

import (
	"context"
	"kr/internal/domain"

	"github.com/google/uuid"
)

func (s *Service) CreateBudget(ctx context.Context, b *domain.Budget) error {
	return s.BudgetRepository.Create(ctx, b)
}

func (s *Service) GetBudgets(ctx context.Context, userID uuid.UUID) ([]domain.Budget, error) {
	return s.BudgetRepository.GetByUserID(ctx, userID)
}

func (s *Service) DeleteBudget(ctx context.Context, id uuid.UUID) error {
	return s.BudgetRepository.Delete(ctx, id)
}
