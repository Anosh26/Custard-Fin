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

func setupTestServer(t *testing.T) (*httpapi.Server, *memory.InMemoryUserRepository, *memory.InMemoryNodeRepository, *memory.InMemoryCatalogRepository) {
	t.Helper()

	userRepo := memory.NewInMemoryUserRepository()
	nodeRepo := memory.NewInMemoryNodeRepository()
	catalogRepo := memory.NewInMemoryCatalogRepository()

	hasher := auth.NewBcryptHasher(bcrypt.MinCost)
	tokenService, err := auth.NewJWTService([]byte("test-secret-key-at-least-32-bytes-long!"))
	if err != nil {
		t.Fatalf("failed to create jwt service: %v", err)
	}

	server := httpapi.NewServer(
		httpapi.Config{Addr: "127.0.0.1:0"},
		userRepo,
		nodeRepo,
		catalogRepo,
		hasher,
		tokenService,
		nil,
	)

	return server, userRepo, nodeRepo, catalogRepo
}

func TestServer_Healthz(t *testing.T) {
	server, _, _, _ := setupTestServer(t)

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
	server, _, _, _ := setupTestServer(t)

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
	server, _, _, _ := setupTestServer(t)

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
	server, _, _, _ := setupTestServer(t)

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
	server, _, _, _ := setupTestServer(t)

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
	server, _, _, catalogRepo := setupTestServer(t)

	// 1. Not Found
	req := httptest.NewRequest(http.MethodGet, "/v1/media/non-existent-id/playback", nil)
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404 for missing media, got %d", w.Code)
	}

	// 2. Existing Media
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

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
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

	// 3. Simulated Repo Error
	catalogRepo.SimulateErr = fmt.Errorf("database query failure")
	req = httptest.NewRequest(http.MethodGet, "/v1/media/"+mediaID+"/playback", nil)
	w = httptest.NewRecorder()
	server.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500 Internal Server Error, got %d", w.Code)
	}
}

