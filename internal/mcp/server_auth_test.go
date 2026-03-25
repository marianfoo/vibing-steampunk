package mcp

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	mcpserver "github.com/mark3labs/mcp-go/server"
)

// --- API Key Middleware Tests ---

func TestAPIKeyMiddleware_ValidKey(t *testing.T) {
	apiKey := "test-secret-key-12345"
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	handler := apiKeyMiddleware(apiKey, next)

	req := httptest.NewRequest(http.MethodPost, "/mcp", nil)
	req.Header.Set("Authorization", "Bearer "+apiKey)
	rw := httptest.NewRecorder()
	handler.ServeHTTP(rw, req)

	if rw.Code != http.StatusOK {
		t.Fatalf("expected 200 for valid API key, got %d: %s", rw.Code, rw.Body.String())
	}
}

func TestAPIKeyMiddleware_InvalidKey(t *testing.T) {
	apiKey := "correct-key"
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	handler := apiKeyMiddleware(apiKey, next)

	req := httptest.NewRequest(http.MethodPost, "/mcp", nil)
	req.Header.Set("Authorization", "Bearer wrong-key")
	rw := httptest.NewRecorder()
	handler.ServeHTTP(rw, req)

	if rw.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for invalid API key, got %d", rw.Code)
	}
}

func TestAPIKeyMiddleware_MissingHeader(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	handler := apiKeyMiddleware("some-key", next)

	req := httptest.NewRequest(http.MethodPost, "/mcp", nil)
	rw := httptest.NewRecorder()
	handler.ServeHTTP(rw, req)

	if rw.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for missing header, got %d", rw.Code)
	}

	// Check WWW-Authenticate header is set
	if rw.Header().Get("WWW-Authenticate") != "Bearer" {
		t.Error("expected WWW-Authenticate: Bearer header")
	}
}

func TestAPIKeyMiddleware_BearerCaseInsensitive(t *testing.T) {
	apiKey := "my-key"
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	handler := apiKeyMiddleware(apiKey, next)

	// Test "bearer" (lowercase) prefix
	req := httptest.NewRequest(http.MethodPost, "/mcp", nil)
	req.Header.Set("Authorization", "bearer "+apiKey)
	rw := httptest.NewRecorder()
	handler.ServeHTTP(rw, req)

	if rw.Code != http.StatusOK {
		t.Fatalf("expected 200 for lowercase bearer prefix, got %d", rw.Code)
	}
}

func TestAPIKeyMiddleware_RawToken(t *testing.T) {
	// Some clients may send the token without "Bearer " prefix
	apiKey := "my-key"
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	handler := apiKeyMiddleware(apiKey, next)

	req := httptest.NewRequest(http.MethodPost, "/mcp", nil)
	req.Header.Set("Authorization", apiKey)
	rw := httptest.NewRecorder()
	handler.ServeHTTP(rw, req)

	if rw.Code != http.StatusOK {
		t.Fatalf("expected 200 for raw token, got %d", rw.Code)
	}
}

func TestAPIKeyMiddleware_TimingResistant(t *testing.T) {
	// Verify that the middleware uses constant-time comparison by ensuring
	// incorrect keys with the same length as the real key are rejected.
	// (Timing attack resistance is ensured by subtle.ConstantTimeCompare)
	apiKey := "abcdefghijklmnop"
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	handler := apiKeyMiddleware(apiKey, next)

	// Same length, different content
	req := httptest.NewRequest(http.MethodPost, "/mcp", nil)
	req.Header.Set("Authorization", "Bearer abcdefghijklmnox")
	rw := httptest.NewRecorder()
	handler.ServeHTTP(rw, req)

	if rw.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for similar key, got %d", rw.Code)
	}
}

// --- Health Endpoint Tests ---

func TestHealthHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rw := httptest.NewRecorder()
	healthHandler(rw, req)

	if rw.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rw.Code)
	}
	if ct := rw.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %q", ct)
	}
	if body := rw.Body.String(); body != `{"status":"ok"}` {
		t.Errorf("expected {\"status\":\"ok\"}, got %q", body)
	}
}

// --- Integration: API Key wired into ServeStreamableHTTP ---

func TestServeStreamableHTTP_WithAPIKey(t *testing.T) {
	cfg := newTestConfig()
	cfg.APIKey = "test-key-for-wiring"
	s := NewServer(cfg)

	var capturedHandler http.Handler

	installTransportHooks(
		t,
		func(_ *mcpserver.MCPServer) error { return nil },
		func(_ *mcpserver.MCPServer, _ ...mcpserver.StreamableHTTPOption) streamableHTTPStarter {
			return &mockStreamableServer{}
		},
	)
	installListenHook(t, func(addr string, handler http.Handler) error {
		capturedHandler = handler
		return nil
	})

	if err := s.ServeStreamableHTTP("127.0.0.1:8080"); err != nil {
		t.Fatalf("ServeStreamableHTTP failed: %v", err)
	}

	if capturedHandler == nil {
		t.Fatal("handler was not captured")
	}

	// Test: request to /mcp without API key → 401
	req := httptest.NewRequest(http.MethodPost, "/mcp", nil)
	rw := httptest.NewRecorder()
	capturedHandler.ServeHTTP(rw, req)
	if rw.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 for /mcp without API key, got %d", rw.Code)
	}

	// Test: request to /mcp with correct API key → should pass auth (may return non-401)
	req = httptest.NewRequest(http.MethodPost, "/mcp", nil)
	req.Header.Set("Authorization", "Bearer test-key-for-wiring")
	rw = httptest.NewRecorder()
	capturedHandler.ServeHTTP(rw, req)
	if rw.Code == http.StatusUnauthorized {
		t.Errorf("expected non-401 for /mcp with correct API key, got %d", rw.Code)
	}

	// Test: health endpoint works without API key
	req = httptest.NewRequest(http.MethodGet, "/health", nil)
	rw = httptest.NewRecorder()
	capturedHandler.ServeHTTP(rw, req)
	if rw.Code != http.StatusOK {
		t.Errorf("expected 200 for /health without auth, got %d", rw.Code)
	}
}

func TestServeStreamableHTTP_WithoutAPIKey(t *testing.T) {
	cfg := newTestConfig()
	// No APIKey set
	s := NewServer(cfg)

	var capturedHandler http.Handler

	installTransportHooks(
		t,
		func(_ *mcpserver.MCPServer) error { return nil },
		func(_ *mcpserver.MCPServer, _ ...mcpserver.StreamableHTTPOption) streamableHTTPStarter {
			return &mockStreamableServer{}
		},
	)
	installListenHook(t, func(addr string, handler http.Handler) error {
		capturedHandler = handler
		return nil
	})

	if err := s.ServeStreamableHTTP("127.0.0.1:8080"); err != nil {
		t.Fatalf("ServeStreamableHTTP failed: %v", err)
	}

	if capturedHandler == nil {
		t.Fatal("handler was not captured")
	}

	// Test: request to /mcp without API key should pass (no auth configured)
	req := httptest.NewRequest(http.MethodPost, "/mcp", nil)
	rw := httptest.NewRecorder()
	capturedHandler.ServeHTTP(rw, req)
	if rw.Code == http.StatusUnauthorized {
		t.Errorf("expected non-401 when no API key is configured, got %d", rw.Code)
	}

	// Test: health endpoint still works
	req = httptest.NewRequest(http.MethodGet, "/health", nil)
	rw = httptest.NewRecorder()
	capturedHandler.ServeHTTP(rw, req)
	if rw.Code != http.StatusOK {
		t.Errorf("expected 200 for /health, got %d", rw.Code)
	}
}

// --- Protected Resource Metadata Tests ---

func TestProtectedResourceMetadataHandler(t *testing.T) {
	handler := protectedResourceMetadataHandler(
		"vsp.example.com:8080",
		"https://login.microsoftonline.com/tenant/v2.0",
	)

	req := httptest.NewRequest(http.MethodGet, "/.well-known/oauth-protected-resource", nil)
	rw := httptest.NewRecorder()
	handler.ServeHTTP(rw, req)

	if rw.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rw.Code)
	}

	var metadata map[string]interface{}
	if err := decodeJSON(rw.Body.Bytes(), &metadata); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	if metadata["resource"] != "vsp.example.com:8080" {
		t.Errorf("expected resource vsp.example.com:8080, got %v", metadata["resource"])
	}

	servers, ok := metadata["authorization_servers"].([]interface{})
	if !ok || len(servers) != 1 {
		t.Fatalf("expected 1 authorization server, got %v", metadata["authorization_servers"])
	}
	if servers[0] != "https://login.microsoftonline.com/tenant/v2.0" {
		t.Errorf("unexpected authorization server: %v", servers[0])
	}
}

// --- Username Mapping Tests ---

func TestLoadUsernameMapping_Empty(t *testing.T) {
	mapping := loadUsernameMapping("")
	if mapping != nil {
		t.Errorf("expected nil for empty path, got %v", mapping)
	}
}

func TestLoadUsernameMapping_NonexistentFile(t *testing.T) {
	mapping := loadUsernameMapping("/nonexistent/path/mapping.json")
	if mapping != nil {
		t.Errorf("expected nil for nonexistent file, got %v", mapping)
	}
}

// helper for JSON decoding in tests
func decodeJSON(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
