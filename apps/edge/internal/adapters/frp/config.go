package frp

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var uuidRegex = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

// ProxyType defines the proxy protocol for FRP tunnel.
type ProxyType string

const (
	ProxyTypeHTTP ProxyType = "http"
	ProxyTypeTCP  ProxyType = "tcp"
)

// Config defines the configuration parameters for the FRP client tunnel.
type Config struct {
	ServerAddr string    // Central relay server address/host
	ServerPort int       // Central relay server port (default: 7000)
	AuthToken  string    // FRP secret authentication token
	NodeID     string    // Edge node UUID (used for proxy name and routing)
	LocalIP    string    // Local IP to forward to (default: "127.0.0.1")
	LocalPort  int       // Local streaming HTTP port (default: 8088)
	ProxyType  ProxyType // "http" (default) or "tcp"
	Subdomain  string    // Subdomain for HTTP vhost routing (defaults to NodeID)
	RemotePort int       // Remote port for TCP proxy (required when ProxyType is "tcp")
	BinaryPath string    // Path to frpc executable (default: "frpc")
	ConfigDir  string    // Directory where temp frpc.toml is generated (default: os.TempDir())
	UseTLS     bool      // Use TLS for outbound tunnel transport (default: true)
}

// Validate ensures all configuration fields meet requirements.
func (c *Config) Validate() error {
	if strings.TrimSpace(c.ServerAddr) == "" {
		return errors.New("frp server address is required")
	}
	if c.ServerPort <= 0 || c.ServerPort > 65535 {
		return fmt.Errorf("invalid frp server port %d: must be 1-65535", c.ServerPort)
	}
	if strings.TrimSpace(c.AuthToken) == "" {
		return errors.New("frp auth token is required")
	}
	if strings.TrimSpace(c.NodeID) == "" {
		return errors.New("frp node id is required")
	}
	if !uuidRegex.MatchString(c.NodeID) {
		return fmt.Errorf("node id %q must be a valid UUID", c.NodeID)
	}
	if c.LocalPort <= 0 || c.LocalPort > 65535 {
		return fmt.Errorf("invalid local port %d: must be 1-65535", c.LocalPort)
	}

	if c.ProxyType == "" {
		c.ProxyType = ProxyTypeHTTP
	}
	if c.ProxyType != ProxyTypeHTTP && c.ProxyType != ProxyTypeTCP {
		return fmt.Errorf("unsupported proxy type %q: must be http or tcp", c.ProxyType)
	}

	if c.ProxyType == ProxyTypeTCP && (c.RemotePort <= 0 || c.RemotePort > 65535) {
		return fmt.Errorf("invalid remote port %d for tcp proxy: must be 1-65535", c.RemotePort)
	}

	if strings.TrimSpace(c.LocalIP) == "" {
		c.LocalIP = "127.0.0.1"
	}
	if strings.TrimSpace(c.Subdomain) == "" {
		c.Subdomain = c.NodeID
	}
	if strings.TrimSpace(c.BinaryPath) == "" {
		c.BinaryPath = "frpc"
	}

	return nil
}

// GenerateTOML produces the official FRP v0.52+ TOML client configuration.
func (c *Config) GenerateTOML() (string, error) {
	if err := c.Validate(); err != nil {
		return "", fmt.Errorf("config validation failed: %w", err)
	}

	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("serverAddr = %q\n", c.ServerAddr))
	sb.WriteString(fmt.Sprintf("serverPort = %d\n", c.ServerPort))
	sb.WriteString("auth.method = \"token\"\n")
	sb.WriteString(fmt.Sprintf("auth.token = %q\n", c.AuthToken))
	sb.WriteString(fmt.Sprintf("transport.tls.enable = %t\n\n", c.UseTLS))

	sb.WriteString("[[proxies]]\n")
	sb.WriteString(fmt.Sprintf("name = \"stream-%s\"\n", c.NodeID))
	sb.WriteString(fmt.Sprintf("type = %q\n", c.ProxyType))
	sb.WriteString(fmt.Sprintf("localIP = %q\n", c.LocalIP))
	sb.WriteString(fmt.Sprintf("localPort = %d\n", c.LocalPort))

	if c.ProxyType == ProxyTypeHTTP {
		sb.WriteString(fmt.Sprintf("customDomains = [%q]\n", c.Subdomain))
	} else if c.ProxyType == ProxyTypeTCP {
		sb.WriteString(fmt.Sprintf("remotePort = %d\n", c.RemotePort))
	}

	return sb.String(), nil
}

// WriteConfigFile creates a temporary TOML file with restricted 0600 permissions.
func (c *Config) WriteConfigFile() (string, error) {
	content, err := c.GenerateTOML()
	if err != nil {
		return "", err
	}

	dir := c.ConfigDir
	if dir == "" {
		dir = os.TempDir()
	}

	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", fmt.Errorf("failed to create frp config dir %s: %w", dir, err)
	}

	tmpFile, err := os.CreateTemp(dir, fmt.Sprintf("frpc_%s_*.toml", c.NodeID[:8]))
	if err != nil {
		return "", fmt.Errorf("failed to create temporary frpc config: %w", err)
	}
	defer tmpFile.Close()

	// Ensure strict 0600 file permissions to protect auth token
	if err := os.Chmod(tmpFile.Name(), 0600); err != nil {
		_ = os.Remove(tmpFile.Name())
		return "", fmt.Errorf("failed to set secure permissions on frpc config: %w", err)
	}

	if _, err := tmpFile.WriteString(content); err != nil {
		_ = os.Remove(tmpFile.Name())
		return "", fmt.Errorf("failed to write frpc config content: %w", err)
	}

	return filepath.Clean(tmpFile.Name()), nil
}
