package frp_test

import (
	"os"
	"strings"
	"testing"

	"github.com/anosh26/custardfin/apps/edge/internal/adapters/frp"
)

func validHTTPConfig() frp.Config {
	return frp.Config{
		ServerAddr: "relay.custardfin.local",
		ServerPort: 7000,
		AuthToken:  "test-secret-token",
		NodeID:     "12345678-1234-1234-1234-123456789abc",
		LocalIP:    "127.0.0.1",
		LocalPort:  8088,
		ProxyType:  frp.ProxyTypeHTTP,
		UseTLS:     true,
	}
}

func TestConfig_Validation_HappyPath(t *testing.T) {
	// 1. HTTP Proxy
	cfg := validHTTPConfig()
	if err := cfg.Validate(); err != nil {
		t.Fatalf("expected valid HTTP config, got error: %v", err)
	}

	// 2. TCP Proxy
	cfgTCP := validHTTPConfig()
	cfgTCP.ProxyType = frp.ProxyTypeTCP
	cfgTCP.RemotePort = 8010
	if err := cfgTCP.Validate(); err != nil {
		t.Fatalf("expected valid TCP config, got error: %v", err)
	}
}

func TestConfig_Validation_Errors(t *testing.T) {
	tests := []struct {
		name      string
		modify    func(c *frp.Config)
		errSubstr string
	}{
		{
			name:      "empty server address",
			modify:    func(c *frp.Config) { c.ServerAddr = "" },
			errSubstr: "server address is required",
		},
		{
			name:      "zero server port",
			modify:    func(c *frp.Config) { c.ServerPort = 0 },
			errSubstr: "invalid frp server port",
		},
		{
			name:      "negative server port",
			modify:    func(c *frp.Config) { c.ServerPort = -1 },
			errSubstr: "invalid frp server port",
		},
		{
			name:      "server port too high",
			modify:    func(c *frp.Config) { c.ServerPort = 70000 },
			errSubstr: "invalid frp server port",
		},
		{
			name:      "empty auth token",
			modify:    func(c *frp.Config) { c.AuthToken = "   " },
			errSubstr: "auth token is required",
		},
		{
			name:      "empty node id",
			modify:    func(c *frp.Config) { c.NodeID = "" },
			errSubstr: "node id is required",
		},
		{
			name:      "invalid node id UUID format",
			modify:    func(c *frp.Config) { c.NodeID = "not-a-valid-uuid" },
			errSubstr: "must be a valid UUID",
		},
		{
			name:      "invalid local port",
			modify:    func(c *frp.Config) { c.LocalPort = 0 },
			errSubstr: "invalid local port",
		},
		{
			name:      "unsupported proxy type",
			modify:    func(c *frp.Config) { c.ProxyType = "udp" },
			errSubstr: "unsupported proxy type",
		},
		{
			name: "tcp proxy missing remote port",
			modify: func(c *frp.Config) {
				c.ProxyType = frp.ProxyTypeTCP
				c.RemotePort = 0
			},
			errSubstr: "invalid remote port",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := validHTTPConfig()
			tt.modify(&cfg)
			err := cfg.Validate()
			if err == nil {
				t.Fatalf("expected error containing %q, got nil", tt.errSubstr)
			}
			if !strings.Contains(err.Error(), tt.errSubstr) {
				t.Errorf("expected error %q to contain %q", err.Error(), tt.errSubstr)
			}
		})
	}
}

func TestConfig_GenerateTOML_HTTP(t *testing.T) {
	cfg := validHTTPConfig()
	content, err := cfg.GenerateTOML()
	if err != nil {
		t.Fatalf("failed to generate TOML: %v", err)
	}

	expectedStrings := []string{
		`serverAddr = "relay.custardfin.local"`,
		`serverPort = 7000`,
		`auth.method = "token"`,
		`auth.token = "test-secret-token"`,
		`transport.tls.enable = true`,
		`[[proxies]]`,
		`name = "stream-12345678-1234-1234-1234-123456789abc"`,
		`type = "http"`,
		`localIP = "127.0.0.1"`,
		`localPort = 8088`,
		`customDomains = ["12345678-1234-1234-1234-123456789abc"]`,
	}

	for _, s := range expectedStrings {
		if !strings.Contains(content, s) {
			t.Errorf("generated TOML missing expected line: %s\nContent:\n%s", s, content)
		}
	}
}

func TestConfig_GenerateTOML_TCP(t *testing.T) {
	cfg := validHTTPConfig()
	cfg.ProxyType = frp.ProxyTypeTCP
	cfg.RemotePort = 8010

	content, err := cfg.GenerateTOML()
	if err != nil {
		t.Fatalf("failed to generate TOML: %v", err)
	}

	if !strings.Contains(content, `type = "tcp"`) {
		t.Errorf("expected type = \"tcp\", got:\n%s", content)
	}
	if !strings.Contains(content, `remotePort = 8010`) {
		t.Errorf("expected remotePort = 8010, got:\n%s", content)
	}
}

func TestConfig_WriteConfigFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "frp_test_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	cfg := validHTTPConfig()
	cfg.ConfigDir = tempDir

	filePath, err := cfg.WriteConfigFile()
	if err != nil {
		t.Fatalf("WriteConfigFile failed: %v", err)
	}

	// Verify file exists
	info, err := os.Stat(filePath)
	if err != nil {
		t.Fatalf("stat on generated config failed: %v", err)
	}

	if info.Size() == 0 {
		t.Errorf("expected non-empty config file")
	}

	// Verify content
	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("reading generated config failed: %v", err)
	}

	if !strings.Contains(string(data), cfg.AuthToken) {
		t.Errorf("expected config file to contain auth token")
	}
}
