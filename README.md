# GoAuth

[![Build Status][build-status-svg]][build-status-url]
[![Go Report Card][goreport-svg]][goreport-url]
[![Used By][used-by-svg]][used-by-url]
[![Docs][docs-godoc-svg]][docs-godoc-url]
[![License][license-svg]][license-url]

GoAuth provides helper libraries for authentication in Go, with a focus on API services. It covers [OAuth 2.0](https://github.com/golang/oauth2), [JWT](https://github.com/golang-jwt/jwt), TLS client authentication and Basic Auth. A primary goal is to be able to create a `*http.Client` from a single JSON application definition.

Major features include:

1. Create `*http.Client` for multiple API services. Use `NewClient()` functions to create `*http.Client` structs for services not supported in `oauth2` like `aha`, `metabase`, `ringcentral`, `salesforce`, `visa`, etc. Generating `*http.Client` structs is especially useful for using with Swagger Codegen auto-generated SDKs to support different auth models.
1. Generically store and retrieve multiple app credentials in a single JSON object via `credentials`. Supports both OAuth 2 and JWT.
1. Create OAuth 2.0 authorization code token from the command line (for test purposes). No website is needed.
1. Retrieve canonical user information via helper libraries to retrieve canonical user information from services. The [SCIM](http://www.simplecloud.info/) user schema is used for a canonical user model. This may be replaced/augmented by OIDC `userinfo` in the future.
1. Transparently handle OAuth 2 for multiple services, e.g. a website that supports Google and Facebook auth. This is demoed in [grokify/beego-oauth2-demo](https://github.com/grokify/beego-oauth2-demo)

## Installation

```
$ go get github.com/grokify/goauth
```

## Usage

### Canonical User Information

`ClientUtil` structs satisfy the interface having `SetClient()` and `GetSCIMUser()` functions.

#### Google

```golang
import(
	"github.com/grokify/goauth/google"
)

// googleOAuth2HTTPClient is *http.Client from Golang OAuth2
googleClientUtil := google.NewClientUtil(googleOAuth2HTTPClient)
scimuser, err := googleClientUtil.GetSCIMUser()
```

#### Facebook

```golang
import(
	"github.com/grokify/goauth/facebook"
)

// fbOAuth2HTTPClient is *http.Client from Golang OAuth2
fbClientUtil := facebook.NewClientUtil(fbOAuth2HTTPClient)
scimuser, err := fbClientUtil.GetSCIMUser()
```

#### RingCentral

```golang
import(
	"github.com/grokify/goauth/ringcentral"
)

// rcOAuth2HTTPClient is *http.Client from Golang OAuth2
rcClientUtil := ringcentral.NewClientUtil(rcOAuth2HTTPClient)
scimuser, err := rcClientUtil.GetSCIMUser()
```

## Test Redirect URL

This repo comes with a generic test OAuth 2 redirect page which can be used with headless (no-UI) apps. To use this test URL, configure the following URL to be your OAuth 2 redirect URI. This will write the Authorization Code in the HTMl which you can then copy and paste into your own app.

The URL is located here:

* [https://grokify.github.io/goauth/oauth2callback/](https://grokify.github.io/goauth/oauth2callback/)

## Example App

See the following repo for a Beego-based demo app:

* https://github.com/grokify/beego-oauth2-demo

 [used-by-svg]: https://sourcegraph.com/github.com/grokify/goauth/-/badge.svg
 [used-by-url]: https://sourcegraph.com/github.com/grokify/goauth?badge
 [build-status-svg]: https://github.com/grokify/goauth/workflows/go%20build/badge.svg
 [build-status-url]: https://github.com/grokify/goauth/actions
 [goreport-svg]: https://goreportcard.com/badge/github.com/grokify/goauth
 [goreport-url]: https://goreportcard.com/report/github.com/grokify/goauth
 [docs-godoc-svg]: https://pkg.go.dev/badge/github.com/grokify/goauth
 [docs-godoc-url]: https://pkg.go.dev/github.com/grokify/goauth
 [loc-svg]: https://tokei.rs/b1/github/grokify/goauth
 [repo-url]: https://github.com/grokify/goauth
 [license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
 [license-url]: https://github.com/grokify/goauth/blob/master/LICENSE.md
