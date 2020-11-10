# Metabase

## Usage

Using default environment variiable names:

```go
import ("github.com/grokify/oauth2more")

httpClient, authResponse, clientConfig, err :=
  metabase.NewClientEnv(nil)
```

## Authentication

To query the Metabase API you need to retrieve a bearer token. You can do this with the following cURL command which is also implemented in the `AuthRequest` function:

```
curl -v -H "Content-Type: application/json" \
  -d '{"username":"myusername","password":"mypassword"}' \
  -XPOST 'http://example.com/api/session'
```

You will receive a response like:

```
{"id":"11112222-3333-4444-5555-666677778888"}
```

You can then use the `id` in the `X-Metabase-Session` header for subsequent API calls. Here's an example:

```
curl -XGET 'https://example.com/api/database' \
  -H 'X-Metabase-Session: 11112222-3333-4444-5555-666677778888'
```
