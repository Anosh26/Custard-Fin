package httpapi

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/anosh26/custardfin/apps/api/internal/domain"
	"github.com/anosh26/custardfin/apps/api/internal/ports"
)

type contextKey string

const (
	claimsContextKey contextKey = "user_claims"
	nodeContextKey   contextKey = "node_entity"
)

// responseWriterInterceptor captures HTTP status code for structured logging.
type responseWriterInterceptor struct {
	http.ResponseWriter
	statusCode int
}

func (w *responseWriterInterceptor) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

// LoggingMiddleware logs each incoming request with method, path, remote addr, duration, and status.
func LoggingMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			interceptor := &responseWriterInterceptor{ResponseWriter: w, statusCode: http.StatusOK}

			next.ServeHTTP(interceptor, r)

			duration := time.Since(start)
			logger.Info("http request",
				"method", r.Method,
				"path", r.URL.Path,
				"remote_addr", r.RemoteAddr,
				"status", interceptor.statusCode,
				"duration_ms", duration.Milliseconds(),
			)
		})
	}
}

// CORSMiddleware enables CORS for UI access.
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Node-Token")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// RecoverMiddleware handles unhandled panics gracefully without terminating the process.
func RecoverMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					logger.Error("panic recovered in http handler", "err", rec)
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)
					_ = json.NewEncoder(w).Encode(map[string]string{
						"code":    "INTERNAL_ERROR",
						"message": "Internal server error occurred",
					})
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

// RequireUserAuth validates JWT Bearer authorization and injects claims into context.
func RequireUserAuth(tokenService ports.TokenService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Missing or invalid Authorization header")
				return
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := tokenService.ValidateToken(tokenStr)
			if err != nil {
				writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid or expired token")
				return
			}

			ctx := context.WithValue(r.Context(), claimsContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireNodeAuth validates X-Node-Token header against node ID in the request path.
func RequireNodeAuth(nodeRepo ports.NodeRepository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			nodeToken := r.Header.Get("X-Node-Token")
			if strings.TrimSpace(nodeToken) == "" {
				writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Missing X-Node-Token header")
				return
			}

			tokenHash := domain.HashNodeToken(nodeToken)
			node, err := nodeRepo.FindByTokenHash(r.Context(), tokenHash)
			if err != nil {
				writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid node token")
				return
			}

			// Validate node ID from path if present
			pathNodeID := r.PathValue("id")
			if pathNodeID != "" && pathNodeID != node.ID {
				writeError(w, http.StatusForbidden, "FORBIDDEN", "Token does not match target node ID")
				return
			}

			ctx := context.WithValue(r.Context(), nodeContextKey, node)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserClaims extracts validated Claims from context.
func GetUserClaims(ctx context.Context) (*ports.Claims, bool) {
	claims, ok := ctx.Value(claimsContextKey).(*ports.Claims)
	return claims, ok
}

// GetNodeEntity extracts validated Node from context.
func GetNodeEntity(ctx context.Context) (*domain.Node, bool) {
	node, ok := ctx.Value(nodeContextKey).(*domain.Node)
	return node, ok
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"code":    code,
		"message": message,
	})
}
