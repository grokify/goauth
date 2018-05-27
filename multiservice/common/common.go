package common

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"golang.org/x/oauth2"
)

type TokenInfo struct {
	ServiceKey  string        `json:"serviceKey,omitempty"`
	ServiceType string        `json:"serviceType,omitempty"`
	Token       *oauth2.Token `json:"token,omitempty"`
}

func FormatKey(key string) string { return strings.TrimSpace(key) }

type TokenSet interface {
	GetTokenInfo(key string) (*TokenInfo, error)
	GetToken(key string) (*oauth2.Token, error)
	SetTokenInfo(key string, tokenInfo *TokenInfo) error
}

func ParseTokenInfo(data []byte) (*TokenInfo, error) {
	tok := &TokenInfo{}
	return tok, json.Unmarshal(data, tok)
}

func NewClientWithTokenSet(
	ctx context.Context, conf *oauth2.Config, token *oauth2.Token,
	tokenSet TokenSet, tokenKey, serviceKey, serviceType string,
) (*http.Client, error) {

	tokenSource := conf.TokenSource(ctx, token)

	if newToken, err := tokenSource.Token(); err != nil {
		return nil, err
	} else if newToken.AccessToken != token.AccessToken {
		if err := tokenSet.SetTokenInfo(tokenKey, &TokenInfo{
			ServiceKey:  serviceKey,
			ServiceType: serviceType,
			Token:       newToken}); err != nil {
			return nil, err
		}
	}
	return oauth2.NewClient(ctx, tokenSource), nil
}
