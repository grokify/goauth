package dpop

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"
)

// VerificationConfig contains configuration for DPoP proof verification.
type VerificationConfig struct {
	// MaxAge is the maximum age of a DPoP proof (based on iat claim).
	// Default: 5 minutes.
	MaxAge time.Duration

	// AllowedClockSkew is the maximum clock skew to allow when validating iat.
	// Default: 30 seconds.
	AllowedClockSkew time.Duration

	// RequireAccessTokenBinding when true requires the ath claim to be present
	// and match the provided access token hash.
	RequireAccessTokenBinding bool

	// NonceValidator is an optional function to validate server-provided nonces.
	// If set and returns an error, verification fails with ErrNonceMismatch.
	NonceValidator func(ctx context.Context, nonce string) error
}

// DefaultVerificationConfig returns the default verification configuration.
func DefaultVerificationConfig() VerificationConfig {
	return VerificationConfig{
		MaxAge:           5 * time.Minute,
		AllowedClockSkew: 30 * time.Second,
	}
}

// Verifier verifies DPoP proofs.
type Verifier struct {
	config VerificationConfig
}

// NewVerifier creates a new DPoP verifier with the given configuration.
func NewVerifier(config VerificationConfig) *Verifier {
	if config.MaxAge == 0 {
		config.MaxAge = 5 * time.Minute
	}
	if config.AllowedClockSkew == 0 {
		config.AllowedClockSkew = 30 * time.Second
	}
	return &Verifier{config: config}
}

// VerificationRequest contains the parameters for verifying a DPoP proof.
type VerificationRequest struct {
	// Proof is the DPoP proof JWT string.
	Proof string
	// Method is the HTTP method of the request (e.g., "POST").
	Method string
	// URI is the HTTP URI of the request (scheme + host + path, no query or fragment).
	URI string
	// AccessToken is the access token to verify binding against (optional).
	// Required if VerificationConfig.RequireAccessTokenBinding is true.
	AccessToken string //nolint:gosec // G117: field holds runtime token value
	// ExpectedNonce is the expected server-provided nonce (optional).
	ExpectedNonce string
}

// VerificationResult contains the result of a successful verification.
type VerificationResult struct {
	// Thumbprint is the JWK thumbprint of the public key used to sign the proof.
	Thumbprint string
	// JTI is the unique identifier of the proof.
	JTI string
	// IssuedAt is when the proof was created.
	IssuedAt time.Time
}

// Verify verifies a DPoP proof against the given request parameters.
func (v *Verifier) Verify(ctx context.Context, req VerificationRequest) (*VerificationResult, error) {
	// Parse and validate signature
	parsed, err := ParseProof(req.Proof)
	if err != nil {
		return nil, err
	}

	// Validate HTTP method
	if !strings.EqualFold(parsed.Claims.HTTPMethod, req.Method) {
		return nil, fmt.Errorf("%w: proof htm=%s, request method=%s",
			ErrMethodMismatch, parsed.Claims.HTTPMethod, req.Method)
	}

	// Validate HTTP URI (normalize before comparison)
	if err := v.validateURI(parsed.Claims.HTTPURI, req.URI); err != nil {
		return nil, err
	}

	// Validate iat (issued at)
	if parsed.Claims.IssuedAt == nil {
		return nil, fmt.Errorf("%w: missing iat claim", ErrInvalidProof)
	}

	now := time.Now()
	iat := parsed.Claims.IssuedAt.Time

	// Check if proof is from the future (with clock skew allowance)
	if iat.After(now.Add(v.config.AllowedClockSkew)) {
		return nil, fmt.Errorf("%w: iat is in the future", ErrInvalidProof)
	}

	// Check if proof is too old
	if now.Sub(iat) > v.config.MaxAge+v.config.AllowedClockSkew {
		return nil, fmt.Errorf("%w: proof iat is %v old, max age is %v",
			ErrProofExpired, now.Sub(iat).Round(time.Second), v.config.MaxAge)
	}

	// Validate jti is present
	if parsed.Claims.ID == "" {
		return nil, fmt.Errorf("%w: missing jti claim", ErrInvalidProof)
	}

	// Validate access token binding if required or if ath is present
	if v.config.RequireAccessTokenBinding {
		if req.AccessToken == "" {
			return nil, fmt.Errorf("%w: access token binding required but no token provided", ErrInvalidProof)
		}
		if parsed.Claims.AccessTokenHash == "" {
			return nil, fmt.Errorf("%w: access token binding required but ath claim missing", ErrInvalidProof)
		}
	}

	if parsed.Claims.AccessTokenHash != "" && req.AccessToken != "" {
		expectedAth := ComputeAccessTokenHash(req.AccessToken)
		if parsed.Claims.AccessTokenHash != expectedAth {
			return nil, fmt.Errorf("%w: expected %s, got %s",
				ErrTokenHashMismatch, expectedAth, parsed.Claims.AccessTokenHash)
		}
	}

	// Validate nonce if validator is configured
	if v.config.NonceValidator != nil {
		if err := v.config.NonceValidator(ctx, parsed.Claims.Nonce); err != nil {
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				return nil, err
			}
			return nil, fmt.Errorf("%w: %v", ErrNonceMismatch, err)
		}
	}

	// Validate expected nonce if provided
	if req.ExpectedNonce != "" && parsed.Claims.Nonce != req.ExpectedNonce {
		return nil, fmt.Errorf("%w: expected %s, got %s",
			ErrNonceMismatch, req.ExpectedNonce, parsed.Claims.Nonce)
	}

	return &VerificationResult{
		Thumbprint: parsed.Thumbprint,
		JTI:        parsed.Claims.ID,
		IssuedAt:   iat,
	}, nil
}

// validateURI validates that the proof's htu matches the request URI.
// URIs are normalized before comparison (scheme and host are case-insensitive).
func (v *Verifier) validateURI(proofURI, requestURI string) error {
	// Parse both URIs
	proofURL, err := url.Parse(proofURI)
	if err != nil {
		return fmt.Errorf("%w: invalid htu: %v", ErrInvalidProof, err)
	}

	requestURL, err := url.Parse(requestURI)
	if err != nil {
		return fmt.Errorf("%w: invalid request URI: %v", ErrURIMismatch, err)
	}

	// Normalize: scheme and host are case-insensitive
	proofScheme := strings.ToLower(proofURL.Scheme)
	requestScheme := strings.ToLower(requestURL.Scheme)

	proofHost := strings.ToLower(proofURL.Host)
	requestHost := strings.ToLower(requestURL.Host)

	// Compare scheme
	if proofScheme != requestScheme {
		return fmt.Errorf("%w: scheme mismatch: proof=%s, request=%s",
			ErrURIMismatch, proofScheme, requestScheme)
	}

	// Compare host
	if proofHost != requestHost {
		return fmt.Errorf("%w: host mismatch: proof=%s, request=%s",
			ErrURIMismatch, proofHost, requestHost)
	}

	// Compare path (case-sensitive)
	// Normalize empty path to "/"
	proofPath := proofURL.Path
	requestPath := requestURL.Path
	if proofPath == "" {
		proofPath = "/"
	}
	if requestPath == "" {
		requestPath = "/"
	}

	if proofPath != requestPath {
		return fmt.Errorf("%w: path mismatch: proof=%s, request=%s",
			ErrURIMismatch, proofPath, requestPath)
	}

	return nil
}

// VerifyTokenBinding verifies that an access token's cnf.jkt claim matches
// the thumbprint from a verified DPoP proof.
func VerifyTokenBinding(tokenThumbprint string, proofThumbprint string) error {
	if tokenThumbprint == "" {
		return fmt.Errorf("%w: token has no cnf.jkt claim", ErrInvalidProof)
	}
	if proofThumbprint == "" {
		return fmt.Errorf("%w: proof has no thumbprint", ErrInvalidProof)
	}
	if tokenThumbprint != proofThumbprint {
		return fmt.Errorf("%w: token jkt=%s, proof jkt=%s",
			ErrTokenHashMismatch, tokenThumbprint, proofThumbprint)
	}
	return nil
}
