package dpop

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestNewVerifier(t *testing.T) {
	v := NewVerifier(VerificationConfig{})

	if v.config.MaxAge != 5*time.Minute {
		t.Errorf("Default MaxAge = %v, want 5m", v.config.MaxAge)
	}
	if v.config.AllowedClockSkew != 30*time.Second {
		t.Errorf("Default AllowedClockSkew = %v, want 30s", v.config.AllowedClockSkew)
	}
}

func TestNewVerifier_CustomConfig(t *testing.T) {
	cfg := VerificationConfig{
		MaxAge:           10 * time.Minute,
		AllowedClockSkew: 1 * time.Minute,
	}
	v := NewVerifier(cfg)

	if v.config.MaxAge != 10*time.Minute {
		t.Errorf("MaxAge = %v, want 10m", v.config.MaxAge)
	}
	if v.config.AllowedClockSkew != 1*time.Minute {
		t.Errorf("AllowedClockSkew = %v, want 1m", v.config.AllowedClockSkew)
	}
}

func TestDefaultVerificationConfig(t *testing.T) {
	cfg := DefaultVerificationConfig()

	if cfg.MaxAge != 5*time.Minute {
		t.Errorf("MaxAge = %v, want 5m", cfg.MaxAge)
	}
	if cfg.AllowedClockSkew != 30*time.Second {
		t.Errorf("AllowedClockSkew = %v, want 30s", cfg.AllowedClockSkew)
	}
	if cfg.RequireAccessTokenBinding {
		t.Error("RequireAccessTokenBinding should be false by default")
	}
}

func TestVerifier_Verify_Success(t *testing.T) {
	kp, _ := GenerateKeyPair()
	v := NewVerifier(DefaultVerificationConfig())

	proof, err := CreateProof(kp, "POST", "https://api.example.com/resource")
	if err != nil {
		t.Fatalf("CreateProof() error: %v", err)
	}

	result, err := v.Verify(context.Background(), VerificationRequest{
		Proof:  proof,
		Method: "POST",
		URI:    "https://api.example.com/resource",
	})
	if err != nil {
		t.Fatalf("Verify() error: %v", err)
	}

	if result.Thumbprint != kp.Thumbprint {
		t.Errorf("Thumbprint = %s, want %s", result.Thumbprint, kp.Thumbprint)
	}
	if result.JTI == "" {
		t.Error("JTI is empty")
	}
	if result.IssuedAt.IsZero() {
		t.Error("IssuedAt is zero")
	}
}

func TestVerifier_Verify_MethodMismatch(t *testing.T) {
	kp, _ := GenerateKeyPair()
	v := NewVerifier(DefaultVerificationConfig())

	proof, _ := CreateProof(kp, "POST", "https://api.example.com/resource")

	_, err := v.Verify(context.Background(), VerificationRequest{
		Proof:  proof,
		Method: "GET", // Mismatch
		URI:    "https://api.example.com/resource",
	})
	if err == nil {
		t.Error("Verify() should fail with method mismatch")
	}
	if !errors.Is(err, ErrMethodMismatch) {
		t.Errorf("Error = %v, want ErrMethodMismatch", err)
	}
}

func TestVerifier_Verify_MethodCaseInsensitive(t *testing.T) {
	kp, _ := GenerateKeyPair()
	v := NewVerifier(DefaultVerificationConfig())

	proof, _ := CreateProof(kp, "POST", "https://api.example.com/resource")

	// Should accept lowercase method in request
	result, err := v.Verify(context.Background(), VerificationRequest{
		Proof:  proof,
		Method: "post", // lowercase
		URI:    "https://api.example.com/resource",
	})
	if err != nil {
		t.Errorf("Verify() should accept case-insensitive method: %v", err)
	}
	if result == nil {
		t.Error("Result should not be nil")
	}
}

func TestVerifier_Verify_URIMismatch(t *testing.T) {
	kp, _ := GenerateKeyPair()
	v := NewVerifier(DefaultVerificationConfig())

	proof, _ := CreateProof(kp, "GET", "https://api.example.com/resource")

	testCases := []struct {
		name string
		uri  string
	}{
		{"different path", "https://api.example.com/other"},
		{"different host", "https://other.example.com/resource"},
		{"different scheme", "http://api.example.com/resource"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := v.Verify(context.Background(), VerificationRequest{
				Proof:  proof,
				Method: "GET",
				URI:    tc.uri,
			})
			if err == nil {
				t.Errorf("Verify() should fail with %s", tc.name)
			}
			if !errors.Is(err, ErrURIMismatch) {
				t.Errorf("Error = %v, want ErrURIMismatch", err)
			}
		})
	}
}

func TestVerifier_Verify_URICaseNormalization(t *testing.T) {
	kp, _ := GenerateKeyPair()
	v := NewVerifier(DefaultVerificationConfig())

	// Create proof with uppercase scheme and host
	proof, _ := CreateProof(kp, "GET", "HTTPS://API.EXAMPLE.COM/resource")

	// Should match lowercase
	result, err := v.Verify(context.Background(), VerificationRequest{
		Proof:  proof,
		Method: "GET",
		URI:    "https://api.example.com/resource",
	})
	if err != nil {
		t.Errorf("Verify() should normalize URI case: %v", err)
	}
	if result == nil {
		t.Error("Result should not be nil")
	}
}

func TestVerifier_Verify_AccessTokenBinding(t *testing.T) {
	kp, _ := GenerateKeyPair()
	accessToken := "my-access-token"

	v := NewVerifier(DefaultVerificationConfig())

	proof, _ := CreateProofWithOptions(kp, "POST", "https://api.example.com", ProofOptions{
		AccessToken: accessToken,
	})

	result, err := v.Verify(context.Background(), VerificationRequest{
		Proof:       proof,
		Method:      "POST",
		URI:         "https://api.example.com",
		AccessToken: accessToken,
	})
	if err != nil {
		t.Fatalf("Verify() error: %v", err)
	}
	if result == nil {
		t.Error("Result should not be nil")
	}
}

func TestVerifier_Verify_AccessTokenMismatch(t *testing.T) {
	kp, _ := GenerateKeyPair()

	v := NewVerifier(DefaultVerificationConfig())

	proof, _ := CreateProofWithOptions(kp, "POST", "https://api.example.com", ProofOptions{
		AccessToken: "token-a",
	})

	_, err := v.Verify(context.Background(), VerificationRequest{
		Proof:       proof,
		Method:      "POST",
		URI:         "https://api.example.com",
		AccessToken: "token-b", // Mismatch
	})
	if err == nil {
		t.Error("Verify() should fail with token mismatch")
	}
	if !errors.Is(err, ErrTokenHashMismatch) {
		t.Errorf("Error = %v, want ErrTokenHashMismatch", err)
	}
}

func TestVerifier_Verify_RequireAccessTokenBinding(t *testing.T) {
	kp, _ := GenerateKeyPair()

	v := NewVerifier(VerificationConfig{
		RequireAccessTokenBinding: true,
	})

	// Proof without ath claim
	proof, _ := CreateProof(kp, "POST", "https://api.example.com")

	_, err := v.Verify(context.Background(), VerificationRequest{
		Proof:       proof,
		Method:      "POST",
		URI:         "https://api.example.com",
		AccessToken: "my-token",
	})
	if err == nil {
		t.Error("Verify() should fail when ath required but missing")
	}
}

func TestVerifier_Verify_RequireAccessTokenBinding_NoToken(t *testing.T) {
	kp, _ := GenerateKeyPair()

	v := NewVerifier(VerificationConfig{
		RequireAccessTokenBinding: true,
	})

	proof, _ := CreateProofWithOptions(kp, "POST", "https://api.example.com", ProofOptions{
		AccessToken: "token",
	})

	_, err := v.Verify(context.Background(), VerificationRequest{
		Proof:  proof,
		Method: "POST",
		URI:    "https://api.example.com",
		// No AccessToken provided
	})
	if err == nil {
		t.Error("Verify() should fail when binding required but no token provided")
	}
}

func TestVerifier_Verify_Nonce(t *testing.T) {
	kp, _ := GenerateKeyPair()

	v := NewVerifier(DefaultVerificationConfig())

	proof, _ := CreateProofWithOptions(kp, "GET", "https://api.example.com", ProofOptions{
		Nonce: "server-nonce-123",
	})

	result, err := v.Verify(context.Background(), VerificationRequest{
		Proof:         proof,
		Method:        "GET",
		URI:           "https://api.example.com",
		ExpectedNonce: "server-nonce-123",
	})
	if err != nil {
		t.Fatalf("Verify() error: %v", err)
	}
	if result == nil {
		t.Error("Result should not be nil")
	}
}

func TestVerifier_Verify_NonceMismatch(t *testing.T) {
	kp, _ := GenerateKeyPair()

	v := NewVerifier(DefaultVerificationConfig())

	proof, _ := CreateProofWithOptions(kp, "GET", "https://api.example.com", ProofOptions{
		Nonce: "nonce-a",
	})

	_, err := v.Verify(context.Background(), VerificationRequest{
		Proof:         proof,
		Method:        "GET",
		URI:           "https://api.example.com",
		ExpectedNonce: "nonce-b",
	})
	if err == nil {
		t.Error("Verify() should fail with nonce mismatch")
	}
	if !errors.Is(err, ErrNonceMismatch) {
		t.Errorf("Error = %v, want ErrNonceMismatch", err)
	}
}

func TestVerifier_Verify_NonceValidator(t *testing.T) {
	kp, _ := GenerateKeyPair()

	validNonces := map[string]bool{"valid-nonce": true}

	v := NewVerifier(VerificationConfig{
		NonceValidator: func(ctx context.Context, nonce string) error {
			if !validNonces[nonce] {
				return errors.New("invalid nonce")
			}
			return nil
		},
	})

	// Valid nonce
	proof, _ := CreateProofWithOptions(kp, "GET", "https://api.example.com", ProofOptions{
		Nonce: "valid-nonce",
	})

	result, err := v.Verify(context.Background(), VerificationRequest{
		Proof:  proof,
		Method: "GET",
		URI:    "https://api.example.com",
	})
	if err != nil {
		t.Fatalf("Verify() with valid nonce error: %v", err)
	}
	if result == nil {
		t.Error("Result should not be nil")
	}

	// Invalid nonce
	proof2, _ := CreateProofWithOptions(kp, "GET", "https://api.example.com", ProofOptions{
		Nonce: "invalid-nonce",
	})

	_, err = v.Verify(context.Background(), VerificationRequest{
		Proof:  proof2,
		Method: "GET",
		URI:    "https://api.example.com",
	})
	if err == nil {
		t.Error("Verify() should fail with invalid nonce")
	}
	if !errors.Is(err, ErrNonceMismatch) {
		t.Errorf("Error = %v, want ErrNonceMismatch", err)
	}
}

func TestVerifier_Verify_NonceValidator_ContextCanceled(t *testing.T) {
	kp, _ := GenerateKeyPair()

	v := NewVerifier(VerificationConfig{
		NonceValidator: func(ctx context.Context, nonce string) error {
			return context.Canceled
		},
	})

	proof, _ := CreateProofWithOptions(kp, "GET", "https://api.example.com", ProofOptions{
		Nonce: "nonce",
	})

	_, err := v.Verify(context.Background(), VerificationRequest{
		Proof:  proof,
		Method: "GET",
		URI:    "https://api.example.com",
	})
	if !errors.Is(err, context.Canceled) {
		t.Errorf("Error = %v, want context.Canceled", err)
	}
}

func TestVerifier_Verify_InvalidProof(t *testing.T) {
	v := NewVerifier(DefaultVerificationConfig())

	_, err := v.Verify(context.Background(), VerificationRequest{
		Proof:  "invalid-proof",
		Method: "GET",
		URI:    "https://api.example.com",
	})
	if err == nil {
		t.Error("Verify() should fail with invalid proof")
	}
}

func TestVerifier_Verify_URIEmptyPath(t *testing.T) {
	kp, _ := GenerateKeyPair()
	v := NewVerifier(DefaultVerificationConfig())

	// Proof with empty path (just host)
	proof, _ := CreateProof(kp, "GET", "https://api.example.com")

	// Request with "/" path should match
	result, err := v.Verify(context.Background(), VerificationRequest{
		Proof:  proof,
		Method: "GET",
		URI:    "https://api.example.com/",
	})
	if err != nil {
		t.Errorf("Verify() should normalize empty path to /: %v", err)
	}
	if result == nil {
		t.Error("Result should not be nil")
	}
}

func TestVerifyTokenBinding(t *testing.T) {
	// Success case
	err := VerifyTokenBinding("thumbprint-abc", "thumbprint-abc")
	if err != nil {
		t.Errorf("VerifyTokenBinding() with matching thumbprints error: %v", err)
	}

	// Mismatch case
	err = VerifyTokenBinding("thumbprint-a", "thumbprint-b")
	if err == nil {
		t.Error("VerifyTokenBinding() should fail with mismatched thumbprints")
	}
	if !errors.Is(err, ErrTokenHashMismatch) {
		t.Errorf("Error = %v, want ErrTokenHashMismatch", err)
	}

	// Empty token thumbprint
	err = VerifyTokenBinding("", "thumbprint")
	if err == nil {
		t.Error("VerifyTokenBinding() should fail with empty token thumbprint")
	}

	// Empty proof thumbprint
	err = VerifyTokenBinding("thumbprint", "")
	if err == nil {
		t.Error("VerifyTokenBinding() should fail with empty proof thumbprint")
	}
}

func TestVerifier_Verify_AllMethods(t *testing.T) {
	kp, _ := GenerateKeyPair()
	v := NewVerifier(DefaultVerificationConfig())

	methods := []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			proof, err := CreateProof(kp, method, "https://api.example.com/test")
			if err != nil {
				t.Fatalf("CreateProof() error: %v", err)
			}

			result, err := v.Verify(context.Background(), VerificationRequest{
				Proof:  proof,
				Method: method,
				URI:    "https://api.example.com/test",
			})
			if err != nil {
				t.Errorf("Verify() error: %v", err)
			}
			if result == nil {
				t.Error("Result should not be nil")
			}
		})
	}
}
