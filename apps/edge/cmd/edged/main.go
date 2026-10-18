package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/anosh26/custardfin/apps/edge/internal/adapters/ffprobe"
	"github.com/anosh26/custardfin/apps/edge/internal/adapters/fs"
	"github.com/anosh26/custardfin/apps/edge/internal/adapters/httpapi"
	"github.com/anosh26/custardfin/apps/edge/internal/adapters/sqlite"
	syncAdapter "github.com/anosh26/custardfin/apps/edge/internal/adapters/sync"
	"github.com/anosh26/custardfin/apps/edge/internal/app"
	"github.com/anosh26/custardfin/apps/edge/internal/ports"
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
	transcodeDir := getEnv("EDGE_TRANSCODE_DIR", "./transcoded")
	dbPath := getEnv("EDGE_DB_PATH", "./data/edge.db")
	listenAddr := getEnv("EDGE_ADDR", ":8090")

	// 2. Initialize Edge SQLite Database
	db, err := sqlite.Open(ctx, dbPath)
	if err != nil {
		slog.Error("failed initializing edge sqlite store", "err", err, "path", dbPath)
		os.Exit(1)
	}
	defer db.Close()
	slog.Info("edge sqlite store initialized", "path", dbPath)

	// 3. Initialize Repositories and Services
	mediaRepo := sqlite.NewMediaRepository(db)
	watchRepo := sqlite.NewWatchProgressRepository(db)

	fileScanner := fs.NewFileScanner()
	hasher := fs.NewHasher()
	prober := ffprobe.NewSubprocessProber("")

	scanService := app.NewScanService(fileScanner, prober, hasher, mediaRepo, logger)

	// 4. Initialize Cloud Manifest Syncer (if configured)
	cloudURL := getEnv("EDGE_CLOUD_URL", "http://localhost:8080")
	nodeID := getEnv("EDGE_NODE_ID", "")
	nodeToken := getEnv("EDGE_NODE_TOKEN", "")

	var syncer ports.ManifestSyncer
	if nodeID != "" && nodeToken != "" {
		var err error
		syncer, err = syncAdapter.NewHTTPManifestSyncer(cloudURL, nodeID, nodeToken, nil)
		if err != nil {
			slog.Warn("failed to initialize manifest syncer", "err", err)
		} else {
			slog.Info("cloud manifest syncer active", "cloud_url", cloudURL, "node_id", nodeID)
		}
	} else {
		slog.Warn("EDGE_NODE_ID or EDGE_NODE_TOKEN not set; running in standalone mode")
	}

	// 5. Initialize and Start Local HTTP Streaming Server
	httpServer := httpapi.NewServer(httpapi.Config{
		Addr:         listenAddr,
		TranscodeDir: transcodeDir,
	}, mediaRepo, watchRepo, syncer, logger)

	if err := httpServer.Start(); err != nil {
		slog.Error("failed starting edge local http server", "err", err)
		os.Exit(1)
	}

	// 5. Run initial media library scan in background
	if _, err := os.Stat(mediaDir); err == nil {
		go func() {
			if _, err := scanService.ScanDirectory(ctx, mediaDir); err != nil {
				slog.Error("media library scan encountered errors", "err", err)
			}
		}()
	} else {
		slog.Warn("media directory does not exist, skipping scan", "path", mediaDir)
	}

	// 6. Wait for shutdown signal
	<-ctx.Done()
	slog.Info("shutdown signal received, commencing graceful shutdown")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		slog.Error("error during http server shutdown", "err", err)
	}

	slog.Info("edged stopped cleanly")
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
