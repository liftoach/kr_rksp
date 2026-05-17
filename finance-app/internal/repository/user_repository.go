package ports

import "kr/finance-app/internal/domain"

type UserRepository interface {
	Create(user domain.User) (int64, error)
	FindByEmail(email string) (*domain.User, error)
	FindByID(id int64) (*domain.User, error)
}
