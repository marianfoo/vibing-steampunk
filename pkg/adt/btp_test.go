package adt

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestParseVCAPServices_NotSet(t *testing.T) {
	os.Unsetenv("VCAP_SERVICES")
	config, err := ParseVCAPServices()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if config != nil {
		t.Errorf("expected nil for unset VCAP_SERVICES, got %v", config)
	}
}

func TestParseVCAPServices_XSUAA(t *testing.T) {
	vcap := `{
		"xsuaa": [{
			"name": "vsp-xsuaa",
			"credentials": {
				"url": "https://tenant.authentication.eu10.hana.ondemand.com",
				"clientid": "sb-vsp-app!t12345",
				"clientsecret": "xsuaa-secret-123",
				"xsappname": "vsp-app!t12345"
			}
		}]
	}`

	os.Setenv("VCAP_SERVICES", vcap)
	defer os.Unsetenv("VCAP_SERVICES")

	config, err := ParseVCAPServices()
	if err != nil {
		t.Fatalf("ParseVCAPServices failed: %v", err)
	}
	if config == nil {
		t.Fatal("expected config, got nil")
	}

	if config.XSUAAUrl != "https://tenant.authentication.eu10.hana.ondemand.com" {
		t.Errorf("unexpected XSUAA URL: %s", config.XSUAAUrl)
	}
	if config.XSUAAClientID != "sb-vsp-app!t12345" {
		t.Errorf("unexpected client ID: %s", config.XSUAAClientID)
	}
	if config.XSUAASecret != "xsuaa-secret-123" {
		t.Errorf("unexpected secret: %s", config.XSUAASecret)
	}
	if config.XSAppName != "vsp-app!t12345" {
		t.Errorf("unexpected app name: %s", config.XSAppName)
	}
}

func TestParseVCAPServices_Destination(t *testing.T) {
	vcap := `{
		"xsuaa": [{
			"name": "xsuaa",
			"credentials": {
				"url": "https://tenant.authentication.eu10.hana.ondemand.com",
				"clientid": "xsuaa-client",
				"clientsecret": "xsuaa-secret"
			}
		}],
		"destination": [{
			"name": "destination",
			"credentials": {
				"uri": "https://destination-configuration.cfapps.eu10.hana.ondemand.com",
				"clientid": "dest-client",
				"clientsecret": "dest-secret",
				"token_service_url": "https://tenant.authentication.eu10.hana.ondemand.com/oauth/token"
			}
		}]
	}`

	os.Setenv("VCAP_SERVICES", vcap)
	defer os.Unsetenv("VCAP_SERVICES")

	config, err := ParseVCAPServices()
	if err != nil {
		t.Fatalf("ParseVCAPServices failed: %v", err)
	}

	if config.DestinationURL != "https://destination-configuration.cfapps.eu10.hana.ondemand.com" {
		t.Errorf("unexpected destination URL: %s", config.DestinationURL)
	}
	if config.DestinationClientID != "dest-client" {
		t.Errorf("unexpected dest client ID: %s", config.DestinationClientID)
	}
}

func TestParseVCAPServices_Invalid(t *testing.T) {
	os.Setenv("VCAP_SERVICES", "not-json")
	defer os.Unsetenv("VCAP_SERVICES")

	_, err := ParseVCAPServices()
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestBTPConfig_ToOAuthConfig(t *testing.T) {
	config := &BTPConfig{
		XSUAAUrl:      "https://tenant.authentication.eu10.hana.ondemand.com",
		XSUAAClientID: "my-client",
		XSUAASecret:   "my-secret",
	}

	oauthCfg := config.ToOAuthConfigForXSUAA()
	if oauthCfg == nil {
		t.Fatal("expected OAuthConfig, got nil")
	}
	if oauthCfg.TokenURL != "https://tenant.authentication.eu10.hana.ondemand.com/oauth/token" {
		t.Errorf("unexpected token URL: %s", oauthCfg.TokenURL)
	}
	if oauthCfg.ClientID != "my-client" {
		t.Errorf("unexpected client ID: %s", oauthCfg.ClientID)
	}
}

func TestBTPConfig_ToOAuthConfig_Empty(t *testing.T) {
	config := &BTPConfig{}
	oauthCfg := config.ToOAuthConfigForXSUAA()
	if oauthCfg != nil {
		t.Errorf("expected nil for empty config, got %v", oauthCfg)
	}
}

func TestDestinationLookup_GetDestination(t *testing.T) {
	// Mock token endpoint
	tokenServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{
			"access_token": "test-token-123",
		})
	}))
	defer tokenServer.Close()

	// Mock destination service
	destServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if auth := r.Header.Get("Authorization"); auth != "Bearer test-token-123" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		result := map[string]interface{}{
			"destinationConfiguration": map[string]string{
				"Name":           "SAP_SYSTEM",
				"Type":           "HTTP",
				"URL":            "https://sap.internal:44300",
				"Authentication": "BasicAuthentication",
				"ProxyType":      "OnPremise",
				"User":           "SAP_USER",
				"Password":       "SAP_PASS",
				"sap-client":     "001",
			},
		}
		json.NewEncoder(w).Encode(result)
	}))
	defer destServer.Close()

	config := &BTPConfig{
		DestinationURL:      destServer.URL,
		DestinationClientID: "dest-client",
		DestinationSecret:   "dest-secret",
		DestinationTokenURL: tokenServer.URL,
	}

	lookup := NewDestinationLookup(config)
	dest, err := lookup.GetDestination(context.Background(), "SAP_SYSTEM")
	if err != nil {
		t.Fatalf("GetDestination failed: %v", err)
	}

	if dest.Name != "SAP_SYSTEM" {
		t.Errorf("expected SAP_SYSTEM, got %s", dest.Name)
	}
	if dest.URL != "https://sap.internal:44300" {
		t.Errorf("unexpected URL: %s", dest.URL)
	}
	if dest.Authentication != "BasicAuthentication" {
		t.Errorf("unexpected auth: %s", dest.Authentication)
	}
	if dest.User != "SAP_USER" {
		t.Errorf("unexpected user: %s", dest.User)
	}
}

func TestDestinationLookup_TokenFailure(t *testing.T) {
	tokenServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "invalid_client", http.StatusUnauthorized)
	}))
	defer tokenServer.Close()

	config := &BTPConfig{
		DestinationURL:      "https://dest.example.com",
		DestinationClientID: "bad-client",
		DestinationSecret:   "bad-secret",
		DestinationTokenURL: tokenServer.URL,
	}

	lookup := NewDestinationLookup(config)
	_, err := lookup.GetDestination(context.Background(), "SAP_SYSTEM")
	if err == nil {
		t.Error("expected error for token failure")
	}
}
