package service

type Service struct {
	BudgetRepository      BudgetRepository
	CategoryRepository    CategoryRepository
	TransactionRepository TransactionRepository
	UserRepository        UserRepository
}
