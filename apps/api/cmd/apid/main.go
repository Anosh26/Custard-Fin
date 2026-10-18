package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/anosh26/custardfin/apps/api/internal/adapters/auth"
	"github.com/anosh26/custardfin/apps/api/internal/adapters/httpapi"
	"github.com/anosh26/custardfin/apps/api/internal/adapters/memory"
	"github.com/anosh26/custardfin/apps/api/internal/adapters/postgres"
	"github.com/anosh26/custardfin/apps/api/internal/ports"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	slog.Info("starting custardfin apid (cloud control api)", "version", "0.1.0")

	// 1. Resolve configuration from environment
	apiAddr := getEnv("API_ADDR", ":8080")
	dbURL := getEnv("DATABASE_URL", "postgres://custard:custard_secret@localhost:5432/custardfin?sslmode=disable")
	jwtSecret := getEnv("JWT_SECRET", "custardfin_development_jwt_secret_key_32_bytes!!")

	// 2. Initialize Auth Services
	hasher := auth.NewBcryptHasher(bcrypt.DefaultCost)
	tokenService, err := auth.NewJWTService([]byte(jwtSecret))
	if err != nil {
		slog.Error("failed to initialize jwt service", "err", err)
		os.Exit(1)
	}

	// 3. Initialize Repositories (PostgreSQL with fallback to In-Memory for dev)
	var (
		userRepo    ports.UserRepository
		nodeRepo    ports.NodeRepository
		catalogRepo ports.CatalogRepository
	)

	pgPool, err := postgres.Connect(ctx, dbURL)
	if err != nil {
		slog.Warn("postgresql unavailable, falling back to in-memory repositories for local development", "err", err)
		userRepo = memory.NewInMemoryUserRepository()
		nodeRepo = memory.NewInMemoryNodeRepository()
		catalogRepo = memory.NewInMemoryCatalogRepository()
	} else {
		defer pgPool.Close()
		slog.Info("postgresql connection pool established")

		// Apply database schema migrations
		if err := postgres.Migrate(ctx, pgPool); err != nil {
			slog.Error("failed to apply postgres schema migrations", "err", err)
			os.Exit(1)
		}
		slog.Info("database schema migrations applied successfully")

		userRepo = postgres.NewUserRepository(pgPool)
		nodeRepo = postgres.NewNodeRepository(pgPool)
		catalogRepo = postgres.NewCatalogRepository(pgPool)
	}

	// 4. Initialize and start HTTP API server
	server := httpapi.NewServer(
		httpapi.Config{Addr: apiAddr},
		userRepo,
		nodeRepo,
		catalogRepo,
		hasher,
		tokenService,
		logger,
	)

	if err := server.Start(); err != nil {
		slog.Error("failed to start cloud api server", "err", err)
		os.Exit(1)
	}

	// 5. Wait for shutdown signal
	<-ctx.Done()
	slog.Info("shutdown signal received, commencing graceful shutdown")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("error during server shutdown", "err", err)
	}

	slog.Info("apid stopped cleanly")
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
