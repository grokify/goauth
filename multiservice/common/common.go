package common

import (
	"encoding/json"
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
