# OAuth 2.0 Util for Go

[![Go Report Card][goreport-svg]][goreport-link]
[![Docs][docs-godoc-svg]][docs-godoc-link]
[![License][license-svg]][license-link]

[OAuth 2.0 - https://github.com/golang/oauth2](https://github.com/golang/oauth2) helper API calls related to OAuth 2.0 user profile information. Currently provices helper libraries to retrieve canonical user information from services. The [SCIM](http://www.simplecloud.info/) user schema is used for a canonical user model.

## Installation

```
$ go get github.com/grokify/oauth2-util-go/...
```

## Usage

### Google

```golang
import(
	"github.com/grokify/oauth2-util-go/googleutil"
)

// googleOAuth2HTTPClient is *http.Client from Golang OAuth2
googleutil := googleutil.GoogleClientUtil(googleOAuth2HTTPClient)
scimuser, err := googleutil.GetSCIMUser()
```

### Facebook

```golang
import(
	"github.com/grokify/oauth2-util-go/facebookutil"
)

// fbOAuth2HTTPClient is *http.Client from Golang OAuth2
fbutil := facebookutil.GoogleClientUtil(fbOAuth2HTTPClient)
scimuser, err := fbutil.GetSCIMUser()
```

 [goreport-svg]: https://goreportcard.com/badge/github.com/grokify/oauth2-util-go
 [goreport-link]: https://goreportcard.com/report/github.com/grokify/oauth2-util-go
 [docs-godoc-svg]: https://img.shields.io/badge/docs-godoc-blue.svg
 [docs-godoc-link]: https://godoc.org/github.com/grokify/oauth2-util-go
 [license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
 [license-link]: https://github.com/grokify/oauth2-util-go/blob/master/LICENSE.md
