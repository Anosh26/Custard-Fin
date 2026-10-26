package httpapi

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/anosh26/custardfin/apps/api/internal/adapters/admission"
	"github.com/anosh26/custardfin/apps/api/internal/adapters/memory"
	"github.com/anosh26/custardfin/apps/api/internal/adapters/prefetch"
	"github.com/anosh26/custardfin/apps/api/internal/domain"
	"github.com/anosh26/custardfin/apps/api/internal/ports"
	"github.com/anosh26/custardfin/libs/contracts"
)

type Config struct {
	Addr          string
	HeartbeatTTL  time.Duration
	BaseURL       string // Cloud API base URL, e.g. "http://localhost:8080"
	SigningSecret []byte // Secret key for HMAC-SHA256 chunk signing (Invariant 9)
}

type Server struct {
	cfg                 Config
	userRepo            ports.UserRepository
	nodeRepo            ports.NodeRepository
	catalogRepo         ports.CatalogRepository
	heartbeatStore      ports.HeartbeatStore
	chunkCache          ports.ChunkCache
	tunnelPuller        ports.TunnelPuller
	saturationEstimator ports.SaturationEstimator
	prefetcher          ports.Prefetcher
	hasher              ports.PasswordHasher
	tokenService        ports.TokenService
	logger              *slog.Logger
	server              *http.Server
	listener            net.Listener
}

func NewServer(
	cfg Config,
	userRepo ports.UserRepository,
	nodeRepo ports.NodeRepository,
	catalogRepo ports.CatalogRepository,
	heartbeatStore ports.HeartbeatStore,
	chunkCache ports.ChunkCache,
	tunnelPuller ports.TunnelPuller,
	hasher ports.PasswordHasher,
	tokenService ports.TokenService,
	logger *slog.Logger,
) *Server {
	if cfg.Addr == "" {
		cfg.Addr = ":8080"
	}
	if cfg.HeartbeatTTL <= 0 {
		cfg.HeartbeatTTL = 20 * time.Second
	}
	if cfg.BaseURL == "" {
		cfg.BaseURL = "http://localhost:8080"
	}
	if len(cfg.SigningSecret) == 0 {
		cfg.SigningSecret = []byte("custardfin-default-origin-shield-secret-32b!")
	}
	if logger == nil {
		logger = slog.Default()
	}
	if chunkCache == nil {
		chunkCache = memory.NewMemoryChunkCache(24 * time.Hour)
	}

	var estimator ports.SaturationEstimator
	if heartbeatStore != nil {
		estimator = admission.NewDefaultEstimator(admission.Config{}, heartbeatStore, logger)
	}

	var prefetcher ports.Prefetcher
	if chunkCache != nil && tunnelPuller != nil {
		if pf, err := prefetch.NewBackgroundPrefetcher(prefetch.Config{}, chunkCache, tunnelPuller, logger); err == nil {
			prefetcher = pf
		}
	}

	s := &Server{
		cfg:                 cfg,
		userRepo:            userRepo,
		nodeRepo:            nodeRepo,
		catalogRepo:         catalogRepo,
		heartbeatStore:      heartbeatStore,
		chunkCache:          chunkCache,
		tunnelPuller:        tunnelPuller,
		saturationEstimator: estimator,
		prefetcher:          prefetcher,
		hasher:              hasher,
		tokenService:        tokenService,
		logger:              logger,
	}

	mux := http.NewServeMux()
	s.registerRoutes(mux)

	// Apply middlewares (Recover -> Logging -> CORS)
	handler := RecoverMiddleware(logger)(
		LoggingMiddleware(logger)(
			CORSMiddleware(mux),
		),
	)

	s.server = &http.Server{
		Addr:              cfg.Addr,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	return s
}

func (s *Server) registerRoutes(mux *http.ServeMux) {
	// Health
	mux.HandleFunc("GET /healthz", s.handleHealthz)

	// Auth
	mux.HandleFunc("POST /v1/auth/register", s.handleRegister)
	mux.HandleFunc("POST /v1/auth/login", s.handleLogin)
	mux.HandleFunc("GET /v1/auth/me", s.handleGetMe)

	// Nodes (Requires JWT Bearer Auth)
	mux.Handle("POST /v1/nodes", RequireUserAuth(s.tokenService)(http.HandlerFunc(s.handleRegisterNode)))

	// Node Heartbeat & Manifest Sync (Requires X-Node-Token Auth)
	nodeAuth := RequireNodeAuth(s.nodeRepo)
	mux.Handle("POST /v1/nodes/{id}/heartbeat", nodeAuth(http.HandlerFunc(s.handleNodeHeartbeat)))
	mux.Handle("PUT /v1/nodes/{id}/media/{mid}/manifest", nodeAuth(http.HandlerFunc(s.handleUpsertManifest)))
	mux.Handle("DELETE /v1/nodes/{id}/media/{mid}/manifest", nodeAuth(http.HandlerFunc(s.handleDeleteManifest)))

	// Node Status & Liveness (Public for Web UI, Dev Console, HUD & Invariant Checks)
	mux.HandleFunc("GET /v1/nodes/{id}/status", s.handleGetNodeStatus)

	// Catalog Search (Public)
	mux.HandleFunc("GET /v1/catalog", s.handleGetCatalog)

	// Playback Route Resolution (Public)
	mux.HandleFunc("GET /v1/media/{id}/playback", s.handleGetPlayback)

	// Thumbnail Poster Image (Public via Origin Shield Pull-Through Cache - Invariant 1, 7)
	mux.HandleFunc("GET /v1/media/{id}/thumbnail", s.handleGetMediaThumbnail)

	// Origin Shield Stream & HLS Chunk Pull-Through (Signed URLs - Invariant 9)
	mux.HandleFunc("GET /v1/stream/{node_id}/{media_id}/{chunk_path...}", s.handleStreamChunk)
}

// Start begins listening on the configured address.
func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.cfg.Addr)
	if err != nil {
		return fmt.Errorf("failed to bind address %s: %w", s.cfg.Addr, err)
	}
	s.listener = ln
	s.logger.Info("custardfin cloud api started", "addr", ln.Addr().String())

	go func() {
		if err := s.server.Serve(ln); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Error("cloud api server error", "err", err)
		}
	}()

	return nil
}

// Addr returns the bound listener address (useful for dynamic :0 ports in tests).
func (s *Server) Addr() string {
	if s.listener != nil {
		return s.listener.Addr().String()
	}
	return s.cfg.Addr
}

// Shutdown gracefully stops the HTTP server.
func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("cloud api server shutting down")
	return s.server.Shutdown(ctx)
}

// ServeHTTP dispatches requests to the internal handler (useful for testing and embedding).
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.server.Handler.ServeHTTP(w, r)
}

// SetSaturationEstimator configures a custom SaturationEstimator (useful for test mocks).
func (s *Server) SetSaturationEstimator(est ports.SaturationEstimator) {
	s.saturationEstimator = est
}

// SetPrefetcher configures a custom Prefetcher (useful for test mocks).
func (s *Server) SetPrefetcher(pf ports.Prefetcher) {
	s.prefetcher = pf
}

// Handlers

func (s *Server) handleHealthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

type authRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type authResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	UserID      string `json:"user_id"`
	Username    string `json:"username"`
}

func (s *Server) handleRegister(w http.ResponseWriter, r *http.Request) {
	var req authRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON payload")
		return
	}

	hashedPwd, err := s.hasher.HashPassword(req.Password)
	if err != nil {
		if errors.Is(err, domain.ErrValidation) {
			writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to process password")
		return
	}

	user, err := domain.NewUser(req.Username, hashedPwd)
	if err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", err.Error())
		return
	}

	if err := s.userRepo.Create(r.Context(), user); err != nil {
		if errors.Is(err, domain.ErrConflict) {
			writeError(w, http.StatusConflict, "CONFLICT", "Username already taken")
			return
		}
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create user")
		return
	}

	// Issue initial token (1 hour TTL)
	token, err := s.tokenService.GenerateToken(user.ID, user.Username, time.Hour)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to generate token")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(authResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   3600,
		UserID:      user.ID,
		Username:    user.Username,
	})
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	var req authRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON payload")
		return
	}

	user, err := s.userRepo.FindByUsername(r.Context(), req.Username)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid username or password")
		return
	}

	if !s.hasher.ComparePassword(req.Password, user.PasswordHash) {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid username or password")
		return
	}

	token, err := s.tokenService.GenerateToken(user.ID, user.Username, time.Hour)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to generate token")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(authResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   3600,
		UserID:      user.ID,
		Username:    user.Username,
	})
}

func (s *Server) handleGetMe(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Missing authorization header")
		return
	}
	token := strings.TrimPrefix(authHeader, "Bearer ")
	claims, err := s.tokenService.ValidateToken(token)
	if err != nil || claims == nil {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid token")
		return
	}
	user, err := s.userRepo.FindByID(r.Context(), claims.UserID)
	if err != nil || user == nil {
		writeError(w, http.StatusNotFound, "NOT_FOUND", "User not found")
		return
	}

	var nodes []*domain.Node
	if s.nodeRepo != nil {
		nodes, _ = s.nodeRepo.ListByUserID(r.Context(), user.ID)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"id":       user.ID,
		"user_id":  user.ID,
		"username": user.Username,
		"nodes":    nodes,
	})
}

type registerNodeRequest struct {
	Name string `json:"name"`
}

type registerNodeResponse struct {
	NodeID       string `json:"node_id"`
	NodeToken    string `json:"node_token"`
	FrpAuthToken string `json:"frp_auth_token"`
}

func (s *Server) handleRegisterNode(w http.ResponseWriter, r *http.Request) {
	claims, ok := GetUserClaims(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User authorization required")
		return
	}

	var req registerNodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON payload")
		return
	}

	rawToken, tokenHash, frpToken, err := domain.GenerateNodeCredentials()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to generate credentials")
		return
	}

	node, err := domain.NewNode(claims.UserID, req.Name, tokenHash, frpToken)
	if err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", err.Error())
		return
	}

	if err := s.nodeRepo.Create(r.Context(), node); err != nil {
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to save node")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(registerNodeResponse{
		NodeID:       node.ID,
		NodeToken:    rawToken, // Sent only once during registration
		FrpAuthToken: frpToken,
	})
}

func (s *Server) handleUpsertManifest(w http.ResponseWriter, r *http.Request) {
	nodeID := r.PathValue("id")
	mediaID := r.PathValue("mid")

	// Read limited body (max 3 KB + 256 B buffer)
	limitedReader := io.LimitReader(r.Body, int64(contracts.MaxManifestSizeBytes+256))
	bodyBytes, err := io.ReadAll(limitedReader)
	if err != nil {
		writeError(w, http.StatusBadRequest, "BAD_REQUEST", "Failed to read request body")
		return
	}

	if len(bodyBytes) > contracts.MaxManifestSizeBytes {
		writeError(w, http.StatusBadRequest, "MANIFEST_TOO_LARGE",
			fmt.Sprintf("Manifest size %d bytes exceeds strict 3 KB limit", len(bodyBytes)))
		return
	}

	var manifest contracts.Manifest
	if err := json.Unmarshal(bodyBytes, &manifest); err != nil {
		writeError(w, http.StatusBadRequest, "BAD_REQUEST", "Invalid manifest JSON format")
		return
	}

	if manifest.NodeID != nodeID || manifest.MediaID != mediaID {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", "Manifest path IDs do not match manifest body IDs")
		return
	}

	if err := manifest.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", err.Error())
		return
	}

	// Map contracts renditions to domain renditions
	var domainRenditions []domain.Rendition
	for _, rend := range manifest.Renditions {
		domainRenditions = append(domainRenditions, domain.Rendition{
			Resolution:     rend.Resolution,
			BitrateBps:     int(rend.BitrateBps),
			PlaylistPath:   rend.PlaylistPath,
			ChunkDurationS: rend.ChunkDurationS,
		})
	}

	catalogItem, err := domain.NewCatalogItem(
		manifest.MediaID,
		manifest.NodeID,
		manifest.Title,
		manifest.Description,
		manifest.DurationS,
		manifest.ContentHash,
		manifest.IsPublic,
		bodyBytes,
		domainRenditions,
	)
	if err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_FAILED", err.Error())
		return
	}

	if err := s.catalogRepo.UpsertManifest(r.Context(), catalogItem); err != nil {
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to sync manifest")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"status":   "synced",
		"media_id": mediaID,
	})
}

func (s *Server) handleDeleteManifest(w http.ResponseWriter, r *http.Request) {
	nodeID := r.PathValue("id")
	mediaID := r.PathValue("mid")

	if err := s.catalogRepo.DeleteManifest(r.Context(), nodeID, mediaID); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			writeError(w, http.StatusNotFound, "NOT_FOUND", "Manifest not found in catalog")
			return
		}
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to delete manifest")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type catalogItemResponse struct {
	MediaID          string `json:"media_id"`
	Title            string `json:"title"`
	Description      string `json:"description"`
	DurationS        int    `json:"duration_s"`
	NodeID           string `json:"node_id"`
	NodeName         string `json:"node_name,omitempty"`
	UploaderUsername string `json:"uploader_username,omitempty"`
	IsPublic         bool   `json:"is_public"`
	IsNodeOnline     bool   `json:"is_node_online"`
	UpdatedAt        string `json:"updated_at"`
}

type catalogListResponse struct {
	Items []catalogItemResponse `json:"items"`
	Total int                   `json:"total"`
}

func (s *Server) handleGetCatalog(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 20
	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
		limit = l
	}

	offset := 0
	if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
		offset = o
	}

	// Authenticated user resolution for private videos
	var allowedNodeIDs []string
	if claims := s.extractUserClaims(r); claims != nil {
		if s.nodeRepo != nil {
			if nodes, err := s.nodeRepo.ListByUserID(r.Context(), claims.UserID); err == nil {
				for _, n := range nodes {
					allowedNodeIDs = append(allowedNodeIDs, n.ID)
				}
			}
		}
	}

	items, total, err := s.catalogRepo.Search(r.Context(), q, limit, offset, allowedNodeIDs)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to query catalog")
		return
	}

	nodeIDs := make([]string, 0, len(items))
	for _, it := range items {
		nodeIDs = append(nodeIDs, it.NodeID)
	}

	onlineMap := make(map[string]bool)
	if s.heartbeatStore != nil && len(nodeIDs) > 0 {
		if om, err := s.heartbeatStore.BatchIsNodeOnline(r.Context(), nodeIDs); err == nil {
			onlineMap = om
		}
	}

	nodeCache := make(map[string]*domain.Node)
	userCache := make(map[string]*domain.User)

	resItems := make([]catalogItemResponse, 0, len(items))
	for _, it := range items {
		node, ok := nodeCache[it.NodeID]
		if !ok && s.nodeRepo != nil {
			node, _ = s.nodeRepo.FindByID(r.Context(), it.NodeID)
			nodeCache[it.NodeID] = node
		}

		var uploaderName string
		var nodeName string
		if node != nil {
			nodeName = node.Name
			user, ok := userCache[node.UserID]
			if !ok && s.userRepo != nil {
				user, _ = s.userRepo.FindByID(r.Context(), node.UserID)
				userCache[node.UserID] = user
			}
			if user != nil {
				uploaderName = user.Username
			}
		}

		resItems = append(resItems, catalogItemResponse{
			MediaID:          it.MediaID,
			Title:            it.Title,
			Description:      it.Description,
			DurationS:        it.DurationS,
			NodeID:           it.NodeID,
			NodeName:         nodeName,
			UploaderUsername: uploaderName,
			IsPublic:         it.IsPublic,
			IsNodeOnline:     onlineMap[it.NodeID],
			UpdatedAt:        it.UpdatedAt.Format(time.RFC3339),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(catalogListResponse{
		Items: resItems,
		Total: total,
	})
}

func (s *Server) handleNodeHeartbeat(w http.ResponseWriter, r *http.Request) {
	nodeID := r.PathValue("id")
	if nodeID == "" {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "Node ID is required")
		return
	}

	var hb contracts.Heartbeat
	if err := json.NewDecoder(r.Body).Decode(&hb); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "Malformed heartbeat payload")
		return
	}

	if hb.NodeID != "" && hb.NodeID != nodeID {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "Node ID in path does not match payload")
		return
	}
	hb.NodeID = nodeID

	if err := hb.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	if s.heartbeatStore != nil {
		if err := s.heartbeatStore.RecordHeartbeat(r.Context(), &hb, s.cfg.HeartbeatTTL); err != nil {
			s.logger.Error("failed to record heartbeat", "node_id", nodeID, "err", err)
			writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to record heartbeat")
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (s *Server) handleGetNodeStatus(w http.ResponseWriter, r *http.Request) {
	nodeID := r.PathValue("id")
	if nodeID == "" || nodeID == "demo" || nodeID == "default" {
		nodeID = "898ea0c7-f7a3-48e8-ba59-b378a219f9ee"
	}

	isOnline := false
	var hb *contracts.Heartbeat
	if s.heartbeatStore != nil {
		var err error
		var exists bool
		hb, exists, err = s.heartbeatStore.GetHeartbeat(r.Context(), nodeID)
		if err == nil && exists && hb != nil {
			isOnline = true
		}
	}

	nvenc := false
	activeStreams := 0
	version := "0.1.0"
	var lastHeartbeat time.Time
	if isOnline && hb != nil {
		if hb.GPUCaps != nil {
			nvenc = hb.GPUCaps.NVENC
		}
		activeStreams = hb.ActiveStreams
		version = hb.Version
		lastHeartbeat = hb.SentAt
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"node_id":        nodeID,
		"status":         map[bool]string{true: "online", false: "offline"}[isOnline],
		"is_online":      isOnline,
		"active_streams": activeStreams,
		"max_streams":    10,
		"nvenc_enabled":  nvenc,
		"version":        version,
		"last_heartbeat": lastHeartbeat,
	})
}

type playbackResponse struct {
	StreamURL        string `json:"stream_url"`
	RouteType        string `json:"route_type"`
	ExpiresAt        string `json:"expires_at"`
	MediaID          string `json:"media_id,omitempty"`
	Title            string `json:"title,omitempty"`
	UploaderUsername string `json:"uploader_username,omitempty"`
	NodeName         string `json:"node_name,omitempty"`
	IsPublic         bool   `json:"is_public"`
}

func (s *Server) extractUserClaims(r *http.Request) *ports.Claims {
	if s.tokenService == nil {
		return nil
	}
	authHeader := r.Header.Get("Authorization")
	token := ""
	if strings.HasPrefix(authHeader, "Bearer ") {
		token = strings.TrimPrefix(authHeader, "Bearer ")
	} else if qToken := r.URL.Query().Get("token"); qToken != "" {
		token = qToken
	}
	if token == "" {
		return nil
	}
	claims, err := s.tokenService.ValidateToken(token)
	if err != nil {
		return nil
	}
	return claims
}

func isDirectLAN(r *http.Request) bool {
	if r.URL.Query().Get("lan") == "true" {
		return true
	}
	if r.Header.Get("X-Playback-Mode") == "direct_lan" {
		return true
	}
	return false
}

func writePlaybackResponse(w http.ResponseWriter, streamURL, routeType string, expiresAt time.Time, item *domain.CatalogItem, uploaderUsername, nodeName string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	res := playbackResponse{
		StreamURL:        streamURL,
		RouteType:        routeType,
		ExpiresAt:        expiresAt.UTC().Format(time.RFC3339),
		IsPublic:         true,
		UploaderUsername: uploaderUsername,
		NodeName:         nodeName,
	}
	if item != nil {
		res.MediaID = item.MediaID
		res.Title = item.Title
		res.IsPublic = item.IsPublic
	}
	_ = json.NewEncoder(w).Encode(res)
}

func (s *Server) handleGetMediaThumbnail(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "Media ID is required")
		return
	}

	item, err := s.catalogRepo.FindByMediaID(r.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			writeError(w, http.StatusNotFound, "NOT_FOUND", "Media not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to retrieve media")
		return
	}

	// Privacy check: if not public and not direct LAN, only the node owner can access
	if !item.IsPublic && !isDirectLAN(r) {
		claims := s.extractUserClaims(r)
		if claims == nil {
			writeError(w, http.StatusForbidden, "FORBIDDEN", "This media thumbnail is private and requires owner authentication")
			return
		}
		if s.nodeRepo != nil {
			node, err := s.nodeRepo.FindByID(r.Context(), item.NodeID)
			if err != nil || node == nil || node.UserID != claims.UserID {
				writeError(w, http.StatusForbidden, "FORBIDDEN", "You do not have permission to view this thumbnail")
				return
			}
		}
	}

	cacheKey := fmt.Sprintf("%s/%s/thumb.jpg", item.NodeID, item.MediaID)

	// 1. Check Origin Shield S3 Cache (Invariant 7)
	if s.chunkCache != nil {
		cachedItem, err := s.chunkCache.Get(r.Context(), cacheKey)
		if err == nil && cachedItem != nil {
			defer cachedItem.Body.Close()
			w.Header().Set("Content-Type", "image/jpeg")
			w.Header().Set("Cache-Control", "public, max-age=86400")
			w.Header().Set("X-Cache", "HIT")
			_, _ = io.Copy(w, cachedItem.Body)
			return
		}
	}

	// 2. Cache Miss: Pull once from Home Edge over FRP reverse tunnel (Invariant 2, 7)
	if s.tunnelPuller == nil {
		writeError(w, http.StatusServiceUnavailable, "TUNNEL_UNAVAILABLE", "Tunnel puller not configured")
		return
	}

	body, _, _, err := s.tunnelPuller.PullChunk(r.Context(), item.NodeID, item.MediaID, "thumb.jpg")
	if err != nil {
		s.logger.Debug("thumbnail not available on edge node", "media_id", id, "node_id", item.NodeID, "err", err)
		writeError(w, http.StatusNotFound, "THUMBNAIL_NOT_FOUND", "Thumbnail not found on edge")
		return
	}
	defer body.Close()

	thumbData, err := io.ReadAll(body)
	if err != nil || len(thumbData) == 0 {
		writeError(w, http.StatusBadGateway, "EDGE_ERROR", "Failed to read thumbnail payload from edge")
		return
	}

	// 3. Populate S3 Origin Shield Cache (Uploads once from home edge!)
	if s.chunkCache != nil {
		_ = s.chunkCache.Put(r.Context(), cacheKey, bytes.NewReader(thumbData), int64(len(thumbData)), "image/jpeg")
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Cache-Control", "public, max-age=86400")
	w.Header().Set("X-Cache", "MISS")
	w.Header().Set("Content-Length", strconv.Itoa(len(thumbData)))
	_, _ = w.Write(thumbData)
}

func (s *Server) handleGetPlayback(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "Media ID is required")
		return
	}

	item, err := s.catalogRepo.FindByMediaID(r.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			writeError(w, http.StatusNotFound, "NOT_FOUND", "Media not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to resolve playback")
		return
	}

	// Resolve node owner and names
	var uploaderUsername string
	var nodeName string
	var nodeOwnerUserID string
	if s.nodeRepo != nil {
		if node, err := s.nodeRepo.FindByID(r.Context(), item.NodeID); err == nil && node != nil {
			nodeName = node.Name
			nodeOwnerUserID = node.UserID
			if s.userRepo != nil {
				if user, err := s.userRepo.FindByID(r.Context(), node.UserID); err == nil && user != nil {
					uploaderUsername = user.Username
				}
			}
		}
	}

	// Privacy check: if not public and not direct LAN, only the node owner can access
	if !item.IsPublic && !isDirectLAN(r) {
		claims := s.extractUserClaims(r)
		if claims == nil {
			writeError(w, http.StatusForbidden, "FORBIDDEN", "This media is private and requires owner authentication")
			return
		}
		if nodeOwnerUserID != "" && claims.UserID != nodeOwnerUserID {
			writeError(w, http.StatusForbidden, "FORBIDDEN", "You do not have permission to view this private media")
			return
		}
	}

	// 1. Direct LAN Detection:
	if isDirectLAN(r) {
		streamURL := fmt.Sprintf("http://localhost:8090/stream/%s/master.m3u8", item.MediaID)
		writePlaybackResponse(w, streamURL, "direct_lan", time.Now().Add(1*time.Hour), item, uploaderUsername, nodeName)
		return
	}

	// 2. Play Decision Tree: Remote viewer
	// Check node liveness and active stream capacity from Redis HeartbeatStore
	var nodeOnline bool = true
	if s.heartbeatStore != nil {
		hb, online, err := s.heartbeatStore.GetHeartbeat(r.Context(), item.NodeID)
		if err != nil {
			s.logger.Error("failed checking node liveness", "node_id", item.NodeID, "err", err)
			writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to check node status")
			return
		}
		nodeOnline = online
		if online && hb != nil && hb.ActiveStreams >= contracts.MaxConcurrentStreamsPerNode {
			w.Header().Set("Retry-After", "30")
			writeError(w, http.StatusTooManyRequests, "NODE_BUSY", "Node capacity limit exceeded (10 active streams)")
			return
		}
	}

	// 3. Offline handling: Check if master playlist exists in cache (FR-7)
	masterKey := fmt.Sprintf("%s/%s/master.m3u8", item.NodeID, item.MediaID)
	cached := false
	if s.chunkCache != nil {
		cached, _ = s.chunkCache.Exists(r.Context(), masterKey)
	}

	if !nodeOnline {
		if !cached {
			writeError(w, http.StatusServiceUnavailable, "NODE_OFFLINE", "Home edge node is currently offline and media is not cached")
			return
		}
		baseURL := s.resolveBaseURL(r)
		expiresAt := time.Now().Add(1 * time.Hour)
		sigQuery := contracts.SignChunkURL(s.cfg.SigningSecret, item.NodeID, item.MediaID, "master.m3u8", expiresAt)
		streamURL := fmt.Sprintf("%s/v1/stream/%s/%s/master.m3u8?%s", baseURL, item.NodeID, item.MediaID, sigQuery)
		writePlaybackResponse(w, streamURL, "cache_only", expiresAt, item, uploaderUsername, nodeName)
		return
	}

	// 4. Remote viewer with active node -> Origin Shield signed stream URL
	// Warm initial video chunk ahead of time to eliminate TTFB buffering (Milestone 8, FR-10)
	if s.prefetcher != nil {
		s.prefetcher.TriggerPrefetch(r.Context(), item.NodeID, item.MediaID, "1080p/segment_000.ts")
	}

	baseURL := s.resolveBaseURL(r)
	expiresAt := time.Now().Add(1 * time.Hour)
	sigQuery := contracts.SignChunkURL(s.cfg.SigningSecret, item.NodeID, item.MediaID, "master.m3u8", expiresAt)
	streamURL := fmt.Sprintf("%s/v1/stream/%s/%s/master.m3u8?%s", baseURL, item.NodeID, item.MediaID, sigQuery)
	writePlaybackResponse(w, streamURL, "origin_shield", expiresAt, item, uploaderUsername, nodeName)
}

func (s *Server) resolveBaseURL(r *http.Request) string {
	if r != nil && r.Host != "" {
		proto := "http"
		if r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https" {
			proto = "https"
		}
		return fmt.Sprintf("%s://%s", proto, r.Host)
	}
	if s.cfg.BaseURL != "" {
		return strings.TrimRight(s.cfg.BaseURL, "/")
	}
	return "http://localhost:8080"
}

func (s *Server) handleStreamChunk(w http.ResponseWriter, r *http.Request) {
	nodeID := r.PathValue("node_id")
	mediaID := r.PathValue("media_id")
	chunkPath := r.PathValue("chunk_path")

	if strings.TrimSpace(nodeID) == "" || strings.TrimSpace(mediaID) == "" || strings.TrimSpace(chunkPath) == "" {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "Node ID, media ID, and chunk path are required")
		return
	}

	// Prevent directory traversal
	cleanPath := path.Clean(strings.TrimPrefix(chunkPath, "/"))
	if strings.HasPrefix(cleanPath, "..") || cleanPath == "." {
		writeError(w, http.StatusBadRequest, "FORBIDDEN_PATH", "Invalid chunk path")
		return
	}

	sig := r.URL.Query().Get("sig")
	exp := r.URL.Query().Get("exp")
	if sig == "" || exp == "" {
		writeError(w, http.StatusUnauthorized, "MISSING_SIGNATURE", "URL signature and expiration are required")
		return
	}

	// Verify HMAC-SHA256 signature (Invariant 9)
	if err := contracts.VerifyChunkSignature(s.cfg.SigningSecret, nodeID, mediaID, cleanPath, sig, exp, time.Now()); err != nil {
		if errors.Is(err, contracts.ErrSignatureExpired) {
			writeError(w, http.StatusUnauthorized, "SIGNATURE_EXPIRED", "URL signature has expired")
			return
		}
		writeError(w, http.StatusUnauthorized, "INVALID_SIGNATURE", "Invalid chunk signature")
		return
	}

	// Evaluate edge node saturation level for adaptive rendition throttling (FR-10)
	loadLevel := domain.LoadNormal
	if s.saturationEstimator != nil {
		if lvl, _, err := s.saturationEstimator.EstimateLoad(r.Context(), nodeID); err == nil {
			loadLevel = lvl
		}
	}

	cacheKey := fmt.Sprintf("%s/%s/%s", nodeID, mediaID, cleanPath)

	// Check Origin Shield Cache (Invariant 7)
	if s.chunkCache != nil {
		cachedItem, err := s.chunkCache.Get(r.Context(), cacheKey)
		if err == nil && cachedItem != nil {
			defer cachedItem.Body.Close()
			w.Header().Set("X-Cache", "HIT")
			s.serveChunkContent(w, r, nodeID, mediaID, cleanPath, cachedItem.Body, cachedItem.ContentType, loadLevel)
			return
		}
	}

	// Cache Miss: Check node liveness and active streams limit (FR-7, Invariant 10)
	if s.heartbeatStore != nil {
		hb, online, err := s.heartbeatStore.GetHeartbeat(r.Context(), nodeID)
		if err != nil {
			s.logger.Error("failed checking node liveness", "node_id", nodeID, "err", err)
			writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to check node status")
			return
		}
		if !online {
			writeError(w, http.StatusServiceUnavailable, "NODE_OFFLINE", "Home edge node is currently offline and requested chunk is not cached")
			return
		}
		if hb != nil && hb.ActiveStreams >= contracts.MaxConcurrentStreamsPerNode {
			w.Header().Set("Retry-After", "30")
			writeError(w, http.StatusTooManyRequests, "NODE_BUSY", "Node capacity limit exceeded (10 active streams)")
			return
		}
	}

	// Pull chunk through FRP tunnel relay (Invariant 2, Invariant 7)
	if s.tunnelPuller == nil {
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Tunnel puller not configured")
		return
	}

	body, contentType, length, err := s.tunnelPuller.PullChunk(r.Context(), nodeID, mediaID, cleanPath)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			writeError(w, http.StatusNotFound, "NOT_FOUND", "Chunk not found on edge node")
			return
		}
		s.logger.Error("failed pulling chunk from edge via tunnel", "node_id", nodeID, "media_id", mediaID, "path", cleanPath, "err", err)
		writeError(w, http.StatusBadGateway, "TUNNEL_ERROR", "Failed to retrieve chunk from home edge node")
		return
	}
	defer body.Close()

	// Read content so we can both store in cache and serve
	data, err := io.ReadAll(body)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed reading chunk body")
		return
	}

	// Save to Origin Shield cache (Invariant 1, Invariant 7)
	if s.chunkCache != nil {
		_ = s.chunkCache.Put(r.Context(), cacheKey, bytes.NewReader(data), int64(len(data)), contentType)
	}

	w.Header().Set("X-Cache", "MISS")
	s.serveChunkContent(w, r, nodeID, mediaID, cleanPath, bytes.NewReader(data), contentType, loadLevel)
	_ = length
}

func (s *Server) serveChunkContent(w http.ResponseWriter, r *http.Request, nodeID, mediaID, chunkPath string, reader io.Reader, contentType string, load domain.LoadLevel) {
	ext := strings.ToLower(path.Ext(chunkPath))
	switch ext {
	case ".m3u8":
		w.Header().Set("Content-Type", "application/vnd.apple.mpegurl")
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		// Playlist rewriting: dynamically sign relative child URLs and filter by load level
		raw, err := io.ReadAll(reader)
		if err != nil {
			http.Error(w, "failed reading playlist", http.StatusInternalServerError)
			return
		}
		playlistDir := path.Dir(chunkPath)
		if playlistDir == "." {
			playlistDir = ""
		}
		expiresAt := time.Now().Add(1 * time.Hour)
		rewritten := rewritePlaylistWithLoad(raw, nodeID, mediaID, playlistDir, s.cfg.SigningSecret, expiresAt, load)
		w.Header().Set("Content-Length", strconv.Itoa(len(rewritten)))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(rewritten)
		return
	case ".ts":
		w.Header().Set("Content-Type", "video/mp2t")
		w.Header().Set("Cache-Control", "public, max-age=86400, immutable")
	default:
		w.Header().Set("Content-Type", contentType)
	}

	w.Header().Set("Accept-Ranges", "bytes")
	w.WriteHeader(http.StatusOK)
	_, _ = io.Copy(w, reader)
}

func rewritePlaylist(content []byte, nodeID, mediaID, playlistDir string, secret []byte, expiresAt time.Time) []byte {
	return rewritePlaylistWithLoad(content, nodeID, mediaID, playlistDir, secret, expiresAt, domain.LoadNormal)
}

type variantStream struct {
	infoLine   string
	uriLine    string
	resolution string
}

func extractResolution(infoLine string) string {
	idx := strings.Index(infoLine, "RESOLUTION=")
	if idx == -1 {
		return ""
	}
	sub := infoLine[idx+len("RESOLUTION="):]
	comma := strings.Index(sub, ",")
	if comma != -1 {
		sub = sub[:comma]
	}
	return strings.TrimSpace(sub)
}

func rewritePlaylistWithLoad(content []byte, nodeID, mediaID, playlistDir string, secret []byte, expiresAt time.Time, load domain.LoadLevel) []byte {
	// If it's not a master playlist, apply regular segment-level URL signing
	if !bytes.Contains(content, []byte("#EXT-X-STREAM-INF")) {
		scanner := bufio.NewScanner(bytes.NewReader(content))
		var out bytes.Buffer

		for scanner.Scan() {
			line := scanner.Text()
			trimmed := strings.TrimSpace(line)

			// Pass through empty lines and comments/tags (#EXT...)
			if trimmed == "" || strings.HasPrefix(trimmed, "#") {
				out.WriteString(line)
				out.WriteString("\n")
				continue
			}

			// Line is a URI (relative path to sub-playlist or segment)
			chunkPath := trimmed
			if playlistDir != "" && playlistDir != "." {
				chunkPath = path.Join(playlistDir, trimmed)
			}

			sigQuery := contracts.SignChunkURL(secret, nodeID, mediaID, chunkPath, expiresAt)
			sep := "?"
			if strings.Contains(trimmed, "?") {
				sep = "&"
			}
			out.WriteString(trimmed + sep + sigQuery)
			out.WriteString("\n")
		}

		return out.Bytes()
	}

	// Master playlist parsing with adaptive rendition filtering
	scanner := bufio.NewScanner(bytes.NewReader(content))
	var headerLines []string
	var variants []variantStream
	var pendingInfo string

	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}

		if strings.HasPrefix(trimmed, "#EXT-X-STREAM-INF") {
			pendingInfo = line
			continue
		}

		if pendingInfo != "" {
			res := extractResolution(pendingInfo)
			variants = append(variants, variantStream{
				infoLine:   pendingInfo,
				uriLine:    line,
				resolution: res,
			})
			pendingInfo = ""
			continue
		}

		if strings.HasPrefix(trimmed, "#") {
			headerLines = append(headerLines, line)
		}
	}

	// Filter variants according to load level
	var acceptedVariants []variantStream
	for _, v := range variants {
		if load.AllowedResolution(v.resolution) {
			acceptedVariants = append(acceptedVariants, v)
		}
	}

	// Safety fallback: if all variants were pruned, keep the lowest rendition to maintain playback
	if len(acceptedVariants) == 0 && len(variants) > 0 {
		acceptedVariants = append(acceptedVariants, variants[len(variants)-1])
	}

	var out bytes.Buffer
	for _, h := range headerLines {
		out.WriteString(h)
		out.WriteString("\n")
	}

	for _, v := range acceptedVariants {
		out.WriteString(v.infoLine)
		out.WriteString("\n")

		chunkPath := strings.TrimSpace(v.uriLine)
		if playlistDir != "" && playlistDir != "." {
			chunkPath = path.Join(playlistDir, chunkPath)
		}

		sigQuery := contracts.SignChunkURL(secret, nodeID, mediaID, chunkPath, expiresAt)
		sep := "?"
		if strings.Contains(v.uriLine, "?") {
			sep = "&"
		}
		out.WriteString(strings.TrimSpace(v.uriLine) + sep + sigQuery)
		out.WriteString("\n")
	}

	return out.Bytes()
}
