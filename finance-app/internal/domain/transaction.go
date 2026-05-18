package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

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
