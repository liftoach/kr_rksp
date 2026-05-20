package service

import (
	auth "kr/pkg/jwt"

	"github.com/google/uuid"
)

type Service struct {
	BudgetRepository      BudgetRepository
	CategoryRepository    CategoryRepository
	TransactionRepository TransactionRepository
	UserRepository        UserRepository
	jwt                   JWTGenerator
}

func NewService(
	budgetRepository BudgetRepository,
	categoryRepository CategoryRepository,
	transactionRepository TransactionRepository,
	userRepository UserRepository,
	jwt *auth.JWTManager,
) *Service {
	return &Service{
		BudgetRepository:      budgetRepository,
		CategoryRepository:    categoryRepository,
		TransactionRepository: transactionRepository,
		UserRepository:        userRepository,
		jwt:                   jwt,
	}
}

type JWTGenerator interface {
	Generate(userID uuid.UUID) (string, error)
}
