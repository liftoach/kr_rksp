package handler

import "github.com/shopspring/decimal"

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

type Summary struct {
	TotalIncome  decimal.Decimal `json:"total_income"`
	TotalExpense decimal.Decimal `json:"total_expense"`
	Balance      decimal.Decimal `json:"balance"`
}
