// Package providers exchanges OAuth 2.0 authorization codes for a normalized
// user profile across providers (currently Google and GitHub), for use in
// "Sign in with ..." login flows.
//
// Provider-specific packages such as google offer richer, provider-native
// APIs; this package trades that depth for one uniform OAuthUser type.
package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// DefaultTimeout bounds requests made by the *WithToken functions, which use
// their own HTTP client. Code-exchange functions use the client derived from
// the oauth2.Config; bound those with a context deadline.
const DefaultTimeout = 30 * time.Second

// maxErrorBody caps how much of a non-200 response body is included in an
// error, so a large or hostile response cannot bloat error messages.
const maxErrorBody = 4096

var defaultClient = &http.Client{Timeout: DefaultTimeout}

// OAuthUser is a user profile normalized across OAuth providers.
type OAuthUser struct {
	ProviderID   string // User ID from the provider
	Email        string
	Name         string
	AvatarURL    string
	Provider     string // "google" or "github"
	AccessToken  string // Provider access token (for API calls)
	RefreshToken string // Provider refresh token (if available)
}

// getJSON GETs url with client and decodes a 200 JSON response into result.
// When accessToken is non-empty it is sent as a Bearer token; otherwise the
// client is expected to authenticate (e.g. one from oauth2.Config.Client).
func getJSON(ctx context.Context, client *http.Client, url, accessToken string, result any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	if accessToken != "" {
		req.Header.Set("Authorization", "Bearer "+accessToken)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("get %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return statusError(url, resp)
	}
	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return fmt.Errorf("decode %s: %w", url, err)
	}
	return nil
}

// statusError describes a non-200 response, including a bounded prefix of
// its body.
func statusError(url string, resp *http.Response) error {
	body, err := io.ReadAll(io.LimitReader(resp.Body, maxErrorBody))
	if err != nil {
		return fmt.Errorf("get %s: status %d (reading body: %w)", url, resp.StatusCode, err)
	}
	return fmt.Errorf("get %s: status %d: %s", url, resp.StatusCode, body)
}
