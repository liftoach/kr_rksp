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

type CategoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) Create(ctx context.Context, c *domain.Category) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}

	query := `
		INSERT INTO categories (id, user_id, name, transaction_type)
		VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.ExecContext(ctx, query,
		c.ID,
		c.UserID,
		c.Name,
		c.TransactionType,
	)

	if err != nil {
		return fmt.Errorf("run sql query: %w", err)
	}

	return nil
}

func (r *CategoryRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Category, error) {
	query := `
		SELECT id, user_id, name, transaction_type
		FROM categories
		WHERE id = $1
	`

	var c domain.Category

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&c.ID,
		&c.UserID,
		&c.Name,
		&c.TransactionType,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, service.ErrNotFound
		}
		return nil, fmt.Errorf("run sql query: %w", err)
	}

	return &c, nil
}

func (r *CategoryRepository) GetAll(ctx context.Context) ([]domain.Category, error) {
	query := `
		SELECT id, user_id, name, transaction_type
		FROM categories
		ORDER BY name ASC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("run sql query: %w", err)
	}
	defer func() {
		err = errors.Join(err, rows.Close())
	}()

	var categories []domain.Category

	for rows.Next() {
		var c domain.Category

		if err = rows.Scan(
			&c.ID,
			&c.UserID,
			&c.Name,
			&c.TransactionType,
		); err != nil {
			return nil, err
		}

		categories = append(categories, c)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("rows error: %w", rows.Err())
	}

	return categories, nil
}

func (r *CategoryRepository) Update(ctx context.Context, c *domain.Category) error {
	query := `
		UPDATE categories
		SET user_id = $1,
		    name = $2,
		    transaction_type = $3
		WHERE id = $4
	`

	res, err := r.db.ExecContext(ctx, query,
		c.UserID,
		c.Name,
		c.TransactionType,
		c.ID,
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

func (r *CategoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM categories
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
