package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"kr/finance-app/internal/domain"
	"kr/finance-app/internal/service"

	"github.com/google/uuid"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Create(ctx context.Context, t *domain.Transaction) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}

	query := `
		INSERT INTO transactions (id, user_id, amount, type, category, description)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.ExecContext(ctx, query,
		t.ID,
		t.UserID,
		t.Amount,
		t.Type,
		t.Category,
		t.Description,
	)

	if err != nil {
		return fmt.Errorf("run sql query: %w", err)
	}

	return nil
}

func (r *TransactionRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Transaction, error) {
	query := `
		SELECT id, user_id, amount, type, category, description, created_at
		FROM transactions
		WHERE id = $1
	`

	var t domain.Transaction

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&t.ID,
		&t.UserID,
		&t.Amount,
		&t.Type,
		&t.Category,
		&t.Description,
		&t.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, service.ErrNotFound
		}
		return nil, fmt.Errorf("run sql query: %w", err)
	}

	return &t, nil
}

func (r *TransactionRepository) GetAll(ctx context.Context) ([]domain.Transaction, error) {
	query := `
		SELECT id, user_id, amount, type, category, description, created_at
		FROM transactions
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("run sql query: %w", err)
	}
	defer func() {
		err = errors.Join(err, rows.Close())
	}()

	var txs []domain.Transaction

	for rows.Next() {
		var t domain.Transaction

		if err = rows.Scan(
			&t.ID,
			&t.UserID,
			&t.Amount,
			&t.Type,
			&t.Category,
			&t.Description,
			&t.CreatedAt,
		); err != nil {
			return nil, err
		}

		txs = append(txs, t)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("rows error: %w", rows.Err())
	}

	return txs, nil
}

func (r *TransactionRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Transaction, error) {
	query := `
		SELECT id, user_id, amount, type, category, description, created_at
		FROM transactions
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

	var txs []domain.Transaction

	for rows.Next() {
		var t domain.Transaction

		if err = rows.Scan(
			&t.ID,
			&t.UserID,
			&t.Amount,
			&t.Type,
			&t.Category,
			&t.Description,
			&t.CreatedAt,
		); err != nil {
			return nil, err
		}

		txs = append(txs, t)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("rows error: %w", rows.Err())
	}

	return txs, nil
}

func (r *TransactionRepository) Update(ctx context.Context, t *domain.Transaction) error {
	query := `
		UPDATE transactions
		SET amount = $1,
		    type = $2,
		    category = $3,
		    description = $4
		WHERE id = $5
	`

	res, err := r.db.ExecContext(ctx, query,
		t.Amount,
		t.Type,
		t.Category,
		t.Description,
		t.ID,
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

func (r *TransactionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM transactions
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
