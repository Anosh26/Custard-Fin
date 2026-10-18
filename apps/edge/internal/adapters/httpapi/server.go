package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/anosh26/custardfin/apps/edge/internal/domain"
	"github.com/anosh26/custardfin/apps/edge/internal/ports"
)

var (
	ErrInvalidPath = errors.New("invalid or escaped file path")
)

// Config holds configuration parameters for the edge HTTP server.
type Config struct {
	Addr         string
	TranscodeDir string
}

// Server provides the local HTTP streaming and metadata API for the home edge daemon.
type Server struct {
	httpServer   *http.Server
	mediaRepo    ports.MediaRepository
	watchRepo    ports.WatchProgressRepository
	transcodeDir string
	logger       *slog.Logger
	listener     net.Listener
}

// NewServer constructs and wires the edge HTTP server.
func NewServer(cfg Config, mediaRepo ports.MediaRepository, watchRepo ports.WatchProgressRepository, logger *slog.Logger) *Server {
	if cfg.Addr == "" {
		cfg.Addr = ":8090"
	}
	if cfg.TranscodeDir == "" {
		cfg.TranscodeDir = "./transcoded"
	}
	if logger == nil {
		logger = slog.Default()
	}

	s := &Server{
		mediaRepo:    mediaRepo,
		watchRepo:    watchRepo,
		transcodeDir: filepath.Clean(cfg.TranscodeDir),
		logger:       logger,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", s.handleHealthz)
	mux.HandleFunc("GET /api/media", s.handleListMedia)
	mux.HandleFunc("GET /api/media/{id}", s.handleGetMedia)
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
			"media_id":    mediaID,
			"position_s":  0,
			"completed":   false,
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
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, HEAD")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Range, Authorization")
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
