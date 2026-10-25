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
	redisAdapter "github.com/anosh26/custardfin/apps/api/internal/adapters/redis"
	s3Adapter "github.com/anosh26/custardfin/apps/api/internal/adapters/s3"
	tunnelAdapter "github.com/anosh26/custardfin/apps/api/internal/adapters/tunnel"
	"github.com/anosh26/custardfin/apps/api/internal/ports"
	redisClient "github.com/redis/go-redis/v9"
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
	redisAddr := getEnv("REDIS_ADDR", "localhost:6379")
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

	if dbURL == "memory" || dbURL == "none" || dbURL == "" {
		slog.Info("using in-memory repositories for local development/testing")
		userRepo = memory.NewInMemoryUserRepository()
		nodeRepo = memory.NewInMemoryNodeRepository()
		catalogRepo = memory.NewInMemoryCatalogRepository()
	} else {
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
	}

	// 4. Initialize Redis Heartbeat Store (Invariant 5) with In-Memory fallback
	var heartbeatStore ports.HeartbeatStore
	if redisAddr == "memory" || redisAddr == "none" || redisAddr == "" {
		slog.Info("using in-memory heartbeat store for local development/testing")
		heartbeatStore = memory.NewInMemoryHeartbeatStore()
	} else {
		rdb := redisClient.NewClient(&redisClient.Options{
			Addr: redisAddr,
		})
		if err := rdb.Ping(ctx).Err(); err != nil {
			slog.Warn("redis unavailable, falling back to in-memory heartbeat store for local development", "err", err)
			heartbeatStore = memory.NewInMemoryHeartbeatStore()
		} else {
			defer rdb.Close()
			slog.Info("connected to redis for node heartbeats", "addr", redisAddr)
			heartbeatStore = redisAdapter.NewHeartbeatStore(rdb)
		}
	}

	// 5. Initialize Origin Shield Object Cache (Invariant 1, Invariant 7)
	var chunkCache ports.ChunkCache
	s3Bucket := getEnv("S3_BUCKET", "custardfin-origin-cache")
	s3Endpoint := getEnv("S3_ENDPOINT", "")
	s3Region := getEnv("S3_REGION", "us-east-1")
	s3AccessKey := getEnv("S3_ACCESS_KEY", "minioadmin")
	s3SecretKey := getEnv("S3_SECRET_KEY", "minioadmin")

	if s3Endpoint != "" || os.Getenv("S3_BUCKET") != "" {
		s3Adapter, err := s3Adapter.NewS3ChunkCache(s3Adapter.Config{
			Endpoint:        s3Endpoint,
			Bucket:          s3Bucket,
			Region:          s3Region,
			AccessKeyID:     s3AccessKey,
			SecretAccessKey: s3SecretKey,
			UsePathStyle:    true,
		})
		if err != nil {
			slog.Warn("failed to initialize s3 chunk cache, falling back to memory cache", "err", err)
			chunkCache = memory.NewMemoryChunkCache(24 * time.Hour)
		} else {
			slog.Info("connected to s3-compatible origin shield cache", "bucket", s3Bucket, "endpoint", s3Endpoint)
			chunkCache = s3Adapter
		}
	} else {
		slog.Info("s3 endpoint not configured, using in-memory origin shield chunk cache")
		chunkCache = memory.NewMemoryChunkCache(24 * time.Hour)
	}

	// 6. Initialize Reverse Tunnel Puller (Invariant 2, Invariant 7)
	tunnelBaseURL := getEnv("TUNNEL_BASE_URL", "http://localhost:8090")
	tunnelPuller, err := tunnelAdapter.NewHTTPTunnelPuller(tunnelAdapter.Config{
		BaseURL: tunnelBaseURL,
		Timeout: 10 * time.Second,
	})
	if err != nil {
		slog.Error("failed to initialize tunnel puller", "err", err)
		os.Exit(1)
	}

	signingSecret := getEnv("SIGNING_SECRET", "custardfin_development_chunk_signing_secret_key_32b!!")
	baseURL := getEnv("API_BASE_URL", "http://localhost:8080")

	// 7. Initialize and start HTTP API server
	server := httpapi.NewServer(
		httpapi.Config{
			Addr:          apiAddr,
			BaseURL:       baseURL,
			SigningSecret: []byte(signingSecret),
		},
		userRepo,
		nodeRepo,
		catalogRepo,
		heartbeatStore,
		chunkCache,
		tunnelPuller,
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
