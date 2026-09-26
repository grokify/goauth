package providers

import (
	"context"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
)

const (
	ProviderGoogle    = "google"
	GoogleUserInfoURL = "https://www.googleapis.com/oauth2/v2/userinfo"
)

// GoogleUserInfo represents the response from Google's userinfo endpoint.
type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
}

// FetchGoogleUser exchanges the authorization code and fetches user info from Google.
func FetchGoogleUser(ctx context.Context, config *oauth2.Config, code string) (*OAuthUser, error) {
	return fetchGoogleUser(ctx, config, code, GoogleUserInfoURL)
}

func fetchGoogleUser(ctx context.Context, config *oauth2.Config, code, userInfoURL string) (*OAuthUser, error) {
	token, err := config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("exchange code: %w", err)
	}

	var userInfo GoogleUserInfo
	if err := getJSON(ctx, config.Client(ctx, token), userInfoURL, "", &userInfo); err != nil {
		return nil, fmt.Errorf("fetch Google user: %w", err)
	}

	return &OAuthUser{
		ProviderID:   userInfo.ID,
		Email:        userInfo.Email,
		Name:         userInfo.Name,
		AvatarURL:    userInfo.Picture,
		Provider:     ProviderGoogle,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}, nil
}

// FetchGoogleUserWithToken fetches user info using an existing access token.
func FetchGoogleUserWithToken(ctx context.Context, accessToken string) (*GoogleUserInfo, error) {
	return fetchGoogleUserWithToken(ctx, defaultClient, GoogleUserInfoURL, accessToken)
}

func fetchGoogleUserWithToken(ctx context.Context, client *http.Client, userInfoURL, accessToken string) (*GoogleUserInfo, error) {
	var userInfo GoogleUserInfo
	if err := getJSON(ctx, client, userInfoURL, accessToken, &userInfo); err != nil {
		return nil, fmt.Errorf("fetch Google user: %w", err)
	}
	return &userInfo, nil
}
