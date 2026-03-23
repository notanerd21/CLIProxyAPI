// Package store provides token storage backends for CLIProxyAPI.
package store

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	cliproxyauth "github.com/router-for-me/CLIProxyAPI/v6/sdk/cliproxy/auth"
)

// MemoryStore is a pure in-memory implementation of the coreauth.Store interface.
// All records are held in RAM and lost on process restart — the web-creator platform
// re-injects tokens at worker startup via the REST injection API.
type MemoryStore struct {
	mu      sync.RWMutex
	records map[string]*cliproxyauth.Auth
}

// NewMemoryStore creates an empty in-memory token store.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		records: make(map[string]*cliproxyauth.Auth),
	}
}

// Save persists an Auth record into the in-memory map.
// The record is keyed by auth.ID. If an existing record has the same ID it is replaced.
// Returns a synthetic "path" string (the record ID) so callers that use the return value
// for logging still get something meaningful.
func (s *MemoryStore) Save(_ context.Context, auth *cliproxyauth.Auth) (string, error) {
	if auth == nil {
		return "", fmt.Errorf("memory store: auth is nil")
	}
	id := strings.TrimSpace(auth.ID)
	if id == "" {
		return "", fmt.Errorf("memory store: auth.ID is required")
	}

	now := time.Now().UTC()
	clone := auth.Clone()
	if clone.CreatedAt.IsZero() {
		clone.CreatedAt = now
	}
	clone.UpdatedAt = now

	s.mu.Lock()
	s.records[id] = clone
	s.mu.Unlock()

	return id, nil
}

// List returns a snapshot of all auth records currently held in memory.
func (s *MemoryStore) List(_ context.Context) ([]*cliproxyauth.Auth, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]*cliproxyauth.Auth, 0, len(s.records))
	for _, rec := range s.records {
		result = append(result, rec.Clone())
	}
	return result, nil
}

// Delete removes the auth record identified by id from the in-memory map.
// Returns nil if the record does not exist (idempotent).
func (s *MemoryStore) Delete(_ context.Context, id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return fmt.Errorf("memory store: id is empty")
	}
	s.mu.Lock()
	delete(s.records, id)
	s.mu.Unlock()
	return nil
}

// InjectToken is a convenience helper used by the REST injection API.
// It builds a minimal Auth record from a raw token string and provider name,
// then persists it via Save. The resulting record ID equals the token itself
// so the injection API can address it directly.
func (s *MemoryStore) InjectToken(ctx context.Context, token, provider string, metadata map[string]any) (*cliproxyauth.Auth, error) {
	token = strings.TrimSpace(token)
	if token == "" {
		return nil, fmt.Errorf("memory store: token is empty")
	}
	provider = strings.TrimSpace(strings.ToLower(provider))
	if provider == "" {
		return nil, fmt.Errorf("memory store: provider is required")
	}

	if metadata == nil {
		metadata = make(map[string]any)
	}

	// Store the raw API key in both the standard attribute and metadata so that
	// the existing auth executors that read Attributes["api_key"] can find it.
	auth := &cliproxyauth.Auth{
		ID:       token,
		Provider: provider,
		Status:   cliproxyauth.StatusActive,
		Attributes: map[string]string{
			"api_key": token,
		},
		Metadata:  metadata,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	if _, err := s.Save(ctx, auth); err != nil {
		return nil, err
	}
	return auth.Clone(), nil
}

// RemoveToken removes a token record identified by the raw token string.
// This is a thin alias for Delete exposed for the REST injection API.
func (s *MemoryStore) RemoveToken(ctx context.Context, token string) error {
	return s.Delete(ctx, strings.TrimSpace(token))
}

// ListTokens returns a masked view of all stored tokens suitable for the REST listing endpoint.
// Each returned entry replaces the middle characters of the token with "..." to avoid leaking
// the full secret over the wire.
func (s *MemoryStore) ListTokens() []TokenSummary {
	s.mu.RLock()
	defer s.mu.RUnlock()

	summaries := make([]TokenSummary, 0, len(s.records))
	for id, rec := range s.records {
		summaries = append(summaries, TokenSummary{
			ID:        maskToken(id),
			Provider:  rec.Provider,
			Status:    string(rec.Status),
			CreatedAt: rec.CreatedAt,
		})
	}
	return summaries
}

// TokenSummary is the masked representation returned by the listing endpoint.
type TokenSummary struct {
	ID        string    `json:"id"`
	Provider  string    `json:"provider"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// maskToken replaces the body of a token with "..." keeping only the first 6 and last 4 chars.
func maskToken(token string) string {
	const keepPrefix = 6
	const keepSuffix = 4
	if len(token) <= keepPrefix+keepSuffix {
		return token[:min(2, len(token))] + "..."
	}
	return token[:keepPrefix] + "..." + token[len(token)-keepSuffix:]
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
