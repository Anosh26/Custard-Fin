package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/anosh26/custardfin/apps/edge/internal/adapters/ffmpeg"
	"github.com/anosh26/custardfin/apps/edge/internal/adapters/ffprobe"
	"github.com/anosh26/custardfin/apps/edge/internal/adapters/frp"
	"github.com/anosh26/custardfin/apps/edge/internal/adapters/fs"
	"github.com/anosh26/custardfin/apps/edge/internal/adapters/heartbeat"
	"github.com/anosh26/custardfin/apps/edge/internal/adapters/httpapi"
	"github.com/anosh26/custardfin/apps/edge/internal/adapters/sqlite"
	syncAdapter "github.com/anosh26/custardfin/apps/edge/internal/adapters/sync"
	"github.com/anosh26/custardfin/apps/edge/internal/app"
	"github.com/anosh26/custardfin/apps/edge/internal/ports"
	"github.com/anosh26/custardfin/libs/contracts"
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

	// 4. Initialize Cloud Manifest Syncer (with demo fallback credentials)
	cloudURL := getEnv("EDGE_CLOUD_URL", getEnv("EDGE_CLOUD_API_URL", "http://localhost:8080"))
	cloudURL = strings.TrimSuffix(cloudURL, "/")
	cloudURL = strings.TrimSuffix(cloudURL, "/v1")
	nodeID := getEnv("EDGE_NODE_ID", "898ea0c7-f7a3-48e8-ba59-b378a219f9ee")
	nodeToken := getEnv("EDGE_NODE_TOKEN", "1c94b7b8fd38a97cd634b4894358fc1feed28f3a180a104aa0b83f9710190de2")
	ffmpegBin := getEnv("EDGE_FFMPEG_BIN", "ffmpeg")
	transcoder := ffmpeg.NewSubprocessTranscoder(ffmpegBin, logger)

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
		MediaDir:     mediaDir,
		NodeID:       nodeID,
	}, mediaRepo, watchRepo, syncer, logger)
	httpServer.SetScanService(scanService)
	httpServer.SetTranscoder(transcoder)

	if err := httpServer.Start(); err != nil {
		slog.Error("failed starting edge local http server", "err", err)
		os.Exit(1)
	}


	// 5. Run initial media library scan and auto-transcoding in background
	if _, err := os.Stat(mediaDir); err == nil {
		go func() {
			if _, err := scanService.ScanDirectory(ctx, mediaDir); err != nil {
				slog.Error("media library scan encountered errors", "err", err)
			}

			// Automatically ensure all scanned media fixtures have 4-second HLS renditions
			mediaList, err := mediaRepo.List(ctx)
			if err == nil {
				for _, m := range mediaList {
					outDir := filepath.Join(transcodeDir, m.ID)
					masterPath := filepath.Join(outDir, "master.m3u8")
					if _, err := os.Stat(masterPath); os.IsNotExist(err) {
						slog.Info("transcoding media fixture to 4s HLS chunks", "title", m.Title, "id", m.ID)
						_, err := transcoder.Transcode(ctx, ports.TranscodeJob{
							MediaID:      m.ID,
							InputPath:    m.Path,
							OutputDir:    outDir,
							SourceWidth:  m.Width,
							SourceHeight: m.Height,
						})
						if err != nil {
							slog.Error("transcoding fixture failed", "title", m.Title, "err", err)
						}
					}

					// Ensure media items are set to public and synced to cloud catalog
					if syncer != nil {
						_ = mediaRepo.SetPublic(ctx, m.ID, true)
						m.IsPublic = true
						renditions := []contracts.Rendition{
							{Resolution: "1080p", BitrateBps: 5000000, PlaylistPath: "1080p.m3u8", ChunkDurationS: 4},
							{Resolution: "720p", BitrateBps: 2800000, PlaylistPath: "720p.m3u8", ChunkDurationS: 4},
							{Resolution: "480p", BitrateBps: 1200000, PlaylistPath: "480p.m3u8", ChunkDurationS: 4},
						}
						if err := syncer.SyncManifest(ctx, m, renditions); err != nil {
							slog.Warn("syncing media manifest to cloud catalog failed", "title", m.Title, "err", err)
						} else {
							slog.Info("manifest synced to cloud catalog", "title", m.Title, "id", m.ID)
						}
					}
				}
			}
		}()

	} else {
		slog.Warn("media directory does not exist, skipping scan", "path", mediaDir)
	}

	// 6. Initialize Outbound FRP Reverse Tunnel Client (if configured)
	frpEnabled := getEnv("EDGE_FRP_ENABLED", "false") == "true"
	frpServerAddr := getEnv("EDGE_FRP_SERVER_ADDR", "localhost")
	frpServerPortStr := getEnv("EDGE_FRP_SERVER_PORT", "7000")
	frpAuthToken := getEnv("EDGE_FRP_AUTH_TOKEN", "")
	frpBinPath := getEnv("EDGE_FRP_BIN", "frpc")

	var tunnelClient ports.TunnelClient
	if frpEnabled && nodeID != "" && frpAuthToken != "" {
		frpPort, _ := strconv.Atoi(frpServerPortStr)
		if frpPort <= 0 {
			frpPort = 7000
		}

		localPort := 8090
		parts := strings.Split(listenAddr, ":")
		if len(parts) > 1 {
			if p, err := strconv.Atoi(parts[len(parts)-1]); err == nil && p > 0 {
				localPort = p
			}
		}

		tunnelClient = frp.NewManagedFRPClient(frp.Config{
			ServerAddr: frpServerAddr,
			ServerPort: frpPort,
			AuthToken:  frpAuthToken,
			NodeID:     nodeID,
			LocalPort:  localPort,
			BinaryPath: frpBinPath,
			UseTLS:     true,
		}, logger)

		if err := tunnelClient.Start(ctx); err != nil {
			slog.Warn("failed starting frp reverse tunnel client", "err", err)
		} else {
			slog.Info("outbound frp reverse tunnel client active", "relay", frpServerAddr, "node_id", nodeID)
		}
	} else if frpEnabled {
		slog.Warn("EDGE_FRP_ENABLED is true but EDGE_NODE_ID or EDGE_FRP_AUTH_TOKEN is missing")
	}

	// 7. Initialize Outbound Periodic Heartbeat Worker (if node credentials configured)
	var hbWorker ports.HeartbeatReporter
	if nodeID != "" && nodeToken != "" {
		heartbeatIntervalStr := getEnv("EDGE_HEARTBEAT_INTERVAL_SEC", "15")
		heartbeatIntervalSec, _ := strconv.Atoi(heartbeatIntervalStr)
		if heartbeatIntervalSec <= 0 {
			heartbeatIntervalSec = 15
		}

		// Probe GPU encoding capabilities (NVENC / QSV / VAAPI)
		caps := transcoder.ProbeCapabilities(ctx)
		var gpuCaps *contracts.GPUCapabilities
		if caps.NVENC || caps.QSV || caps.VAAPI {
			gpuCaps = &caps
		}

		worker, err := heartbeat.NewWorker(heartbeat.Config{
			CloudBaseURL:  cloudURL,
			NodeID:        nodeID,
			NodeToken:     nodeToken,
			Version:       "0.1.0",
			Interval:      time.Duration(heartbeatIntervalSec) * time.Second,
			GPUCaps:       gpuCaps,
			StreamTracker: httpServer,
			Logger:        logger,
		})
		if err != nil {
			slog.Warn("failed initializing edge heartbeat worker", "err", err)
		} else {
			if err := worker.Start(ctx); err != nil {
				slog.Warn("failed starting edge heartbeat worker", "err", err)
			} else {
				hbWorker = worker
				slog.Info("outbound heartbeat worker active",
					"interval_s", heartbeatIntervalSec,
					"cloud_url", cloudURL,
					"nvenc", gpuCaps != nil && gpuCaps.NVENC,
				)
			}
		}
	}

	// 8. Wait for shutdown signal
	<-ctx.Done()
	slog.Info("shutdown signal received, commencing graceful shutdown")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if hbWorker != nil {
		if err := hbWorker.Stop(); err != nil {
			slog.Error("error during heartbeat worker shutdown", "err", err)
		}
	}

	if tunnelClient != nil {
		if err := tunnelClient.Stop(); err != nil {
			slog.Error("error during frp tunnel client shutdown", "err", err)
		}
	}

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
