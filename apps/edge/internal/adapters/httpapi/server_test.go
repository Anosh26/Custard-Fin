package httpapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/anosh26/custardfin/apps/edge/internal/adapters/ffmpeg"
	"github.com/anosh26/custardfin/apps/edge/internal/domain"
	"github.com/anosh26/custardfin/apps/edge/internal/ports"
	"github.com/anosh26/custardfin/libs/contracts"
)

// FakeMediaRepository for HTTP handler testing
type FakeMediaRepo struct {
	mu        sync.RWMutex
	items     map[string]*domain.Media
	failList  bool
	failGetID bool
}

func (r *FakeMediaRepo) FindByPath(ctx context.Context, path string) (*domain.Media, error) {
	return nil, nil
}

func (r *FakeMediaRepo) FindByID(ctx context.Context, id string) (*domain.Media, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if r.failGetID {
		return nil, fmt.Errorf("simulated db failure")
	}
	return r.items[id], nil
}

func (r *FakeMediaRepo) FindByHash(ctx context.Context, hash string) (*domain.Media, error) {
	return nil, nil
}

func (r *FakeMediaRepo) Save(ctx context.Context, media *domain.Media) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.items[media.ID] = media
	return nil
}

func (r *FakeMediaRepo) List(ctx context.Context) ([]*domain.Media, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if r.failList {
		return nil, fmt.Errorf("simulated db list failure")
	}
	var list []*domain.Media
	for _, m := range r.items {
		list = append(list, m)
	}
	return list, nil
}

func (r *FakeMediaRepo) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.items, id)
	return nil
}

func (r *FakeMediaRepo) SetPublic(ctx context.Context, id string, isPublic bool) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	m, exists := r.items[id]
	if !exists {
		return fmt.Errorf("media not found")
	}
	m.IsPublic = isPublic
	return nil
}

type FakeSyncer struct {
	synced   map[string]bool
	unsynced map[string]bool
	failSync bool
}

func (s *FakeSyncer) SyncManifest(ctx context.Context, media *domain.Media, renditions []contracts.Rendition) error {
	if s.failSync {
		return fmt.Errorf("simulated cloud sync failure")
	}
	s.synced[media.ID] = true
	delete(s.unsynced, media.ID)
	return nil
}

func (s *FakeSyncer) UnsyncManifest(ctx context.Context, mediaID string) error {
	s.unsynced[mediaID] = true
	delete(s.synced, mediaID)
	return nil
}

// FakeWatchRepo for HTTP handler testing
type FakeWatchRepo struct {
	mu    sync.RWMutex
	items map[string]*domain.WatchProgress
}

func (r *FakeWatchRepo) GetWatchProgress(ctx context.Context, mediaID string) (*domain.WatchProgress, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.items[mediaID], nil
}

func (r *FakeWatchRepo) SaveWatchProgress(ctx context.Context, wp *domain.WatchProgress) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.items[wp.MediaID] = wp
	return nil
}

func (r *FakeWatchRepo) DeleteWatchProgress(ctx context.Context, mediaID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.items, mediaID)
	return nil
}

func setupTestServer(t *testing.T) (*Server, *FakeMediaRepo, *FakeWatchRepo, *FakeSyncer, string) {
	tempDir := t.TempDir()
	transcodeDir := filepath.Join(tempDir, "transcoded")
	_ = os.MkdirAll(transcodeDir, 0755)

	mediaRepo := &FakeMediaRepo{items: make(map[string]*domain.Media)}
	watchRepo := &FakeWatchRepo{items: make(map[string]*domain.WatchProgress)}
	syncer := &FakeSyncer{synced: make(map[string]bool), unsynced: make(map[string]bool)}

	cfg := Config{
		Addr:         ":0",
		TranscodeDir: transcodeDir,
	}

	server := NewServer(cfg, mediaRepo, watchRepo, syncer, nil)
	return server, mediaRepo, watchRepo, syncer, transcodeDir
}

func TestServer_Healthz(t *testing.T) {
	server, _, _, _, _ := setupTestServer(t)
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()

	server.httpServer.Handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), `"status":"ok"`) {
		t.Errorf("expected status ok body, got %s", w.Body.String())
	}
}

func TestServer_CORS_Middleware(t *testing.T) {
	server, _, _, _, _ := setupTestServer(t)

	// Test OPTIONS preflight
	req := httptest.NewRequest(http.MethodOptions, "/api/media", nil)
	w := httptest.NewRecorder()

	server.httpServer.Handler.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected status 204 No Content for OPTIONS, got %d", w.Code)
	}
	if w.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Errorf("missing Access-Control-Allow-Origin header")
	}
	if !strings.Contains(w.Header().Get("Access-Control-Allow-Methods"), "GET") {
		t.Errorf("missing allowed methods in header")
	}
}

func TestServer_MediaEndpoints(t *testing.T) {
	server, mediaRepo, _, _, _ := setupTestServer(t)
	ctx := context.Background()

	m1, _ := domain.NewMedia("/v/movie1.mp4", 1000, time.Now(), "hash1", 120, 1920, 1080, "h264", "aac")
	_ = mediaRepo.Save(ctx, m1)

	// 1. List Media
	reqList := httptest.NewRequest(http.MethodGet, "/api/media", nil)
	wList := httptest.NewRecorder()
	server.httpServer.Handler.ServeHTTP(wList, reqList)

	if wList.Code != http.StatusOK {
		t.Fatalf("expected 200 for list media, got %d", wList.Code)
	}
	var listResp []*domain.Media
	if err := json.NewDecoder(wList.Body).Decode(&listResp); err != nil {
		t.Fatal(err)
	}
	if len(listResp) != 1 || listResp[0].ID != m1.ID {
		t.Errorf("unexpected list response: %+v", listResp)
	}

	// 2. Get Media Found
	reqGet := httptest.NewRequest(http.MethodGet, "/api/media/"+m1.ID, nil)
	wGet := httptest.NewRecorder()
	server.httpServer.Handler.ServeHTTP(wGet, reqGet)

	if wGet.Code != http.StatusOK {
		t.Fatalf("expected 200 for get media, got %d", wGet.Code)
	}

	// 3. Get Media Not Found
	reqNotFound := httptest.NewRequest(http.MethodGet, "/api/media/non-existent-id", nil)
	wNotFound := httptest.NewRecorder()
	server.httpServer.Handler.ServeHTTP(wNotFound, reqNotFound)

	if wNotFound.Code != http.StatusNotFound {
		t.Errorf("expected 404 for missing media, got %d", wNotFound.Code)
	}
}

func TestServer_Streaming_HeadersAndByteRanges(t *testing.T) {
	server, mediaRepo, _, _, transcodeDir := setupTestServer(t)
	ctx := context.Background()

	m, _ := domain.NewMedia("/v/movie.mp4", 5000, time.Now(), "hash", 60, 1920, 1080, "h264", "aac")
	_ = mediaRepo.Save(ctx, m)

	// Create test HLS assets in transcodeDir/m.ID/
	mediaDir := filepath.Join(transcodeDir, m.ID)
	renditionDir := filepath.Join(mediaDir, "1080p")
	_ = os.MkdirAll(renditionDir, 0755)

	playlistFile := filepath.Join(mediaDir, "master.m3u8")
	_ = os.WriteFile(playlistFile, []byte("#EXTM3U\n#EXT-X-VERSION:3\n"), 0644)

	chunkFile := filepath.Join(renditionDir, "segment_0000.ts")
	chunkBytes := []byte("fake MPEG-TS chunk segment bytes for testing byte ranges")
	_ = os.WriteFile(chunkFile, chunkBytes, 0644)

	// 1. Fetch Playlist (.m3u8)
	reqPlaylist := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/stream/%s/master.m3u8", m.ID), nil)
	wPlaylist := httptest.NewRecorder()
	server.httpServer.Handler.ServeHTTP(wPlaylist, reqPlaylist)

	if wPlaylist.Code != http.StatusOK {
		t.Fatalf("expected 200 for playlist, got %d", wPlaylist.Code)
	}
	if wPlaylist.Header().Get("Content-Type") != "application/vnd.apple.mpegurl" {
		t.Errorf("expected m3u8 mime type, got %s", wPlaylist.Header().Get("Content-Type"))
	}
	if !strings.Contains(wPlaylist.Header().Get("Cache-Control"), "no-cache") {
		t.Errorf("expected no-cache for playlist, got %s", wPlaylist.Header().Get("Cache-Control"))
	}

	// 2. Fetch Segment (.ts) with Full Request
	reqChunk := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/stream/%s/1080p/segment_0000.ts", m.ID), nil)
	wChunk := httptest.NewRecorder()
	server.httpServer.Handler.ServeHTTP(wChunk, reqChunk)

	if wChunk.Code != http.StatusOK {
		t.Fatalf("expected 200 for chunk, got %d", wChunk.Code)
	}
	if wChunk.Header().Get("Content-Type") != "video/mp2t" {
		t.Errorf("expected video/mp2t mime type, got %s", wChunk.Header().Get("Content-Type"))
	}
	if wChunk.Header().Get("Accept-Ranges") != "bytes" {
		t.Errorf("expected Accept-Ranges: bytes header")
	}

	// 3. Fetch Segment with Range Request
	reqRange := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/stream/%s/1080p/segment_0000.ts", m.ID), nil)
	reqRange.Header.Set("Range", "bytes=0-9")
	wRange := httptest.NewRecorder()
	server.httpServer.Handler.ServeHTTP(wRange, reqRange)

	if wRange.Code != http.StatusPartialContent {
		t.Fatalf("expected 206 Partial Content, got %d", wRange.Code)
	}
	if wRange.Body.Len() != 10 {
		t.Errorf("expected 10 bytes returned, got %d", wRange.Body.Len())
	}
	if !strings.Contains(wRange.Header().Get("Content-Range"), "bytes 0-9/") {
		t.Errorf("expected Content-Range bytes 0-9/..., got %s", wRange.Header().Get("Content-Range"))
	}

	// 4. Path Traversal Attempt (via Mux and Direct Handler)
	reqTraversal := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/stream/%s/../../secret.txt", m.ID), nil)
	wTraversal := httptest.NewRecorder()
	server.httpServer.Handler.ServeHTTP(wTraversal, reqTraversal)

	// Go's ServeMux intercepts relative dots and redirects (307/301) to cleaned URL, or returns 400/404
	if wTraversal.Code != http.StatusBadRequest && wTraversal.Code != http.StatusNotFound && wTraversal.Code != http.StatusTemporaryRedirect {
		t.Errorf("expected 400, 404, or 307 on path traversal, got %d", wTraversal.Code)
	}

	// Direct handler path traversal test
	reqDirect := httptest.NewRequest(http.MethodGet, "/stream/"+m.ID+"/test", nil)
	reqDirect.SetPathValue("id", m.ID)
	reqDirect.SetPathValue("path", "../secret.txt")
	wDirect := httptest.NewRecorder()
	server.handleStream(wDirect, reqDirect)
	if wDirect.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for direct handler path traversal, got %d", wDirect.Code)
	}
}

func TestServer_WatchProgress_Endpoints(t *testing.T) {
	server, _, _, _, _ := setupTestServer(t)

	// 1. Initial GET on unwatched media returns default 0s
	reqGet := httptest.NewRequest(http.MethodGet, "/api/watch/media-xyz", nil)
	wGet := httptest.NewRecorder()
	server.httpServer.Handler.ServeHTTP(wGet, reqGet)

	if wGet.Code != http.StatusOK {
		t.Fatalf("expected 200 for default watch progress, got %d", wGet.Code)
	}
	var wpResp domain.WatchProgress
	if err := json.NewDecoder(wGet.Body).Decode(&wpResp); err != nil {
		t.Fatal(err)
	}
	if wpResp.PositionS != 0 {
		t.Errorf("expected default position 0, got %d", wpResp.PositionS)
	}

	// 2. POST update progress
	payload := `{"position_s": 320, "completed": false}`
	reqPost := httptest.NewRequest(http.MethodPost, "/api/watch/media-xyz", bytes.NewBufferString(payload))
	reqPost.Header.Set("Content-Type", "application/json")
	wPost := httptest.NewRecorder()
	server.httpServer.Handler.ServeHTTP(wPost, reqPost)

	if wPost.Code != http.StatusOK {
		t.Fatalf("expected 200 for post watch progress, got %d (%s)", wPost.Code, wPost.Body.String())
	}

	// 3. GET after update returns saved position
	reqGet2 := httptest.NewRequest(http.MethodGet, "/api/watch/media-xyz", nil)
	wGet2 := httptest.NewRecorder()
	server.httpServer.Handler.ServeHTTP(wGet2, reqGet2)

	var wpUpdated domain.WatchProgress
	_ = json.NewDecoder(wGet2.Body).Decode(&wpUpdated)
	if wpUpdated.PositionS != 320 {
		t.Errorf("expected position 320, got %d", wpUpdated.PositionS)
	}

	// 4. POST with invalid payload returns 400
	reqBadJSON := httptest.NewRequest(http.MethodPost, "/api/watch/media-xyz", bytes.NewBufferString("invalid json"))
	wBadJSON := httptest.NewRecorder()
	server.httpServer.Handler.ServeHTTP(wBadJSON, reqBadJSON)
	if wBadJSON.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for bad json, got %d", wBadJSON.Code)
	}

	// 5. POST with negative position returns 400
	reqNeg := httptest.NewRequest(http.MethodPost, "/api/watch/media-xyz", bytes.NewBufferString(`{"position_s": -50}`))
	wNeg := httptest.NewRecorder()
	server.httpServer.Handler.ServeHTTP(wNeg, reqNeg)
	if wNeg.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for negative position, got %d", wNeg.Code)
	}
}

func TestServer_StartAndShutdown(t *testing.T) {
	server, _, _, _, _ := setupTestServer(t)

	if err := server.Start(); err != nil {
		t.Fatalf("server start failed: %v", err)
	}

	addr := server.Addr()
	if !strings.Contains(addr, ":") {
		t.Errorf("expected valid host:port address, got %s", addr)
	}

	// Verify server is accepting connections
	resp, err := http.Get("http://" + addr + "/healthz")
	if err != nil {
		t.Fatalf("healthz request failed: %v", err)
	}
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		t.Fatalf("server shutdown returned error: %v", err)
	}
}

func TestServer_FullLocalStreamingIntegration(t *testing.T) {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		t.Skip("ffmpeg not found, skipping full streaming integration test")
	}

	tempDir := t.TempDir()
	sourceFile := filepath.Join(tempDir, "source.mp4")
	transcodeDir := filepath.Join(tempDir, "transcoded")
	_ = os.MkdirAll(transcodeDir, 0755)

	// 1. Generate test video clip
	cmd := exec.Command("ffmpeg", "-y",
		"-f", "lavfi", "-i", "testsrc=duration=4:size=320x240:rate=30",
		"-f", "lavfi", "-i", "sine=frequency=1000:duration=4",
		"-c:v", "libx264", "-c:a", "aac",
		sourceFile,
	)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("failed generating test clip: %v (%s)", err, string(out))
	}

	// 2. Pre-encode to HLS using SubprocessTranscoder
	transcoder := ffmpeg.NewSubprocessTranscoder("", nil)
	m, _ := domain.NewMedia(sourceFile, 10000, time.Now(), "hash_stream_e2e", 4, 320, 240, "h264", "aac")

	mediaTranscodeDir := filepath.Join(transcodeDir, m.ID)
	result, err := transcoder.Transcode(context.Background(), ports.TranscodeJob{
		MediaID:       m.ID,
		InputPath:     sourceFile,
		OutputDir:     mediaTranscodeDir,
		SourceWidth:   320,
		SourceHeight:  240,
		PreferredMode: ports.EncoderModeCPU,
	})
	if err != nil {
		t.Fatalf("transcode failed: %v", err)
	}

	// 3. Setup Server and Register Media
	mediaRepo := &FakeMediaRepo{items: map[string]*domain.Media{m.ID: m}}
	watchRepo := &FakeWatchRepo{items: make(map[string]*domain.WatchProgress)}

	server := NewServer(Config{Addr: ":0", TranscodeDir: transcodeDir}, mediaRepo, watchRepo, nil, nil)
	if err := server.Start(); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = server.Shutdown(context.Background()) }()

	baseURL := "http://" + server.Addr()

	// 4. Fetch Master Playlist over HTTP
	masterURL := fmt.Sprintf("%s/stream/%s/master.m3u8", baseURL, m.ID)
	resp, err := http.Get(masterURL)
	if err != nil {
		t.Fatalf("failed to fetch master playlist: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 for master playlist, got %d", resp.StatusCode)
	}
	masterContent, _ := io.ReadAll(resp.Body)
	if !strings.Contains(string(masterContent), "480p/index.m3u8") {
		t.Errorf("master playlist missing rendition path: %s", string(masterContent))
	}

	// 5. Fetch Child Playlist over HTTP
	childURL := fmt.Sprintf("%s/stream/%s/%s", baseURL, m.ID, result.Renditions[0].PlaylistPath)
	respChild, err := http.Get(childURL)
	if err != nil {
		t.Fatalf("failed to fetch child playlist: %v", err)
	}
	defer respChild.Body.Close()

	if respChild.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 for child playlist, got %d", respChild.StatusCode)
	}

	// 6. Fetch First Chunk Segment over HTTP with byte range
	chunkURL := fmt.Sprintf("%s/stream/%s/%s/segment_0000.ts", baseURL, m.ID, result.Renditions[0].Resolution)
	reqChunk, _ := http.NewRequest(http.MethodGet, chunkURL, nil)
	reqChunk.Header.Set("Range", "bytes=0-187") // 188-byte MPEG-TS packet
	respChunk, err := http.DefaultClient.Do(reqChunk)
	if err != nil {
		t.Fatalf("failed fetching chunk with range: %v", err)
	}
	defer respChunk.Body.Close()

	if respChunk.StatusCode != http.StatusPartialContent {
		t.Fatalf("expected 206 Partial Content for chunk range, got %d", respChunk.StatusCode)
	}
	if respChunk.Header.Get("Content-Type") != "video/mp2t" {
		t.Errorf("expected video/mp2t, got %s", respChunk.Header.Get("Content-Type"))
	}
}

func TestServer_TogglePublic_Endpoint(t *testing.T) {
	server, mediaRepo, _, syncer, _ := setupTestServer(t)
	ctx := context.Background()

	m, _ := domain.NewMedia("/media/open_movie.mp4", 1000, time.Now(), "hash_toggle", 120, 1920, 1080, "h264", "aac")
	_ = mediaRepo.Save(ctx, m)

	// 1. Toggle Public = true
	bodyTrue := bytes.NewBufferString(`{"is_public": true}`)
	reqTrue := httptest.NewRequest(http.MethodPut, "/api/media/"+m.ID+"/public", bodyTrue)
	wTrue := httptest.NewRecorder()
	server.httpServer.Handler.ServeHTTP(wTrue, reqTrue)

	if wTrue.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for toggle public true, got %d: %s", wTrue.Code, wTrue.Body.String())
	}
	if !syncer.synced[m.ID] {
		t.Errorf("expected syncer.SyncManifest to be called for media %s", m.ID)
	}

	// 2. Toggle Public = false (Unpublish)
	bodyFalse := bytes.NewBufferString(`{"is_public": false}`)
	reqFalse := httptest.NewRequest(http.MethodPut, "/api/media/"+m.ID+"/public", bodyFalse)
	wFalse := httptest.NewRecorder()
	server.httpServer.Handler.ServeHTTP(wFalse, reqFalse)

	if wFalse.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for toggle public false, got %d: %s", wFalse.Code, wFalse.Body.String())
	}
	if !syncer.unsynced[m.ID] {
		t.Errorf("expected syncer.UnsyncManifest to be called for media %s", m.ID)
	}

	// 3. Non-existent media -> 404
	req404 := httptest.NewRequest(http.MethodPut, "/api/media/non-existent-id/public", bytes.NewBufferString(`{"is_public": true}`))
	w404 := httptest.NewRecorder()
	server.httpServer.Handler.ServeHTTP(w404, req404)
	if w404.Code != http.StatusNotFound {
		t.Errorf("expected 404 for non-existent media, got %d", w404.Code)
	}

	// 4. Invalid JSON -> 400
	req400 := httptest.NewRequest(http.MethodPut, "/api/media/"+m.ID+"/public", bytes.NewBufferString(`invalid json`))
	w400 := httptest.NewRecorder()
	server.httpServer.Handler.ServeHTTP(w400, req400)
	if w400.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for bad json, got %d", w400.Code)
	}
}

func TestServer_ActiveStreamsTracker(t *testing.T) {
	tmpDir := t.TempDir()
	mediaRepo := &FakeMediaRepo{items: make(map[string]*domain.Media)}
	watchRepo := &FakeWatchRepo{items: make(map[string]*domain.WatchProgress)}
	server := NewServer(Config{Addr: "127.0.0.1:0", TranscodeDir: tmpDir}, mediaRepo, watchRepo, nil, nil)

	// Check initial active streams
	if n := server.ActiveStreams(); n != 0 {
		t.Fatalf("expected 0 active streams initially, got %d", n)
	}

	// Create test file in transcode dir
	m := &domain.Media{
		ID:        "stream-test-media-1",
		Title:     "Stream Test",
		DurationS: 60,
	}
	_ = mediaRepo.Save(context.Background(), m)
	mediaDir := filepath.Join(tmpDir, m.ID)
	_ = os.MkdirAll(mediaDir, 0755)
	_ = os.WriteFile(filepath.Join(mediaDir, "chunk_000.ts"), []byte("ts-content"), 0644)

	// Stream from client 1
	req1 := httptest.NewRequest(http.MethodGet, "/stream/"+m.ID+"/chunk_000.ts", nil)
	req1.RemoteAddr = "192.168.1.10:54321"
	w1 := httptest.NewRecorder()
	server.httpServer.Handler.ServeHTTP(w1, req1)
	if w1.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w1.Code)
	}

	if n := server.ActiveStreams(); n != 1 {
		t.Fatalf("expected 1 active stream after client 1, got %d", n)
	}

	// Stream from client 2
	req2 := httptest.NewRequest(http.MethodGet, "/stream/"+m.ID+"/chunk_000.ts", nil)
	req2.RemoteAddr = "192.168.1.20:54322"
	w2 := httptest.NewRecorder()
	server.httpServer.Handler.ServeHTTP(w2, req2)

	if n := server.ActiveStreams(); n != 2 {
		t.Fatalf("expected 2 active streams after client 2, got %d", n)
	}

	// Repeated stream from client 1 does not inflate count
	w1Again := httptest.NewRecorder()
	server.httpServer.Handler.ServeHTTP(w1Again, req1)
	if n := server.ActiveStreams(); n != 2 {
		t.Fatalf("expected still 2 active streams for repeated requests, got %d", n)
	}
}

func TestServer_DeleteMedia(t *testing.T) {
	tempDir := t.TempDir()
	sourceFile := filepath.Join(tempDir, "video_to_delete.mp4")
	_ = os.WriteFile(sourceFile, []byte("fake-mp4-data"), 0644)

	transcodeDir := filepath.Join(tempDir, "transcoded")
	m, _ := domain.NewMedia(sourceFile, 1000, time.Now(), "hash_del_test", 60, 1280, 720, "h264", "aac")

	mediaTranscodeDir := filepath.Join(transcodeDir, m.ID)
	_ = os.MkdirAll(mediaTranscodeDir, 0755)
	_ = os.WriteFile(filepath.Join(mediaTranscodeDir, "master.m3u8"), []byte("#EXTM3U"), 0644)

	mediaRepo := &FakeMediaRepo{
		items: map[string]*domain.Media{
			m.ID: m,
		},
	}
	syncer := &FakeSyncer{
		synced:   map[string]bool{m.ID: true},
		unsynced: make(map[string]bool),
	}

	cfg := Config{
		Addr:         ":0",
		NodeID:       "test-node",
		TranscodeDir: transcodeDir,
	}
	server := NewServer(cfg, mediaRepo, nil, syncer, nil)

	// 1. Delete media with ?delete_file=true
	req := httptest.NewRequest(http.MethodDelete, "/api/media/"+m.ID+"?delete_file=true", nil)
	w := httptest.NewRecorder()
	server.httpServer.Handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}

	// Verify SQLite repo deleted it
	remaining, err := mediaRepo.FindByID(context.Background(), m.ID)
	if err != nil || remaining != nil {
		t.Errorf("expected media to be deleted from repo, found: %v", remaining)
	}

	// Verify cloud unsyncer was called
	if !syncer.unsynced[m.ID] {
		t.Errorf("expected syncer.UnsyncManifest to be called for %s", m.ID)
	}

	// Verify transcoded directory was purged
	if _, err := os.Stat(mediaTranscodeDir); !os.IsNotExist(err) {
		t.Errorf("expected transcoded dir to be removed, but still exists")
	}

	// Verify physical source file was deleted
	if _, err := os.Stat(sourceFile); !os.IsNotExist(err) {
		t.Errorf("expected source video file to be removed from disk, but still exists")
	}

	// 2. Subsequent DELETE returns 404
	req404 := httptest.NewRequest(http.MethodDelete, "/api/media/"+m.ID, nil)
	w404 := httptest.NewRecorder()
	server.httpServer.Handler.ServeHTTP(w404, req404)
	if w404.Code != http.StatusNotFound {
		t.Errorf("expected 404 on deleting non-existent media, got %d", w404.Code)
	}
}

