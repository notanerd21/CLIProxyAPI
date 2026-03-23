// Package usage previously provided usage tracking and logging functionality.
// Usage statistics have been removed in the notanerd21 fork.
// Web-creator tracks its own usage metrics.
//
// All public functions are retained as no-ops for backward compatibility.
package usage

import (
	"context"

	"github.com/gin-gonic/gin"
	coreusage "github.com/router-for-me/CLIProxyAPI/v6/sdk/cliproxy/usage"
)

// init is intentionally empty — we do NOT register the plugin.
func init() {}

// SetStatisticsEnabled is a no-op. Usage statistics are disabled.
func SetStatisticsEnabled(_ bool) {}

// StatisticsEnabled always returns false.
func StatisticsEnabled() bool { return false }

// RequestStatistics is a stub type for backward compatibility.
type RequestStatistics struct{}

// StatisticsSnapshot is a stub type for backward compatibility.
type StatisticsSnapshot struct {
	TotalRequests  int64              `json:"total_requests"`
	SuccessCount   int64              `json:"success_count"`
	FailureCount   int64              `json:"failure_count"`
	TotalTokens    int64              `json:"total_tokens"`
	APIs           map[string]any     `json:"apis"`
	RequestsByDay  map[string]int64   `json:"requests_by_day"`
	RequestsByHour map[string]int64   `json:"requests_by_hour"`
	TokensByDay    map[string]int64   `json:"tokens_by_day"`
	TokensByHour   map[string]int64   `json:"tokens_by_hour"`
}

// MergeResult is a stub type for backward compatibility.
type MergeResult struct {
	TotalRequests int64 `json:"total_requests"`
	SuccessCount  int64 `json:"success_count"`
	FailureCount  int64 `json:"failure_count"`
	TotalTokens   int64 `json:"total_tokens"`
	Added         int64 `json:"added"`
	Skipped       int64 `json:"skipped"`
}

// GetRequestStatistics returns a no-op stub.
func GetRequestStatistics() *RequestStatistics { return &RequestStatistics{} }

// Snapshot returns an empty snapshot.
func (s *RequestStatistics) Snapshot() StatisticsSnapshot { return StatisticsSnapshot{} }

// MergeSnapshot is a no-op. Returns empty result.
func (s *RequestStatistics) MergeSnapshot(_ StatisticsSnapshot) MergeResult { return MergeResult{} }

// Record is a no-op.
func (s *RequestStatistics) Record(_ context.Context, _ coreusage.Record) {}

// LoggerPlugin is a stub for backward compatibility.
type LoggerPlugin struct{}

// NewLoggerPlugin returns a no-op plugin.
func NewLoggerPlugin() *LoggerPlugin { return &LoggerPlugin{} }

// HandleUsage is a no-op.
func (p *LoggerPlugin) HandleUsage(_ context.Context, _ coreusage.Record) {}

// --- Management handler stubs (referenced by management handler) ---

// GetUsageStatisticsJSON returns empty JSON for the management endpoint.
func GetUsageStatisticsJSON(_ *gin.Context) StatisticsSnapshot {
	return StatisticsSnapshot{}
}
