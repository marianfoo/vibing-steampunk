package adt

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// BTPConfig holds SAP BTP deployment configuration.
// Parsed from VCAP_SERVICES environment variable when running on Cloud Foundry.
type BTPConfig struct {
	// XSUAA credentials (from VCAP_SERVICES binding)
	XSUAAUrl      string // UAA URL (e.g., "https://tenant.authentication.eu10.hana.ondemand.com")
	XSUAAClientID string // OAuth2 client ID
	XSUAASecret   string // OAuth2 client secret
	XSAppName     string // Application name for scope checks

	// Destination service credentials
	DestinationURL      string // Destination service URL
	DestinationClientID string
	DestinationSecret   string
	DestinationTokenURL string

	// Connectivity service (Cloud Connector proxy)
	ConnectivityProxyHost string // e.g., "connectivityproxy.internal.cf.us10-001.hana.ondemand.com"
	ConnectivityProxyPort string // e.g., "20003"
	ConnectivityClientID  string
	ConnectivitySecret    string
	ConnectivityTokenURL  string
}

// VCAPServices represents the VCAP_SERVICES JSON structure.
type VCAPServices struct {
	XSUAA        []vcapBinding `json:"xsuaa"`
	Destination  []vcapBinding `json:"destination"`
	Connectivity []vcapBinding `json:"connectivity"`
}

type vcapBinding struct {
	Name        string          `json:"name"`
	Credentials json.RawMessage `json:"credentials"`
}

type xsuaaCredentials struct {
	URL          string `json:"url"`
	ClientID     string `json:"clientid"`
	ClientSecret string `json:"clientsecret"`
	XSAppName    string `json:"xsappname"`
}

type destinationCredentials struct {
	URI          string `json:"uri"`
	URL          string `json:"url"`
	ClientID     string `json:"clientid"`
	ClientSecret string `json:"clientsecret"`
	TokenURL     string `json:"token_service_url"`
}

type connectivityCredentials struct {
	ClientID           string `json:"clientid"`
	ClientSecret       string `json:"clientsecret"`
	TokenURL           string `json:"token_service_url"`
	URL                string `json:"url"` // XSUAA URL fallback
	OnPremiseProxyHost string `json:"onpremise_proxy_host"`
	OnPremiseProxyPort string `json:"onpremise_proxy_http_port"`
}

// ParseVCAPServices parses BTP service bindings from the VCAP_SERVICES environment variable.
// Returns nil if VCAP_SERVICES is not set (not running on BTP).
func ParseVCAPServices() (*BTPConfig, error) {
	vcapJSON := os.Getenv("VCAP_SERVICES")
	if vcapJSON == "" {
		return nil, nil // Not running on BTP
	}

	var vcap VCAPServices
	if err := json.Unmarshal([]byte(vcapJSON), &vcap); err != nil {
		return nil, fmt.Errorf("parsing VCAP_SERVICES: %w", err)
	}

	config := &BTPConfig{}

	// Parse XSUAA binding
	if len(vcap.XSUAA) > 0 {
		var creds xsuaaCredentials
		if err := json.Unmarshal(vcap.XSUAA[0].Credentials, &creds); err != nil {
			return nil, fmt.Errorf("parsing XSUAA credentials: %w", err)
		}
		config.XSUAAUrl = creds.URL
		config.XSUAAClientID = creds.ClientID
		config.XSUAASecret = creds.ClientSecret
		config.XSAppName = creds.XSAppName
	}

	// Parse Destination binding
	if len(vcap.Destination) > 0 {
		var creds destinationCredentials
		if err := json.Unmarshal(vcap.Destination[0].Credentials, &creds); err != nil {
			return nil, fmt.Errorf("parsing Destination credentials: %w", err)
		}
		config.DestinationURL = creds.URI
		if config.DestinationURL == "" {
			config.DestinationURL = creds.URL
		}
		config.DestinationClientID = creds.ClientID
		config.DestinationSecret = creds.ClientSecret
		config.DestinationTokenURL = creds.TokenURL
		// Fallback: construct token URL from XSUAA URL if token_service_url is not present
		if config.DestinationTokenURL == "" && creds.URL != "" {
			config.DestinationTokenURL = strings.TrimSuffix(creds.URL, "/") + "/oauth/token"
		}
	}

	// Parse Connectivity binding
	if len(vcap.Connectivity) > 0 {
		var creds connectivityCredentials
		if err := json.Unmarshal(vcap.Connectivity[0].Credentials, &creds); err != nil {
			return nil, fmt.Errorf("parsing Connectivity credentials: %w", err)
		}
		config.ConnectivityProxyHost = creds.OnPremiseProxyHost
		config.ConnectivityProxyPort = creds.OnPremiseProxyPort
		config.ConnectivityClientID = creds.ClientID
		config.ConnectivitySecret = creds.ClientSecret
		config.ConnectivityTokenURL = creds.TokenURL
		// Fallback: construct token URL from XSUAA URL
		if config.ConnectivityTokenURL == "" && creds.URL != "" {
			config.ConnectivityTokenURL = strings.TrimSuffix(creds.URL, "/") + "/oauth/token"
		} else if config.ConnectivityTokenURL != "" && !strings.HasSuffix(config.ConnectivityTokenURL, "/oauth/token") {
			// token_service_url is the base XSUAA URL, need to append /oauth/token
			config.ConnectivityTokenURL = strings.TrimSuffix(config.ConnectivityTokenURL, "/") + "/oauth/token"
		}
	}

	return config, nil
}

// ToOAuthConfigForXSUAA converts BTP XSUAA credentials to OAuthConfig.
func (b *BTPConfig) ToOAuthConfigForXSUAA() *OAuthConfig {
	if b.XSUAAUrl == "" || b.XSUAAClientID == "" {
		return nil
	}
	return &OAuthConfig{
		TokenURL:     b.XSUAAUrl + "/oauth/token",
		ClientID:     b.XSUAAClientID,
		ClientSecret: b.XSUAASecret,
	}
}

// DestinationLookup retrieves a destination from the BTP Destination Service.
// Returns the destination's URL and authentication details.
type DestinationLookup struct {
	config     *BTPConfig
	httpClient *http.Client
}

// Destination holds a resolved BTP destination.
type Destination struct {
	Name           string `json:"Name"`
	Type           string `json:"Type"`
	URL            string `json:"URL"`
	Authentication string `json:"Authentication"`
	ProxyType      string `json:"ProxyType"`
	User           string `json:"User"`
	Password       string `json:"Password"`
	SAPClient      string `json:"sap-client"`

	// Token-based auth (from destination service token exchange)
	AuthTokens []DestinationAuthToken `json:"authTokens"`
}

// DestinationAuthToken represents an auth token from the destination service.
type DestinationAuthToken struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

// NewDestinationLookup creates a client for the BTP Destination Service API.
func NewDestinationLookup(config *BTPConfig) *DestinationLookup {
	return &DestinationLookup{
		config:     config,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// GetDestination retrieves a destination by name from the BTP Destination Service.
// See: https://help.sap.com/docs/connectivity/sap-btp-connectivity-cf/calling-destination-service-rest-api
func (d *DestinationLookup) GetDestination(ctx context.Context, name string) (*Destination, error) {
	// First, get a token for the destination service
	token, err := d.getDestinationToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting destination service token: %w", err)
	}

	// Call the destination service API
	destURL := strings.TrimSuffix(d.config.DestinationURL, "/") +
		"/destination-configuration/v1/destinations/" + url.PathEscape(name)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, destURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating destination request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := d.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("calling destination service: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading destination response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("destination service returned HTTP %d: %s", resp.StatusCode, string(body))
	}

	// The response wraps the destination in a destinationConfiguration object
	var result struct {
		DestinationConfiguration Destination `json:"destinationConfiguration"`
		AuthTokens               []DestinationAuthToken `json:"authTokens"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("parsing destination response: %w", err)
	}

	dest := &result.DestinationConfiguration
	dest.AuthTokens = result.AuthTokens
	return dest, nil
}

// getDestinationToken obtains an OAuth2 token for the Destination Service API.
func (d *DestinationLookup) getDestinationToken(ctx context.Context) (string, error) {
	tokenURL := d.config.DestinationTokenURL
	if tokenURL == "" {
		// Fallback: use XSUAA token URL
		tokenURL = d.config.XSUAAUrl + "/oauth/token"
	}

	data := url.Values{
		"grant_type":    {"client_credentials"},
		"client_id":     {d.config.DestinationClientID},
		"client_secret": {d.config.DestinationSecret},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenURL,
		strings.NewReader(data.Encode()))
	if err != nil {
		return "", fmt.Errorf("creating token request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := d.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("fetching token: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("reading token response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("token endpoint returned HTTP %d: %s", resp.StatusCode, string(body))
	}

	var token struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.Unmarshal(body, &token); err != nil {
		return "", fmt.Errorf("parsing token: %w", err)
	}

	return token.AccessToken, nil
}

// ResolveDestination looks up a BTP destination by name and returns the SAP URL,
// username, and password extracted from the destination configuration.
// This is used on startup to override SAP_URL / SAP_USER / SAP_PASSWORD
// when running on BTP with a configured destination.
func (d *DestinationLookup) ResolveDestination(ctx context.Context, name string) (sapURL, user, password string, err error) {
	dest, err := d.GetDestination(ctx, name)
	if err != nil {
		return "", "", "", fmt.Errorf("looking up destination %q: %w", name, err)
	}
	return dest.URL, dest.User, dest.Password, nil
}

// ConnectivityProxyTransport returns an http.RoundTripper that routes requests through
// the BTP Connectivity Service proxy (Cloud Connector). It obtains a JWT token from the
// connectivity service's XSUAA and sets it as Proxy-Authorization header.
func (b *BTPConfig) ConnectivityProxyTransport(base *http.Transport) (http.RoundTripper, error) {
	if b.ConnectivityProxyHost == "" {
		return base, nil // No connectivity proxy configured
	}

	proxyURL := &url.URL{
		Scheme: "http",
		Host:   b.ConnectivityProxyHost + ":" + b.ConnectivityProxyPort,
	}
	base.Proxy = http.ProxyURL(proxyURL)

	return &connectivityProxyRoundTripper{
		base:   base,
		config: b,
	}, nil
}

// connectivityProxyRoundTripper injects a Proxy-Authorization header with a JWT
// obtained from the connectivity service's XSUAA. The token is cached and refreshed.
type connectivityProxyRoundTripper struct {
	base   http.RoundTripper
	config *BTPConfig

	token     string
	expiresAt time.Time
}

func (rt *connectivityProxyRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	token, err := rt.getToken(req.Context())
	if err != nil {
		return nil, fmt.Errorf("getting connectivity proxy token: %w", err)
	}
	req = req.Clone(req.Context())
	req.Header.Set("Proxy-Authorization", "Bearer "+token)
	// SAP-Connectivity-SCC-Location_ID header can be added here if needed
	return rt.base.RoundTrip(req)
}

func (rt *connectivityProxyRoundTripper) getToken(ctx context.Context) (string, error) {
	if rt.token != "" && time.Now().Before(rt.expiresAt) {
		return rt.token, nil
	}

	tokenURL := rt.config.ConnectivityTokenURL
	if tokenURL == "" {
		return "", fmt.Errorf("connectivity token URL not configured")
	}

	data := url.Values{
		"grant_type":    {"client_credentials"},
		"client_id":     {rt.config.ConnectivityClientID},
		"client_secret": {rt.config.ConnectivitySecret},
	}

	tokenReq, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenURL,
		strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}
	tokenReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Use a direct client (no proxy) for token requests
	directClient := &http.Client{Timeout: 30 * time.Second}
	resp, err := directClient.Do(tokenReq)
	if err != nil {
		return "", fmt.Errorf("fetching connectivity token: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("connectivity token endpoint returned HTTP %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	rt.token = result.AccessToken
	// Refresh 60 seconds before expiry
	rt.expiresAt = time.Now().Add(time.Duration(result.ExpiresIn-60) * time.Second)
	return rt.token, nil
}
