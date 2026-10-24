package e2e_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

func findProjectRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed getting wd: %v", err)
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.work")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatalf("could not locate project root with go.work")
		}
		dir = parent
	}
}

func ensureBinariesBuilt(t *testing.T, root string) (string, string) {
	t.Helper()
	exeExt := ""
	if runtime.GOOS == "windows" {
		exeExt = ".exe"
	}
	apidPath := filepath.Join(root, "bin", "apid"+exeExt)
	edgedPath := filepath.Join(root, "bin", "edged"+exeExt)

	cmd := exec.Command("go", "build", "-o", apidPath, "./apps/api/cmd/apid")
	cmd.Dir = root
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("failed to build apid: %v\nOutput: %s", err, string(out))
	}

	cmd = exec.Command("go", "build", "-o", edgedPath, "./apps/edge/cmd/edged")
	cmd.Dir = root
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("failed to build edged: %v\nOutput: %s", err, string(out))
	}

	return apidPath, edgedPath
}

func waitForHealth(url string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		resp, err := http.Get(url)
		if err == nil {
			_ = resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				return nil
			}
		}
		time.Sleep(50 * time.Millisecond)
	}
	return fmt.Errorf("timeout waiting for %s", url)
}

func TestCustardFin_E2E_FullLifecycle(t *testing.T) {
	root := findProjectRoot(t)
	apidBin, edgedBin := ensureBinariesBuilt(t, root)

	tempDir, err := os.MkdirTemp("", "custardfin_e2e_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	edgeMediaDir := filepath.Join(tempDir, "media")
	edgeTranscodeDir := filepath.Join(tempDir, "transcoded")
	edgeDBPath := filepath.Join(tempDir, "edge.db")
	_ = os.MkdirAll(edgeMediaDir, 0755)
	_ = os.MkdirAll(edgeTranscodeDir, 0755)

	// Pre-generate 4s test video
	videoPath := filepath.Join(edgeMediaDir, "sample.mp4")
	ffmpegCmd := exec.Command("ffmpeg",
		"-y",
		"-f", "lavfi",
		"-i", "testsrc=duration=4:size=640x360:rate=30",
		"-c:v", "libx264",
		"-pix_fmt", "yuv420p",
		videoPath,
	)
	if out, err := ffmpegCmd.CombinedOutput(); err != nil {
		t.Fatalf("failed generating test video with ffmpeg: %v\nOutput: %s", err, string(out))
	}
	t.Logf("[E2E] Generated test video: %s", videoPath)

	// Ports for testing
	apiPort := "18080"
	edgePort := "18090"
	apiBaseURL := "http://127.0.0.1:" + apiPort
	edgeBaseURL := "http://127.0.0.1:" + edgePort

	// -------------------------------------------------------------
	// 1. Launch Cloud API Daemon (apid)
	// -------------------------------------------------------------
	apiCtx, apiCancel := context.WithCancel(context.Background())
	defer apiCancel()

	apidCmd := exec.CommandContext(apiCtx, apidBin)
	apidCmd.Env = append(os.Environ(),
		"API_ADDR=127.0.0.1:"+apiPort,
		"API_BASE_URL="+apiBaseURL,
		"TUNNEL_BASE_URL="+edgeBaseURL,
		"SIGNING_SECRET=e2e_testing_chunk_signing_secret_key_32b!!",
		"DATABASE_URL=memory",
		"REDIS_ADDR=memory",
		"JWT_SECRET=e2e_testing_secret_key_at_least_32_bytes_long!!",
	)
	apidCmd.Stdout = os.Stdout
	apidCmd.Stderr = os.Stderr

	if err := apidCmd.Start(); err != nil {
		t.Fatalf("failed starting apid daemon: %v", err)
	}
	defer func() {
		apiCancel()
		_ = apidCmd.Wait()
	}()

	if err := waitForHealth(apiBaseURL+"/healthz", 10*time.Second); err != nil {
		t.Fatalf("apid failed to become healthy: %v", err)
	}
	t.Logf("[E2E] Cloud API is healthy on %s", apiBaseURL)

	// -------------------------------------------------------------
	// 2. Register User & Edge Node on Cloud API
	// -------------------------------------------------------------
	// Register user
	regBody := `{"username":"alice_streamer","password":"SuperSecret123!"}`
	regResp, err := http.Post(apiBaseURL+"/v1/auth/register", "application/json", strings.NewReader(regBody))
	if err != nil {
		t.Fatalf("user registration failed: %v", err)
	}
	defer regResp.Body.Close()

	if regResp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201 Created from user register, got %d", regResp.StatusCode)
	}

	var authRes struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(regResp.Body).Decode(&authRes); err != nil {
		t.Fatalf("failed decoding auth response: %v", err)
	}

	// Register node with JWT bearer auth
	nodeReq, err := http.NewRequest(http.MethodPost, apiBaseURL+"/v1/nodes", strings.NewReader(`{"name":"alice-rtx-edge"}`))
	if err != nil {
		t.Fatalf("failed creating node register request: %v", err)
	}
	nodeReq.Header.Set("Authorization", "Bearer "+authRes.AccessToken)
	nodeReq.Header.Set("Content-Type", "application/json")

	nodeResp, err := http.DefaultClient.Do(nodeReq)
	if err != nil {
		t.Fatalf("node registration failed: %v", err)
	}
	defer nodeResp.Body.Close()

	if nodeResp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201 Created from node register, got %d", nodeResp.StatusCode)
	}

	var nodeRes struct {
		NodeID       string `json:"node_id"`
		NodeToken    string `json:"node_token"`
		FRPAuthToken string `json:"frp_auth_token"`
	}
	if err := json.NewDecoder(nodeResp.Body).Decode(&nodeRes); err != nil {
		t.Fatalf("failed decoding node registration response: %v", err)
	}
	t.Logf("[E2E] Node successfully registered: ID=%s", nodeRes.NodeID)

	// -------------------------------------------------------------
	// 3. Launch Edge Daemon (edged) Wired to Cloud API
	// -------------------------------------------------------------
	edgeCtx, edgeCancel := context.WithCancel(context.Background())
	defer edgeCancel()

	edgedCmd := exec.CommandContext(edgeCtx, edgedBin)
	edgedCmd.Env = append(os.Environ(),
		"EDGE_ADDR=127.0.0.1:"+edgePort,
		"EDGE_MEDIA_DIR="+edgeMediaDir,
		"EDGE_TRANSCODE_DIR="+edgeTranscodeDir,
		"EDGE_DB_PATH="+edgeDBPath,
		"EDGE_CLOUD_URL="+apiBaseURL,
		"EDGE_NODE_ID="+nodeRes.NodeID,
		"EDGE_NODE_TOKEN="+nodeRes.NodeToken,
		"EDGE_FRP_ENABLED=false",
	)
	edgedCmd.Stdout = os.Stdout
	edgedCmd.Stderr = os.Stderr

	if err := edgedCmd.Start(); err != nil {
		t.Fatalf("failed starting edged daemon: %v", err)
	}
	defer func() {
		edgeCancel()
		_ = edgedCmd.Wait()
	}()

	if err := waitForHealth(edgeBaseURL+"/healthz", 5*time.Second); err != nil {
		t.Fatalf("edged failed to become healthy: %v", err)
	}
	t.Logf("[E2E] Edge Daemon is healthy on %s", edgeBaseURL)

	// -------------------------------------------------------------
	// 4. Wait for Edge Library Scan & Discover Media
	// -------------------------------------------------------------
	var mediaID string
	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		resp, err := http.Get(edgeBaseURL + "/api/media")
		if err == nil {
			var listRes []struct {
				ID       string `json:"id"`
				Path     string `json:"path"`
				IsPublic bool   `json:"is_public"`
			}
			if err := json.NewDecoder(resp.Body).Decode(&listRes); err == nil && len(listRes) > 0 {
				_ = resp.Body.Close()
				mediaID = listRes[0].ID
				t.Logf("[E2E] Edge scanner discovered media: ID=%s", mediaID)
				break
			}
			_ = resp.Body.Close()
		}
		time.Sleep(200 * time.Millisecond)
	}

	if mediaID == "" {
		t.Fatalf("edge scanner did not discover sample.mp4 within 10s")
	}

	// -------------------------------------------------------------
	// 5. Pre-create 4s HLS Segments in transcodeDir for Media
	// -------------------------------------------------------------
	mediaTranscodeDir := filepath.Join(edgeTranscodeDir, mediaID, "1080p")
	_ = os.MkdirAll(mediaTranscodeDir, 0755)

	// Transcode sample.mp4 into strict 4-second HLS segments
	hlsCmd := exec.Command("ffmpeg",
		"-y",
		"-i", videoPath,
		"-c:v", "libx264",
		"-hls_time", "4",
		"-hls_list_size", "0",
		"-hls_segment_filename", filepath.Join(mediaTranscodeDir, "segment_%03d.ts"),
		filepath.Join(mediaTranscodeDir, "index.m3u8"),
	)
	if out, err := hlsCmd.CombinedOutput(); err != nil {
		t.Fatalf("failed generating HLS segments: %v\nOutput: %s", err, string(out))
	}

	// Write Master Playlist
	masterContent := fmt.Sprintf("#EXTM3U\n#EXT-X-VERSION:3\n#EXT-X-STREAM-INF:BANDWIDTH=4500000,RESOLUTION=1920x1080\n1080p/index.m3u8\n")
	if err := os.WriteFile(filepath.Join(edgeTranscodeDir, mediaID, "master.m3u8"), []byte(masterContent), 0644); err != nil {
		t.Fatalf("failed writing master playlist: %v", err)
	}
	t.Logf("[E2E] Generated 4s HLS segments & master playlist for media %s", mediaID)

	// -------------------------------------------------------------
	// 6. Verify Public Catalog Initially Empty
	// -------------------------------------------------------------
	catResp, err := http.Get(apiBaseURL + "/v1/catalog")
	if err != nil {
		t.Fatalf("catalog query failed: %v", err)
	}
	defer catResp.Body.Close()

	var catRes struct {
		Total int `json:"total"`
	}
	_ = json.NewDecoder(catResp.Body).Decode(&catRes)
	if catRes.Total != 0 {
		t.Fatalf("expected initial catalog total=0, got %d", catRes.Total)
	}

	// -------------------------------------------------------------
	// 7. Toggle Media to Public on Edge Daemon
	// -------------------------------------------------------------
	toggleReq, err := http.NewRequest(
		http.MethodPut,
		fmt.Sprintf("%s/api/media/%s/public", edgeBaseURL, mediaID),
		strings.NewReader(`{"is_public": true}`),
	)
	if err != nil {
		t.Fatalf("failed creating toggle request: %v", err)
	}
	toggleReq.Header.Set("Content-Type", "application/json")

	toggleResp, err := http.DefaultClient.Do(toggleReq)
	if err != nil {
		t.Fatalf("toggle public request failed: %v", err)
	}
	defer toggleResp.Body.Close()

	if toggleResp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(toggleResp.Body)
		t.Fatalf("expected 200 OK from toggle public, got %d: %s", toggleResp.StatusCode, string(body))
	}
	t.Logf("[E2E] Successfully marked media %s as Public on Edge", mediaID)

	// -------------------------------------------------------------
	// 8. Verify Manifest in Cloud API Catalog
	// -------------------------------------------------------------
	catResp, err = http.Get(apiBaseURL + "/v1/catalog")
	if err != nil {
		t.Fatalf("catalog query failed: %v", err)
	}
	defer catResp.Body.Close()

	var pubCatRes struct {
		Items []struct {
			MediaID      string `json:"media_id"`
			Title        string `json:"title"`
			NodeID       string `json:"node_id"`
			IsNodeOnline bool   `json:"is_node_online"`
		} `json:"items"`
		Total int `json:"total"`
	}
	if err := json.NewDecoder(catResp.Body).Decode(&pubCatRes); err != nil {
		t.Fatalf("failed decoding catalog response: %v", err)
	}

	if pubCatRes.Total != 1 || len(pubCatRes.Items) != 1 {
		t.Fatalf("expected 1 item in public catalog, got total=%d", pubCatRes.Total)
	}

	item := pubCatRes.Items[0]
	if item.MediaID != mediaID {
		t.Errorf("expected media ID %s, got %s", mediaID, item.MediaID)
	}
	if item.NodeID != nodeRes.NodeID {
		t.Errorf("expected node ID %s, got %s", nodeRes.NodeID, item.NodeID)
	}
	if !item.IsNodeOnline {
		t.Errorf("expected node %s to be reported online in catalog via heartbeat", nodeRes.NodeID)
	}
	t.Logf("[E2E] Verified item in Cloud Catalog: MediaID=%s, NodeID=%s, IsOnline=%v", item.MediaID, item.NodeID, item.IsNodeOnline)

	// -------------------------------------------------------------
	// 9. Resolve Playback Route from Cloud API
	// -------------------------------------------------------------
	// 9a. Remote playback resolution -> origin_shield with signed URL
	pbResp, err := http.Get(fmt.Sprintf("%s/v1/media/%s/playback", apiBaseURL, mediaID))
	if err != nil {
		t.Fatalf("playback route resolution failed: %v", err)
	}
	defer pbResp.Body.Close()

	if pbResp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK from playback route, got %d", pbResp.StatusCode)
	}

	var pbRes struct {
		StreamURL string `json:"stream_url"`
		RouteType string `json:"route_type"`
		ExpiresAt string `json:"expires_at"`
	}
	if err := json.NewDecoder(pbResp.Body).Decode(&pbRes); err != nil {
		t.Fatalf("failed decoding playback response: %v", err)
	}

	if pbRes.RouteType != "origin_shield" {
		t.Errorf("expected origin_shield route type, got %s", pbRes.RouteType)
	}
	if !strings.Contains(pbRes.StreamURL, "sig=") {
		t.Errorf("expected signed stream URL, got %s", pbRes.StreamURL)
	}
	t.Logf("[E2E] Resolved playback: URL=%s (Route: %s)", pbRes.StreamURL, pbRes.RouteType)

	// 9b. Direct LAN playback resolution (?lan=true) -> direct_lan
	lanResp, err := http.Get(fmt.Sprintf("%s/v1/media/%s/playback?lan=true", apiBaseURL, mediaID))
	if err != nil {
		t.Fatalf("lan playback route resolution failed: %v", err)
	}
	defer lanResp.Body.Close()

	var lanRes struct {
		StreamURL string `json:"stream_url"`
		RouteType string `json:"route_type"`
	}
	if err := json.NewDecoder(lanResp.Body).Decode(&lanRes); err != nil {
		t.Fatalf("failed decoding lan playback response: %v", err)
	}
	if lanRes.RouteType != "direct_lan" {
		t.Errorf("expected direct_lan route type, got %s", lanRes.RouteType)
	}

	// 9c. Stream via Origin Shield Pull-Through (Cache MISS then Cache HIT)
	shieldResp1, err := http.Get(pbRes.StreamURL)
	if err != nil {
		t.Fatalf("failed fetching stream via origin shield: %v", err)
	}
	defer shieldResp1.Body.Close()
	if shieldResp1.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK from origin shield master playlist, got %d", shieldResp1.StatusCode)
	}
	if xCache := shieldResp1.Header.Get("X-Cache"); xCache != "MISS" {
		t.Errorf("expected X-Cache: MISS on initial pull, got %s", xCache)
	}
	shieldBytes1, _ := io.ReadAll(shieldResp1.Body)
	if !bytes.Contains(shieldBytes1, []byte("sig=")) {
		t.Errorf("expected rewritten playlist with signed URLs from origin shield:\n%s", string(shieldBytes1))
	}

	// Second request should HIT cache
	shieldResp2, err := http.Get(pbRes.StreamURL)
	if err != nil {
		t.Fatalf("failed fetching cached stream via origin shield: %v", err)
	}
	defer shieldResp2.Body.Close()
	if xCache := shieldResp2.Header.Get("X-Cache"); xCache != "HIT" {
		t.Errorf("expected X-Cache: HIT on second pull, got %s", xCache)
	}
	t.Logf("[E2E] Verified Origin Shield pull-through cache: initial MISS, subsequent HIT")

	// -------------------------------------------------------------
	// 10. End-to-End Stream Retrieval from Edge Server
	// -------------------------------------------------------------
	// A. Master Playlist
	masterResp, err := http.Get(fmt.Sprintf("%s/stream/%s/master.m3u8", edgeBaseURL, mediaID))
	if err != nil {
		t.Fatalf("failed fetching master playlist: %v", err)
	}
	defer masterResp.Body.Close()

	if masterResp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK from master playlist, got %d", masterResp.StatusCode)
	}

	masterBytes, _ := io.ReadAll(masterResp.Body)
	if !bytes.Contains(masterBytes, []byte("#EXTM3U")) {
		t.Fatalf("invalid master playlist content:\n%s", string(masterBytes))
	}

	// B. Rendition Playlist
	renditionResp, err := http.Get(fmt.Sprintf("%s/stream/%s/1080p/index.m3u8", edgeBaseURL, mediaID))
	if err != nil {
		t.Fatalf("failed fetching rendition playlist: %v", err)
	}
	defer renditionResp.Body.Close()

	if renditionResp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK from rendition playlist, got %d", renditionResp.StatusCode)
	}

	renditionBytes, _ := io.ReadAll(renditionResp.Body)
	if !bytes.Contains(renditionBytes, []byte("#EXT-X-TARGETDURATION:4")) {
		t.Fatalf("expected 4s chunk duration invariant in playlist:\n%s", string(renditionBytes))
	}

	// C. First Video Chunk
	chunkResp, err := http.Get(fmt.Sprintf("%s/stream/%s/1080p/segment_000.ts", edgeBaseURL, mediaID))
	if err != nil {
		t.Fatalf("failed fetching video segment: %v", err)
	}
	defer chunkResp.Body.Close()

	if chunkResp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK from video segment, got %d", chunkResp.StatusCode)
	}

	chunkData, _ := io.ReadAll(chunkResp.Body)
	if len(chunkData) == 0 {
		t.Fatalf("expected non-empty video chunk")
	}
	t.Logf("[E2E] Streamed video segment: %d bytes (MPEG-TS)", len(chunkData))

	// D. Byte Range Request (Accept-Ranges: bytes)
	rangeReq, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/stream/%s/1080p/segment_000.ts", edgeBaseURL, mediaID), nil)
	rangeReq.Header.Set("Range", "bytes=0-187")
	rangeResp, err := http.DefaultClient.Do(rangeReq)
	if err != nil {
		t.Fatalf("byte range request failed: %v", err)
	}
	defer rangeResp.Body.Close()

	if rangeResp.StatusCode != http.StatusPartialContent {
		t.Fatalf("expected 206 Partial Content, got %d", rangeResp.StatusCode)
	}
	rangeData, _ := io.ReadAll(rangeResp.Body)
	if len(rangeData) != 188 {
		t.Fatalf("expected 188 bytes, got %d", len(rangeData))
	}

	// -------------------------------------------------------------
	// 11. Toggle Public to False (Unpublish)
	// -------------------------------------------------------------
	unpubReq, _ := http.NewRequest(
		http.MethodPut,
		fmt.Sprintf("%s/api/media/%s/public", edgeBaseURL, mediaID),
		strings.NewReader(`{"is_public": false}`),
	)
	unpubReq.Header.Set("Content-Type", "application/json")

	unpubResp, err := http.DefaultClient.Do(unpubReq)
	if err != nil {
		t.Fatalf("unpublish request failed: %v", err)
	}
	defer unpubResp.Body.Close()

	if unpubResp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK from unpublish, got %d", unpubResp.StatusCode)
	}

	// Verify Catalog is now empty
	catResp, err = http.Get(apiBaseURL + "/v1/catalog")
	if err != nil {
		t.Fatalf("catalog query failed: %v", err)
	}
	defer catResp.Body.Close()

	var emptyCatRes struct {
		Total int `json:"total"`
	}
	_ = json.NewDecoder(catResp.Body).Decode(&emptyCatRes)
	if emptyCatRes.Total != 0 {
		t.Fatalf("expected 0 items in catalog after unpublishing, got %d", emptyCatRes.Total)
	}

	// Invariant Check: Local private media remains fully available on Edge!
	edgeMediaResp, err := http.Get(fmt.Sprintf("%s/api/media/%s", edgeBaseURL, mediaID))
	if err != nil {
		t.Fatalf("local edge media check failed: %v", err)
	}
	defer edgeMediaResp.Body.Close()

	if edgeMediaResp.StatusCode != http.StatusOK {
		t.Fatalf("expected local media to remain available on edge, got %d", edgeMediaResp.StatusCode)
	}

	t.Log("[E2E] SUCCESS: Full cross-subsystem end-to-end integration verified!")
}
