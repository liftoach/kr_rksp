package service

import (
	auth "kr/pkg/jwt"
)

type Service struct {
	BudgetRepository      BudgetRepository
	CategoryRepository    CategoryRepository
	TransactionRepository TransactionRepository
	UserRepository        UserRepository
	jwt                   *auth.JWTManager
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
