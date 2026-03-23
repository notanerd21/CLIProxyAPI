// Package util provides utility functions for the CLI Proxy API server.
// Proxy support has been removed in the notanerd21 fork — SetProxy is a no-op.
// All connections go direct to AI providers. Less code, less attack surface.
package util

import (
	"net/http"

	"github.com/router-for-me/CLIProxyAPI/v6/sdk/config"
)

// SetProxy is a no-op. Proxy support has been removed.
// Returns the httpClient unmodified for backward compatibility.
func SetProxy(_ *config.SDKConfig, httpClient *http.Client) *http.Client {
	return httpClient
}
