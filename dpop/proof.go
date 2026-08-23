package dpop

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

var (
	// ErrInvalidProof is returned when a DPoP proof is malformed or invalid.
	ErrInvalidProof = errors.New("invalid DPoP proof")
	// ErrMethodMismatch is returned when the htm claim doesn't match the request method.
	ErrMethodMismatch = errors.New("HTTP method mismatch")
	// ErrURIMismatch is returned when the htu claim doesn't match the request URI.
	ErrURIMismatch = errors.New("HTTP URI mismatch")
	// ErrTokenHashMismatch is returned when the ath claim doesn't match the access token hash.
	ErrTokenHashMismatch = errors.New("access token hash mismatch")
	// ErrProofExpired is returned when a DPoP proof's iat is too old.
	ErrProofExpired = errors.New("DPoP proof expired")
	// ErrNonceMismatch is returned when the nonce claim doesn't match the expected nonce.
	ErrNonceMismatch = errors.New("nonce mismatch")
)

// CreateProof creates a DPoP proof JWT for the given request.
// The proof is signed with the key pair's private key and embeds the public key in the header.
func CreateProof(kp *KeyPair, method, uri string) (string, error) {
	return CreateProofWithOptions(kp, method, uri, ProofOptions{})
}

// ProofOptions contains optional parameters for creating a DPoP proof.
type ProofOptions struct {
	// AccessToken is the access token to bind to the proof (for ath claim).
	AccessToken string //nolint:gosec // G117: field holds runtime token value
	// Nonce is a server-provided nonce for replay protection.
	Nonce string
}

// CreateProofWithOptions creates a DPoP proof JWT with optional parameters.
func CreateProofWithOptions(kp *KeyPair, method, uri string, opts ProofOptions) (string, error) {
	if kp == nil || kp.PrivateKey == nil {
		return "", ErrInvalidKey
	}

	if !IsValidHTTPMethod(method) {
		return "", fmt.Errorf("%w: invalid HTTP method %s", ErrInvalidProof, method)
	}

	// Create claims
	claims := NewProofClaims(method, uri, opts.AccessToken)
	if opts.Nonce != "" {
		claims.WithNonce(opts.Nonce)
	}

	// Create JWK for header
	jwk, err := ToJWK(kp.PublicKey())
	if err != nil {
		return "", fmt.Errorf("creating JWK: %w", err)
	}

	// Determine signing method based on curve
	var signingMethod jwt.SigningMethod
	switch kp.PrivateKey.Curve {
	case elliptic.P256():
		signingMethod = jwt.SigningMethodES256
	case elliptic.P384():
		signingMethod = jwt.SigningMethodES384
	case elliptic.P521():
		signingMethod = jwt.SigningMethodES512
	default:
		return "", fmt.Errorf("%w: unsupported curve", ErrUnsupportedAlgorithm)
	}

	// Create token with custom header
	token := jwt.NewWithClaims(signingMethod, claims)

	// Set custom header fields
	token.Header["typ"] = DPoPTokenType
	token.Header["jwk"] = jwk

	// Sign the token
	return token.SignedString(kp.PrivateKey)
}

// ParsedProof contains the parsed and validated DPoP proof.
type ParsedProof struct {
	// Claims contains the proof claims.
	Claims *ProofClaims
	// PublicKey is the public key extracted from the jwk header.
	PublicKey *ecdsa.PublicKey
	// Thumbprint is the JWK thumbprint of the public key.
	Thumbprint string
}

// ParseProof parses a DPoP proof JWT without verifying it against request parameters.
// Use this for initial parsing before validation.
func ParseProof(proofString string) (*ParsedProof, error) {
	// Split the JWT to get header
	parts := strings.Split(proofString, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("%w: invalid JWT format", ErrInvalidProof)
	}

	// Decode header
	headerJSON, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, fmt.Errorf("%w: invalid header encoding: %v", ErrInvalidProof, err)
	}

	var header struct {
		Alg string `json:"alg"`
		Typ string `json:"typ"`
		JWK *JWK   `json:"jwk"`
	}
	if err := json.Unmarshal(headerJSON, &header); err != nil {
		return nil, fmt.Errorf("%w: invalid header JSON: %v", ErrInvalidProof, err)
	}

	// Validate typ header
	if header.Typ != DPoPTokenType {
		return nil, fmt.Errorf("%w: typ must be %s, got %s", ErrInvalidProof, DPoPTokenType, header.Typ)
	}

	// Validate jwk header is present
	if header.JWK == nil {
		return nil, fmt.Errorf("%w: missing jwk header", ErrInvalidProof)
	}

	// Convert JWK to public key
	publicKey, err := header.JWK.ToPublicKey()
	if err != nil {
		return nil, fmt.Errorf("%w: invalid jwk: %v", ErrInvalidProof, err)
	}

	// Compute thumbprint
	thumbprint, err := header.JWK.Thumbprint()
	if err != nil {
		return nil, fmt.Errorf("%w: failed to compute thumbprint: %v", ErrInvalidProof, err)
	}

	// Determine signing method from alg header
	var signingMethod jwt.SigningMethod
	switch header.Alg {
	case "ES256":
		signingMethod = jwt.SigningMethodES256
	case "ES384":
		signingMethod = jwt.SigningMethodES384
	case "ES512":
		signingMethod = jwt.SigningMethodES512
	default:
		return nil, fmt.Errorf("%w: unsupported algorithm %s", ErrInvalidProof, header.Alg)
	}

	// Parse and verify the token signature
	claims := &ProofClaims{}
	token, err := jwt.ParseWithClaims(proofString, claims, func(token *jwt.Token) (any, error) {
		// Verify the signing method matches
		if token.Method.Alg() != signingMethod.Alg() {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("%w: signature verification failed: %v", ErrInvalidProof, err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("%w: token not valid", ErrInvalidProof)
	}

	return &ParsedProof{
		Claims:     claims,
		PublicKey:  publicKey,
		Thumbprint: thumbprint,
	}, nil
}
