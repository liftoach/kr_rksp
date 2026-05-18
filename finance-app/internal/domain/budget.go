package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Budget struct {
	ID         uuid.UUID
	UserID     uuid.UUID
	CategoryID uuid.UUID
	Limit      decimal.Decimal
	Period     string
	CreatedAt  time.Time
}

type Category struct {
	ID              uuid.UUID
	UserID          uuid.UUID
	Name            string
	TransactionType string
}

type Transaction struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Amount      decimal.Decimal
	Type        string
	Category    string
	CreatedAt   time.Time
	Description string
}

const (
	TransactionTypeOut = "OUT"
	TransactionTypeIn  = "IN"
)

type User struct {
	ID        uuid.UUID
	Email     string
	Password  string
	CreatedAt time.Time
}
