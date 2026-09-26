# Google

# Configure App and OAuth URLs

Use the following to configure a Google Cloud app for OAuth2. This includes configuring the OAuth 2.0 Redirect URL.

1. Go to https://console.cloud.google.com/
1. Select your project.
1. Click "APIs and Services" > "Credentials"

## Running the Examples

Two runnable examples load credentials from a [credentials set](https://github.com/grokify/goauth#credentials-set-multiple-accounts)
file and fetch the signed-in user's Google profile.

**OAuth 2.0 web flow** (`google/cmd/oauth2web`), using an `oauth2` credential
with `"service": "google"`:

```bash
go run ./google/cmd/oauth2web --creds credentials.json --account my-google-app
```

**Service account** (`google/cmd/serviceaccount`), using a `gcpsa` credential:

```bash
go run ./google/cmd/serviceaccount --creds credentials.json --account my-service-account
```

Both print only token metadata and the profile, never credentials or tokens.
Keep credentials files out of version control; this repository ignores files
starting with `_`, so a name like `_credentials.json` is not committed.

For "Sign in with Google" in your own application, see
[Sign in with Google or GitHub](guides/login-providers.md).
