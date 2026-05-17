package service

import (
	"context"
	"kr/finance-app/internal/domain"
	auth "kr/finance-app/pkg/jwt"

	"github.com/google/uuid"
)

type Service struct {
	BudgetRepository      BudgetRepository
	CategoryRepository    CategoryRepository
	TransactionRepository TransactionRepository
	UserRepository        UserRepository
	jwt                   *auth.JWTManager
}

type BudgetRepository interface {
	Create(ctx context.Context, b *domain.Budget) error

	GetByID(ctx context.Context, id uuid.UUID) (*domain.Budget, error)
	GetAll(ctx context.Context) ([]domain.Budget, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Budget, error)

	Update(ctx context.Context, b *domain.Budget) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type TransactionRepository interface {
	Create(ctx context.Context, t *domain.Transaction) error

	GetByID(ctx context.Context, id uuid.UUID) (*domain.Transaction, error)
	GetAll(ctx context.Context) ([]domain.Transaction, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Transaction, error)

	Update(ctx context.Context, t *domain.Transaction) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error

	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	GetAll(ctx context.Context) ([]domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)

	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type CategoryRepository interface {
	Create(ctx context.Context, c *domain.Category) error

	GetByID(ctx context.Context, id uuid.UUID) (*domain.Category, error)
	GetAll(ctx context.Context) ([]domain.Category, error)

	Update(ctx context.Context, c *domain.Category) error
	Delete(ctx context.Context, id uuid.UUID) error
}
