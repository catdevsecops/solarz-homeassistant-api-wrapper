// Package service contains business logic.
package service

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/catdevsecops/solarz-api/internal/model"
)

const (
	// defaultSolarzEndpoint is the default trusted Solarz API endpoint.
	defaultSolarzEndpoint = "https://app.solarz.com.br/shareable/generation/period"
	// httpClientTimeout is the timeout for HTTP requests (prevent hanging connections).
	httpClientTimeout = 10 * time.Second
)

// allowedSolarzHosts contains the list of trusted hosts for Solarz API.
var allowedSolarzHosts = map[string]bool{
	"app.solarz.com.br": true,
}

// getSecureHTTPClient returns a configured HTTP client with security restrictions.
func getSecureHTTPClient() *http.Client {
	return &http.Client{
		Timeout: httpClientTimeout,
		Transport: &http.Transport{
			// Disable redirects to prevent open redirect attacks
			DisableKeepAlives:     true,
			MaxIdleConns:          5,
			MaxIdleConnsPerHost:   1,
			MaxConnsPerHost:       1,
			IdleConnTimeout:       5 * time.Second,
			TLSHandshakeTimeout:   5 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			// Prevent redirects
			return http.ErrUseLastResponse
		},
	}
}

func isValidSolarzURL(urlString string) error {
	// Parse the URL to validate its format.
	parsedURL, err := url.Parse(urlString)
	if err != nil {
		return err
	}

	// Validate scheme is HTTPS (no HTTP, file://, etc).
	if parsedURL.Scheme != "https" {
		return errors.New("invalid URL scheme: only HTTPS is allowed")
	}

	// Extract hostname without port.
	hostname := parsedURL.Hostname()
	if hostname == "" {
		return errors.New("invalid URL: missing hostname")
	}

	// Check against whitelist of allowed hosts.
	if !allowedSolarzHosts[hostname] {
		return errors.New("untrusted host: " + hostname + " not in whitelist")
	}

	// Prevent requests to private/local IPs (RFC 1918, localhost, etc).
	ip := net.ParseIP(hostname)
	if ip != nil {
		if ip.IsLoopback() || ip.IsPrivate() || ip.IsUnspecified() {
			return errors.New("invalid host: private/local IP not allowed")
		}
	}

	// Prevent requests to cloud metadata endpoints.
	if hostname == "169.254.169.254" || hostname == "metadata.google.internal" {
		return errors.New("invalid host: cloud metadata endpoint not allowed")
	}

	return nil
}

// GetData returns all items from Solarz API.
func GetData() ([]model.Item, error) {
	return GetDataWithContext(context.Background())
}

// GetDataWithContext returns all items from Solarz API with context.
func GetDataWithContext(ctx context.Context) ([]model.Item, error) {
	solarzURL := os.Getenv("SOLARZ_ENDPOINT")
	if solarzURL == "" {
		solarzURL = defaultSolarzEndpoint
	}

	// Validate URL to prevent SSRF attacks.
	if err := isValidSolarzURL(solarzURL); err != nil {
		log.Printf("invalid SOLARZ_ENDPOINT: %v", err)
		return nil, err
	}

	// #nosec G704 - URL is validated by isValidSolarzURL() function above
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, solarzURL, nil)
	if err != nil {
		return nil, err
	}

	// Use secure HTTP client with timeouts and redirect prevention
	secureClient := getSecureHTTPClient()
	// #nosec G704 - Client uses secure configuration with validation above
	resp, err := secureClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("failed to close response body: %v", err)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var solarzResp model.SolarzResponse
	if err := json.Unmarshal(body, &solarzResp); err != nil {
		return nil, err
	}

	// Finds the item with the most recent date.
	result := make([]model.Item, 0)
	var latestDado *model.DadoGeracao

	for i := range solarzResp.Dados {
		if latestDado == nil || solarzResp.Dados[i].Data > latestDado.Data {
			latestDado = &solarzResp.Dados[i]
		}
	}

	if latestDado != nil {
		item := model.Item{
			ID:    latestDado.Data,
			Name:  latestDado.Data + " - " + latestDado.Denominacao,
			Value: formatFloat(latestDado.Quantidade),
		}
		result = append(result, item)
	}

	return result, nil
}

// FormatFloat converts float64 to formatted string.
func FormatFloat(val float64) string {
	return strconv.FormatFloat(val, 'f', 2, 64)
}

// formatFloat is an internal wrapper for FormatFloat.
func formatFloat(val float64) string {
	return FormatFloat(val)
}
