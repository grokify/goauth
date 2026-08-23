// Package dpop implements Demonstrating Proof of Possession (DPoP) per RFC 9449.
// DPoP binds access tokens to cryptographic key pairs held by clients,
// preventing stolen tokens from being used without the private key.
package dpop

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
)

var (
	// ErrInvalidKey is returned when a key is malformed or invalid.
	ErrInvalidKey = errors.New("invalid key")
	// ErrUnsupportedAlgorithm is returned for unsupported cryptographic algorithms.
	ErrUnsupportedAlgorithm = errors.New("unsupported algorithm")
)

// KeyPair represents an ES256 key pair for DPoP.
type KeyPair struct {
	// PrivateKey is the ECDSA private key for signing proofs.
	PrivateKey *ecdsa.PrivateKey
	// Thumbprint is the JWK thumbprint (RFC 7638) of the public key.
	Thumbprint string
}

// GenerateKeyPair creates a new ES256 (P-256/secp256r1) key pair for DPoP.
// The thumbprint is computed per RFC 7638 using SHA-256.
func GenerateKeyPair() (*KeyPair, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("generating ECDSA key: %w", err)
	}

	thumbprint, err := ComputeThumbprint(&privateKey.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("computing thumbprint: %w", err)
	}

	return &KeyPair{
		PrivateKey: privateKey,
		Thumbprint: thumbprint,
	}, nil
}

// PublicKey returns the public key portion of the key pair.
func (kp *KeyPair) PublicKey() *ecdsa.PublicKey {
	if kp.PrivateKey == nil {
		return nil
	}
	return &kp.PrivateKey.PublicKey
}

// ComputeThumbprint computes the JWK thumbprint per RFC 7638 for an ECDSA public key.
// For EC keys, the required members are: crv, kty, x, y (in lexicographic order).
func ComputeThumbprint(publicKey *ecdsa.PublicKey) (string, error) {
	if publicKey == nil {
		return "", ErrInvalidKey
	}

	// Determine curve name
	var crv string
	switch publicKey.Curve {
	case elliptic.P256():
		crv = "P-256"
	case elliptic.P384():
		crv = "P-384"
	case elliptic.P521():
		crv = "P-521"
	default:
		return "", fmt.Errorf("%w: unsupported curve", ErrUnsupportedAlgorithm)
	}

	// Get byte size for the curve (for padding)
	byteSize := (publicKey.Curve.Params().BitSize + 7) / 8

	// Use PublicKey.Bytes() which returns uncompressed SEC 1 format: 0x04 || X || Y
	// Each coordinate is already padded to byteSize
	pubBytes, err := publicKey.Bytes()
	if err != nil {
		return "", fmt.Errorf("%w: failed to encode public key: %v", ErrInvalidKey, err)
	}
	x := base64URLEncode(pubBytes[1 : byteSize+1])
	y := base64URLEncode(pubBytes[byteSize+1:])

	// RFC 7638: Create JSON with required members in lexicographic order
	// For EC keys: crv, kty, x, y
	thumbprintInput := fmt.Sprintf(`{"crv":"%s","kty":"EC","x":"%s","y":"%s"}`, crv, x, y)

	// SHA-256 hash
	hash := sha256.Sum256([]byte(thumbprintInput))

	// Base64url encode without padding
	return base64URLEncode(hash[:]), nil
}

// JWK represents a JSON Web Key for embedding in DPoP proof headers.
type JWK struct {
	Kty string `json:"kty"`           // Key type (EC)
	Crv string `json:"crv"`           // Curve (P-256)
	X   string `json:"x"`             // X coordinate
	Y   string `json:"y"`             // Y coordinate
	Alg string `json:"alg,omitempty"` // Algorithm (ES256)
}

// ToJWK converts an ECDSA public key to a JWK representation.
func ToJWK(publicKey *ecdsa.PublicKey) (*JWK, error) {
	if publicKey == nil {
		return nil, ErrInvalidKey
	}

	var crv, alg string
	switch publicKey.Curve {
	case elliptic.P256():
		crv = "P-256"
		alg = "ES256"
	case elliptic.P384():
		crv = "P-384"
		alg = "ES384"
	case elliptic.P521():
		crv = "P-521"
		alg = "ES512"
	default:
		return nil, fmt.Errorf("%w: unsupported curve", ErrUnsupportedAlgorithm)
	}

	byteSize := (publicKey.Curve.Params().BitSize + 7) / 8

	// Use PublicKey.Bytes() which returns uncompressed SEC 1 format: 0x04 || X || Y
	pubBytes, err := publicKey.Bytes()
	if err != nil {
		return nil, fmt.Errorf("%w: failed to encode public key: %v", ErrInvalidKey, err)
	}

	return &JWK{
		Kty: "EC",
		Crv: crv,
		X:   base64URLEncode(pubBytes[1 : byteSize+1]),
		Y:   base64URLEncode(pubBytes[byteSize+1:]),
		Alg: alg,
	}, nil
}

// Thumbprint computes and returns the JWK thumbprint for this key.
func (j *JWK) Thumbprint() (string, error) {
	// RFC 7638: required members in lexicographic order for EC: crv, kty, x, y
	thumbprintInput := fmt.Sprintf(`{"crv":"%s","kty":"%s","x":"%s","y":"%s"}`, j.Crv, j.Kty, j.X, j.Y)
	hash := sha256.Sum256([]byte(thumbprintInput))
	return base64URLEncode(hash[:]), nil
}

// ToPublicKey converts a JWK back to an ECDSA public key.
func (j *JWK) ToPublicKey() (*ecdsa.PublicKey, error) {
	if j.Kty != "EC" {
		return nil, fmt.Errorf("%w: expected EC key type", ErrInvalidKey)
	}

	var curve elliptic.Curve
	switch j.Crv {
	case "P-256":
		curve = elliptic.P256()
	case "P-384":
		curve = elliptic.P384()
	case "P-521":
		curve = elliptic.P521()
	default:
		return nil, fmt.Errorf("%w: unsupported curve %s", ErrUnsupportedAlgorithm, j.Crv)
	}

	xBytes, err := base64URLDecode(j.X)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid x coordinate: %v", ErrInvalidKey, err)
	}

	yBytes, err := base64URLDecode(j.Y)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid y coordinate: %v", ErrInvalidKey, err)
	}

	// Build the uncompressed SEC 1 point (0x04 || X || Y) with each coordinate
	// left-padded to the curve's field size, then let ParseUncompressedPublicKey
	// validate that the point is well-formed and on the curve.
	uncompressed, err := uncompressedPoint(curve, xBytes, yBytes)
	if err != nil {
		return nil, err
	}

	pub, err := ecdsa.ParseUncompressedPublicKey(curve, uncompressed)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidKey, err)
	}
	return pub, nil
}

// uncompressedPoint assembles an uncompressed SEC 1 encoded point
// (0x04 || X || Y) from big-endian coordinate bytes, left-padding each
// coordinate to the curve's field size.
func uncompressedPoint(curve elliptic.Curve, xBytes, yBytes []byte) ([]byte, error) {
	byteSize := (curve.Params().BitSize + 7) / 8
	if len(xBytes) > byteSize || len(yBytes) > byteSize {
		return nil, fmt.Errorf("%w: coordinate too large for curve", ErrInvalidKey)
	}

	out := make([]byte, 1+2*byteSize)
	out[0] = 0x04
	copy(out[1+byteSize-len(xBytes):1+byteSize], xBytes)
	copy(out[1+2*byteSize-len(yBytes):], yBytes)
	return out, nil
}

// padScalar left-pads a big-endian scalar to the curve's field size. A scalar
// longer than the field size (after accounting for it being a fixed-width
// value) is rejected with a message that states the expected length.
func padScalar(curve elliptic.Curve, b []byte) ([]byte, error) {
	byteSize := (curve.Params().BitSize + 7) / 8
	if len(b) > byteSize {
		return nil, fmt.Errorf("%w: scalar is %d bytes, want at most %d for %s", ErrInvalidKey, len(b), byteSize, curve.Params().Name)
	}
	if len(b) == byteSize {
		return b, nil
	}
	out := make([]byte, byteSize)
	copy(out[byteSize-len(b):], b)
	return out, nil
}

// Serialize serializes the key pair for storage.
// The private key is encoded in PKCS#8 format and base64url encoded.
type SerializedKeyPair struct {
	PrivateKeyD string `json:"d"`          // Private key scalar
	PublicKeyX  string `json:"x"`          // Public key X coordinate
	PublicKeyY  string `json:"y"`          // Public key Y coordinate
	Curve       string `json:"crv"`        // Curve name
	Thumbprint  string `json:"thumbprint"` // JWK thumbprint
}

// Serialize converts the key pair to a serializable format.
func (kp *KeyPair) Serialize() (*SerializedKeyPair, error) {
	if kp.PrivateKey == nil {
		return nil, ErrInvalidKey
	}

	var crv string
	switch kp.PrivateKey.Curve {
	case elliptic.P256():
		crv = "P-256"
	case elliptic.P384():
		crv = "P-384"
	case elliptic.P521():
		crv = "P-521"
	default:
		return nil, fmt.Errorf("%w: unsupported curve", ErrUnsupportedAlgorithm)
	}

	byteSize := (kp.PrivateKey.Curve.Params().BitSize + 7) / 8

	// Use PrivateKey.Bytes() for D (already padded to byteSize)
	// Use PublicKey.Bytes() for X/Y (uncompressed SEC 1 format: 0x04 || X || Y)
	privBytes, err := kp.PrivateKey.Bytes()
	if err != nil {
		return nil, fmt.Errorf("%w: failed to encode private key: %v", ErrInvalidKey, err)
	}
	pubBytes, err := kp.PrivateKey.PublicKey.Bytes()
	if err != nil {
		return nil, fmt.Errorf("%w: failed to encode public key: %v", ErrInvalidKey, err)
	}

	return &SerializedKeyPair{
		PrivateKeyD: base64URLEncode(privBytes),
		PublicKeyX:  base64URLEncode(pubBytes[1 : byteSize+1]),
		PublicKeyY:  base64URLEncode(pubBytes[byteSize+1:]),
		Curve:       crv,
		Thumbprint:  kp.Thumbprint,
	}, nil
}

// SerializeJSON serializes the key pair to JSON bytes.
func (kp *KeyPair) SerializeJSON() ([]byte, error) {
	serialized, err := kp.Serialize()
	if err != nil {
		return nil, err
	}
	return json.Marshal(serialized)
}

// DeserializeKeyPair reconstructs a key pair from its serialized form.
func DeserializeKeyPair(s *SerializedKeyPair) (*KeyPair, error) {
	if s == nil {
		return nil, ErrInvalidKey
	}

	var curve elliptic.Curve
	switch s.Curve {
	case "P-256":
		curve = elliptic.P256()
	case "P-384":
		curve = elliptic.P384()
	case "P-521":
		curve = elliptic.P521()
	default:
		return nil, fmt.Errorf("%w: unsupported curve %s", ErrUnsupportedAlgorithm, s.Curve)
	}

	dBytes, err := base64URLDecode(s.PrivateKeyD)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid private key: %v", ErrInvalidKey, err)
	}

	// ParseRawPrivateKey requires the scalar to be exactly the curve's field
	// size, but a scalar with leading zero bytes may have been stored trimmed.
	// Left-pad so trimmed scalars round-trip rather than failing with an opaque
	// "invalid scalar length" error.
	dBytes, err = padScalar(curve, dBytes)
	if err != nil {
		return nil, err
	}

	// ParseRawPrivateKey derives the public key from the scalar and validates
	// that it is well-formed.
	privateKey, err := ecdsa.ParseRawPrivateKey(curve, dBytes)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidKey, err)
	}

	// Verify the stored X/Y coordinates are consistent with the public key
	// derived from the private scalar. This detects tampering or corruption of
	// the serialized blob.
	xBytes, err := base64URLDecode(s.PublicKeyX)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid x coordinate: %v", ErrInvalidKey, err)
	}
	yBytes, err := base64URLDecode(s.PublicKeyY)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid y coordinate: %v", ErrInvalidKey, err)
	}

	pubBytes, err := privateKey.PublicKey.Bytes()
	if err != nil {
		return nil, fmt.Errorf("%w: failed to encode public key: %v", ErrInvalidKey, err)
	}
	byteSize := (curve.Params().BitSize + 7) / 8
	if !bytes.Equal(xBytes, pubBytes[1:byteSize+1]) || !bytes.Equal(yBytes, pubBytes[byteSize+1:]) {
		return nil, fmt.Errorf("%w: public key coordinates do not match private key", ErrInvalidKey)
	}

	// Recompute thumbprint to verify integrity
	thumbprint, err := ComputeThumbprint(&privateKey.PublicKey)
	if err != nil {
		return nil, err
	}

	// Verify thumbprint matches if provided
	if s.Thumbprint != "" && s.Thumbprint != thumbprint {
		return nil, fmt.Errorf("%w: thumbprint mismatch", ErrInvalidKey)
	}

	return &KeyPair{
		PrivateKey: privateKey,
		Thumbprint: thumbprint,
	}, nil
}

// DeserializeKeyPairJSON reconstructs a key pair from JSON bytes.
func DeserializeKeyPairJSON(data []byte) (*KeyPair, error) {
	var s SerializedKeyPair
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("%w: invalid JSON: %v", ErrInvalidKey, err)
	}
	return DeserializeKeyPair(&s)
}

// Signer returns a crypto.Signer interface for the private key.
func (kp *KeyPair) Signer() crypto.Signer {
	return kp.PrivateKey
}

// base64URLEncode encodes bytes to base64url without padding.
func base64URLEncode(data []byte) string {
	return base64.RawURLEncoding.EncodeToString(data)
}

// base64URLDecode decodes base64url (with or without padding).
func base64URLDecode(s string) ([]byte, error) {
	// Try without padding first
	data, err := base64.RawURLEncoding.DecodeString(s)
	if err != nil {
		// Try with padding
		data, err = base64.URLEncoding.DecodeString(s)
	}
	return data, err
}

// padBytes pads bytes to the specified length with leading zeros.
func padBytes(b []byte, size int) []byte {
	if len(b) >= size {
		return b
	}
	padded := make([]byte, size)
	copy(padded[size-len(b):], b)
	return padded
}
