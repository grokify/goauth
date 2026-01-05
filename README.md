# GoAuth

[![Build Status][build-status-svg]][build-status-url]
[![Lint Status][lint-status-svg]][lint-status-url]
[![Go Report Card][goreport-svg]][goreport-url]
[![Docs][docs-godoc-svg]][docs-godoc-url]
[![License][license-svg]][license-url]

GoAuth is a comprehensive Go authentication library designed to simplify OAuth 2.0, JWT, and other authentication methods for API services. It provides a unified configuration system for handling multiple authentication types and services, making it easy to create authenticated `*http.Client` instances from a single JSON configuration file.

## Features

- **Unified Credentials Management**: Single JSON configuration format for multiple authentication types and services
- **Multiple Authentication Types**: OAuth 2.0, JWT, Basic Auth, GCP Service Account, and custom header/query authentication
- **40+ OAuth 2.0 Providers**: Pre-configured endpoints for popular services
- **Multiple Grant Types**: Authorization Code, Client Credentials, Password, JWT Bearer, SAML2 Bearer, and Refresh Token
- **PKCE Support**: Proof Key for Code Exchange for enhanced security
- **SCIM User Model**: Canonical user information retrieval across services using [SCIM](http://www.simplecloud.info/) schema
- **CLI Tools**: Command-line utilities for token generation and API requests
- **Multi-Service OAuth**: Support for applications using multiple OAuth providers (e.g., "Login with Google" and "Login with Facebook")

## Installation

```bash
go get github.com/grokify/goauth
```

**Requirements**: Go 1.24+

## Supported OAuth 2.0 Providers

GoAuth includes pre-configured OAuth 2.0 endpoints for the following services:

| Service | Service Key | Notes |
|---------|-------------|-------|
| Aha | `aha` | Requires subdomain |
| Asana | `asana` | |
| Atlassian | `atlassian` | |
| eBay | `ebay` | Production |
| eBay Sandbox | `ebaysandbox` | Sandbox |
| Facebook | `facebook` | |
| GitHub | `github` | |
| Google | `google` | |
| HubSpot | `hubspot` | |
| Instagram | `instagram` | |
| Lyft | `lyft` | |
| Mailchimp | `mailchimp` | |
| Monday.com | `monday` | |
| PagerDuty | `pagerduty` | |
| PayPal | `paypal` | Production |
| PayPal Sandbox | `paypalsandbox` | Sandbox |
| Pipedrive | `pipedrive` | |
| Practicesuite | `practicesuite` | |
| RingCentral | `ringcentral` | Production |
| RingCentral Sandbox | `ringcentralsandbox` | Sandbox |
| Shippo | `shippo` | |
| Shopify | `shopify` | Requires subdomain |
| Slack | `slack` | |
| Stack Overflow | `stackoverflow` | |
| Stripe | `stripe` | |
| Todoist | `todoist` | |
| Uber | `uber` | |
| WePay | `wepay` | Production |
| WePay Sandbox | `wepaysandbox` | Sandbox |
| Wrike | `wrike` | |
| Wunderlist | `wunderlist` | |
| Zoom | `zoom` | |

Additional service-specific packages are available for: Auth0, Metabase, Salesforce, SuccessFactors, Visa, SparkPost, and Zendesk.

## Authentication Types

GoAuth supports the following authentication types:

| Type | Type Key | Description |
|------|----------|-------------|
| Basic Auth | `basic` | HTTP Basic Authentication |
| OAuth 2.0 | `oauth2` | OAuth 2.0 with multiple grant types |
| JWT | `jwt` | JSON Web Token generation |
| GCP Service Account | `gcpsa` | Google Cloud Platform Service Account |
| Google OAuth 2.0 | `googleoauth2` | Google-specific OAuth 2.0 |
| Header/Query | `headerquery` | Custom header or query parameter authentication |

### OAuth 2.0 Grant Types

- `authorization_code` - Authorization Code flow
- `client_credentials` - Client Credentials flow
- `password` - Resource Owner Password Credentials
- `urn:ietf:params:oauth:grant-type:jwt-bearer` - JWT Bearer
- `urn:ietf:params:oauth:grant-type:saml2-bearer` - SAML2 Bearer
- `refresh_token` - Refresh Token
- `account_credentials` - Account Credentials (Zoom Server-to-Server)

## Configuration

### Credentials Set (Multiple Accounts)

GoAuth uses a JSON configuration format that supports multiple credentials:

```json
{
  "credentials": {
    "my-google-app": {
      "service": "google",
      "type": "oauth2",
      "oauth2": {
        "clientID": "your-client-id",
        "clientSecret": "your-client-secret",
        "redirectURL": "https://example.com/callback",
        "scope": ["email", "profile"],
        "grantType": "authorization_code"
      }
    },
    "my-ringcentral-app": {
      "service": "ringcentral",
      "type": "oauth2",
      "oauth2": {
        "clientID": "your-client-id",
        "clientSecret": "your-client-secret",
        "grantType": "password",
        "username": "your-username",
        "password": "your-password"
      }
    },
    "my-api-key": {
      "type": "headerquery",
      "headerquery": {
        "serverURL": "https://api.example.com",
        "header": {
          "X-API-Key": "your-api-key"
        }
      }
    }
  }
}
```

### Credential Types

#### OAuth 2.0 Credentials

```json
{
  "service": "github",
  "type": "oauth2",
  "oauth2": {
    "serverURL": "https://api.github.com",
    "clientID": "your-client-id",
    "clientSecret": "your-client-secret",
    "redirectURL": "https://example.com/callback",
    "scope": ["repo", "user"],
    "grantType": "authorization_code",
    "pkce": false
  }
}
```

#### Basic Auth Credentials

```json
{
  "type": "basic",
  "basic": {
    "username": "your-username",
    "password": "your-password",
    "serverURL": "https://api.example.com",
    "allowInsecure": false
  }
}
```

#### JWT Credentials

```json
{
  "type": "jwt",
  "jwt": {
    "issuer": "your-issuer",
    "privateKey": "your-private-key",
    "signingMethod": "HS256"
  }
}
```

Supported signing methods: `ES256`, `ES384`, `ES512`, `HS256`, `HS384`, `HS512`

#### Header/Query Credentials

```json
{
  "type": "headerquery",
  "headerquery": {
    "serverURL": "https://api.example.com",
    "header": {
      "Authorization": "Bearer your-token",
      "X-Custom-Header": "value"
    },
    "query": {
      "api_key": "your-api-key"
    }
  }
}
```

## Usage

### Creating an HTTP Client

```go
package main

import (
    "context"
    "github.com/grokify/goauth"
)

func main() {
    ctx := context.Background()

    // From credentials file with account key
    client, err := goauth.NewClient(ctx, "credentials.json", "my-google-app")
    if err != nil {
        panic(err)
    }

    // Use client for API requests
    resp, err := client.Get("https://api.example.com/resource")
}
```

### Loading Credentials Set

```go
package main

import (
    "context"
    "github.com/grokify/goauth"
)

func main() {
    ctx := context.Background()

    // Load credentials set from file
    set, err := goauth.ReadFileCredentialsSet("credentials.json", true)
    if err != nil {
        panic(err)
    }

    // Get specific credentials
    creds, err := set.Get("my-google-app")
    if err != nil {
        panic(err)
    }

    // Create client from credentials
    client, err := creds.NewClient(ctx)
    if err != nil {
        panic(err)
    }

    // List all account keys
    accounts := set.Accounts()
}
```

### CLI-based Token Retrieval

For authorization code flow without a web server:

```go
package main

import (
    "context"
    "github.com/grokify/goauth"
)

func main() {
    ctx := context.Background()

    creds, _ := goauth.NewCredentialsFromSetFile("credentials.json", "my-app", false)

    // This will print the authorization URL and prompt for the code
    client, err := creds.NewClientCLI(ctx, "random-state")
    if err != nil {
        panic(err)
    }
}
```

### Canonical User Information (SCIM)

GoAuth provides `ClientUtil` implementations that satisfy the `OAuth2Util` interface for retrieving canonical user information:

```go
type OAuth2Util interface {
    SetClient(*http.Client)
    GetSCIMUser() (scim.User, error)
}
```

#### Google

```go
import "github.com/grokify/goauth/google"

googleClientUtil := google.NewClientUtil(googleOAuth2HTTPClient)
scimUser, err := googleClientUtil.GetSCIMUser()
```

#### Facebook

```go
import "github.com/grokify/goauth/facebook"

fbClientUtil := facebook.NewClientUtil(fbOAuth2HTTPClient)
scimUser, err := fbClientUtil.GetSCIMUser()
```

#### RingCentral

```go
import "github.com/grokify/goauth/ringcentral"

rcClientUtil := ringcentral.NewClientUtil(rcOAuth2HTTPClient)
scimUser, err := rcClientUtil.GetSCIMUser()
```

Also available for: Aha, Zoom, Metabase, Zendesk, and Salesforce.

## CLI Tools

GoAuth includes command-line tools for authentication tasks:

### goauth

Main token retrieval tool supporting all authentication types:

```bash
go run cmd/goauth/main.go --credentials credentials.json --account my-app
```

### goapi

Make authenticated API requests:

```bash
go run cmd/goapi/main.go --credentials credentials.json --account my-app --url https://api.example.com/resource
```

## Package Structure

| Package | Description |
|---------|-------------|
| `goauth` | Core credentials management and client creation |
| `authutil` | Low-level authentication utilities (BasicAuth, OAuth2, JWT, scope management) |
| `endpoints` | Pre-configured OAuth 2.0 endpoints for 30+ services |
| `scim` | SCIM schema user/group models for canonical user representation |
| `multiservice` | Multi-provider OAuth2 management for applications |
| `google` | Google-specific OAuth2 and GCP service account handling |
| `ringcentral` | RingCentral API integration |
| `facebook` | Facebook OAuth2 and user data retrieval |
| `aha`, `zoom`, `metabase`, `zendesk`, `salesforce`, `hubspot` | Service-specific implementations |

## Test Redirect URL

This repo includes a generic test OAuth 2 redirect page for headless (no-UI) applications. Configure your OAuth 2 redirect URI to:

**[https://grokify.github.io/goauth/oauth2callback/](https://grokify.github.io/goauth/oauth2callback/)**

This page displays the Authorization Code which you can copy and paste into your CLI application.

## Example Applications

- **Multi-Service OAuth Demo**: [github.com/grokify/beegoutil](https://github.com/grokify/beegoutil) - Beego-based demo showing Google and Facebook authentication
- **Examples Directory**: See `examples/` and service-specific `cmd/` directories for usage examples

## Contributing

Contributions are welcome. Please submit pull requests or create issues for bugs and feature requests.

## License

GoAuth is available under the [MIT License](LICENSE).

 [build-status-svg]: https://github.com/grokify/goauth/actions/workflows/ci.yaml/badge.svg?branch=master
 [build-status-url]: https://github.com/grokify/goauth/actions/workflows/ci.yaml
 [lint-status-svg]: https://github.com/grokify/goauth/actions/workflows/lint.yaml/badge.svg?branch=master
 [lint-status-url]: https://github.com/grokify/goauth/actions/workflows/lint.yaml
 [goreport-svg]: https://goreportcard.com/badge/github.com/grokify/goauth
 [goreport-url]: https://goreportcard.com/report/github.com/grokify/goauth
 [docs-godoc-svg]: https://pkg.go.dev/badge/github.com/grokify/goauth
 [docs-godoc-url]: https://pkg.go.dev/github.com/grokify/goauth
 [license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
 [license-url]: https://github.com/grokify/goauth/blob/master/LICENSE
