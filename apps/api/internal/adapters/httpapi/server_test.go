package httpapi_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/anosh26/custardfin/apps/api/internal/adapters/auth"
	"github.com/anosh26/custardfin/apps/api/internal/adapters/httpapi"
	"github.com/anosh26/custardfin/apps/api/internal/adapters/memory"
	"github.com/anosh26/custardfin/apps/api/internal/domain"
	"github.com/anosh26/custardfin/apps/api/internal/ports"
	"github.com/anosh26/custardfin/libs/contracts"
	"golang.org/x/crypto/bcrypt"
)

var testSigningSecret = []byte("test-signing-secret-at-least-32b!")

func setupTestServer(t *testing.T) (*httpapi.Server, *memory.InMemoryUserRepository, *memory.InMemoryNodeRepository, *memory.InMemoryCatalogRepository, *memory.InMemoryHeartbeatStore) {
	s, u, n, c, h, _ := setupTestServerFull(t, nil, nil)
	return s, u, n, c, h
}

func setupTestServerFull(
	t *testing.T,
	cache ports.ChunkCache,
	puller ports.TunnelPuller,
) (*httpapi.Server, *memory.InMemoryUserRepository, *memory.InMemoryNodeRepository, *memory.InMemoryCatalogRepository, *memory.InMemoryHeartbeatStore, *memory.ChunkCache) {
	t.Helper()

	userRepo := memory.NewInMemoryUserRepository()
	nodeRepo := memory.NewInMemoryNodeRepository()
	catalogRepo := memory.NewInMemoryCatalogRepository()
	heartbeatStore := memory.NewInMemoryHeartbeatStore()

	var memCache *memory.ChunkCache
	if cache == nil {
		memCache = memory.NewMemoryChunkCache(24 * time.Hour)
		cache = memCache
	} else if mc, ok := cache.(*memory.ChunkCache); ok {
		memCache = mc
	}

	hasher := auth.NewBcryptHasher(bcrypt.MinCost)
	tokenService, err := auth.NewJWTService([]byte("test-secret-key-at-least-32-bytes-long!"))
	if err != nil {
		t.Fatalf("failed to create jwt service: %v", err)
	}

	server := httpapi.NewServer(
		httpapi.Config{
			Addr:          "127.0.0.1:0",
			HeartbeatTTL:  20 * time.Second,
			BaseURL:       "http://localhost:8080",
			SigningSecret: testSigningSecret,
		},
		userRepo,
		nodeRepo,
		catalogRepo,
		heartbeatStore,
		cache,
		puller,
		hasher,
		tokenService,
		nil,
	)

	return server, userRepo, nodeRepo, catalogRepo, heartbeatStore, memCache
}

func TestServer_Healthz(t *testing.T) {
	server, _, _, _, _ := setupTestServer(t)

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var res map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("failed to parse json: %v", err)
	}
	if res["status"] != "ok" {
		t.Errorf("expected status 'ok', got %s", res["status"])
	}
}

func TestServer_CORS_Middleware(t *testing.T) {
	server, _, _, _, _ := setupTestServer(t)

	req := httptest.NewRequest(http.MethodOptions, "/v1/auth/login", nil)
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected status 204 for OPTIONS, got %d", w.Code)
	}

	if origin := w.Header().Get("Access-Control-Allow-Origin"); origin != "*" {
		t.Errorf("expected Access-Control-Allow-Origin: *, got %s", origin)
	}
	if methods := w.Header().Get("Access-Control-Allow-Methods"); !strings.Contains(methods, "POST") {
		t.Errorf("expected Access-Control-Allow-Methods to contain POST, got %s", methods)
	}
}

func TestServer_Auth_RegisterAndLogin(t *testing.T) {
	server, _, _, _, _ := setupTestServer(t)

	// 1. Weak password registration (< 8 chars)
	weakPayload := `{"username":"alice","password":"123"}`
	req := httptest.NewRequest(http.MethodPost, "/v1/auth/register", bytes.NewBufferString(weakPayload))
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for weak password, got %d", w.Code)
	}

	// 2. Successful Registration
	validPayload := `{"username":"alice","password":"password123"}`
	req = httptest.NewRequest(http.MethodPost, "/v1/auth/register", bytes.NewBufferString(validPayload))
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created, got %d: %s", w.Code, w.Body.String())
	}

	var authRes struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &authRes); err != nil {
		t.Fatalf("failed to decode auth response: %v", err)
	}
	if authRes.AccessToken == "" || authRes.TokenType != "Bearer" {
		t.Errorf("invalid auth response: %+v", authRes)
	}

	// 3. Duplicate Registration
	req = httptest.NewRequest(http.MethodPost, "/v1/auth/register", bytes.NewBufferString(validPayload))
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	if w.Code != http.StatusConflict {
		t.Errorf("expected 409 Conflict for duplicate username, got %d", w.Code)
	}

	// 4. Successful Login
	req = httptest.NewRequest(http.MethodPost, "/v1/auth/login", bytes.NewBufferString(validPayload))
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for valid login, got %d", w.Code)
	}

	// 5. Wrong Password Login
	wrongPwdPayload := `{"username":"alice","password":"wrong_password"}`
	req = httptest.NewRequest(http.MethodPost, "/v1/auth/login", bytes.NewBufferString(wrongPwdPayload))
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 Unauthorized for wrong password, got %d", w.Code)
	}
}

func TestServer_Nodes_PairingAndManifestSync(t *testing.T) {
	server, _, _, _, _ := setupTestServer(t)

	// Step 1: Register User & Obtain JWT
	regPayload := `{"username":"bob_node_owner","password":"secure_password_123"}`
	req := httptest.NewRequest(http.MethodPost, "/v1/auth/register", bytes.NewBufferString(regPayload))
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)

	var authRes struct {
		AccessToken string `json:"access_token"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &authRes)
	jwtToken := authRes.AccessToken

	// Step 2: Register Edge Node without JWT -> 401 Unauthorized
	nodePayload := `{"name":"Bob Home RTX"}`
	req = httptest.NewRequest(http.MethodPost, "/v1/nodes", bytes.NewBufferString(nodePayload))
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 without JWT, got %d", w.Code)
	}

	// Step 3: Register Edge Node with JWT -> 201 Created
	req = httptest.NewRequest(http.MethodPost, "/v1/nodes", bytes.NewBufferString(nodePayload))
	req.Header.Set("Authorization", "Bearer "+jwtToken)
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created for node registration, got %d: %s", w.Code, w.Body.String())
	}

	var nodeRes struct {
		NodeID       string `json:"node_id"`
		NodeToken    string `json:"node_token"`
		FrpAuthToken string `json:"frp_auth_token"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &nodeRes); err != nil {
		t.Fatalf("failed to decode node registration response: %v", err)
	}
	if nodeRes.NodeID == "" || nodeRes.NodeToken == "" || nodeRes.FrpAuthToken == "" {
		t.Fatalf("expected non-empty node credentials: %+v", nodeRes)
	}

	// Step 4: Edge Node Syncs Manifest using X-Node-Token
	mediaID := "00000000-0000-0000-0000-000000000100"
	validManifest := contracts.Manifest{
		SchemaVersion: "1.0",
		MediaID:       mediaID,
		NodeID:        nodeRes.NodeID,
		Title:         "Cosmos Laundromat",
		Description:   "First cycle of open movie project gooseberry",
		DurationS:     720,
		ContentHash:   "a1b2c3d4e5f60718293a4b5c6d7e8f90a1b2c3d4e5f60718293a4b5c6d7e8f90",
		Renditions: []contracts.Rendition{
			{Resolution: "1080p", BitrateBps: 5000000, PlaylistPath: "1080p.m3u8", ChunkDurationS: 4},
			{Resolution: "720p", BitrateBps: 2800000, PlaylistPath: "720p.m3u8", ChunkDurationS: 4},
		},
		IsPublic:  true,
		UpdatedAt: time.Now().UTC(),
	}

	manifestBytes, _ := json.Marshal(validManifest)
	syncURL := fmt.Sprintf("/v1/nodes/%s/media/%s/manifest", nodeRes.NodeID, mediaID)

	// Sync without X-Node-Token -> 401 Unauthorized
	req = httptest.NewRequest(http.MethodPut, syncURL, bytes.NewBuffer(manifestBytes))
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 without X-Node-Token, got %d", w.Code)
	}

	// Sync with valid X-Node-Token -> 200 OK
	req = httptest.NewRequest(http.MethodPut, syncURL, bytes.NewBuffer(manifestBytes))
	req.Header.Set("X-Node-Token", nodeRes.NodeToken)
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for manifest sync, got %d: %s", w.Code, w.Body.String())
	}

	// Step 5: Test Invariant 1 - Manifest > 3 KB is strictly rejected
	oversizedBytes := make([]byte, contracts.MaxManifestSizeBytes+10)
	req = httptest.NewRequest(http.MethodPut, syncURL, bytes.NewBuffer(oversizedBytes))
	req.Header.Set("X-Node-Token", nodeRes.NodeToken)
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 Bad Request for manifest exceeding 3 KB cap, got %d", w.Code)
	}

	// Step 6: Public Catalog Browse & Search
	req = httptest.NewRequest(http.MethodGet, "/v1/catalog?q=cosmos", nil)
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for catalog search, got %d", w.Code)
	}

	var catalogRes struct {
		Items []struct {
			MediaID string `json:"media_id"`
			Title   string `json:"title"`
		} `json:"items"`
		Total int `json:"total"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &catalogRes); err != nil {
		t.Fatalf("failed to decode catalog response: %v", err)
	}
	if catalogRes.Total != 1 || len(catalogRes.Items) != 1 {
		t.Errorf("expected 1 catalog item, got total=%d len=%d", catalogRes.Total, len(catalogRes.Items))
	}
	if catalogRes.Items[0].Title != "Cosmos Laundromat" {
		t.Errorf("expected Cosmos Laundromat, got %s", catalogRes.Items[0].Title)
	}

	// Step 7: Delete Manifest (Unpublish)
	req = httptest.NewRequest(http.MethodDelete, syncURL, nil)
	req.Header.Set("X-Node-Token", nodeRes.NodeToken)
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204 No Content for delete manifest, got %d", w.Code)
	}

	// Confirm catalog is now empty
	req = httptest.NewRequest(http.MethodGet, "/v1/catalog", nil)
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	_ = json.Unmarshal(w.Body.Bytes(), &catalogRes)
	if catalogRes.Total != 0 {
		t.Errorf("expected 0 items in catalog after deletion, got %d", catalogRes.Total)
	}
}

func TestServer_StartAndShutdown(t *testing.T) {
	server, _, _, _, _ := setupTestServer(t)

	if err := server.Start(); err != nil {
		t.Fatalf("failed to start server: %v", err)
	}

	addr := server.Addr()
	if addr == "" {
		t.Fatalf("expected non-empty bound address")
	}

	// Ping healthz over real TCP
	resp, err := http.Get("http://" + addr + "/healthz")
	if err != nil {
		t.Fatalf("failed to ping server: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200 OK from running server, got %d", resp.StatusCode)
	}

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		t.Errorf("expected clean shutdown, got %v", err)
	}
}

func TestServer_Playback(t *testing.T) {
	server, _, _, catalogRepo, heartbeatStore := setupTestServer(t)

	// 1. Not Found
	req := httptest.NewRequest(http.MethodGet, "/v1/media/non-existent-id/playback", nil)
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404 for missing media, got %d", w.Code)
	}

	// 2. Existing Media - Node is Offline (no heartbeat recorded)
	mediaID := "77777777-7777-7777-7777-777777777777"
	nodeID := "11111111-1111-1111-1111-111111111111"

	err := catalogRepo.UpsertManifest(context.Background(), &domain.CatalogItem{
		MediaID:     mediaID,
		NodeID:      nodeID,
		Title:       "Test Movie",
		Description: "A great movie",
		DurationS:   120,
		ContentHash: "hash-12345",
		IsPublic:    true,
		UpdatedAt:   time.Now().UTC(),
	})
	if err != nil {
		t.Fatalf("failed to insert test catalog item: %v", err)
	}

	req = httptest.NewRequest(http.MethodGet, "/v1/media/"+mediaID+"/playback", nil)
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)

	if w.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503 Service Unavailable when node is offline, got %d: %s", w.Code, w.Body.String())
	}

	// 3. Existing Media - Node is Online with healthy capacity
	hb := contracts.Heartbeat{
		NodeID:        nodeID,
		SentAt:        time.Now().UTC(),
		Version:       "1.0.0",
		ActiveStreams: 2,
	}
	if err := heartbeatStore.RecordHeartbeat(context.Background(), &hb, 30*time.Second); err != nil {
		t.Fatalf("failed to record test heartbeat: %v", err)
	}

	// 3a. Remote request defaults to Origin Shield signed stream URL
	req = httptest.NewRequest(http.MethodGet, "/v1/media/"+mediaID+"/playback", nil)
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK when node is online, got %d: %s", w.Code, w.Body.String())
	}

	var pbRes struct {
		StreamURL string `json:"stream_url"`
		RouteType string `json:"route_type"`
		ExpiresAt string `json:"expires_at"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &pbRes); err != nil {
		t.Fatalf("failed to parse playback response: %v", err)
	}

	if !strings.Contains(pbRes.StreamURL, mediaID) || !strings.Contains(pbRes.StreamURL, "sig=") {
		t.Errorf("expected stream url to contain media ID and signature, got %s", pbRes.StreamURL)
	}
	if pbRes.RouteType != "origin_shield" {
		t.Errorf("expected origin_shield route type, got %s", pbRes.RouteType)
	}
	if pbRes.ExpiresAt == "" {
		t.Errorf("expected non-empty expires_at")
	}

	// 3b. Direct LAN request (?lan=true) routes to direct_lan
	reqLAN := httptest.NewRequest(http.MethodGet, "/v1/media/"+mediaID+"/playback?lan=true", nil)
	wLAN := httptest.NewRecorder()
	server.ServeHTTP(wLAN, reqLAN)

	if wLAN.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on direct lan request, got %d", wLAN.Code)
	}
	var pbResLAN struct {
		StreamURL string `json:"stream_url"`
		RouteType string `json:"route_type"`
	}
	_ = json.Unmarshal(wLAN.Body.Bytes(), &pbResLAN)
	if pbResLAN.RouteType != "direct_lan" {
		t.Errorf("expected direct_lan route type, got %s", pbResLAN.RouteType)
	}

	// 4. Invariant 10: Node Busy (ActiveStreams >= 10)
	busyHb := contracts.Heartbeat{
		NodeID:        nodeID,
		SentAt:        time.Now().UTC(),
		Version:       "1.0.0",
		ActiveStreams: 10,
	}
	if err := heartbeatStore.RecordHeartbeat(context.Background(), &busyHb, 30*time.Second); err != nil {
		t.Fatalf("failed to record busy heartbeat: %v", err)
	}

	req = httptest.NewRequest(http.MethodGet, "/v1/media/"+mediaID+"/playback", nil)
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)

	if w.Code != http.StatusTooManyRequests {
		t.Fatalf("expected 429 Too Many Requests when node reaches 10-stream cap, got %d: %s", w.Code, w.Body.String())
	}
	if retryAfter := w.Header().Get("Retry-After"); retryAfter != "30" {
		t.Errorf("expected Retry-After: 30 header, got %s", retryAfter)
	}

	// 5. Simulated HeartbeatStore Error
	heartbeatStore.SimulateErr = fmt.Errorf("redis connection dropped")
	req = httptest.NewRequest(http.MethodGet, "/v1/media/"+mediaID+"/playback", nil)
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500 when heartbeat store fails, got %d", w.Code)
	}
	heartbeatStore.SimulateErr = nil

	// 6. Simulated Repo Error
	catalogRepo.SimulateErr = fmt.Errorf("database query failure")
	req = httptest.NewRequest(http.MethodGet, "/v1/media/"+mediaID+"/playback", nil)
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500 Internal Server Error, got %d", w.Code)
	}
}

type mockTunnelPuller struct {
	pulledPaths []string
	content     map[string][]byte
	contentType map[string]string
	err         error
}

func (m *mockTunnelPuller) PullChunk(ctx context.Context, nodeID, mediaID, chunkPath string) (io.ReadCloser, string, int64, error) {
	if m.err != nil {
		return nil, "", 0, m.err
	}
	m.pulledPaths = append(m.pulledPaths, chunkPath)
	data, ok := m.content[chunkPath]
	if !ok {
		return nil, "", 0, domain.ErrNotFound
	}
	ct := m.contentType[chunkPath]
	if ct == "" {
		ct = "application/octet-stream"
	}
	return io.NopCloser(bytes.NewReader(data)), ct, int64(len(data)), nil
}

func TestServer_OriginShield_Streaming(t *testing.T) {
	nodeID := "11111111-1111-1111-1111-111111111111"
	mediaID := "22222222-2222-2222-2222-222222222222"

	rawMaster := "#EXTM3U\n#EXT-X-VERSION:3\n#EXT-X-STREAM-INF:BANDWIDTH=5000000,RESOLUTION=1920x1080\n1080p/index.m3u8\n"
	rawIndex := "#EXTM3U\n#EXT-X-VERSION:3\n#EXT-X-TARGETDURATION:4\n#EXTINF:4.000000,\nsegment_0000.ts\n#EXT-X-ENDLIST\n"
	rawSegment := []byte("binary-ts-segment-payload")

	puller := &mockTunnelPuller{
		content: map[string][]byte{
			"master.m3u8":           []byte(rawMaster),
			"1080p/index.m3u8":      []byte(rawIndex),
			"1080p/segment_0000.ts": rawSegment,
		},
		contentType: map[string]string{
			"master.m3u8":           "application/vnd.apple.mpegurl",
			"1080p/index.m3u8":      "application/vnd.apple.mpegurl",
			"1080p/segment_0000.ts": "video/mp2t",
		},
	}

	memCache := memory.NewMemoryChunkCache(24 * time.Hour)
	server, _, _, catalogRepo, heartbeatStore, _ := setupTestServerFull(t, memCache, puller)

	// Register catalog item
	_ = catalogRepo.UpsertManifest(context.Background(), &domain.CatalogItem{
		MediaID:     mediaID,
		NodeID:      nodeID,
		Title:       "Sample Film",
		DurationS:   100,
		ContentHash: "hash-film",
		IsPublic:    true,
		UpdatedAt:   time.Now().UTC(),
	})

	// Register online heartbeat
	_ = heartbeatStore.RecordHeartbeat(context.Background(), &contracts.Heartbeat{
		NodeID:        nodeID,
		SentAt:        time.Now().UTC(),
		Version:       "1.0.0",
		ActiveStreams: 1,
	}, 30*time.Second)

	now := time.Now()
	expFuture := now.Add(1 * time.Hour)
	expPast := now.Add(-10 * time.Minute)

	// 1. Missing Signature -> 401
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/stream/%s/%s/master.m3u8", nodeID, mediaID), nil)
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 for missing signature, got %d", w.Code)
	}

	// 2. Invalid / Tampered Signature -> 401
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/stream/%s/%s/master.m3u8?sig=tampered&exp=%d", nodeID, mediaID, expFuture.Unix()), nil)
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 for tampered signature, got %d", w.Code)
	}

	// 3. Expired Signature -> 401
	expiredSig := contracts.SignChunkURL(testSigningSecret, nodeID, mediaID, "master.m3u8", expPast)
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/stream/%s/%s/master.m3u8?%s", nodeID, mediaID, expiredSig), nil)
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 for expired signature, got %d", w.Code)
	}

	// 4. Directory Traversal Attempt -> 400
	cleanMalicious := path.Clean("../passwd")
	travSig := contracts.SignChunkURL(testSigningSecret, nodeID, mediaID, cleanMalicious, expFuture)
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/stream/%s/%s/%%2e%%2e/passwd?%s", nodeID, mediaID, travSig), nil)
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest && w.Code != http.StatusUnauthorized {
		t.Errorf("expected 400 or 401 for forbidden directory traversal path, got %d", w.Code)
	}

	// 5. Valid request on Cache MISS -> Pulls through tunnel, caches, returns X-Cache: MISS
	validSigMaster := contracts.SignChunkURL(testSigningSecret, nodeID, mediaID, "master.m3u8", expFuture)
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/stream/%s/%s/master.m3u8?%s", nodeID, mediaID, validSigMaster), nil)
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on valid stream chunk request, got %d: %s", w.Code, w.Body.String())
	}
	if cacheHdr := w.Header().Get("X-Cache"); cacheHdr != "MISS" {
		t.Errorf("expected X-Cache: MISS on initial fetch, got %s", cacheHdr)
	}

	// Verify master.m3u8 URLs were rewritten with HMAC signatures
	masterBody := w.Body.String()
	if !strings.Contains(masterBody, "1080p/index.m3u8?sig=") {
		t.Errorf("expected rewritten playlist with signed sub-playlist URL, got:\n%s", masterBody)
	}

	// 6. Second request on Cache HIT -> Serves from cache without pulling tunnel again, X-Cache: HIT
	pullCountBefore := len(puller.pulledPaths)
	w2 := httptest.NewRecorder()
	server.ServeHTTP(w2, req)

	if w2.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on cached stream request, got %d", w2.Code)
	}
	if cacheHdr := w2.Header().Get("X-Cache"); cacheHdr != "HIT" {
		t.Errorf("expected X-Cache: HIT on second fetch, got %s", cacheHdr)
	}
	if len(puller.pulledPaths) != pullCountBefore {
		t.Errorf("expected no additional tunnel pull on cache hit")
	}

	// 7. Requesting rendition index playlist (1080p/index.m3u8)
	validSigIndex := contracts.SignChunkURL(testSigningSecret, nodeID, mediaID, "1080p/index.m3u8", expFuture)
	reqIndex := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/stream/%s/%s/1080p/index.m3u8?%s", nodeID, mediaID, validSigIndex), nil)
	wIndex := httptest.NewRecorder()
	server.ServeHTTP(wIndex, reqIndex)

	if wIndex.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for rendition index, got %d: %s", wIndex.Code, wIndex.Body.String())
	}
	indexBody := wIndex.Body.String()
	if !strings.Contains(indexBody, "segment_0000.ts?sig=") {
		t.Errorf("expected rewritten segment URL with signature, got:\n%s", indexBody)
	}

	// 8. Requesting media segment .ts (1080p/segment_0000.ts)
	validSigSeg := contracts.SignChunkURL(testSigningSecret, nodeID, mediaID, "1080p/segment_0000.ts", expFuture)
	reqSeg := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/stream/%s/%s/1080p/segment_0000.ts?%s", nodeID, mediaID, validSigSeg), nil)
	wSeg := httptest.NewRecorder()
	server.ServeHTTP(wSeg, reqSeg)

	if wSeg.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for .ts segment, got %d: %s", wSeg.Code, wSeg.Body.String())
	}
	if ct := wSeg.Header().Get("Content-Type"); ct != "video/mp2t" {
		t.Errorf("expected Content-Type video/mp2t, got %s", ct)
	}
	if !bytes.Equal(wSeg.Body.Bytes(), rawSegment) {
		t.Errorf("segment bytes mismatch")
	}

	// 9. Node Offline Handling:
	// If node goes offline, cached segments (like master.m3u8) still return 200 OK (Cache HIT)
	heartbeatStore.DeleteHeartbeat(nodeID)

	wOfflineCached := httptest.NewRecorder()
	server.ServeHTTP(wOfflineCached, req) // master.m3u8 is cached!
	if wOfflineCached.Code != http.StatusOK {
		t.Errorf("expected cached chunk to play when node is offline (FR-7), got %d", wOfflineCached.Code)
	}

	// But an uncached chunk returns 503 NODE_OFFLINE when node is offline
	uncachedSig := contracts.SignChunkURL(testSigningSecret, nodeID, mediaID, "480p/segment_9999.ts", expFuture)
	reqUncached := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/stream/%s/%s/480p/segment_9999.ts?%s", nodeID, mediaID, uncachedSig), nil)
	wUncached := httptest.NewRecorder()
	server.ServeHTTP(wUncached, reqUncached)

	if wUncached.Code != http.StatusServiceUnavailable {
		t.Errorf("expected 503 Service Unavailable when node is offline and chunk uncached, got %d", wUncached.Code)
	}
}

func TestServer_Node_Heartbeat(t *testing.T) {
	server, _, _, catalogRepo, heartbeatStore := setupTestServer(t)

	// Step 1: Register User & Obtain JWT
	regPayload := `{"username":"carol_node_owner","password":"secure_password_123"}`
	req := httptest.NewRequest(http.MethodPost, "/v1/auth/register", bytes.NewBufferString(regPayload))
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)

	var authRes struct {
		AccessToken string `json:"access_token"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &authRes)

	// Step 2: Register Edge Node
	nodePayload := `{"name":"Carol Node RTX"}`
	req = httptest.NewRequest(http.MethodPost, "/v1/nodes", bytes.NewBufferString(nodePayload))
	req.Header.Set("Authorization", "Bearer "+authRes.AccessToken)
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)

	var nodeRes struct {
		NodeID    string `json:"node_id"`
		NodeToken string `json:"node_token"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &nodeRes)

	hbURL := fmt.Sprintf("/v1/nodes/%s/heartbeat", nodeRes.NodeID)

	validPayload := fmt.Sprintf(`{
		"node_id": %q,
		"sent_at": %q,
		"version": "1.0.0",
		"gpu_caps": {"nvenc": true, "qsv": false},
		"active_streams": 3
	}`, nodeRes.NodeID, time.Now().UTC().Format(time.RFC3339))

	// Subtest 1: Missing X-Node-Token -> 401
	req = httptest.NewRequest(http.MethodPost, hbURL, bytes.NewBufferString(validPayload))
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 without X-Node-Token, got %d", w.Code)
	}

	// Subtest 2: Invalid X-Node-Token -> 401
	req = httptest.NewRequest(http.MethodPost, hbURL, bytes.NewBufferString(validPayload))
	req.Header.Set("X-Node-Token", "invalid-token")
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 with invalid X-Node-Token, got %d", w.Code)
	}

	// Subtest 3a: Path node_id does not match authenticated token node -> 403 Forbidden
	mismatchPathURL := "/v1/nodes/00000000-0000-0000-0000-000000000999/heartbeat"
	req = httptest.NewRequest(http.MethodPost, mismatchPathURL, bytes.NewBufferString(validPayload))
	req.Header.Set("X-Node-Token", nodeRes.NodeToken)
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403 for path node mismatch with token, got %d", w.Code)
	}

	// Subtest 3b: Body node_id does not match path node_id -> 400 Bad Request
	mismatchBodyPayload := fmt.Sprintf(`{
		"node_id": "00000000-0000-0000-0000-000000000999",
		"sent_at": %q,
		"version": "1.0.0",
		"active_streams": 3
	}`, time.Now().UTC().Format(time.RFC3339))
	req = httptest.NewRequest(http.MethodPost, hbURL, bytes.NewBufferString(mismatchBodyPayload))
	req.Header.Set("X-Node-Token", nodeRes.NodeToken)
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for path/body node mismatch, got %d", w.Code)
	}

	// Subtest 4: Malformed JSON payload -> 400
	req = httptest.NewRequest(http.MethodPost, hbURL, bytes.NewBufferString("{invalid-json"))
	req.Header.Set("X-Node-Token", nodeRes.NodeToken)
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for malformed json, got %d", w.Code)
	}

	// Subtest 5: Exceeding 10 active streams in payload -> 400 Validation Error
	overCapPayload := fmt.Sprintf(`{
		"node_id": %q,
		"sent_at": %q,
		"version": "1.0.0",
		"active_streams": 15
	}`, nodeRes.NodeID, time.Now().UTC().Format(time.RFC3339))
	req = httptest.NewRequest(http.MethodPost, hbURL, bytes.NewBufferString(overCapPayload))
	req.Header.Set("X-Node-Token", nodeRes.NodeToken)
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for active_streams > 10, got %d: %s", w.Code, w.Body.String())
	}

	// Subtest 6: Valid Heartbeat -> 200 OK
	req = httptest.NewRequest(http.MethodPost, hbURL, bytes.NewBufferString(validPayload))
	req.Header.Set("X-Node-Token", nodeRes.NodeToken)
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for valid heartbeat, got %d: %s", w.Code, w.Body.String())
	}

	// Subtest 7: Verify catalog online indicator reflects heartbeat
	_ = catalogRepo.UpsertManifest(context.Background(), &domain.CatalogItem{
		MediaID:     "88888888-8888-8888-8888-888888888888",
		NodeID:      nodeRes.NodeID,
		Title:       "Carol Video",
		Description: "Video from online node",
		DurationS:   60,
		ContentHash: "hash-888",
		IsPublic:    true,
		UpdatedAt:   time.Now().UTC(),
	})

	req = httptest.NewRequest(http.MethodGet, "/v1/catalog", nil)
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for catalog, got %d", w.Code)
	}

	var catRes struct {
		Items []struct {
			MediaID      string `json:"media_id"`
			IsNodeOnline bool   `json:"is_node_online"`
		} `json:"items"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &catRes); err != nil {
		t.Fatalf("failed to decode catalog response: %v", err)
	}

	var foundOnline bool
	for _, it := range catRes.Items {
		if it.MediaID == "88888888-8888-8888-8888-888888888888" && it.IsNodeOnline {
			foundOnline = true
			break
		}
	}
	if !foundOnline {
		t.Errorf("expected catalog item to report is_node_online = true")
	}

	// Subtest 8: Simulated HeartbeatStore error -> 500
	heartbeatStore.SimulateErr = fmt.Errorf("storage failure")
	req = httptest.NewRequest(http.MethodPost, hbURL, bytes.NewBufferString(validPayload))
	req.Header.Set("X-Node-Token", nodeRes.NodeToken)
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500 for storage error, got %d", w.Code)
	}
}

func TestServer_AdaptiveRenditionThrottling(t *testing.T) {
	nodeID := "33333333-3333-3333-3333-333333333333"
	mediaID := "44444444-4444-4444-4444-444444444444"

	multiMaster := "#EXTM3U\n#EXT-X-VERSION:3\n#EXT-X-STREAM-INF:BANDWIDTH=5000000,RESOLUTION=1920x1080\n1080p/index.m3u8\n#EXT-X-STREAM-INF:BANDWIDTH=2500000,RESOLUTION=1280x720\n720p/index.m3u8\n#EXT-X-STREAM-INF:BANDWIDTH=1000000,RESOLUTION=854x480\n480p/index.m3u8\n"

	puller := &mockTunnelPuller{
		content: map[string][]byte{
			"master.m3u8": []byte(multiMaster),
		},
		contentType: map[string]string{
			"master.m3u8": "application/vnd.apple.mpegurl",
		},
	}

	server, _, _, _, heartbeatStore, _ := setupTestServerFull(t, nil, puller)

	// Subtest 1: LoadNormal (2 active streams) -> All renditions present
	_ = heartbeatStore.RecordHeartbeat(context.Background(), &contracts.Heartbeat{
		NodeID:        nodeID,
		SentAt:        time.Now(),
		Version:       "1.0.0",
		ActiveStreams: 2,
	}, time.Minute)

	exp := time.Now().Add(time.Hour)
	sig := contracts.SignChunkURL(testSigningSecret, nodeID, mediaID, "master.m3u8", exp)
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/stream/%s/%s/master.m3u8?%s", nodeID, mediaID, sig), nil)
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK under LoadNormal, got %d", w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, "1080p/index.m3u8") || !strings.Contains(body, "720p/index.m3u8") || !strings.Contains(body, "480p/index.m3u8") {
		t.Errorf("expected all 3 renditions under LoadNormal, got:\n%s", body)
	}

	// Subtest 2: LoadModerate (5 active streams) -> 1080p pruned, 720p & 480p retained
	_ = heartbeatStore.RecordHeartbeat(context.Background(), &contracts.Heartbeat{
		NodeID:        nodeID,
		SentAt:        time.Now(),
		Version:       "1.0.0",
		ActiveStreams: 5,
	}, time.Minute)

	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/stream/%s/%s/master.m3u8?%s", nodeID, mediaID, sig), nil)
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK under LoadModerate, got %d", w.Code)
	}
	body = w.Body.String()
	if strings.Contains(body, "1080p/index.m3u8") {
		t.Errorf("expected 1080p to be throttled out under LoadModerate, got:\n%s", body)
	}
	if !strings.Contains(body, "720p/index.m3u8") || !strings.Contains(body, "480p/index.m3u8") {
		t.Errorf("expected 720p and 480p under LoadModerate, got:\n%s", body)
	}

	// Subtest 3: LoadSaturated (8 active streams) -> 1080p and 720p pruned, 480p retained
	_ = heartbeatStore.RecordHeartbeat(context.Background(), &contracts.Heartbeat{
		NodeID:        nodeID,
		SentAt:        time.Now(),
		Version:       "1.0.0",
		ActiveStreams: 8,
	}, time.Minute)

	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/stream/%s/%s/master.m3u8?%s", nodeID, mediaID, sig), nil)
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK under LoadSaturated, got %d", w.Code)
	}
	body = w.Body.String()
	if strings.Contains(body, "1080p/index.m3u8") || strings.Contains(body, "720p/index.m3u8") {
		t.Errorf("expected 1080p and 720p to be throttled out under LoadSaturated, got:\n%s", body)
	}
	if !strings.Contains(body, "480p/index.m3u8") {
		t.Errorf("expected 480p under LoadSaturated, got:\n%s", body)
	}

	// Subtest 4: Safety Fallback - Single 1080p rendition under Saturated load is preserved
	singleMaster := "#EXTM3U\n#EXT-X-VERSION:3\n#EXT-X-STREAM-INF:BANDWIDTH=5000000,RESOLUTION=1920x1080\n1080p/index.m3u8\n"
	pullerSingle := &mockTunnelPuller{
		content: map[string][]byte{
			"master.m3u8": []byte(singleMaster),
		},
		contentType: map[string]string{
			"master.m3u8": "application/vnd.apple.mpegurl",
		},
	}
	serverSingle, _, _, _, heartbeatSingle, _ := setupTestServerFull(t, nil, pullerSingle)
	_ = heartbeatSingle.RecordHeartbeat(context.Background(), &contracts.Heartbeat{
		NodeID:        nodeID,
		SentAt:        time.Now(),
		Version:       "1.0.0",
		ActiveStreams: 8,
	}, time.Minute)

	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/stream/%s/%s/master.m3u8?%s", nodeID, mediaID, sig), nil)
	w = httptest.NewRecorder()
	serverSingle.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for fallback, got %d", w.Code)
	}
	body = w.Body.String()
	if !strings.Contains(body, "1080p/index.m3u8") {
		t.Errorf("expected fallback to preserve 1080p rather than producing empty playlist, got:\n%s", body)
	}
}

type mockPrefetcher struct {
	triggeredNode  string
	triggeredMedia string
	triggeredPath  string
	calledCount    int
}

func (m *mockPrefetcher) TriggerPrefetch(ctx context.Context, nodeID, mediaID, initialChunkPath string) {
	m.calledCount++
	m.triggeredNode = nodeID
	m.triggeredMedia = mediaID
	m.triggeredPath = initialChunkPath
}

func TestServer_FirstChunkPrefetchOnPlayback(t *testing.T) {
	nodeID := "55555555-5555-5555-5555-555555555555"
	mediaID := "66666666-6666-6666-6666-666666666666"

	server, _, _, catalogRepo, heartbeatStore, _ := setupTestServerFull(t, nil, nil)

	mockPf := &mockPrefetcher{}
	server.SetPrefetcher(mockPf)

	// Setup online node
	_ = heartbeatStore.RecordHeartbeat(context.Background(), &contracts.Heartbeat{
		NodeID:        nodeID,
		SentAt:        time.Now(),
		Version:       "1.0.0",
		ActiveStreams: 1,
	}, time.Minute)

	// Setup catalog media
	_ = catalogRepo.UpsertManifest(context.Background(), &domain.CatalogItem{
		MediaID:     mediaID,
		NodeID:      nodeID,
		Title:       "Prefetch Test Video",
		DurationS:   120,
		ContentHash: "hash-prefetch",
		IsPublic:    true,
		UpdatedAt:   time.Now().UTC(),
	})

	// Request remote playback
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/media/%s/playback", mediaID), nil)
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for playback, got %d: %s", w.Code, w.Body.String())
	}

	if mockPf.calledCount != 1 {
		t.Errorf("expected prefetcher to be triggered once, called %d times", mockPf.calledCount)
	}
	if mockPf.triggeredNode != nodeID || mockPf.triggeredMedia != mediaID || mockPf.triggeredPath != "1080p/segment_000.ts" {
		t.Errorf("unexpected prefetch params: node=%s, media=%s, path=%s", mockPf.triggeredNode, mockPf.triggeredMedia, mockPf.triggeredPath)
	}
}
