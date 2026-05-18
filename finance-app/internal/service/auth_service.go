package service

import (
	"context"
	"errors"
	"fmt"
	"kr/internal/domain"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error

	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	GetAll(ctx context.Context) ([]domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)

	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id uuid.UUID) error
}

func (s *Service) Register(ctx context.Context, email, password string) error {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	user := &domain.User{
		ID:        uuid.New(),
		Email:     email,
		Password:  string(hash),
		CreatedAt: time.Now(),
	}

	return s.UserRepository.Create(ctx, user)
}

func (s *Service) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.UserRepository.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return "", ErrNotFound
		}
		return "", fmt.Errorf("userRepository.GetByEmail: %w", err)
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return "", errors.New("invalid password")
	}
	return s.jwt.Generate(user.ID)
}
