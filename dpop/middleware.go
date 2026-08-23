package dpop

import (
	"context"
	"net/http"
	"strings"
)

// contextKey is a type for context keys to avoid collisions.
type contextKey string

const (
	// ContextKeyVerificationResult is the context key for the verification result.
	ContextKeyVerificationResult contextKey = "dpop_verification_result"
	// ContextKeyThumbprint is the context key for the JWK thumbprint.
	ContextKeyThumbprint contextKey = "dpop_thumbprint"
)

// Header names for DPoP.
const (
	// HeaderDPoP is the header name for DPoP proofs.
	HeaderDPoP = "DPoP"
	// HeaderAuthorization is the standard Authorization header.
	HeaderAuthorization = "Authorization"
	// AuthSchemeBearer is the Bearer authorization scheme.
	AuthSchemeBearer = "Bearer"
	// AuthSchemeDPoP is the DPoP authorization scheme.
	AuthSchemeDPoP = "DPoP"
)

// MiddlewareConfig contains configuration for the DPoP middleware.
type MiddlewareConfig struct {
	// Verifier is the DPoP verifier to use.
	Verifier *Verifier
	// RequireDPoP when true rejects requests without DPoP proofs.
	// When false, requests without DPoP proofs are allowed through.
	RequireDPoP bool
	// ExtractAccessToken is a function to extract the access token from the request.
	// If nil, the middleware extracts from the Authorization header.
	ExtractAccessToken func(r *http.Request) string
	// OnError is called when verification fails.
	// If nil, a 401 Unauthorized response is sent.
	OnError func(w http.ResponseWriter, r *http.Request, err error)
}

// Middleware creates HTTP middleware that validates DPoP proofs.
func Middleware(config MiddlewareConfig) func(http.Handler) http.Handler {
	if config.Verifier == nil {
		config.Verifier = NewVerifier(DefaultVerificationConfig())
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract DPoP proof from header
			proof := r.Header.Get(HeaderDPoP)

			// If no proof and not required, continue
			if proof == "" {
				if config.RequireDPoP {
					handleError(w, r, config, ErrInvalidProof)
					return
				}
				next.ServeHTTP(w, r)
				return
			}

			// Extract access token
			var accessToken string
			if config.ExtractAccessToken != nil {
				accessToken = config.ExtractAccessToken(r)
			} else {
				accessToken = extractBearerToken(r)
			}

			// Build request URI (scheme + host + path)
			uri := buildRequestURI(r)

			// Verify the proof
			result, err := config.Verifier.Verify(r.Context(), VerificationRequest{
				Proof:       proof,
				Method:      r.Method,
				URI:         uri,
				AccessToken: accessToken,
			})
			if err != nil {
				handleError(w, r, config, err)
				return
			}

			// Add verification result to context
			ctx := context.WithValue(r.Context(), ContextKeyVerificationResult, result)
			ctx = context.WithValue(ctx, ContextKeyThumbprint, result.Thumbprint)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// extractBearerToken extracts a bearer token from the Authorization header.
func extractBearerToken(r *http.Request) string {
	auth := r.Header.Get(HeaderAuthorization)
	if auth == "" {
		return ""
	}

	// Check for "Bearer " or "DPoP " prefix
	if token, found := strings.CutPrefix(auth, AuthSchemeBearer+" "); found {
		return token
	}
	if token, found := strings.CutPrefix(auth, AuthSchemeDPoP+" "); found {
		return token
	}

	return ""
}

// buildRequestURI constructs the URI for DPoP verification.
// Per RFC 9449, this should be scheme + host + path (no query or fragment).
func buildRequestURI(r *http.Request) string {
	scheme := "https"
	if r.TLS == nil {
		scheme = "http"
	}

	// Use X-Forwarded-Proto if behind a proxy
	if proto := r.Header.Get("X-Forwarded-Proto"); proto != "" {
		scheme = proto
	}

	host := r.Host
	// Use X-Forwarded-Host if behind a proxy
	if forwardedHost := r.Header.Get("X-Forwarded-Host"); forwardedHost != "" {
		host = forwardedHost
	}

	return scheme + "://" + host + r.URL.Path
}

// handleError handles verification errors.
func handleError(w http.ResponseWriter, r *http.Request, config MiddlewareConfig, err error) {
	if config.OnError != nil {
		config.OnError(w, r, err)
		return
	}

	// Default error handling: 401 Unauthorized
	w.Header().Set("WWW-Authenticate", `DPoP error="invalid_dpop_proof"`)
	http.Error(w, err.Error(), http.StatusUnauthorized)
}

// GetVerificationResult retrieves the DPoP verification result from the request context.
func GetVerificationResult(ctx context.Context) *VerificationResult {
	result, _ := ctx.Value(ContextKeyVerificationResult).(*VerificationResult)
	return result
}

// GetThumbprint retrieves the DPoP thumbprint from the request context.
func GetThumbprint(ctx context.Context) string {
	thumbprint, _ := ctx.Value(ContextKeyThumbprint).(string)
	return thumbprint
}

// RequireDPoP is a helper function to create middleware that requires DPoP.
func RequireDPoP(verifier *Verifier) func(http.Handler) http.Handler {
	return Middleware(MiddlewareConfig{
		Verifier:    verifier,
		RequireDPoP: true,
	})
}

// OptionalDPoP is a helper function to create middleware that accepts but doesn't require DPoP.
func OptionalDPoP(verifier *Verifier) func(http.Handler) http.Handler {
	return Middleware(MiddlewareConfig{
		Verifier:    verifier,
		RequireDPoP: false,
	})
}
