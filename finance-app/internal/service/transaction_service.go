package service

import (
	"context"
	"kr/finance-app/internal/domain"

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
