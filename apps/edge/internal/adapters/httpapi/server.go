package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/anosh26/custardfin/apps/edge/internal/app"
	"github.com/anosh26/custardfin/apps/edge/internal/domain"
	"github.com/anosh26/custardfin/apps/edge/internal/ports"
	"github.com/anosh26/custardfin/libs/contracts"
)

var (
	ErrInvalidPath = errors.New("invalid or escaped file path")
)

var _ ports.StreamTracker = (*Server)(nil)

// Config holds configuration parameters for the edge HTTP server.
type Config struct {
	Addr         string
	TranscodeDir string
	MediaDir     string
	NodeID       string
}

// Server provides the local HTTP streaming and metadata API for the home edge daemon.
type Server struct {
	httpServer   *http.Server
	mediaRepo    ports.MediaRepository
	watchRepo    ports.WatchProgressRepository
	syncer       ports.ManifestSyncer
	transcodeDir string
	mediaDir     string
	nodeID       string
	scanService  *app.ScanService
	transcoder   ports.Transcoder
	logger       *slog.Logger
	listener     net.Listener
	streamMu     sync.Mutex
	sessions     map[string]time.Time
}

// NewServer constructs and wires the edge HTTP server.
func NewServer(cfg Config, mediaRepo ports.MediaRepository, watchRepo ports.WatchProgressRepository, syncer ports.ManifestSyncer, logger *slog.Logger) *Server {
	if cfg.Addr == "" {
		cfg.Addr = ":8090"
	}
	if cfg.TranscodeDir == "" {
		cfg.TranscodeDir = "./transcoded"
	}
	if cfg.MediaDir == "" {
		cfg.MediaDir = "./media"
	}
	if logger == nil {
		logger = slog.Default()
	}

	s := &Server{
		mediaRepo:    mediaRepo,
		watchRepo:    watchRepo,
		syncer:       syncer,
		transcodeDir: filepath.Clean(cfg.TranscodeDir),
		mediaDir:     filepath.Clean(cfg.MediaDir),
		nodeID:       cfg.NodeID,
		logger:       logger,
		sessions:     make(map[string]time.Time),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", s.handleHealthz)
	mux.HandleFunc("GET /api/status", s.handleStatus)
	mux.HandleFunc("GET /api/media", s.handleListMedia)
	mux.HandleFunc("GET /api/media/{id}", s.handleGetMedia)
	mux.HandleFunc("PUT /api/media/{id}/public", s.handleTogglePublic)
	mux.HandleFunc("POST /api/media/upload", s.handleUploadMedia)
	mux.HandleFunc("POST /api/media/scan", s.handleScanMedia)
	mux.HandleFunc("POST /api/media/generate-sample", s.handleGenerateSample)
	mux.HandleFunc("GET /stream/{id}/{path...}", s.handleStream)
	mux.HandleFunc("GET /api/watch/{id}", s.handleGetWatchProgress)
	mux.HandleFunc("POST /api/watch/{id}", s.handleSaveWatchProgress)

	// Apply CORS and logging middleware
	handler := s.withRecovery(s.withCORS(s.withLogging(mux)))

	s.httpServer = &http.Server{
		Addr:              cfg.Addr,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	return s
}

// Start binds to the configured network address and begins serving requests.
func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.httpServer.Addr)
	if err != nil {
		return fmt.Errorf("failed to bind listener on %s: %w", s.httpServer.Addr, err)
	}
	s.listener = ln
	s.logger.Info("edge local streaming server started", "addr", ln.Addr().String())

	go func() {
		if err := s.httpServer.Serve(ln); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Error("edge http server encountered fatal error", "err", err)
		}
	}()

	return nil
}

// Addr returns the bound network address of the server.
func (s *Server) Addr() string {
	if s.listener != nil {
		return s.listener.Addr().String()
	}
	return s.httpServer.Addr
}

// Shutdown gracefully terminates the HTTP server.
func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("edge local streaming server shutting down")
	return s.httpServer.Shutdown(ctx)
}

func (s *Server) handleHealthz(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) handleListMedia(w http.ResponseWriter, r *http.Request) {
	list, err := s.mediaRepo.List(r.Context())
	if err != nil {
		s.logger.Error("failed listing media", "err", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	if list == nil {
		list = []*domain.Media{}
	}
	writeJSON(w, http.StatusOK, list)
}

func (s *Server) handleGetMedia(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	m, err := s.mediaRepo.FindByID(r.Context(), id)
	if err != nil {
		s.logger.Error("failed fetching media", "id", id, "err", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	if m == nil {
		http.Error(w, "media not found", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, m)
}

func (s *Server) handleStream(w http.ResponseWriter, r *http.Request) {
	mediaID := r.PathValue("id")
	subPath := r.PathValue("path")

	if strings.TrimSpace(mediaID) == "" || strings.TrimSpace(subPath) == "" {
		http.Error(w, "invalid stream path", http.StatusBadRequest)
		return
	}

	// Verify media exists
	m, err := s.mediaRepo.FindByID(r.Context(), mediaID)
	if err != nil || m == nil {
		http.Error(w, "media not found", http.StatusNotFound)
		return
	}

	// Prevent directory traversal
	cleanSub := filepath.Clean(filepath.FromSlash(subPath))
	if strings.HasPrefix(cleanSub, "..") || filepath.IsAbs(cleanSub) {
		http.Error(w, "forbidden path", http.StatusBadRequest)
		return
	}

	mediaTranscodeDir := filepath.Join(s.transcodeDir, mediaID)
	targetFile := filepath.Join(mediaTranscodeDir, cleanSub)

	rel, err := filepath.Rel(mediaTranscodeDir, targetFile)
	if err != nil || strings.HasPrefix(rel, "..") {
		http.Error(w, "forbidden path traversal", http.StatusBadRequest)
		return
	}

	info, err := os.Stat(targetFile)
	if err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "file not found", http.StatusNotFound)
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	// Set proper HLS MIME types and caching headers
	ext := strings.ToLower(filepath.Ext(targetFile))
	switch ext {
	case ".m3u8":
		w.Header().Set("Content-Type", "application/vnd.apple.mpegurl")
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	case ".ts":
		w.Header().Set("Content-Type", "video/mp2t")
		// Immutable 4-second chunk segments can be cached indefinitely
		w.Header().Set("Cache-Control", "public, max-age=86400, immutable")
	}

	w.Header().Set("Accept-Ranges", "bytes")

	// Track active stream session per remote client and media ID
	remoteHost, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		remoteHost = r.RemoteAddr
	}
	s.recordStreamActivity(remoteHost + ":" + mediaID)

	http.ServeFile(w, r, targetFile)
	_ = info
}

func (s *Server) handleGetWatchProgress(w http.ResponseWriter, r *http.Request) {
	mediaID := r.PathValue("id")
	wp, err := s.watchRepo.GetWatchProgress(r.Context(), mediaID)
	if err != nil {
		s.logger.Error("failed fetching watch progress", "mediaID", mediaID, "err", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	if wp == nil {
		// Default zero progress
		writeJSON(w, http.StatusOK, map[string]any{
			"media_id":   mediaID,
			"position_s": 0,
			"completed":  false,
		})
		return
	}
	writeJSON(w, http.StatusOK, wp)
}

type saveWatchRequest struct {
	PositionS int  `json:"position_s"`
	Completed bool `json:"completed"`
}

func (s *Server) handleSaveWatchProgress(w http.ResponseWriter, r *http.Request) {
	mediaID := r.PathValue("id")
	var req saveWatchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	wp, err := domain.NewWatchProgress(mediaID, req.PositionS, req.Completed)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := s.watchRepo.SaveWatchProgress(r.Context(), wp); err != nil {
		s.logger.Error("failed saving watch progress", "mediaID", mediaID, "err", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, wp)
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

// Middleware: CORS
func (s *Server) withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS, HEAD")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Range, Authorization, X-Node-Token")
		w.Header().Set("Access-Control-Expose-Headers", "Content-Length, Content-Range, Accept-Ranges")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}


// Middleware: Logging
func (s *Server) withLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		s.logger.Debug("http request",
			"method", r.Method,
			"path", r.URL.Path,
			"duration", time.Since(start),
		)
	})
}

// Middleware: Panic Recovery
func (s *Server) withRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				s.logger.Error("recovered from panic in http handler", "err", rec)
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

type togglePublicRequest struct {
	IsPublic bool `json:"is_public"`
}

func (s *Server) handleTogglePublic(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Missing media id"})
		return
	}

	var req togglePublicRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid JSON payload"})
		return
	}

	media, err := s.mediaRepo.FindByID(r.Context(), id)
	if err != nil || media == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "Media not found"})
		return
	}

	// Update SQLite
	if err := s.mediaRepo.SetPublic(r.Context(), id, req.IsPublic); err != nil {
		s.logger.Error("failed updating is_public", "id", id, "err", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to update media status"})
		return
	}
	media.IsPublic = req.IsPublic

	// If syncer is configured, sync to Cloud API
	if s.syncer != nil {
		if req.IsPublic {
			// Provide renditions (1080p, 720p, 480p) conforming to Invariant 6
			renditions := []contracts.Rendition{
				{Resolution: "1080p", BitrateBps: 5000000, PlaylistPath: "1080p.m3u8", ChunkDurationS: 4},
				{Resolution: "720p", BitrateBps: 2800000, PlaylistPath: "720p.m3u8", ChunkDurationS: 4},
				{Resolution: "480p", BitrateBps: 1200000, PlaylistPath: "480p.m3u8", ChunkDurationS: 4},
			}
			if err := s.syncer.SyncManifest(r.Context(), media, renditions); err != nil {
				s.logger.Error("cloud manifest sync error", "id", id, "err", err)
				writeJSON(w, http.StatusBadGateway, map[string]string{"error": fmt.Sprintf("Saved locally, but cloud sync failed: %v", err)})
				return
			}
		} else {
			if err := s.syncer.UnsyncManifest(r.Context(), id); err != nil {
				s.logger.Error("cloud manifest unsync error", "id", id, "err", err)
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"media_id":  id,
		"is_public": req.IsPublic,
		"status":    "updated",
	})
}

func (s *Server) recordStreamActivity(sessionKey string) {
	s.streamMu.Lock()
	defer s.streamMu.Unlock()
	s.sessions[sessionKey] = time.Now()
}

// ActiveStreams returns the number of active distinct stream sessions within the last 15 seconds.
func (s *Server) ActiveStreams() int {
	s.streamMu.Lock()
	defer s.streamMu.Unlock()

	now := time.Now()
	activeThreshold := 15 * time.Second
	for key, lastSeen := range s.sessions {
		if now.Sub(lastSeen) > activeThreshold {
			delete(s.sessions, key)
		}
	}
	return len(s.sessions)
}

// SetScanService assigns the scan service for directory library rescans.
func (s *Server) SetScanService(svc *app.ScanService) {
	s.scanService = svc
}

// SetTranscoder assigns the transcoder for on-demand/upload transcoding.
func (s *Server) SetTranscoder(t ports.Transcoder) {
	s.transcoder = t
}

func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	mediaList, _ := s.mediaRepo.List(r.Context())
	count := 0
	if mediaList != nil {
		count = len(mediaList)
	}
	hasNVENC := false
	if s.transcoder != nil {
		enc, _ := s.transcoder.ProbeEncoder(r.Context())
		hasNVENC = strings.Contains(enc, "nvenc")
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"node_id":        s.nodeID,
		"status":         "online",
		"media_dir":      s.mediaDir,
		"transcode_dir":  s.transcodeDir,
		"media_count":    count,
		"active_streams": s.ActiveStreams(),
		"max_streams":    10,
		"nvenc_enabled":  hasNVENC,
	})
}

func (s *Server) handleUploadMedia(w http.ResponseWriter, r *http.Request) {
	// Limit request payload to 500 MB
	r.Body = http.MaxBytesReader(w, r.Body, 500<<20)
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Failed to parse form or file exceeds 500 MB limit"})
		return
	}

	file, header, err := r.FormFile("video")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Missing 'video' file in form payload"})
		return
	}
	defer file.Close()

	cleanName := filepath.Base(header.Filename)
	ext := strings.ToLower(filepath.Ext(cleanName))
	if ext != ".mp4" && ext != ".mkv" && ext != ".mov" && ext != ".webm" && ext != ".avi" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Unsupported video format; please upload .mp4, .mkv, .mov, or .webm"})
		return
	}

	if s.mediaDir == "" {
		s.mediaDir = "./media"
	}
	_ = os.MkdirAll(s.mediaDir, 0755)

	targetPath := filepath.Join(s.mediaDir, cleanName)
	dst, err := os.Create(targetPath)
	if err != nil {
		s.logger.Error("failed creating target media file", "path", targetPath, "err", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed creating file on disk"})
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		s.logger.Error("failed writing uploaded video payload", "err", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed writing video payload"})
		return
	}

	s.logger.Info("video file uploaded successfully to edge node", "filename", cleanName, "size_bytes", header.Size)

	// Trigger library scan, transcode, and catalog sync in background
	if s.scanService != nil {
		go func() {
			bgCtx := context.Background()
			_, _ = s.scanService.ScanDirectory(bgCtx, s.mediaDir)
			mediaItem, err := s.mediaRepo.FindByPath(bgCtx, targetPath)
			if err == nil && mediaItem != nil && s.transcoder != nil {
				outDir := filepath.Join(s.transcodeDir, mediaItem.ID)
				_, _ = s.transcoder.Transcode(bgCtx, ports.TranscodeJob{
					MediaID:      mediaItem.ID,
					InputPath:    mediaItem.Path,
					OutputDir:    outDir,
					SourceWidth:  mediaItem.Width,
					SourceHeight: mediaItem.Height,
				})
				if s.syncer != nil {
					_ = s.mediaRepo.SetPublic(bgCtx, mediaItem.ID, true)
					mediaItem.IsPublic = true
					renditions := []contracts.Rendition{
						{Resolution: "1080p", BitrateBps: 5000000, PlaylistPath: "1080p.m3u8", ChunkDurationS: 4},
						{Resolution: "720p", BitrateBps: 2800000, PlaylistPath: "720p.m3u8", ChunkDurationS: 4},
						{Resolution: "480p", BitrateBps: 1200000, PlaylistPath: "480p.m3u8", ChunkDurationS: 4},
					}
					_ = s.syncer.SyncManifest(bgCtx, mediaItem, renditions)
				}
			}
		}()
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"status":   "ok",
		"filename": cleanName,
		"size":     header.Size,
		"message":  "Video uploaded to edge node. Transcoding and catalog publishing in progress.",
	})
}

func (s *Server) handleScanMedia(w http.ResponseWriter, r *http.Request) {
	if s.scanService == nil {
		writeJSON(w, http.StatusOK, map[string]any{"status": "ok", "message": "No scan service attached"})
		return
	}

	if s.mediaDir == "" {
		s.mediaDir = "./media"
	}
	result, err := s.scanService.ScanDirectory(r.Context(), s.mediaDir)
	if err != nil {
		s.logger.Error("manual scan failed", "err", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	go func() {
		bgCtx := context.Background()
		items, _ := s.mediaRepo.List(bgCtx)
		for _, m := range items {
			outDir := filepath.Join(s.transcodeDir, m.ID)
			masterPath := filepath.Join(outDir, "master.m3u8")
			if _, err := os.Stat(masterPath); os.IsNotExist(err) && s.transcoder != nil {
				_, _ = s.transcoder.Transcode(bgCtx, ports.TranscodeJob{
					MediaID:      m.ID,
					InputPath:    m.Path,
					OutputDir:    outDir,
					SourceWidth:  m.Width,
					SourceHeight: m.Height,
				})
			}
			if s.syncer != nil && m.IsPublic {
				renditions := []contracts.Rendition{
					{Resolution: "1080p", BitrateBps: 5000000, PlaylistPath: "1080p.m3u8", ChunkDurationS: 4},
					{Resolution: "720p", BitrateBps: 2800000, PlaylistPath: "720p.m3u8", ChunkDurationS: 4},
					{Resolution: "480p", BitrateBps: 1200000, PlaylistPath: "480p.m3u8", ChunkDurationS: 4},
				}
				_ = s.syncer.SyncManifest(bgCtx, m, renditions)
			}
		}
	}()

	writeJSON(w, http.StatusOK, map[string]any{
		"status":    "ok",
		"total":     result.TotalScanned,
		"added":     result.Added,
		"updated":   result.Updated,
		"unchanged": result.Unchanged,
	})
}

func (s *Server) handleGenerateSample(w http.ResponseWriter, r *http.Request) {
	if s.mediaDir == "" {
		s.mediaDir = "./media"
	}
	_ = os.MkdirAll(s.mediaDir, 0755)

	timestamp := time.Now().Format("20060102_150405")
	fileName := fmt.Sprintf("sample_demo_%s.mp4", timestamp)
	targetPath := filepath.Join(s.mediaDir, fileName)

	cmd := exec.CommandContext(r.Context(), "ffmpeg", "-y", "-v", "error",
		"-f", "lavfi", "-i", "testsrc2=duration=15:size=1280x720:rate=30",
		"-f", "lavfi", "-i", "sine=frequency=1000:duration=15",
		"-pix_fmt", "yuv420p", "-c:v", "libx264", "-preset", "ultrafast",
		"-c:a", "aac", "-b:a", "128k", "-shortest",
		targetPath,
	)

	if err := cmd.Run(); err != nil {
		s.logger.Error("failed generating sample test video clip", "err", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed synthesizing sample clip with ffmpeg"})
		return
	}

	if s.scanService != nil {
		_, _ = s.scanService.ScanDirectory(r.Context(), s.mediaDir)
		mediaItem, err := s.mediaRepo.FindByPath(r.Context(), targetPath)
		if err == nil && mediaItem != nil && s.transcoder != nil {
			outDir := filepath.Join(s.transcodeDir, mediaItem.ID)
			_, _ = s.transcoder.Transcode(r.Context(), ports.TranscodeJob{
				MediaID:      mediaItem.ID,
				InputPath:    mediaItem.Path,
				OutputDir:    outDir,
				SourceWidth:  mediaItem.Width,
				SourceHeight: mediaItem.Height,
			})
			if s.syncer != nil {
				_ = s.mediaRepo.SetPublic(r.Context(), mediaItem.ID, true)
				mediaItem.IsPublic = true
				renditions := []contracts.Rendition{
					{Resolution: "1080p", BitrateBps: 5000000, PlaylistPath: "1080p.m3u8", ChunkDurationS: 4},
					{Resolution: "720p", BitrateBps: 2800000, PlaylistPath: "720p.m3u8", ChunkDurationS: 4},
					{Resolution: "480p", BitrateBps: 1200000, PlaylistPath: "480p.m3u8", ChunkDurationS: 4},
				}
				_ = s.syncer.SyncManifest(r.Context(), mediaItem, renditions)
			}
		}
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"status":   "ok",
		"filename": fileName,
		"message":  "Sample clip synthesized, transcoded to 4s HLS, and published to community catalog!",
	})
}

