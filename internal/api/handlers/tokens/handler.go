// Package tokens provides the REST injection API for the in-memory token store.
// These endpoints are used by the web-creator platform to inject, remove, and list
// API tokens at worker startup.
//
// All endpoints are protected by the X-Inject-Secret header which must match the
// CLIPROXY_INJECT_SECRET environment variable. If that variable is empty the endpoints
// are disabled entirely and return 404.
package tokens

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/router-for-me/CLIProxyAPI/v6/internal/store"
)

// Handler holds a reference to the shared MemoryStore and the configured inject secret.
type Handler struct {
	store  *store.MemoryStore
	secret string
}

// NewHandler creates a token injection handler.
// secret is the value of CLIPROXY_INJECT_SECRET; if empty the handler is disabled.
func NewHandler(ms *store.MemoryStore, secret string) *Handler {
	return &Handler{
		store:  ms,
		secret: strings.TrimSpace(secret),
	}
}

// Enabled reports whether the injection API is active (i.e. a secret was configured).
func (h *Handler) Enabled() bool {
	return h != nil && h.secret != ""
}

// Middleware returns a Gin middleware that validates the X-Inject-Secret header.
// If the handler is disabled the middleware always aborts with 404.
func (h *Handler) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !h.Enabled() {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "token injection API is disabled"})
			return
		}
		provided := strings.TrimSpace(c.GetHeader("X-Inject-Secret"))
		if provided == "" || provided != h.secret {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or missing X-Inject-Secret"})
			return
		}
		c.Next()
	}
}

// injectRequest is the payload accepted by POST /api/tokens.
type injectRequest struct {
	Token    string         `json:"token" binding:"required"`
	Provider string         `json:"provider" binding:"required"`
	Metadata map[string]any `json:"metadata"`
}

// injectResponse is returned on successful injection.
type injectResponse struct {
	ID        string    `json:"id"`        // masked token
	Provider  string    `json:"provider"`
	CreatedAt time.Time `json:"created_at"`
}

// InjectToken handles POST /api/tokens.
// Injects a new token into the in-memory store.
func (h *Handler) InjectToken(c *gin.Context) {
	var req injectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request: " + err.Error()})
		return
	}

	auth, err := h.store.InjectToken(c.Request.Context(), req.Token, req.Provider, req.Metadata)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	summaries := h.store.ListTokens()
	var masked string
	for _, s := range summaries {
		if strings.HasPrefix(req.Token, s.ID[:min(6, len(s.ID))]) {
			masked = s.ID
			break
		}
	}
	if masked == "" {
		masked = req.Token[:min(6, len(req.Token))] + "..."
	}

	c.JSON(http.StatusCreated, injectResponse{
		ID:        masked,
		Provider:  auth.Provider,
		CreatedAt: auth.CreatedAt,
	})
}

// RemoveToken handles DELETE /api/tokens/:token.
// Removes a token from the in-memory store.
func (h *Handler) RemoveToken(c *gin.Context) {
	token := strings.TrimSpace(c.Param("token"))
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "token param is required"})
		return
	}

	if err := h.store.RemoveToken(c.Request.Context(), token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "token removed"})
}

// ListTokens handles GET /api/tokens.
// Returns a masked list of all tokens currently in the store.
func (h *Handler) ListTokens(c *gin.Context) {
	summaries := h.store.ListTokens()
	c.JSON(http.StatusOK, gin.H{"tokens": summaries})
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
