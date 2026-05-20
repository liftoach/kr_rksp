package service

import (
	"context"
	"kr/internal/domain"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type mockBudgetRepository struct {
	budgets map[uuid.UUID]*domain.Budget
	err     error
}

func (m *mockBudgetRepository) Create(
	ctx context.Context,
	b *domain.Budget,
) error {
	if m.err != nil {
		return m.err
	}

	if m.budgets == nil {
		m.budgets = make(map[uuid.UUID]*domain.Budget)
	}

	m.budgets[b.ID] = b
	return nil
}

func (m *mockBudgetRepository) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (*domain.Budget, error) {
	b, ok := m.budgets[id]
	if !ok {
		return nil, ErrNotFound
	}

	return b, nil
}

func (m *mockBudgetRepository) GetAll(
	ctx context.Context,
) ([]domain.Budget, error) {
	var result []domain.Budget

	for _, b := range m.budgets {
		result = append(result, *b)
	}

	return result, nil
}

func (m *mockBudgetRepository) GetByUserID(
	ctx context.Context,
	userID uuid.UUID,
) ([]domain.Budget, error) {
	var result []domain.Budget

	for _, b := range m.budgets {
		if b.UserID == userID {
			result = append(result, *b)
		}
	}

	return result, m.err
}

func (m *mockBudgetRepository) Update(
	ctx context.Context,
	b *domain.Budget,
) error {
	m.budgets[b.ID] = b
	return nil
}

func (m *mockBudgetRepository) Delete(
	ctx context.Context,
	id uuid.UUID,
) error {
	delete(m.budgets, id)
	return m.err
}

func FuzzService_CreateBudget(f *testing.F) {
	f.Add("monthly", "1000.50")
	f.Add("", "0")
	f.Add("yearly", "-500")

	f.Fuzz(func(t *testing.T, period, limitStr string) {
		limit, err := decimal.NewFromString(limitStr)
		if err != nil {
			t.Skip()
		}

		repo := &mockBudgetRepository{}

		svc := &Service{
			BudgetRepository: repo,
		}

		budget := &domain.Budget{
			ID:         uuid.New(),
			UserID:     uuid.New(),
			CategoryID: uuid.New(),
			Limit:      limit,
			Period:     period,
			CreatedAt:  time.Now(),
		}

		err = svc.CreateBudget(context.Background(), budget)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		saved, ok := repo.budgets[budget.ID]
		if !ok {
			t.Fatalf("budget not saved")
		}

		if !saved.Limit.Equal(limit) {
			t.Fatalf(
				"expected limit %s, got %s",
				limit,
				saved.Limit,
			)
		}

		if saved.Period != period {
			t.Fatalf(
				"expected period %q, got %q",
				period,
				saved.Period,
			)
		}
	})
}

func FuzzService_GetBudgets(f *testing.F) {
	f.Add("monthly")
	f.Add("")
	f.Add("weekly")

	f.Fuzz(func(t *testing.T, period string) {
		userID := uuid.New()

		budget := &domain.Budget{
			ID:         uuid.New(),
			UserID:     userID,
			CategoryID: uuid.New(),
			Limit:      decimal.NewFromInt(1000),
			Period:     period,
			CreatedAt:  time.Now(),
		}

		repo := &mockBudgetRepository{
			budgets: map[uuid.UUID]*domain.Budget{
				budget.ID: budget,
			},
		}

		svc := &Service{
			BudgetRepository: repo,
		}

		budgets, err := svc.GetBudgets(
			context.Background(),
			userID,
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(budgets) != 1 {
			t.Fatalf(
				"expected 1 budget, got %d",
				len(budgets),
			)
		}

		if budgets[0].Period != period {
			t.Fatalf(
				"expected period %q, got %q",
				period,
				budgets[0].Period,
			)
		}
	})
}

func FuzzService_DeleteBudget(f *testing.F) {
	f.Add("monthly")

	f.Fuzz(func(t *testing.T, period string) {
		id := uuid.New()

		repo := &mockBudgetRepository{
			budgets: map[uuid.UUID]*domain.Budget{
				id: {
					ID:         id,
					UserID:     uuid.New(),
					CategoryID: uuid.New(),
					Limit:      decimal.NewFromInt(500),
					Period:     period,
					CreatedAt:  time.Now(),
				},
			},
		}

		svc := &Service{
			BudgetRepository: repo,
		}

		err := svc.DeleteBudget(
			context.Background(),
			id,
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if _, exists := repo.budgets[id]; exists {
			t.Fatalf("budget was not deleted")
		}
	})
}
