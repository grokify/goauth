package providers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"golang.org/x/oauth2"
)

const (
	ProviderGitHub  = "github"
	GitHubUserURL   = "https://api.github.com/user"
	GitHubEmailsURL = "https://api.github.com/user/emails"
)

// ErrNoVerifiedEmail is returned when a GitHub user has no verified email.
var ErrNoVerifiedEmail = errors.New("no verified email found")

// GitHubUserInfo represents the response from GitHub's user endpoint.
type GitHubUserInfo struct {
	ID        int64  `json:"id"`
	Login     string `json:"login"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}

// GitHubEmail represents an email from GitHub's emails endpoint.
type GitHubEmail struct {
	Email    string `json:"email"`
	Primary  bool   `json:"primary"`
	Verified bool   `json:"verified"`
}

type gitHubEndpoints struct {
	user, emails string
}

var defaultGitHubEndpoints = gitHubEndpoints{user: GitHubUserURL, emails: GitHubEmailsURL}

// FetchGitHubUser exchanges the authorization code and fetches user info from
// GitHub. When the profile email is private, the primary verified email is
// fetched from the emails endpoint, which requires the "user:email" scope.
func FetchGitHubUser(ctx context.Context, config *oauth2.Config, code string) (*OAuthUser, error) {
	return fetchGitHubUser(ctx, config, code, defaultGitHubEndpoints)
}

func fetchGitHubUser(ctx context.Context, config *oauth2.Config, code string, ep gitHubEndpoints) (*OAuthUser, error) {
	token, err := config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("exchange code: %w", err)
	}
	client := config.Client(ctx, token)

	var userInfo GitHubUserInfo
	if err := getJSON(ctx, client, ep.user, "", &userInfo); err != nil {
		return nil, fmt.Errorf("fetch GitHub user: %w", err)
	}

	email := userInfo.Email
	if email == "" {
		var emails []GitHubEmail
		if err := getJSON(ctx, client, ep.emails, "", &emails); err != nil {
			return nil, fmt.Errorf("fetch GitHub emails: %w", err)
		}
		if email, err = PrimaryVerifiedEmail(emails); err != nil {
			return nil, err
		}
	}

	name := userInfo.Name
	if name == "" {
		name = userInfo.Login
	}

	return &OAuthUser{
		ProviderID:   strconv.FormatInt(userInfo.ID, 10),
		Email:        email,
		Name:         name,
		AvatarURL:    userInfo.AvatarURL,
		Provider:     ProviderGitHub,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}, nil
}

// PrimaryVerifiedEmail returns the primary verified email, falling back to
// any verified email, or ErrNoVerifiedEmail.
func PrimaryVerifiedEmail(emails []GitHubEmail) (string, error) {
	for _, e := range emails {
		if e.Primary && e.Verified {
			return e.Email, nil
		}
	}
	for _, e := range emails {
		if e.Verified {
			return e.Email, nil
		}
	}
	return "", ErrNoVerifiedEmail
}

// FetchGitHubUserWithToken fetches user info using an existing access token.
func FetchGitHubUserWithToken(ctx context.Context, accessToken string) (*GitHubUserInfo, error) {
	return fetchGitHubUserWithToken(ctx, defaultClient, GitHubUserURL, accessToken)
}

func fetchGitHubUserWithToken(ctx context.Context, client *http.Client, userURL, accessToken string) (*GitHubUserInfo, error) {
	var userInfo GitHubUserInfo
	if err := getJSON(ctx, client, userURL, accessToken, &userInfo); err != nil {
		return nil, fmt.Errorf("fetch GitHub user: %w", err)
	}
	return &userInfo, nil
}
