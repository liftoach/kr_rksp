package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"kr/internal/domain"
	"kr/internal/service"

	"github.com/google/uuid"
)

type BudgetRepository struct {
	db *sql.DB
}

func NewBudgetRepository(db *sql.DB) *BudgetRepository {
	return &BudgetRepository{db: db}
}

func (r *BudgetRepository) Create(ctx context.Context, b *domain.Budget) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}

	query := `
		INSERT INTO budgets (id, user_id, category_id, budget_limit, period)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.ExecContext(ctx, query,
		b.ID,
		b.UserID,
		b.CategoryID,
		b.Limit,
		b.Period,
	)

	if err != nil {
		return fmt.Errorf("run sql query: %w", err)
	}

	return nil
}

func (r *BudgetRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Budget, error) {
	query := `
		SELECT id, user_id, category_id, budget_limit, period, created_at
		FROM budgets
		WHERE id = $1
	`

	var b domain.Budget

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&b.ID,
		&b.UserID,
		&b.CategoryID,
		&b.Limit,
		&b.Period,
		&b.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, service.ErrNotFound
		}
		return nil, fmt.Errorf("run sql query: %w", err)
	}

	return &b, nil
}

func (r *BudgetRepository) GetAll(ctx context.Context) ([]domain.Budget, error) {
	query := `
		SELECT id, user_id, category_id, budget_limit, period, created_at
		FROM budgets
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("run sql query: %w", err)
	}
	defer func() {
		err = errors.Join(err, rows.Close())
	}()

	var budgets []domain.Budget

	for rows.Next() {
		var b domain.Budget

		if err = rows.Scan(
			&b.ID,
			&b.UserID,
			&b.CategoryID,
			&b.Limit,
			&b.Period,
			&b.CreatedAt,
		); err != nil {
			return nil, err
		}

		budgets = append(budgets, b)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("rows error: %w", rows.Err())
	}

	return budgets, nil
}

func (r *BudgetRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Budget, error) {
	query := `
		SELECT id, user_id, category_id, budget_limit, period, created_at
		FROM budgets
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("run sql query: %w", err)
	}
	defer func() {
		err = errors.Join(err, rows.Close())
	}()

	var budgets []domain.Budget

	for rows.Next() {
		var b domain.Budget

		if err = rows.Scan(
			&b.ID,
			&b.UserID,
			&b.CategoryID,
			&b.Limit,
			&b.Period,
			&b.CreatedAt,
		); err != nil {
			return nil, err
		}

		budgets = append(budgets, b)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("rows error: %w", rows.Err())
	}

	return budgets, nil
}

func (r *BudgetRepository) Update(ctx context.Context, b *domain.Budget) error {
	query := `
		UPDATE budgets
		SET user_id = $1,
		    category_id = $2,
		    budget_limit = $3,
		    period = $4
		WHERE id = $5
	`

	res, err := r.db.ExecContext(ctx, query,
		b.UserID,
		b.CategoryID,
		b.Limit,
		b.Period,
		b.ID,
	)
	if err != nil {
		return err
	}

	aff, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if aff == 0 {
		return service.ErrNotFound
	}

	return nil
}

func (r *BudgetRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM budgets
		WHERE id = $1
	`

	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	aff, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if aff == 0 {
		return service.ErrNotFound
	}

	return nil
}
