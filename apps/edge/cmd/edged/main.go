package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/anosh26/custardfin/apps/edge/internal/adapters/ffprobe"
	"github.com/anosh26/custardfin/apps/edge/internal/adapters/fs"
	"github.com/anosh26/custardfin/apps/edge/internal/adapters/sqlite"
	"github.com/anosh26/custardfin/apps/edge/internal/app"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	slog.Info("starting custardfin edged (edge node daemon)", "version", "0.1.0")

	// 1. Resolve configuration from environment
	mediaDir := getEnv("EDGE_MEDIA_DIR", "./media")
	dbPath := getEnv("EDGE_DB_PATH", "./data/edge.db")

	// 2. Initialize Edge SQLite Database
	db, err := sqlite.Open(ctx, dbPath)
	if err != nil {
		slog.Error("failed initializing edge sqlite store", "err", err, "path", dbPath)
		os.Exit(1)
	}
	defer db.Close()
	slog.Info("edge sqlite store initialized", "path", dbPath)

	// 3. Initialize Adapters and Services
	mediaRepo := sqlite.NewMediaRepository(db)
	fileScanner := fs.NewFileScanner()
	hasher := fs.NewHasher()
	prober := ffprobe.NewSubprocessProber("")

	scanService := app.NewScanService(fileScanner, prober, hasher, mediaRepo, logger)

	// 4. Run initial media library scan
	if _, err := os.Stat(mediaDir); err == nil {
		go func() {
			if _, err := scanService.ScanDirectory(ctx, mediaDir); err != nil {
				slog.Error("media library scan encountered errors", "err", err)
			}
		}()
	} else {
		slog.Warn("media directory does not exist, skipping scan", "path", mediaDir)
	}

	// 5. Wait for shutdown signal
	<-ctx.Done()
	slog.Info("edged shutting down gracefully")
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
