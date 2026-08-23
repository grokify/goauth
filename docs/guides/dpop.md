# DPoP (RFC 9449)

The `dpop` package implements [RFC 9449 — OAuth 2.0 Demonstrating Proof of Possession](https://www.rfc-editor.org/rfc/rfc9449). DPoP binds an access token to a client-held key pair: each request carries a signed `DPoP` proof, so a stolen token is useless without the corresponding private key.

```go
import "github.com/grokify/goauth/dpop"
```

## Concepts

- **Key pair** — an ES256 (P-256) key the client generates and keeps. Its public key is identified by a **JWK thumbprint** (RFC 7638).
- **Proof** — a short-lived JWT signed by the client's private key, carrying the HTTP method + URI and, optionally, a hash of the access token (`ath`) and a server nonce.
- **Verifier** — server-side validation of a proof against the request, with age/clock-skew checks and optional token binding.

## Client

### Generate a key pair

```go
kp, err := dpop.GenerateKeyPair()
if err != nil {
    return err
}
thumbprint := kp.Thumbprint // JWK thumbprint identifying the public key
```

### Create a proof for a request

```go
proof, err := dpop.CreateProof(kp, http.MethodGet, "https://api.example.com/resource")
if err != nil {
    return err
}
req.Header.Set("DPoP", proof)
```

The URI is the scheme + host + path — **no query string or fragment**.

### Bind the proof to an access token

For resource requests, bind the proof to the access token (adds the `ath` claim) and send the token with the `DPoP` auth scheme:

```go
proof, err := dpop.CreateProofWithOptions(kp, http.MethodGet, uri, dpop.ProofOptions{
    AccessToken: accessToken,
    Nonce:       serverNonce, // optional, if the server issued one
})
if err != nil {
    return err
}
req.Header.Set("Authorization", "DPoP "+accessToken)
req.Header.Set("DPoP", proof)
```

### Persist a key pair

A DPoP key pair is long-lived (it binds tokens across requests). Serialize it to store between sessions:

```go
data, err := kp.SerializeJSON()   // []byte, safe to store
// ...later...
kp, err = dpop.DeserializeKeyPairJSON(data)
```

## Server

### Verify a proof directly

```go
verifier := dpop.NewVerifier(dpop.DefaultVerificationConfig())

result, err := verifier.Verify(ctx, dpop.VerificationRequest{
    Proof:       r.Header.Get("DPoP"),
    Method:      r.Method,
    URI:         "https://api.example.com" + r.URL.Path,
    AccessToken: accessToken, // required if RequireAccessTokenBinding is set
})
if err != nil {
    http.Error(w, "invalid DPoP proof", http.StatusUnauthorized)
    return
}
// result.Thumbprint identifies the client key; compare it to the token's binding.
```

`DefaultVerificationConfig` allows a 5-minute proof age and 30s clock skew. Set `RequireAccessTokenBinding: true` to require the `ath` claim, or a `NonceValidator` to enforce server nonces.

### Use the middleware

`Middleware` validates the `DPoP` proof and stores the result in the request context:

```go
verifier := dpop.NewVerifier(dpop.DefaultVerificationConfig())

mw := dpop.Middleware(dpop.MiddlewareConfig{
    Verifier:    verifier,
    RequireDPoP: true, // reject requests without a valid proof
})

handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    thumbprint := dpop.GetThumbprint(r.Context())
    _ = thumbprint // bind/authorize against the token's cnf.jkt
    w.WriteHeader(http.StatusOK)
}))
```

Convenience wrappers over a verifier:

- `dpop.RequireDPoP(verifier)` — reject requests without a valid proof.
- `dpop.OptionalDPoP(verifier)` — validate when present, allow through when absent.

Retrieve results downstream with `dpop.GetVerificationResult(ctx)` and `dpop.GetThumbprint(ctx)`.

## Notes

- Proofs are single-use and short-lived; generate a fresh proof per request.
- The package is provider-independent — it depends only on the standard library plus `golang-jwt/jwt/v5` and `google/uuid`.
- This package previously lived in `systemforge/session/dpop`; the API is unchanged (see the [v0.24.0 release notes](../releases/v0.24.0.md)).

## Links

- [RFC 9449 — OAuth 2.0 Demonstrating Proof of Possession](https://www.rfc-editor.org/rfc/rfc9449)
- [Go package reference](https://pkg.go.dev/github.com/grokify/goauth/dpop)
