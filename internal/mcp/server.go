// Package mcp provides the MCP server implementation for ABAP ADT tools.
package mcp

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/oisee/vibing-steampunk/pkg/adt"
)

const (
	// DefaultStreamableHTTPAddr is the default listen address for streamable HTTP transport.
	DefaultStreamableHTTPAddr = "127.0.0.1:8080"
	// DefaultStreamableHTTPPath is the MCP endpoint path for streamable HTTP transport.
	DefaultStreamableHTTPPath = "/mcp"
)

type streamableHTTPStarter interface {
	http.Handler
	Start(addr string) error
}

var serveStdioFunc = func(mcpServer *server.MCPServer) error {
	return server.ServeStdio(mcpServer)
}

var newStreamableHTTPServerFunc = func(mcpServer *server.MCPServer, opts ...server.StreamableHTTPOption) streamableHTTPStarter {
	return server.NewStreamableHTTPServer(mcpServer, opts...)
}

var listenAndServeFunc = func(addr string, handler http.Handler) error {
	return (&http.Server{Addr: addr, Handler: handler}).ListenAndServe()
}

// AsyncTask represents a background task status.
type AsyncTask struct {
	ID        string      `json:"id"`
	Type      string      `json:"type"`       // "report", "export", etc.
	Status    string      `json:"status"`     // "running", "completed", "error"
	StartedAt time.Time   `json:"started_at"`
	EndedAt   *time.Time  `json:"ended_at,omitempty"`
	Result    interface{} `json:"result,omitempty"`
	Error     string      `json:"error,omitempty"`
}

// Server wraps the MCP server with ADT client.
type Server struct {
	mcpServer      *server.MCPServer
	adtClient      *adt.Client
	amdpWSClient   *adt.AMDPWebSocketClient   // WebSocket-based AMDP client (ZADT_VSP)
	debugWSClient  *adt.DebugWebSocketClient  // WebSocket-based debug client (ZADT_VSP)
	config         *Config                    // Server configuration for session manager creation
	featureProber  *adt.FeatureProber         // Feature detection system (safety network)
	featureConfig  adt.FeatureConfig          // Feature configuration

	// Async task management
	asyncTasks   map[string]*AsyncTask
	asyncTasksMu sync.RWMutex
	asyncTaskID  int64
}

// Config holds MCP server configuration.
type Config struct {
	// SAP connection settings
	BaseURL            string
	Username           string
	Password           string
	Client             string
	Language           string
	InsecureSkipVerify bool

	// Cookie authentication (alternative to basic auth)
	Cookies map[string]string

	// Verbose output
	Verbose bool

	// Mode: focused or expert (default: focused)
	Mode string

	// DisabledGroups disables groups of tools using short codes:
	// 5/U = UI5/BSP tools, T = Test tools, H = HANA/AMDP debugger, D = ABAP Debugger
	// Example: "TH" disables Tests and HANA debugger tools
	DisabledGroups string

	// Transport mode: stdio (default) or http-streamable
	Transport string

	// HTTPAddr is the listen address for the http-streamable transport (default: 127.0.0.1:8080).
	// Set to 0.0.0.0:8080 when running in Docker or to accept remote connections.
	HTTPAddr string

	// Safety configuration
	ReadOnly         bool
	BlockFreeSQL     bool
	AllowedOps       string
	DisallowedOps    string
	AllowedPackages  []string
	EnableTransports        bool     // Explicitly enable transport management (default: disabled)
	TransportReadOnly       bool     // Only allow read operations on transports (list, get)
	AllowedTransports       []string // Whitelist specific transports (supports wildcards like "A4HK*")
	AllowTransportableEdits bool     // Allow editing objects that require transport requests

	// Feature configuration (safety network)
	// Values: "auto" (default, probe system), "on" (force enabled), "off" (force disabled)
	FeatureHANA      string // HANA database detection (required for some AMDP features)
	FeatureAbapGit   string // abapGit integration
	FeatureRAP       string // RAP/OData development (DDLS, BDEF, SRVD, SRVB)
	FeatureAMDP      string // AMDP/HANA debugger
	FeatureUI5       string // UI5/Fiori BSP management
	FeatureTransport string // CTS transport management (distinct from EnableTransports safety)

	// Debugger configuration
	TerminalID string // SAP GUI terminal ID for cross-tool breakpoint sharing

	// Granular tool visibility (from .vsp.json)
	// Key: tool name, Value: true=enabled, false=disabled
	// Takes highest priority over mode and disabled groups
	ToolsConfig map[string]bool
}

// NewServer creates a new MCP server for ABAP ADT tools.
func NewServer(cfg *Config) *Server {
	// Create ADT client
	opts := []adt.Option{
		adt.WithClient(cfg.Client),
		adt.WithLanguage(cfg.Language),
	}
	if cfg.InsecureSkipVerify {
		opts = append(opts, adt.WithInsecureSkipVerify())
	}
	if len(cfg.Cookies) > 0 {
		opts = append(opts, adt.WithCookies(cfg.Cookies))
	}
	if cfg.Verbose {
		opts = append(opts, adt.WithVerbose())
	}

	// Configure safety settings
	safety := adt.UnrestrictedSafetyConfig() // Default: unrestricted for backwards compatibility
	if cfg.ReadOnly || cfg.Mode == "readonly" {
		// readonly mode implies ReadOnly safety (belt-and-suspenders)
		safety.ReadOnly = true
	}
	if cfg.BlockFreeSQL {
		safety.BlockFreeSQL = true
	}
	if cfg.AllowedOps != "" {
		safety.AllowedOps = cfg.AllowedOps
	}
	if cfg.DisallowedOps != "" {
		safety.DisallowedOps = cfg.DisallowedOps
	}
	if len(cfg.AllowedPackages) > 0 {
		safety.AllowedPackages = cfg.AllowedPackages
	}
	if cfg.EnableTransports {
		safety.EnableTransports = true
	}
	if cfg.TransportReadOnly {
		safety.TransportReadOnly = true
	}
	if len(cfg.AllowedTransports) > 0 {
		safety.AllowedTransports = cfg.AllowedTransports
	}
	if cfg.AllowTransportableEdits {
		safety.AllowTransportableEdits = true
	}
	opts = append(opts, adt.WithSafety(safety))

	adtClient := adt.NewClient(cfg.BaseURL, cfg.Username, cfg.Password, opts...)

	// Set terminal ID for debugger operations
	// Priority: 1) Custom ID (SAP GUI), 2) User-based ID
	if cfg.TerminalID != "" {
		adt.SetTerminalID(cfg.TerminalID)
	}
	adt.SetTerminalIDUser(cfg.Username)

	// Configure feature detection (safety network)
	featureConfig := adt.FeatureConfig{
		HANA:      parseFeatureMode(cfg.FeatureHANA),
		AbapGit:   parseFeatureMode(cfg.FeatureAbapGit),
		RAP:       parseFeatureMode(cfg.FeatureRAP),
		AMDP:      parseFeatureMode(cfg.FeatureAMDP),
		UI5:       parseFeatureMode(cfg.FeatureUI5),
		Transport: parseFeatureMode(cfg.FeatureTransport),
	}

	// Create feature prober
	featureProber := adt.NewFeatureProber(adtClient, featureConfig, cfg.Verbose)

	// Create MCP server
	mcpServer := server.NewMCPServer(
		"mcp-abap-adt-go",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
	)

	s := &Server{
		mcpServer:     mcpServer,
		adtClient:     adtClient,
		config:        cfg,
		featureProber: featureProber,
		featureConfig: featureConfig,
		asyncTasks:    make(map[string]*AsyncTask),
	}

	// Register tools based on mode, disabled groups, and granular tool config
	s.registerTools(cfg.Mode, cfg.DisabledGroups, cfg.ToolsConfig)

	return s
}

// parseFeatureMode converts string to FeatureMode
func parseFeatureMode(s string) adt.FeatureMode {
	switch strings.ToLower(s) {
	case "on", "true", "1", "yes", "enabled":
		return adt.FeatureModeOn
	case "off", "false", "0", "no", "disabled":
		return adt.FeatureModeOff
	default:
		return adt.FeatureModeAuto
	}
}

// Serve starts the MCP server using the selected transport.
// For http-streamable, the listen address is taken from Config.HTTPAddr (falling back to DefaultStreamableHTTPAddr).
func (s *Server) Serve(transport string) error {
	switch strings.ToLower(strings.TrimSpace(transport)) {
	case "", "stdio":
		return s.ServeStdio()
	case "http-streamable":
		addr := s.config.HTTPAddr
		if strings.TrimSpace(addr) == "" {
			addr = DefaultStreamableHTTPAddr
		}
		return s.ServeStreamableHTTP(addr)
	default:
		return fmt.Errorf("unsupported transport: %s (must be 'stdio' or 'http-streamable')", transport)
	}
}

// ServeStdio starts the MCP server on stdin/stdout.
func (s *Server) ServeStdio() error {
	return serveStdioFunc(s.mcpServer)
}

// ServeStreamableHTTP starts the MCP server using streamable HTTP transport.
// It validates the Origin header on all incoming connections to prevent DNS rebinding attacks,
// as required by the MCP specification:
// https://modelcontextprotocol.io/specification/2025-03-26/basic/transports
func (s *Server) ServeStreamableHTTP(addr string) error {
	if strings.TrimSpace(addr) == "" {
		addr = DefaultStreamableHTTPAddr
	}

	mcpHandler := newStreamableHTTPServerFunc(
		s.mcpServer,
		server.WithEndpointPath(DefaultStreamableHTTPPath),
	)

	mux := http.NewServeMux()
	mux.Handle(DefaultStreamableHTTPPath, originValidationMiddleware(addr, mcpHandler))
	return listenAndServeFunc(addr, mux)
}

// originValidationMiddleware returns an HTTP handler that validates the Origin header
// on incoming requests to prevent DNS rebinding attacks per the MCP specification.
// Requests without an Origin header are allowed (same-origin browser requests omit it).
// Requests with an Origin whose host does not match the server's host are rejected with 403.
func originValidationMiddleware(serverAddr string, next http.Handler) http.Handler {
	serverHost, _, err := net.SplitHostPort(serverAddr)
	if err != nil {
		serverHost = serverAddr
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			u, parseErr := url.Parse(origin)
			if parseErr != nil || !isSameOriginHost(u.Hostname(), serverHost) {
				http.Error(w, "Forbidden: invalid Origin header", http.StatusForbidden)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

// isSameOriginHost reports whether originHost is the same logical host as serverHost,
// treating 127.0.0.1, ::1 and localhost as equivalent.
func isSameOriginHost(originHost, serverHost string) bool {
	normalize := func(h string) string {
		if h == "127.0.0.1" || h == "::1" {
			return "localhost"
		}
		return h
	}
	return normalize(originHost) == normalize(serverHost)
}

// newToolResultError creates an error result for tool execution failures.
func newToolResultError(message string) *mcp.CallToolResult {
	result := mcp.NewToolResultText(message)
	result.IsError = true
	return result
}

// ensureWSConnected ensures the WebSocket client is connected, creating it if needed.
// Returns error result if connection fails, nil on success.
func (s *Server) ensureWSConnected(ctx context.Context, toolName string) *mcp.CallToolResult {
	if s.amdpWSClient == nil || !s.amdpWSClient.IsConnected() {
		s.amdpWSClient = adt.NewAMDPWebSocketClient(
			s.config.BaseURL, s.config.Client, s.config.Username, s.config.Password, s.config.InsecureSkipVerify,
		)
		if err := s.amdpWSClient.Connect(ctx); err != nil {
			s.amdpWSClient = nil
			return newToolResultError(fmt.Sprintf("%s: WebSocket connect failed: %v", toolName, err))
		}
	}
	return nil
}

// requireActiveAMDPSession checks if there's an active AMDP debug session.
// Returns error result if no session, nil if session is active.
func (s *Server) requireActiveAMDPSession() *mcp.CallToolResult {
	if s.amdpWSClient == nil || !s.amdpWSClient.IsActive() {
		return newToolResultError("No active AMDP session. Use AMDPDebuggerStart first.")
	}
	return nil
}

// Tool handlers are in separate files:
// - handlers_read.go: GetProgram, GetClass, GetTable, etc.
// - handlers_system.go: GetSystemInfo, GetFeatures, etc.
// - handlers_analysis.go: GetCallGraph, TraceExecution, etc.
// - handlers_codeintel.go: FindDefinition, FindReferences, CodeCompletion, etc.
// - handlers_devtools.go: SyntaxCheck, Activate, ATC, etc.
// - handlers_crud.go: Lock, Create, Update, Delete, etc.
// - handlers_debugger.go: SetBreakpoint, DebuggerListen, etc.
// - handlers_amdp.go: AMDPDebugger* handlers
// - handlers_ui5.go: UI5ListApps, UI5GetApp, etc.
// - handlers_git.go: GitTypes, GitExport
// - handlers_report.go: RunReport, GetVariants, etc.
// - handlers_install.go: InstallZADTVSP, InstallAbapGit, etc.
// - handlers_transport.go: ListTransports, GetTransport, etc.
//
// Tool registration is in:
// - tools_register.go: registerTools() and all register*Tools() methods
// - tools_groups.go: toolGroups() - group definitions for --disabled-groups
// - tools_focused.go: focusedToolSet() - focused mode whitelist
// - tools_aliases.go: registerToolAliases() - short alias names
