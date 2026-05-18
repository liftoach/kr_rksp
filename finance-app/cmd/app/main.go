package main

import (
	"context"
	"kr/internal/config"
	"kr/internal/handler"
	"kr/internal/repository"
	"kr/internal/service"
	auth "kr/pkg/jwt"
	postgres "kr/pkg/pgx"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer func() {
		slog.Info("shutting down")
		stop()
		slog.Info("shut down")
	}()
	slog.Info("Starting finance-app")

	cfg := config.Load()

	slog.Info("config loaded")

	jwtManager := auth.NewJWTManager(
		cfg.JWTSecret,
		5*time.Hour,
	)

	db, err := postgres.New(ctx, postgres.Config{
		Host:     cfg.DBHost,
		Port:     cfg.DBPort,
		User:     cfg.DBUser,
		Password: cfg.DBPassword,
		DBName:   cfg.DBName,
		SSLMode:  cfg.DBSSLMode,
	})
	if err != nil {
		log.Fatalf("db connection failed: %v", err)
	}
	defer db.Close()

	slog.Info("connected to db")

	userRepo := repository.NewUserRepository(db)
	budgetRepo := repository.NewBudgetRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)

	slog.Info("initilzed repo")

	svc := service.NewService(
		budgetRepo,
		categoryRepo,
		transactionRepo,
		userRepo,
		jwtManager,
	)

	slog.Info("service initialized")

	app := handler.NewApp(svc, jwtManager)

	slog.Info("app initialized")

	go func() {
		log.Printf("server starting on %s", cfg.Port)

		if err = app.Run(":" + cfg.Port); err != nil {
			log.Fatalf("server failed: %v", err)
		}
		slog.Info("server starting on port %s", cfg.Port)
	}()

	<-ctx.Done()

	log.Println("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = app.Shutdown(shutdownCtx); err != nil {
		log.Printf("server shutdown error: %v", err)
	}

	log.Println("server stopped cleanly")
}
