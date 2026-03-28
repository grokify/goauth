// Package middleware provides HTTP middleware for authentication.
package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/grokify/goauth/jwt"
)

// contextKey is a custom type for context keys to avoid collisions.
type contextKey string

const (
	// UserContextKey is the context key for the authenticated user.
	UserContextKey contextKey = "authenticated_user"
)

// AuthenticatedUser represents the authenticated user in the request context.
type AuthenticatedUser struct {
	ID              uuid.UUID
	Email           string
	IsPlatformAdmin bool
}

// BearerAuth creates an authentication middleware that extracts JWT from Bearer token.
// It adds the authenticated user to the request context if a valid token is present.
// Requests without a token or with an invalid token continue without user context.
func BearerAuth(jwtService *jwt.Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract token from Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				next.ServeHTTP(w, r)
				return
			}

			// Check for Bearer prefix
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
				next.ServeHTTP(w, r)
				return
			}

			tokenString := parts[1]

			// Validate token
			claims, err := jwtService.ValidateAccessToken(tokenString)
			if err != nil {
				// Token is invalid or expired, continue without user context
				next.ServeHTTP(w, r)
				return
			}

			// Add user to context
			user := &AuthenticatedUser{
				ID:              claims.UserID,
				Email:           claims.Email,
				IsPlatformAdmin: claims.IsPlatformAdmin,
			}

			ctx := context.WithValue(r.Context(), UserContextKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// UserFromContext extracts the authenticated user from the context.
func UserFromContext(ctx context.Context) *AuthenticatedUser {
	user, ok := ctx.Value(UserContextKey).(*AuthenticatedUser)
	if !ok {
		return nil
	}
	return user
}

// RequireAuth is a middleware that requires authentication.
// Returns 401 Unauthorized if no valid token is present.
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := UserFromContext(r.Context())
		if user == nil {
			http.Error(w, `{"title":"Unauthorized","status":401,"detail":"Authentication required"}`, http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// RequirePlatformAdmin is a middleware that requires platform admin access.
// Returns 403 Forbidden if the user is not a platform admin.
func RequirePlatformAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := UserFromContext(r.Context())
		if user == nil {
			http.Error(w, `{"title":"Unauthorized","status":401,"detail":"Authentication required"}`, http.StatusUnauthorized)
			return
		}
		if !user.IsPlatformAdmin {
			http.Error(w, `{"title":"Forbidden","status":403,"detail":"Platform admin access required"}`, http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
