package domain

import "github.com/google/uuid"

type Category struct {
	ID              uuid.UUID
	UserID          uuid.UUID
	Name            string
	TransactionType string
}
