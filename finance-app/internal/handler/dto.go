package handler

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

type CreateTransactionRequest struct {
	Amount      decimal.Decimal `json:"amount"`
	Type        string          `json:"type"` // IN / OUT
	Category    string          `json:"category"`
	Description string          `json:"description"`
}

type CreateCategoryRequest struct {
	Name string `json:"name"`
	Type string `json:"type"` // IN / OUT
}

type CreateBudgetRequest struct {
	CategoryID uuid.UUID `json:"category_id"`
	Limit      float64   `json:"limit"`
	Period     string    `json:"period"`
}

type Summary struct {
	TotalIncome  decimal.Decimal `json:"total_income"`
	TotalExpense decimal.Decimal `json:"total_expense"`
	Balance      decimal.Decimal `json:"balance"`
}
