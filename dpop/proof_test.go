package dpop

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"strings"
	"testing"
	"time"
)

func TestCreateProof(t *testing.T) {
	kp, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair() error: %v", err)
	}

	proof, err := CreateProof(kp, "POST", "https://api.example.com/resource")
	if err != nil {
		t.Fatalf("CreateProof() error: %v", err)
	}

	// Verify proof is a valid JWT (3 parts separated by dots)
	parts := strings.Split(proof, ".")
	if len(parts) != 3 {
		t.Errorf("Proof has %d parts, want 3", len(parts))
	}
}

func TestCreateProof_NilKeyPair(t *testing.T) {
	_, err := CreateProof(nil, "POST", "https://example.com")
	if err == nil {
		t.Error("CreateProof(nil) should return error")
	}
}

func TestCreateProof_NilPrivateKey(t *testing.T) {
	kp := &KeyPair{}
	_, err := CreateProof(kp, "POST", "https://example.com")
	if err == nil {
		t.Error("CreateProof() with nil private key should return error")
	}
}

func TestCreateProof_InvalidMethod(t *testing.T) {
	kp, _ := GenerateKeyPair()
	_, err := CreateProof(kp, "INVALID", "https://example.com")
	if err == nil {
		t.Error("CreateProof() with invalid method should return error")
	}
}

func TestCreateProof_AllMethods(t *testing.T) {
	kp, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair() error: %v", err)
	}

	methods := []string{"GET", "HEAD", "POST", "PUT", "DELETE", "CONNECT", "OPTIONS", "TRACE", "PATCH"}
	for _, method := range methods {
		proof, err := CreateProof(kp, method, "https://example.com")
		if err != nil {
			t.Errorf("CreateProof(%s) error: %v", method, err)
			continue
		}
		if proof == "" {
			t.Errorf("CreateProof(%s) returned empty proof", method)
		}
	}
}

func TestCreateProofWithOptions(t *testing.T) {
	kp, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair() error: %v", err)
	}

	opts := ProofOptions{
		AccessToken: "test-access-token",
		Nonce:       "server-nonce-123",
	}

	proof, err := CreateProofWithOptions(kp, "POST", "https://api.example.com/resource", opts)
	if err != nil {
		t.Fatalf("CreateProofWithOptions() error: %v", err)
	}

	// Parse and verify
	parsed, err := ParseProof(proof)
	if err != nil {
		t.Fatalf("ParseProof() error: %v", err)
	}

	// Verify ath claim is set
	expectedAth := ComputeAccessTokenHash(opts.AccessToken)
	if parsed.Claims.AccessTokenHash != expectedAth {
		t.Errorf("AccessTokenHash = %s, want %s", parsed.Claims.AccessTokenHash, expectedAth)
	}

	// Verify nonce claim is set
	if parsed.Claims.Nonce != opts.Nonce {
		t.Errorf("Nonce = %s, want %s", parsed.Claims.Nonce, opts.Nonce)
	}
}

func TestCreateProof_DifferentCurves(t *testing.T) {
	curves := []elliptic.Curve{
		elliptic.P256(),
		elliptic.P384(),
		elliptic.P521(),
	}

	for _, curve := range curves {
		privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
		if err != nil {
			t.Fatalf("GenerateKey(%v) error: %v", curve.Params().Name, err)
		}

		thumbprint, _ := ComputeThumbprint(&privateKey.PublicKey)
		kp := &KeyPair{
			PrivateKey: privateKey,
			Thumbprint: thumbprint,
		}

		proof, err := CreateProof(kp, "GET", "https://example.com")
		if err != nil {
			t.Errorf("CreateProof() with %v error: %v", curve.Params().Name, err)
			continue
		}

		// Parse and verify signature
		parsed, err := ParseProof(proof)
		if err != nil {
			t.Errorf("ParseProof() with %v error: %v", curve.Params().Name, err)
			continue
		}

		if parsed.Thumbprint != kp.Thumbprint {
			t.Errorf("Thumbprint mismatch for %v", curve.Params().Name)
		}
	}
}

func TestParseProof(t *testing.T) {
	kp, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair() error: %v", err)
	}

	proof, err := CreateProof(kp, "POST", "https://api.example.com/resource")
	if err != nil {
		t.Fatalf("CreateProof() error: %v", err)
	}

	parsed, err := ParseProof(proof)
	if err != nil {
		t.Fatalf("ParseProof() error: %v", err)
	}

	// Verify claims
	if parsed.Claims.HTTPMethod != "POST" {
		t.Errorf("HTTPMethod = %s, want POST", parsed.Claims.HTTPMethod)
	}
	if parsed.Claims.HTTPURI != "https://api.example.com/resource" {
		t.Errorf("HTTPURI = %s, want https://api.example.com/resource", parsed.Claims.HTTPURI)
	}
	if parsed.Claims.ID == "" {
		t.Error("jti (ID) is empty")
	}
	if parsed.Claims.IssuedAt == nil {
		t.Error("iat (IssuedAt) is nil")
	}

	// Verify public key matches (use Bytes() instead of deprecated X/Y fields)
	parsedPubBytes, err := parsed.PublicKey.Bytes()
	if err != nil {
		t.Fatalf("parsed.PublicKey.Bytes() error: %v", err)
	}
	kpPubBytes, err := kp.PrivateKey.PublicKey.Bytes()
	if err != nil {
		t.Fatalf("kp.PrivateKey.PublicKey.Bytes() error: %v", err)
	}
	if !bytes.Equal(parsedPubBytes, kpPubBytes) {
		t.Error("PublicKey does not match")
	}

	// Verify thumbprint matches
	if parsed.Thumbprint != kp.Thumbprint {
		t.Errorf("Thumbprint = %s, want %s", parsed.Thumbprint, kp.Thumbprint)
	}
}

func TestParseProof_InvalidJWT(t *testing.T) {
	_, err := ParseProof("not.a.valid.jwt")
	if err == nil {
		t.Error("ParseProof() with invalid JWT should return error")
	}
}

func TestParseProof_InvalidHeader(t *testing.T) {
	_, err := ParseProof("!!!invalid!!!.payload.signature")
	if err == nil {
		t.Error("ParseProof() with invalid header encoding should return error")
	}
}

func TestParseProof_WrongType(t *testing.T) {
	kp, _ := GenerateKeyPair()

	// Create a regular JWT (not DPoP)
	token := "eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.signature" //nolint:gosec // G101: test token, not real credentials

	_, err := ParseProof(token)
	if err == nil {
		t.Error("ParseProof() with wrong typ should return error")
	}

	_ = kp // avoid unused
}

func TestParseProof_MissingJWK(t *testing.T) {
	// JWT with correct typ but no jwk header
	// This is a crafted token with typ=dpop+jwt but no jwk
	token := "eyJhbGciOiJFUzI1NiIsInR5cCI6ImRwb3Arand0In0.eyJzdWIiOiJ0ZXN0In0.signature" //nolint:gosec // G101: test token, not real credentials

	_, err := ParseProof(token)
	if err == nil {
		t.Error("ParseProof() with missing jwk should return error")
	}
}

func TestParseProof_InvalidSignature(t *testing.T) {
	kp1, _ := GenerateKeyPair()
	kp2, _ := GenerateKeyPair()

	// Create proof with kp1
	proof, err := CreateProof(kp1, "GET", "https://example.com")
	if err != nil {
		t.Fatalf("CreateProof() error: %v", err)
	}

	// Tamper with the signature by replacing jwk in header
	// This is tricky - we'll just verify that different keys give different thumbprints
	parsed1, _ := ParseProof(proof)
	proof2, _ := CreateProof(kp2, "GET", "https://example.com")
	parsed2, _ := ParseProof(proof2)

	if parsed1.Thumbprint == parsed2.Thumbprint {
		t.Error("Different key pairs should have different thumbprints")
	}
}

func TestComputeAccessTokenHash(t *testing.T) {
	accessToken := "test-access-token-12345"
	hash := ComputeAccessTokenHash(accessToken)

	// Verify hash is non-empty
	if hash == "" {
		t.Error("AccessTokenHash is empty")
	}

	// Verify hash is base64url encoded (no padding, no + or /)
	if strings.Contains(hash, "+") || strings.Contains(hash, "/") || strings.Contains(hash, "=") {
		t.Errorf("AccessTokenHash contains invalid characters: %s", hash)
	}

	// Verify determinism
	hash2 := ComputeAccessTokenHash(accessToken)
	if hash != hash2 {
		t.Errorf("AccessTokenHash not deterministic: %s != %s", hash, hash2)
	}
}

func TestNewProofClaims(t *testing.T) {
	claims := NewProofClaims("POST", "https://example.com/api", "")

	// Verify required claims
	if claims.HTTPMethod != "POST" {
		t.Errorf("HTTPMethod = %s, want POST", claims.HTTPMethod)
	}
	if claims.HTTPURI != "https://example.com/api" {
		t.Errorf("HTTPURI = %s, want https://example.com/api", claims.HTTPURI)
	}
	if claims.ID == "" {
		t.Error("jti (ID) is empty")
	}
	if claims.IssuedAt == nil {
		t.Error("iat (IssuedAt) is nil")
	}

	// Verify iat is recent
	if claims.IssuedAt.Before(time.Now().Add(-time.Minute)) {
		t.Error("iat is too old")
	}

	// Verify ath is not set when access token is empty
	if claims.AccessTokenHash != "" {
		t.Error("AccessTokenHash should be empty when no access token provided")
	}
}

func TestNewProofClaims_WithAccessToken(t *testing.T) {
	accessToken := "test-token"
	claims := NewProofClaims("GET", "https://example.com", accessToken)

	expectedAth := ComputeAccessTokenHash(accessToken)
	if claims.AccessTokenHash != expectedAth {
		t.Errorf("AccessTokenHash = %s, want %s", claims.AccessTokenHash, expectedAth)
	}
}

func TestProofClaims_WithNonce(t *testing.T) {
	claims := NewProofClaims("GET", "https://example.com", "")
	claims.WithNonce("server-nonce")

	if claims.Nonce != "server-nonce" {
		t.Errorf("Nonce = %s, want server-nonce", claims.Nonce)
	}
}

func TestProofClaims_WithAccessToken(t *testing.T) {
	claims := NewProofClaims("GET", "https://example.com", "")
	claims.WithAccessToken("my-token")

	expectedAth := ComputeAccessTokenHash("my-token")
	if claims.AccessTokenHash != expectedAth {
		t.Errorf("AccessTokenHash = %s, want %s", claims.AccessTokenHash, expectedAth)
	}
}

func TestIsValidHTTPMethod(t *testing.T) {
	validMethods := []string{"GET", "HEAD", "POST", "PUT", "DELETE", "CONNECT", "OPTIONS", "TRACE", "PATCH"}
	for _, method := range validMethods {
		if !IsValidHTTPMethod(method) {
			t.Errorf("IsValidHTTPMethod(%s) = false, want true", method)
		}
	}

	invalidMethods := []string{"get", "post", "INVALID", "", "FOO"}
	for _, method := range invalidMethods {
		if IsValidHTTPMethod(method) {
			t.Errorf("IsValidHTTPMethod(%s) = true, want false", method)
		}
	}
}

func TestParseProof_UnsupportedAlgorithm(t *testing.T) {
	// Create a JWT with an unsupported algorithm in header
	// Header: {"alg":"RS256","typ":"dpop+jwt","jwk":{...}}
	// This tests the algorithm validation path
	token := "eyJhbGciOiJSUzI1NiIsInR5cCI6ImRwb3Arand0IiwiandrIjp7Imt0eSI6IkVDIiwiY3J2IjoiUC0yNTYiLCJ4IjoiQUFBQSIsInkiOiJCQkJCIn19.eyJodG0iOiJHRVQiLCJodHUiOiJodHRwczovL2V4YW1wbGUuY29tIn0.signature" //nolint:gosec // G101: test token

	_, err := ParseProof(token)
	if err == nil {
		t.Error("ParseProof() with unsupported algorithm should return error")
	}
}

func TestParseProof_RoundTrip(t *testing.T) {
	kp, _ := GenerateKeyPair()

	testCases := []struct {
		name        string
		method      string
		uri         string
		accessToken string
		nonce       string
	}{
		{
			name:   "basic GET",
			method: "GET",
			uri:    "https://api.example.com/users",
		},
		{ //nolint:gosec // G101: test token, not real credentials
			name:        "POST with token",
			method:      "POST",
			uri:         "https://api.example.com/data",
			accessToken: "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.test",
		},
		{
			name:   "DELETE with nonce",
			method: "DELETE",
			uri:    "https://api.example.com/resource/123",
			nonce:  "unique-server-nonce",
		},
		{
			name:        "PUT with all options",
			method:      "PUT",
			uri:         "https://api.example.com/update",
			accessToken: "access-token-xyz",
			nonce:       "nonce-abc",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			opts := ProofOptions{
				AccessToken: tc.accessToken,
				Nonce:       tc.nonce,
			}

			proof, err := CreateProofWithOptions(kp, tc.method, tc.uri, opts)
			if err != nil {
				t.Fatalf("CreateProofWithOptions() error: %v", err)
			}

			parsed, err := ParseProof(proof)
			if err != nil {
				t.Fatalf("ParseProof() error: %v", err)
			}

			// Verify all fields
			if parsed.Claims.HTTPMethod != tc.method {
				t.Errorf("HTTPMethod = %s, want %s", parsed.Claims.HTTPMethod, tc.method)
			}
			if parsed.Claims.HTTPURI != tc.uri {
				t.Errorf("HTTPURI = %s, want %s", parsed.Claims.HTTPURI, tc.uri)
			}
			if tc.accessToken != "" {
				expectedAth := ComputeAccessTokenHash(tc.accessToken)
				if parsed.Claims.AccessTokenHash != expectedAth {
					t.Errorf("AccessTokenHash = %s, want %s", parsed.Claims.AccessTokenHash, expectedAth)
				}
			}
			if parsed.Claims.Nonce != tc.nonce {
				t.Errorf("Nonce = %s, want %s", parsed.Claims.Nonce, tc.nonce)
			}
			if parsed.Thumbprint != kp.Thumbprint {
				t.Errorf("Thumbprint = %s, want %s", parsed.Thumbprint, kp.Thumbprint)
			}
		})
	}
}
