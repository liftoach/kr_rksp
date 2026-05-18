package service

import (
	"context"
	"kr/internal/domain"

	"github.com/google/uuid"
)

func (s *Service) CreateCategory(ctx context.Context, c *domain.Category) error {
	return s.CategoryRepository.Create(ctx, c)
}

func (s *Service) GetCategories(ctx context.Context, userID uuid.UUID) ([]domain.Category, error) {
	return s.CategoryRepository.GetAll(ctx)
}

func (s *Service) DeleteCategory(ctx context.Context, id uuid.UUID) error {
	return s.CategoryRepository.Delete(ctx, id)
}
