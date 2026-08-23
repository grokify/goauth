# GoAuth

**OAuth 2.0, JWT, and multi-provider authentication for Go.**

GoAuth provides a unified credentials model and helpers for authenticating to many services, plus building blocks like DPoP (RFC 9449) for sender-constrained tokens.

## Features

- **Unified credentials** — a single JSON configuration format across authentication types and services
- **Multiple auth types** — OAuth 2.0, JWT, Basic Auth, GCP Service Account, and custom header/query auth
- **40+ OAuth 2.0 providers** — pre-configured endpoints for popular services
- **Multiple grant types** — Authorization Code, Client Credentials, Password, JWT Bearer, SAML2 Bearer, Refresh Token
- **PKCE** — Proof Key for Code Exchange
- **DPoP (RFC 9449)** — Demonstrating Proof of Possession for sender-constrained access tokens

## Installation

```bash
go get github.com/grokify/goauth
```

## Guides

- [DPoP (RFC 9449)](guides/dpop.md) — sender-constrained tokens: client proofs, server verification, HTTP middleware
- [Google](google.md) — configuring Google OAuth apps and URLs

## Documentation

- [Changelog](changelog.md)
- [Go package reference](https://pkg.go.dev/github.com/grokify/goauth)

!!! note "Work in progress"
    This documentation site is thin and growing. Contributions and expansions are welcome.
