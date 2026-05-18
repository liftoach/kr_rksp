package handler

import (
	"context"
	"kr/internal/domain"
	"kr/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Handler struct {
	svc service.Service
}

func NewHandler(svc *service.Service) *Handler {
	return &Handler{
		svc: *svc,
	}
}

func (h *Handler) Me(c *fiber.Ctx) error {
	userID := c.Locals("user_id")

	return c.JSON(fiber.Map{
		"user_id": userID,
	})
}

func (h *Handler) Register(c *fiber.Ctx) error {
	var req RegisterRequest

	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}

	if req.Email == "" || req.Password == "" {
		return fiber.ErrBadRequest
	}

	ctx := c.UserContext()

	err := h.svc.Register(ctx, req.Email, req.Password)
	if err != nil {
		return fiber.ErrBadRequest
	}

	return c.SendStatus(fiber.StatusCreated)
}

func (h *Handler) Login(c *fiber.Ctx) error {
	var req LoginRequest

	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}

	if req.Email == "" || req.Password == "" {
		return fiber.ErrBadRequest
	}

	ctx := c.UserContext()

	token, err := h.svc.Login(ctx, req.Email, req.Password)
	if err != nil {
		return fiber.ErrUnauthorized
	}

	return c.JSON(AuthResponse{
		Token: token,
	})
}

func (h *Handler) CreateTransaction(c *fiber.Ctx) error {
	var req CreateTransactionRequest

	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}

	userIDStr := c.Locals("user_id").(string)
	userID, _ := uuid.Parse(userIDStr)

	tx := &domain.Transaction{
		ID:          uuid.New(),
		UserID:      userID,
		Amount:      req.Amount,
		Type:        req.Type,
		Category:    req.Category,
		Description: req.Description,
	}

	err := h.svc.CreateTransaction(c.UserContext(), tx)
	if err != nil {
		return fiber.ErrBadRequest
	}

	return c.Status(fiber.StatusCreated).JSON(tx)
}

func (h *Handler) GetTransactions(c *fiber.Ctx) error {
	userIDStr := c.Locals("user_id").(string)
	userID, _ := uuid.Parse(userIDStr)

	txs, err := h.svc.GetTransactionsByUser(c.UserContext(), userID)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	return c.JSON(txs)
}

func (h *Handler) GetTransactionByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.ErrBadRequest
	}

	tx, err := h.svc.GetTransactionByID(c.UserContext(), id)
	if err != nil {
		return fiber.ErrNotFound
	}

	return c.JSON(tx)
}

func (h *Handler) DeleteTransaction(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.ErrBadRequest
	}

	err = h.svc.DeleteTransaction(c.UserContext(), id)
	if err != nil {
		return fiber.ErrBadRequest
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *Handler) CreateCategory(c *fiber.Ctx) error {
	var req CreateCategoryRequest

	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}

	userID, _ := uuid.Parse(c.Locals("user_id").(string))

	category := &domain.Category{
		ID:              uuid.New(),
		UserID:          userID,
		Name:            req.Name,
		TransactionType: req.Type,
	}

	err := h.svc.CreateCategory(c.UserContext(), category)
	if err != nil {
		return fiber.ErrBadRequest
	}

	return c.Status(fiber.StatusCreated).JSON(category)
}

func (h *Handler) GetCategories(c *fiber.Ctx) error {
	userID, _ := uuid.Parse(c.Locals("user_id").(string))

	cats, err := h.svc.GetCategories(c.UserContext(), userID)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	return c.JSON(cats)
}
func (h *Handler) DeleteCategory(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.ErrBadRequest
	}

	err = h.svc.DeleteCategory(c.UserContext(), id)
	if err != nil {
		return fiber.ErrBadRequest
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *Handler) CreateBudget(c *fiber.Ctx) error {
	var req CreateBudgetRequest

	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}

	userID, _ := uuid.Parse(c.Locals("user_id").(string))

	budget := &domain.Budget{
		ID:         uuid.New(),
		UserID:     userID,
		CategoryID: req.CategoryID,
		Period:     req.Period,
	}

	budget.Limit = decimal.NewFromFloat(req.Limit)

	err := h.svc.CreateBudget(c.UserContext(), budget)
	if err != nil {
		return fiber.ErrBadRequest
	}

	return c.Status(fiber.StatusCreated).JSON(budget)
}

func (h *Handler) GetBudgets(c *fiber.Ctx) error {
	userID, _ := uuid.Parse(c.Locals("user_id").(string))

	budgets, err := h.svc.GetBudgets(c.UserContext(), userID)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	return c.JSON(budgets)
}

func (h *Handler) GetSummary(c *fiber.Ctx) error {
	userID, _ := uuid.Parse(c.Locals("user_id").(string))

	summary, err := h.summary(c.UserContext(), userID)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	return c.JSON(summary)
}

func (h *Handler) summary(ctx context.Context, userID uuid.UUID) (Summary, error) {
	txs, err := h.svc.TransactionRepository.GetByUserID(ctx, userID)
	if err != nil {
		return Summary{}, err
	}

	income := decimal.Zero
	expense := decimal.Zero

	for _, t := range txs {
		switch t.Type {
		case domain.TransactionTypeIn:
			income = income.Add(t.Amount)

		case domain.TransactionTypeOut:
			expense = expense.Add(t.Amount)
		}
	}

	return Summary{
		TotalIncome:  income,
		TotalExpense: expense,
		Balance:      income.Sub(expense),
	}, nil
}
