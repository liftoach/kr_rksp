package service

import (
	"context"
	"kr/internal/domain"

	"github.com/google/uuid"
)

type TransactionRepository interface {
	Create(ctx context.Context, t *domain.Transaction) error

	GetByID(ctx context.Context, id uuid.UUID) (*domain.Transaction, error)
	GetAll(ctx context.Context) ([]domain.Transaction, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Transaction, error)

	Update(ctx context.Context, t *domain.Transaction) error
	Delete(ctx context.Context, id uuid.UUID) error
}

func (s *Service) CreateTransaction(ctx context.Context, t *domain.Transaction) error {
	return s.TransactionRepository.Create(ctx, t)
}

func (s *Service) GetTransactionsByUser(ctx context.Context, userID uuid.UUID) ([]domain.Transaction, error) {
	return s.TransactionRepository.GetByUserID(ctx, userID)
}

func (s *Service) GetTransactionByID(ctx context.Context, id uuid.UUID) (*domain.Transaction, error) {
	return s.TransactionRepository.GetByID(ctx, id)
}

func (s *Service) DeleteTransaction(ctx context.Context, id uuid.UUID) error {
	return s.TransactionRepository.Delete(ctx, id)
}
