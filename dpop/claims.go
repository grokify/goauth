package dpop

import (
	"crypto/sha256"
	"encoding/base64"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// ProofClaims represents the claims for a DPoP proof JWT per RFC 9449.
type ProofClaims struct {
	jwt.RegisteredClaims

	// HTTPMethod is the HTTP method of the request (htm claim).
	// REQUIRED. The value of the HTTP method of the request to which the JWT is attached.
	HTTPMethod string `json:"htm"`

	// HTTPURI is the HTTP URI of the request (htu claim).
	// REQUIRED. The HTTP target URI, without query and fragment parts.
	HTTPURI string `json:"htu"`

	// AccessTokenHash is the base64url-encoded SHA-256 hash of the access token (ath claim).
	// REQUIRED when the DPoP proof is sent with a request for a protected resource.
	AccessTokenHash string `json:"ath,omitempty"`

	// Nonce is the server-provided nonce value (nonce claim).
	// Used for replay protection when the server issues nonces.
	Nonce string `json:"nonce,omitempty"`
}

// NewProofClaims creates a new DPoP proof claims structure.
// The method and uri parameters are required.
// The accessToken parameter is optional - if provided, the ath claim will be computed.
func NewProofClaims(method, uri string, accessToken string) *ProofClaims {
	now := time.Now()
	claims := &ProofClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt: jwt.NewNumericDate(now),
			ID:       uuid.NewString(), // jti - unique identifier
		},
		HTTPMethod: method,
		HTTPURI:    uri,
	}

	// If access token is provided, compute ath claim
	if accessToken != "" {
		claims.AccessTokenHash = ComputeAccessTokenHash(accessToken)
	}

	return claims
}

// WithNonce adds a server-provided nonce to the claims.
func (c *ProofClaims) WithNonce(nonce string) *ProofClaims {
	c.Nonce = nonce
	return c
}

// WithAccessToken computes and sets the ath claim from an access token.
func (c *ProofClaims) WithAccessToken(accessToken string) *ProofClaims {
	c.AccessTokenHash = ComputeAccessTokenHash(accessToken)
	return c
}

// ComputeAccessTokenHash computes the base64url-encoded SHA-256 hash of an access token.
// This is used for the ath (access token hash) claim in DPoP proofs.
func ComputeAccessTokenHash(accessToken string) string {
	hash := sha256.Sum256([]byte(accessToken))
	return base64.RawURLEncoding.EncodeToString(hash[:])
}

// ProofHeader represents the JWT header for a DPoP proof.
// The header MUST include the public key (jwk) and MUST have typ=dpop+jwt.
type ProofHeader struct {
	Algorithm string `json:"alg"` // Algorithm (ES256, ES384, ES512)
	Type      string `json:"typ"` // MUST be "dpop+jwt"
	JWK       *JWK   `json:"jwk"` // Public key
}

const (
	// DPoPTokenType is the required typ header value for DPoP proofs.
	DPoPTokenType = "dpop+jwt"
)

// ValidHTTPMethods contains the valid HTTP methods for DPoP proofs.
var ValidHTTPMethods = map[string]bool{
	"GET":     true,
	"HEAD":    true,
	"POST":    true,
	"PUT":     true,
	"DELETE":  true,
	"CONNECT": true,
	"OPTIONS": true,
	"TRACE":   true,
	"PATCH":   true,
}

// IsValidHTTPMethod checks if the given method is a valid HTTP method.
func IsValidHTTPMethod(method string) bool {
	return ValidHTTPMethods[method]
}
