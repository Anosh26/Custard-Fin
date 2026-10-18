package httpapi_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/anosh26/custardfin/apps/api/internal/adapters/auth"
	"github.com/anosh26/custardfin/apps/api/internal/adapters/httpapi"
	"github.com/anosh26/custardfin/apps/api/internal/adapters/memory"
	"github.com/anosh26/custardfin/apps/api/internal/domain"
	"github.com/anosh26/custardfin/libs/contracts"
	"golang.org/x/crypto/bcrypt"
)

func setupTestServer(t *testing.T) (*httpapi.Server, *memory.InMemoryUserRepository, *memory.InMemoryNodeRepository, *memory.InMemoryCatalogRepository, *memory.InMemoryHeartbeatStore) {
	t.Helper()

	userRepo := memory.NewInMemoryUserRepository()
	nodeRepo := memory.NewInMemoryNodeRepository()
	catalogRepo := memory.NewInMemoryCatalogRepository()
	heartbeatStore := memory.NewInMemoryHeartbeatStore()

	hasher := auth.NewBcryptHasher(bcrypt.MinCost)
	tokenService, err := auth.NewJWTService([]byte("test-secret-key-at-least-32-bytes-long!"))
	if err != nil {
		t.Fatalf("failed to create jwt service: %v", err)
	}

	server := httpapi.NewServer(
		httpapi.Config{
			Addr:         "127.0.0.1:0",
			HeartbeatTTL: 20 * time.Second,
		},
		userRepo,
		nodeRepo,
		catalogRepo,
		heartbeatStore,
		hasher,
		tokenService,
		nil,
	)

	return server, userRepo, nodeRepo, catalogRepo, heartbeatStore
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

	if !strings.Contains(pbRes.StreamURL, mediaID) {
		t.Errorf("expected stream url to contain media ID, got %s", pbRes.StreamURL)
	}
	if pbRes.RouteType != "direct_lan" {
		t.Errorf("expected direct_lan route type, got %s", pbRes.RouteType)
	}
	if pbRes.ExpiresAt == "" {
		t.Errorf("expected non-empty expires_at")
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

