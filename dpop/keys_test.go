package dpop

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/json"
	"strings"
	"testing"
)

func TestGenerateKeyPair(t *testing.T) {
	kp, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair() error: %v", err)
	}

	// Verify private key is set
	if kp.PrivateKey == nil {
		t.Error("PrivateKey is nil")
	}

	// Verify curve is P-256
	if kp.PrivateKey.Curve != elliptic.P256() {
		t.Errorf("Curve = %v, want P-256", kp.PrivateKey.Curve)
	}

	// Verify thumbprint is set and non-empty
	if kp.Thumbprint == "" {
		t.Error("Thumbprint is empty")
	}

	// Verify thumbprint is base64url encoded (no padding, no + or /)
	if strings.Contains(kp.Thumbprint, "+") || strings.Contains(kp.Thumbprint, "/") {
		t.Errorf("Thumbprint contains invalid base64url characters: %s", kp.Thumbprint)
	}
	if strings.Contains(kp.Thumbprint, "=") {
		t.Errorf("Thumbprint contains padding: %s", kp.Thumbprint)
	}

	// SHA-256 hash base64url encoded should be 43 characters
	if len(kp.Thumbprint) != 43 {
		t.Errorf("Thumbprint length = %d, want 43 (SHA-256 base64url)", len(kp.Thumbprint))
	}
}

func TestGenerateKeyPair_Unique(t *testing.T) {
	kp1, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair() 1 error: %v", err)
	}

	kp2, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair() 2 error: %v", err)
	}

	// Verify different key pairs have different thumbprints
	if kp1.Thumbprint == kp2.Thumbprint {
		t.Error("Two generated key pairs have the same thumbprint")
	}

	// Verify different private keys (use Bytes() instead of deprecated D field)
	priv1, err := kp1.PrivateKey.Bytes()
	if err != nil {
		t.Fatalf("kp1.PrivateKey.Bytes() error: %v", err)
	}
	priv2, err := kp2.PrivateKey.Bytes()
	if err != nil {
		t.Fatalf("kp2.PrivateKey.Bytes() error: %v", err)
	}
	if bytes.Equal(priv1, priv2) {
		t.Error("Two generated key pairs have the same private key")
	}
}

func TestKeyPair_PublicKey(t *testing.T) {
	kp, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair() error: %v", err)
	}

	pubKey := kp.PublicKey()
	if pubKey == nil {
		t.Fatal("PublicKey() is nil")
	}

	// Verify it's the same as the embedded public key (use Bytes() instead of deprecated X/Y fields)
	pubBytes1, err := pubKey.Bytes()
	if err != nil {
		t.Fatalf("pubKey.Bytes() error: %v", err)
	}
	pubBytes2, err := kp.PrivateKey.PublicKey.Bytes()
	if err != nil {
		t.Fatalf("kp.PrivateKey.PublicKey.Bytes() error: %v", err)
	}
	if !bytes.Equal(pubBytes1, pubBytes2) {
		t.Error("PublicKey() does not match embedded public key")
	}
}

func TestKeyPair_PublicKey_Nil(t *testing.T) {
	kp := &KeyPair{}
	if kp.PublicKey() != nil {
		t.Error("PublicKey() should be nil for empty KeyPair")
	}
}

func TestComputeThumbprint(t *testing.T) {
	// Generate a key pair
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("GenerateKey() error: %v", err)
	}

	thumbprint, err := ComputeThumbprint(&privateKey.PublicKey)
	if err != nil {
		t.Fatalf("ComputeThumbprint() error: %v", err)
	}

	// Verify thumbprint is base64url encoded
	if thumbprint == "" {
		t.Error("Thumbprint is empty")
	}

	// Computing again should give the same result (deterministic)
	thumbprint2, err := ComputeThumbprint(&privateKey.PublicKey)
	if err != nil {
		t.Fatalf("ComputeThumbprint() second call error: %v", err)
	}

	if thumbprint != thumbprint2 {
		t.Errorf("Thumbprint not deterministic: %s != %s", thumbprint, thumbprint2)
	}
}

func TestComputeThumbprint_NilKey(t *testing.T) {
	_, err := ComputeThumbprint(nil)
	if err == nil {
		t.Error("ComputeThumbprint(nil) should return error")
	}
	if err != ErrInvalidKey {
		t.Errorf("Error = %v, want ErrInvalidKey", err)
	}
}

func TestComputeThumbprint_DifferentCurves(t *testing.T) {
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

		thumbprint, err := ComputeThumbprint(&privateKey.PublicKey)
		if err != nil {
			t.Errorf("ComputeThumbprint(%v) error: %v", curve.Params().Name, err)
			continue
		}

		if thumbprint == "" {
			t.Errorf("ComputeThumbprint(%v) returned empty thumbprint", curve.Params().Name)
		}
	}
}

func TestToJWK(t *testing.T) {
	kp, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair() error: %v", err)
	}

	jwk, err := ToJWK(kp.PublicKey())
	if err != nil {
		t.Fatalf("ToJWK() error: %v", err)
	}

	// Verify JWK fields
	if jwk.Kty != "EC" {
		t.Errorf("Kty = %s, want EC", jwk.Kty)
	}
	if jwk.Crv != "P-256" {
		t.Errorf("Crv = %s, want P-256", jwk.Crv)
	}
	if jwk.Alg != "ES256" {
		t.Errorf("Alg = %s, want ES256", jwk.Alg)
	}
	if jwk.X == "" {
		t.Error("X is empty")
	}
	if jwk.Y == "" {
		t.Error("Y is empty")
	}
}

func TestToJWK_NilKey(t *testing.T) {
	_, err := ToJWK(nil)
	if err == nil {
		t.Error("ToJWK(nil) should return error")
	}
}

func TestJWK_Thumbprint(t *testing.T) {
	kp, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair() error: %v", err)
	}

	jwk, err := ToJWK(kp.PublicKey())
	if err != nil {
		t.Fatalf("ToJWK() error: %v", err)
	}

	thumbprint, err := jwk.Thumbprint()
	if err != nil {
		t.Fatalf("JWK.Thumbprint() error: %v", err)
	}

	// Should match the key pair's thumbprint
	if thumbprint != kp.Thumbprint {
		t.Errorf("JWK thumbprint = %s, want %s", thumbprint, kp.Thumbprint)
	}
}

func TestJWK_ToPublicKey(t *testing.T) {
	kp, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair() error: %v", err)
	}

	jwk, err := ToJWK(kp.PublicKey())
	if err != nil {
		t.Fatalf("ToJWK() error: %v", err)
	}

	// Convert back to public key
	pubKey, err := jwk.ToPublicKey()
	if err != nil {
		t.Fatalf("JWK.ToPublicKey() error: %v", err)
	}

	// Verify coordinates match (use Bytes() instead of deprecated X/Y fields)
	pubBytes1, err := pubKey.Bytes()
	if err != nil {
		t.Fatalf("pubKey.Bytes() error: %v", err)
	}
	pubBytes2, err := kp.PrivateKey.PublicKey.Bytes()
	if err != nil {
		t.Fatalf("kp.PrivateKey.PublicKey.Bytes() error: %v", err)
	}
	if !bytes.Equal(pubBytes1, pubBytes2) {
		t.Error("Public key mismatch after JWK round-trip")
	}
}

func TestJWK_ToPublicKey_InvalidCurve(t *testing.T) {
	jwk := &JWK{
		Kty: "EC",
		Crv: "unknown",
		X:   "AAAA",
		Y:   "BBBB",
	}

	_, err := jwk.ToPublicKey()
	if err == nil {
		t.Error("ToPublicKey() with invalid curve should return error")
	}
}

func TestJWK_ToPublicKey_InvalidKty(t *testing.T) {
	jwk := &JWK{
		Kty: "RSA",
		Crv: "P-256",
		X:   "AAAA",
		Y:   "BBBB",
	}

	_, err := jwk.ToPublicKey()
	if err == nil {
		t.Error("ToPublicKey() with non-EC key type should return error")
	}
}

func TestKeyPair_Serialize(t *testing.T) {
	kp, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair() error: %v", err)
	}

	serialized, err := kp.Serialize()
	if err != nil {
		t.Fatalf("Serialize() error: %v", err)
	}

	// Verify all fields are set
	if serialized.PrivateKeyD == "" {
		t.Error("PrivateKeyD is empty")
	}
	if serialized.PublicKeyX == "" {
		t.Error("PublicKeyX is empty")
	}
	if serialized.PublicKeyY == "" {
		t.Error("PublicKeyY is empty")
	}
	if serialized.Curve != "P-256" {
		t.Errorf("Curve = %s, want P-256", serialized.Curve)
	}
	if serialized.Thumbprint != kp.Thumbprint {
		t.Errorf("Thumbprint = %s, want %s", serialized.Thumbprint, kp.Thumbprint)
	}
}

func TestKeyPair_SerializeJSON(t *testing.T) {
	kp, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair() error: %v", err)
	}

	data, err := kp.SerializeJSON()
	if err != nil {
		t.Fatalf("SerializeJSON() error: %v", err)
	}

	// Verify it's valid JSON
	var s SerializedKeyPair
	if err := json.Unmarshal(data, &s); err != nil {
		t.Fatalf("Invalid JSON: %v", err)
	}

	// Verify roundtrip
	if s.Thumbprint != kp.Thumbprint {
		t.Errorf("Thumbprint after JSON roundtrip = %s, want %s", s.Thumbprint, kp.Thumbprint)
	}
}

func TestDeserializeKeyPair(t *testing.T) {
	// Generate original key pair
	original, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair() error: %v", err)
	}

	// Serialize
	serialized, err := original.Serialize()
	if err != nil {
		t.Fatalf("Serialize() error: %v", err)
	}

	// Deserialize
	restored, err := DeserializeKeyPair(serialized)
	if err != nil {
		t.Fatalf("DeserializeKeyPair() error: %v", err)
	}

	// Verify key matches (use Bytes() instead of deprecated D/X/Y fields)
	restoredPriv, err := restored.PrivateKey.Bytes()
	if err != nil {
		t.Fatalf("restored.PrivateKey.Bytes() error: %v", err)
	}
	originalPriv, err := original.PrivateKey.Bytes()
	if err != nil {
		t.Fatalf("original.PrivateKey.Bytes() error: %v", err)
	}
	if !bytes.Equal(restoredPriv, originalPriv) {
		t.Error("Private key does not match")
	}

	restoredPub, err := restored.PrivateKey.PublicKey.Bytes()
	if err != nil {
		t.Fatalf("restored.PrivateKey.PublicKey.Bytes() error: %v", err)
	}
	originalPub, err := original.PrivateKey.PublicKey.Bytes()
	if err != nil {
		t.Fatalf("original.PrivateKey.PublicKey.Bytes() error: %v", err)
	}
	if !bytes.Equal(restoredPub, originalPub) {
		t.Error("Public key does not match")
	}

	// Verify thumbprint matches
	if restored.Thumbprint != original.Thumbprint {
		t.Errorf("Thumbprint = %s, want %s", restored.Thumbprint, original.Thumbprint)
	}
}

func TestDeserializeKeyPair_Nil(t *testing.T) {
	_, err := DeserializeKeyPair(nil)
	if err == nil {
		t.Error("DeserializeKeyPair(nil) should return error")
	}
}

func TestDeserializeKeyPair_ThumbprintMismatch(t *testing.T) {
	kp, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair() error: %v", err)
	}

	serialized, err := kp.Serialize()
	if err != nil {
		t.Fatalf("Serialize() error: %v", err)
	}

	// Corrupt the thumbprint
	serialized.Thumbprint = "corrupted-thumbprint"

	_, err = DeserializeKeyPair(serialized)
	if err == nil {
		t.Error("DeserializeKeyPair() with mismatched thumbprint should return error")
	}
}

func TestDeserializeKeyPairJSON(t *testing.T) {
	// Generate original key pair
	original, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair() error: %v", err)
	}

	// Serialize to JSON
	data, err := original.SerializeJSON()
	if err != nil {
		t.Fatalf("SerializeJSON() error: %v", err)
	}

	// Deserialize from JSON
	restored, err := DeserializeKeyPairJSON(data)
	if err != nil {
		t.Fatalf("DeserializeKeyPairJSON() error: %v", err)
	}

	// Verify thumbprint matches
	if restored.Thumbprint != original.Thumbprint {
		t.Errorf("Thumbprint = %s, want %s", restored.Thumbprint, original.Thumbprint)
	}
}

func TestDeserializeKeyPairJSON_InvalidJSON(t *testing.T) {
	_, err := DeserializeKeyPairJSON([]byte("not valid json"))
	if err == nil {
		t.Error("DeserializeKeyPairJSON() with invalid JSON should return error")
	}
}

func TestKeyPair_Signer(t *testing.T) {
	kp, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair() error: %v", err)
	}

	signer := kp.Signer()
	if signer == nil {
		t.Error("Signer() is nil")
	}

	// Verify it's the same key (use Bytes() instead of deprecated D field)
	ecdsaSigner, ok := signer.(*ecdsa.PrivateKey)
	if !ok {
		t.Fatal("Signer() is not *ecdsa.PrivateKey")
	}
	signerBytes, err := ecdsaSigner.Bytes()
	if err != nil {
		t.Fatalf("ecdsaSigner.Bytes() error: %v", err)
	}
	kpBytes, err := kp.PrivateKey.Bytes()
	if err != nil {
		t.Fatalf("kp.PrivateKey.Bytes() error: %v", err)
	}
	if !bytes.Equal(signerBytes, kpBytes) {
		t.Error("Signer() is not the same key")
	}
}

func TestBase64URLEncode(t *testing.T) {
	tests := []struct {
		input    []byte
		expected string
	}{
		{[]byte{}, ""},
		{[]byte{0x00}, "AA"},
		{[]byte{0xff}, "_w"},
		{[]byte{0xfb, 0xff}, "-_8"},
	}

	for _, tt := range tests {
		result := base64URLEncode(tt.input)
		if result != tt.expected {
			t.Errorf("base64URLEncode(%v) = %s, want %s", tt.input, result, tt.expected)
		}
	}
}

func TestBase64URLDecode(t *testing.T) {
	tests := []struct {
		input    string
		expected []byte
	}{
		{"", []byte{}},
		{"AA", []byte{0x00}},
		{"_w", []byte{0xff}},
		{"-_8", []byte{0xfb, 0xff}},
		// With padding
		{"AA==", []byte{0x00}},
	}

	for _, tt := range tests {
		result, err := base64URLDecode(tt.input)
		if err != nil {
			t.Errorf("base64URLDecode(%s) error: %v", tt.input, err)
			continue
		}
		if len(result) != len(tt.expected) {
			t.Errorf("base64URLDecode(%s) = %v, want %v", tt.input, result, tt.expected)
			continue
		}
		for i := range result {
			if result[i] != tt.expected[i] {
				t.Errorf("base64URLDecode(%s) = %v, want %v", tt.input, result, tt.expected)
				break
			}
		}
	}
}

func TestPadBytes(t *testing.T) {
	tests := []struct {
		input    []byte
		size     int
		expected []byte
	}{
		{[]byte{0x01}, 1, []byte{0x01}},
		{[]byte{0x01}, 2, []byte{0x00, 0x01}},
		{[]byte{0x01}, 4, []byte{0x00, 0x00, 0x00, 0x01}},
		{[]byte{0x01, 0x02}, 2, []byte{0x01, 0x02}},
		{[]byte{0x01, 0x02, 0x03}, 2, []byte{0x01, 0x02, 0x03}}, // larger than size
	}

	for _, tt := range tests {
		result := padBytes(tt.input, tt.size)
		if len(result) != len(tt.expected) {
			t.Errorf("padBytes(%v, %d) = %v, want %v", tt.input, tt.size, result, tt.expected)
			continue
		}
		for i := range result {
			if result[i] != tt.expected[i] {
				t.Errorf("padBytes(%v, %d) = %v, want %v", tt.input, tt.size, result, tt.expected)
				break
			}
		}
	}
}

func TestToJWK_P384(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		t.Fatalf("GenerateKey() error: %v", err)
	}

	jwk, err := ToJWK(&privateKey.PublicKey)
	if err != nil {
		t.Fatalf("ToJWK() error: %v", err)
	}

	if jwk.Crv != "P-384" {
		t.Errorf("Crv = %s, want P-384", jwk.Crv)
	}
	if jwk.Alg != "ES384" {
		t.Errorf("Alg = %s, want ES384", jwk.Alg)
	}
}

func TestToJWK_P521(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	if err != nil {
		t.Fatalf("GenerateKey() error: %v", err)
	}

	jwk, err := ToJWK(&privateKey.PublicKey)
	if err != nil {
		t.Fatalf("ToJWK() error: %v", err)
	}

	if jwk.Crv != "P-521" {
		t.Errorf("Crv = %s, want P-521", jwk.Crv)
	}
	if jwk.Alg != "ES512" {
		t.Errorf("Alg = %s, want ES512", jwk.Alg)
	}
}

func TestJWK_ToPublicKey_P384(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		t.Fatalf("GenerateKey() error: %v", err)
	}

	jwk, err := ToJWK(&privateKey.PublicKey)
	if err != nil {
		t.Fatalf("ToJWK() error: %v", err)
	}

	pubKey, err := jwk.ToPublicKey()
	if err != nil {
		t.Fatalf("ToPublicKey() error: %v", err)
	}

	if pubKey.Curve != elliptic.P384() {
		t.Error("Curve mismatch for P-384")
	}
}

func TestJWK_ToPublicKey_P521(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	if err != nil {
		t.Fatalf("GenerateKey() error: %v", err)
	}

	jwk, err := ToJWK(&privateKey.PublicKey)
	if err != nil {
		t.Fatalf("ToJWK() error: %v", err)
	}

	pubKey, err := jwk.ToPublicKey()
	if err != nil {
		t.Fatalf("ToPublicKey() error: %v", err)
	}

	if pubKey.Curve != elliptic.P521() {
		t.Error("Curve mismatch for P-521")
	}
}

func TestJWK_ToPublicKey_InvalidX(t *testing.T) {
	jwk := &JWK{
		Kty: "EC",
		Crv: "P-256",
		X:   "invalid!!!",
		Y:   "BBBB",
	}

	_, err := jwk.ToPublicKey()
	if err == nil {
		t.Error("ToPublicKey() with invalid X should return error")
	}
}

func TestJWK_ToPublicKey_InvalidY(t *testing.T) {
	// Get a valid X from a real key
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	jwk, _ := ToJWK(&privateKey.PublicKey)

	jwk.Y = "invalid!!!"

	_, err := jwk.ToPublicKey()
	if err == nil {
		t.Error("ToPublicKey() with invalid Y should return error")
	}
}

func TestSerialize_P384(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		t.Fatalf("GenerateKey() error: %v", err)
	}

	thumbprint, _ := ComputeThumbprint(&privateKey.PublicKey)
	kp := &KeyPair{
		PrivateKey: privateKey,
		Thumbprint: thumbprint,
	}

	serialized, err := kp.Serialize()
	if err != nil {
		t.Fatalf("Serialize() error: %v", err)
	}

	if serialized.Curve != "P-384" {
		t.Errorf("Curve = %s, want P-384", serialized.Curve)
	}
}

func TestSerialize_P521(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	if err != nil {
		t.Fatalf("GenerateKey() error: %v", err)
	}

	thumbprint, _ := ComputeThumbprint(&privateKey.PublicKey)
	kp := &KeyPair{
		PrivateKey: privateKey,
		Thumbprint: thumbprint,
	}

	serialized, err := kp.Serialize()
	if err != nil {
		t.Fatalf("Serialize() error: %v", err)
	}

	if serialized.Curve != "P-521" {
		t.Errorf("Curve = %s, want P-521", serialized.Curve)
	}
}

func TestSerialize_NilKey(t *testing.T) {
	kp := &KeyPair{}
	_, err := kp.Serialize()
	if err == nil {
		t.Error("Serialize() with nil key should return error")
	}
}

func TestDeserializeKeyPair_P384(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		t.Fatalf("GenerateKey() error: %v", err)
	}

	thumbprint, _ := ComputeThumbprint(&privateKey.PublicKey)
	original := &KeyPair{
		PrivateKey: privateKey,
		Thumbprint: thumbprint,
	}

	serialized, _ := original.Serialize()
	restored, err := DeserializeKeyPair(serialized)
	if err != nil {
		t.Fatalf("DeserializeKeyPair() error: %v", err)
	}

	if restored.PrivateKey.Curve != elliptic.P384() {
		t.Error("Curve mismatch for P-384")
	}
}

func TestDeserializeKeyPair_P521(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	if err != nil {
		t.Fatalf("GenerateKey() error: %v", err)
	}

	thumbprint, _ := ComputeThumbprint(&privateKey.PublicKey)
	original := &KeyPair{
		PrivateKey: privateKey,
		Thumbprint: thumbprint,
	}

	serialized, _ := original.Serialize()
	restored, err := DeserializeKeyPair(serialized)
	if err != nil {
		t.Fatalf("DeserializeKeyPair() error: %v", err)
	}

	if restored.PrivateKey.Curve != elliptic.P521() {
		t.Error("Curve mismatch for P-521")
	}
}

func TestDeserializeKeyPair_InvalidCurve(t *testing.T) {
	s := &SerializedKeyPair{
		Curve:       "unknown",
		PrivateKeyD: "AAAA",
		PublicKeyX:  "BBBB",
		PublicKeyY:  "CCCC",
	}

	_, err := DeserializeKeyPair(s)
	if err == nil {
		t.Error("DeserializeKeyPair() with invalid curve should return error")
	}
}

func TestDeserializeKeyPair_InvalidPrivateKey(t *testing.T) {
	s := &SerializedKeyPair{
		Curve:       "P-256",
		PrivateKeyD: "invalid!!!",
		PublicKeyX:  "BBBB",
		PublicKeyY:  "CCCC",
	}

	_, err := DeserializeKeyPair(s)
	if err == nil {
		t.Error("DeserializeKeyPair() with invalid private key should return error")
	}
}

func TestDeserializeKeyPair_InvalidPublicKeyX(t *testing.T) {
	s := &SerializedKeyPair{
		Curve:       "P-256",
		PrivateKeyD: "AAAA",
		PublicKeyX:  "invalid!!!",
		PublicKeyY:  "CCCC",
	}

	_, err := DeserializeKeyPair(s)
	if err == nil {
		t.Error("DeserializeKeyPair() with invalid X should return error")
	}
}

func TestDeserializeKeyPair_InvalidPublicKeyY(t *testing.T) {
	// Get a valid key to get valid D and X
	kp, _ := GenerateKeyPair()
	serialized, _ := kp.Serialize()

	serialized.PublicKeyY = "invalid!!!"

	_, err := DeserializeKeyPair(serialized)
	if err == nil {
		t.Error("DeserializeKeyPair() with invalid Y should return error")
	}
}

func TestDeserializeKeyPair_EmptyThumbprint(t *testing.T) {
	// Generate and serialize a key pair
	kp, _ := GenerateKeyPair()
	serialized, _ := kp.Serialize()

	// Clear thumbprint to test recomputation
	serialized.Thumbprint = ""

	restored, err := DeserializeKeyPair(serialized)
	if err != nil {
		t.Fatalf("DeserializeKeyPair() error: %v", err)
	}

	// Thumbprint should be recomputed
	if restored.Thumbprint == "" {
		t.Error("Thumbprint should be recomputed")
	}
	if restored.Thumbprint != kp.Thumbprint {
		t.Errorf("Thumbprint = %s, want %s", restored.Thumbprint, kp.Thumbprint)
	}
}

// TestRFC7638Example tests against a known RFC 7638 example (adapted for EC).
// RFC 7638 Section 3.1 provides an RSA example, but we verify our EC implementation
// produces consistent, deterministic thumbprints.
func TestRFC7638Consistency(t *testing.T) {
	kp, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair() error: %v", err)
	}

	// Compute thumbprint multiple times
	for i := range 10 {
		thumbprint, err := ComputeThumbprint(kp.PublicKey())
		if err != nil {
			t.Fatalf("ComputeThumbprint() iteration %d error: %v", i, err)
		}
		if thumbprint != kp.Thumbprint {
			t.Errorf("Thumbprint iteration %d = %s, want %s", i, thumbprint, kp.Thumbprint)
		}
	}
}
