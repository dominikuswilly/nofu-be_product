package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/dominikuswilly/nofu-be_product/internal/config"
	"github.com/dominikuswilly/nofu-be_product/internal/handler"
	"github.com/dominikuswilly/nofu-be_product/internal/repository"
	"github.com/dominikuswilly/nofu-be_product/internal/server"
	"github.com/dominikuswilly/nofu-be_product/internal/usecase"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func main() {
	// 1. Logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// 2. Config
	cfg := config.Load()

	// 3. Database
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	logger.Info("Connecting to database", zap.String("dsn_masked", fmt.Sprintf("host=%s user=%s dbname=%s", cfg.DBHost, cfg.DBUser, cfg.DBName)))

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		logger.Fatal("Failed to open database connection", zap.Error(err))
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		logger.Fatal("Failed to ping database", zap.Error(err))
	}

	// Create table if not exists (Simple migration for now)
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS products (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		description TEXT,
		price DECIMAL(10, 2) NOT NULL,
		created_at TIMESTAMP NOT NULL,
		updated_at TIMESTAMP NOT NULL
	);`
	if _, err := db.Exec(createTableSQL); err != nil {
		logger.Fatal("Failed to run migration", zap.Error(err))
	}

	// 4. Layers Setup
	repo := repository.NewPostgresProductRepository(db)
	uc := usecase.NewProductUsecase(repo)
	h := handler.NewProductHandler(uc, logger)

	// 5. Server
	srv := server.NewServer(cfg, h, logger)

	// 6. Graceful Shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := srv.Start(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server start failed", zap.Error(err))
		}
	}()

	<-ctx.Done()
	logger.Info("Shutting down gracefully...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exiting")
}
