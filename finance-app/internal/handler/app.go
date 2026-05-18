package handler

import (
	"context"
	"kr/finance-app/internal/service"
	auth "kr/finance-app/pkg/jwt"

	"github.com/gofiber/fiber/v2"
)

type App struct {
	fiber         *fiber.App
	handler       *Handler
	jwtMiddleware fiber.Handler
}

func NewApp(svc *service.Service, jwtManager *auth.JWTManager) *App {
	app := fiber.New()

	h := NewHandler(svc)
	jwt := AuthMiddleware(jwtManager)

	a := &App{
		fiber:         app,
		handler:       h,
		jwtMiddleware: jwt,
	}

	a.registerMiddlewares()
	a.registerRoutes()

	return a
}

func (a *App) registerMiddlewares() {
	a.fiber.Use(CORSMiddleware())
}

func (a *App) registerRoutes() {
	a.fiber.Post("/auth/register", a.handler.Register)
	a.fiber.Post("/auth/login", a.handler.Login)

	api := a.fiber.Group("/api", a.jwtMiddleware)
	api.Get("/me", a.handler.Me)
	api.Post("/transactions", a.handler.CreateTransaction)
	api.Get("/transactions", a.handler.GetTransactions)
	api.Get("/transactions/:id", a.handler.GetTransactionByID)
	api.Delete("/transactions/:id", a.handler.DeleteTransaction)

	api.Post("/categories", a.handler.CreateCategory)
	api.Get("/categories", a.handler.GetCategories)
	api.Delete("/categories/:id", a.handler.DeleteCategory)

	api.Post("/budgets", a.handler.CreateBudget)
	api.Get("/budgets", a.handler.GetBudgets)

	api.Get("/analytics/summary", a.handler.GetSummary)
}

func (a *App) Run(addr string) error {
	return a.fiber.Listen(addr)
}

func (a *App) Shutdown(ctx context.Context) error {
	return a.fiber.ShutdownWithContext(ctx)
}
