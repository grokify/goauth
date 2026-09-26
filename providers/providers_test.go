package providers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"golang.org/x/oauth2"
)

const testAccessToken = "test-access-token" //nolint:gosec // G101: test fixture, not a credential

// newProviderServer serves a token endpoint at /token plus the given JSON
// responses, and fails the test if a provider API call lacks the token.
func newProviderServer(t *testing.T, responses map[string]any) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path == "/token" {
			writeJSON(t, w, map[string]string{
				"access_token":  testAccessToken,
				"refresh_token": "test-refresh-token",
				"token_type":    "Bearer",
			})
			return
		}
		if got := r.Header.Get("Authorization"); got != "Bearer "+testAccessToken {
			t.Errorf("%s: Authorization = %q", r.URL.Path, got)
		}
		body, ok := responses[r.URL.Path]
		if !ok {
			http.Error(w, `{"message":"not found"}`, http.StatusNotFound)
			return
		}
		writeJSON(t, w, body)
	}))
	t.Cleanup(srv.Close)
	return srv
}

func writeJSON(t *testing.T, w http.ResponseWriter, v any) {
	t.Helper()
	if err := json.NewEncoder(w).Encode(v); err != nil {
		t.Errorf("encode response: %v", err)
	}
}

func testConfig(srv *httptest.Server) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     "client-id",
		ClientSecret: "client-secret",
		Endpoint: oauth2.Endpoint{
			TokenURL:  srv.URL + "/token",
			AuthStyle: oauth2.AuthStyleInParams,
		},
	}
}

func TestFetchGitHubUser(t *testing.T) {
	cases := []struct {
		name      string
		user      GitHubUserInfo
		emails    []GitHubEmail
		wantEmail string
		wantName  string
		wantErr   error
	}{{
		name:      "public profile email",
		user:      GitHubUserInfo{ID: 42, Login: "octo", Name: "Octo Cat", Email: "octo@example.com"},
		wantEmail: "octo@example.com",
		wantName:  "Octo Cat",
	}, {
		name: "private email falls back to emails endpoint; name falls back to login",
		user: GitHubUserInfo{ID: 42, Login: "octo"},
		emails: []GitHubEmail{
			{Email: "old@example.com", Verified: true},
			{Email: "primary@example.com", Primary: true, Verified: true},
		},
		wantEmail: "primary@example.com",
		wantName:  "octo",
	}, {
		name:    "no verified email",
		user:    GitHubUserInfo{ID: 42, Login: "octo"},
		emails:  []GitHubEmail{{Email: "x@example.com", Primary: true}},
		wantErr: ErrNoVerifiedEmail,
	}}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			srv := newProviderServer(t, map[string]any{"/user": tc.user, "/user/emails": tc.emails})
			ep := gitHubEndpoints{user: srv.URL + "/user", emails: srv.URL + "/user/emails"}

			got, err := fetchGitHubUser(context.Background(), testConfig(srv), "code", ep)
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("err = %v, want %v", err, tc.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			want := OAuthUser{
				ProviderID:   "42",
				Email:        tc.wantEmail,
				Name:         tc.wantName,
				Provider:     ProviderGitHub,
				AccessToken:  testAccessToken,
				RefreshToken: "test-refresh-token",
			}
			if *got != want {
				t.Errorf("got %+v, want %+v", *got, want)
			}
		})
	}
}

func TestFetchGoogleUser(t *testing.T) {
	info := GoogleUserInfo{ID: "g-1", Email: "a@example.com", Name: "A B", Picture: "https://example.com/a.png"}
	srv := newProviderServer(t, map[string]any{"/userinfo": info})

	got, err := fetchGoogleUser(context.Background(), testConfig(srv), "code", srv.URL+"/userinfo")
	if err != nil {
		t.Fatal(err)
	}
	want := OAuthUser{
		ProviderID:   "g-1",
		Email:        "a@example.com",
		Name:         "A B",
		AvatarURL:    "https://example.com/a.png",
		Provider:     ProviderGoogle,
		AccessToken:  testAccessToken,
		RefreshToken: "test-refresh-token",
	}
	if *got != want {
		t.Errorf("got %+v, want %+v", *got, want)
	}
}

func TestFetchUserWithToken(t *testing.T) {
	srv := newProviderServer(t, map[string]any{
		"/user":     GitHubUserInfo{ID: 7, Login: "octo"},
		"/userinfo": GoogleUserInfo{ID: "g-7"},
	})
	ctx := context.Background()

	gh, err := fetchGitHubUserWithToken(ctx, srv.Client(), srv.URL+"/user", testAccessToken)
	if err != nil {
		t.Fatal(err)
	}
	if gh.ID != 7 || gh.Login != "octo" {
		t.Errorf("GitHub user = %+v", gh)
	}

	g, err := fetchGoogleUserWithToken(ctx, srv.Client(), srv.URL+"/userinfo", testAccessToken)
	if err != nil {
		t.Fatal(err)
	}
	if g.ID != "g-7" {
		t.Errorf("Google user = %+v", g)
	}
}

func TestGetJSONNonOK(t *testing.T) {
	srv := newProviderServer(t, nil)
	var v map[string]any
	err := getJSON(context.Background(), srv.Client(), srv.URL+"/missing", testAccessToken, &v)
	if err == nil {
		t.Fatal("expected error for 404")
	}
	if !strings.Contains(err.Error(), "status 404") || !strings.Contains(err.Error(), "not found") {
		t.Errorf("error lacks status/body: %v", err)
	}
}

func TestPrimaryVerifiedEmail(t *testing.T) {
	cases := []struct {
		name   string
		emails []GitHubEmail
		want   string
	}{
		{"primary verified wins", []GitHubEmail{{Email: "a", Verified: true}, {Email: "b", Primary: true, Verified: true}}, "b"},
		{"unverified primary skipped", []GitHubEmail{{Email: "a", Primary: true}, {Email: "b", Verified: true}}, "b"},
		{"none verified", []GitHubEmail{{Email: "a", Primary: true}}, ""},
		{"empty", nil, ""},
	}
	for _, tc := range cases {
		got, err := PrimaryVerifiedEmail(tc.emails)
		if tc.want == "" {
			if !errors.Is(err, ErrNoVerifiedEmail) {
				t.Errorf("%s: err = %v, want ErrNoVerifiedEmail", tc.name, err)
			}
			continue
		}
		if err != nil || got != tc.want {
			t.Errorf("%s: got (%q, %v), want %q", tc.name, got, err, tc.want)
		}
	}
}
