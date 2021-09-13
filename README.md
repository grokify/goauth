# OAuth 2.0 More for Go

[![Build Status][build-status-svg]][build-status-url]
[![Go Report Card][goreport-svg]][goreport-url]
[![Used By][used-by-svg]][used-by-url]
[![Docs][docs-godoc-svg]][docs-godoc-url]
[![License][license-svg]][license-url]

More [OAuth 2.0 - https://github.com/golang/oauth2](https://github.com/golang/oauth2) functionality. Currently provides:

* `NewClient()` functions to create `*http.Client` structs for services not supported in `oauth2` like `aha`, `metabase`, `ringcentral`, `salesforce`, `visa`, etc. Generating `*http.Client` structs is especially useful for using with Swagger Codegen auto-generated SDKs to support different auth models.
* Helper libraries to retrieve canonical user information from services. The [SCIM](http://www.simplecloud.info/) user schema is used for a canonical user model.
* Multi-service libraries to more transparently handle OAuth 2 for multiple services, e.g. a website that supports Google and Facebook auth. This is demoed in [grokify/beego-oauth2-demo](https://github.com/grokify/beego-oauth2-demo)

## Installation

```
$ go get github.com/grokify/oauth2more
```

## Usage

### Canonical User Information

`ClientUtil` structs satisfy the interface having `SetClient()` and `GetSCIMUser()` functions.

#### Google

```golang
import(
	"github.com/grokify/oauth2more/google"
)

// googleOAuth2HTTPClient is *http.Client from Golang OAuth2
googleClientUtil := google.NewClientUtil(googleOAuth2HTTPClient)
scimuser, err := googleClientUtil.GetSCIMUser()
```

#### Facebook

```golang
import(
	"github.com/grokify/oauth2more/facebook"
)

// fbOAuth2HTTPClient is *http.Client from Golang OAuth2
fbClientUtil := facebook.NewClientUtil(fbOAuth2HTTPClient)
scimuser, err := fbClientUtil.GetSCIMUser()
```

#### RingCentral

```golang
import(
	"github.com/grokify/oauth2more/ringcentral"
)

// rcOAuth2HTTPClient is *http.Client from Golang OAuth2
rcClientUtil := ringcentral.NewClientUtil(rcOAuth2HTTPClient)
scimuser, err := rcClientUtil.GetSCIMUser()
```

## Test Redirect URL

This repo comes with a generic test OAuth 2 redirect page which can be used with headless (no-UI) apps. To use this test URL, configure the following URL to be your OAuth 2 redirect URI. This will write the Authorization Code in the HTMl which you can then copy and paste into your own app.

The URL is located here:

* [https://grokify.github.io/oauth2more/oauth2callback/](https://grokify.github.io/oauth2more/oauth2callback/)

## Example App

See the following repo for a Beego-based demo app:

* https://github.com/grokify/beego-oauth2-demo

 [used-by-svg]: https://sourcegraph.com/github.com/grokify/oauth2more/-/badge.svg
 [used-by-url]: https://sourcegraph.com/github.com/grokify/oauth2more?badge
 [build-status-svg]: https://github.com/grokify/oauth2more/workflows/go%20build/badge.svg
 [build-status-url]: https://github.com/grokify/oauth2more/actions
 [goreport-svg]: https://goreportcard.com/badge/github.com/grokify/oauth2more
 [goreport-url]: https://goreportcard.com/report/github.com/grokify/oauth2more
 [docs-godoc-svg]: https://pkg.go.dev/badge/github.com/grokify/oauth2more
 [docs-godoc-url]: https://pkg.go.dev/github.com/grokify/oauth2more
 [license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
 [license-url]: https://github.com/grokify/oauth2more/blob/master/LICENSE.md
