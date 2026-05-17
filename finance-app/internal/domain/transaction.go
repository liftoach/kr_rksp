package domain

import "time"

type Transaction struct {
	ID        uuid.uuid
	UserID    int64
	Amount    decimal.Decimal
	Type      string
	Category  string
	CreatedAt time.Time
}
