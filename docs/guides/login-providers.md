# Sign in with Google or GitHub

The `providers` package turns the last step of an OAuth 2.0 authorization-code
login, exchanging the code and looking up who signed in, into one call. It
returns the same `OAuthUser` type for every provider, so your login handler
doesn't branch on provider-specific response formats.

```go
import "github.com/grokify/goauth/providers"
```

## The `OAuthUser` type

```go
type OAuthUser struct {
    ProviderID   string // user ID from the provider
    Email        string
    Name         string
    AvatarURL    string
    Provider     string // providers.ProviderGoogle or providers.ProviderGitHub
    AccessToken  string // provider access token, for further API calls
    RefreshToken string // provider refresh token, if issued
}
```

Key your user records on `Provider` + `ProviderID`, not on email: a user's email
address can change, and two providers can report the same address.

## Exchanging the authorization code

In your OAuth callback handler, pass the `code` query parameter to the fetch
function for that provider:

```go
import (
    "golang.org/x/oauth2"
    "golang.org/x/oauth2/github"
)

var githubConfig = &oauth2.Config{
    ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
    ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
    RedirectURL:  "https://example.com/auth/github/callback",
    Scopes:       []string{"read:user", "user:email"},
    Endpoint:     github.Endpoint,
}

func githubCallback(w http.ResponseWriter, r *http.Request) {
    // Verify the "state" parameter against the session first (not shown).
    ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
    defer cancel()

    user, err := providers.FetchGitHubUser(ctx, githubConfig, r.URL.Query().Get("code"))
    if err != nil {
        http.Error(w, "login failed", http.StatusUnauthorized)
        return
    }
    // user.Provider == "github", user.ProviderID, user.Email, ...
}
```

`providers.FetchGoogleUser(ctx, googleConfig, code)` works the same way; request
the `openid`, `email`, and `profile` scopes.

Code-exchange requests use the HTTP client derived from your `oauth2.Config`,
so bound them with a context deadline as shown.

### GitHub email handling

A GitHub user can keep their profile email private, in which case the `/user`
endpoint returns no email. `FetchGitHubUser` then calls `/user/emails` and
uses the **primary verified** address, falling back to any verified address.
This requires the `user:email` scope. If the user has no verified email, the
call fails with `providers.ErrNoVerifiedEmail`:

```go
if errors.Is(err, providers.ErrNoVerifiedEmail) {
    // ask the user to verify an email address on GitHub
}
```

If `name` is not set on the GitHub profile, `OAuthUser.Name` falls back to the
user's login.

## Fetching a profile with an existing token

When you already hold an access token, for example a stored one, fetch the
provider-native profile directly:

```go
gh, err := providers.FetchGitHubUserWithToken(ctx, accessToken) // *GitHubUserInfo
g, err := providers.FetchGoogleUserWithToken(ctx, accessToken)  // *GoogleUserInfo
```

These use their own HTTP client, bounded by `providers.DefaultTimeout` (30s).

## Errors

- A failed code exchange is wrapped as `exchange code: ...`.
- A non-200 provider response reports the URL, the status code, and up to
  4 KB of the response body.

Don't show these errors to end users verbatim; log them and return a generic
failure message.

## Relationship to other packages

For provider-specific features, use the provider packages directly: for
example, `google.ClientUtil` for Google userinfo and SCIM mapping. The
`providers` package trades that depth for one uniform type across providers.
