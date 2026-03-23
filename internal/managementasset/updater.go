// Package managementasset previously provided auto-updating management panel downloads.
// This functionality has been removed in the notanerd21 fork for security reasons:
// - Auto-downloads from third-party domains (cpamc.router-for.me) eliminated
// - No more 3-hour polling to GitHub releases
// - GITSTORE_GIT_TOKEN no longer leaked in management asset requests
//
// Management is handled by web-creator's backpanel UI instead.
package managementasset

import (
	"context"

	"github.com/router-for-me/CLIProxyAPI/v6/internal/config"
)

// ManagementFileName kept for backward compatibility with any code referencing it.
const ManagementFileName = "management.html"

// SetCurrentConfig is a no-op. Previously stored config for auto-updater decisions.
func SetCurrentConfig(_ *config.Config) {}

// StartAutoUpdater is a no-op. Auto-updating has been removed.
func StartAutoUpdater(_ context.Context, _ string) {}

// EnsureLatestManagementHTML is a no-op. Always returns false (no panel available).
func EnsureLatestManagementHTML(_ context.Context, _, _, _ string) bool { return false }

// StaticDir always returns empty. No management panel assets are served.
func StaticDir(_ string) string { return "" }

// FilePath always returns empty. No management panel assets are served.
func FilePath(_ string) string { return "" }
