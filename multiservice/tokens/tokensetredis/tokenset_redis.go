package tokensetredis

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/grokify/goauth/multiservice/tokens"
	"github.com/grokify/sogo/database/kvs"
	"github.com/grokify/sogo/database/kvs/redis"
	"golang.org/x/oauth2"
)

type TokenSet struct {
	gsClient kvs.Client
}

func NewTokenSet(client *redis.Client) *TokenSet {
	return &TokenSet{gsClient: client}
}

func (toks *TokenSet) GetToken(ctx context.Context, key string) (*oauth2.Token, error) {
	if tokInfo, err := toks.GetTokenInfo(ctx, key); err != nil {
		return nil, err
	} else {
		return tokInfo.Token, nil
	}
}

func (toks *TokenSet) GetTokenInfo(ctx context.Context, key string) (*tokens.TokenInfo, error) {
	key = tokens.FormatKey(key)
	data := toks.gsClient.GetOrDefaultString(ctx, key, "")
	if len(strings.TrimSpace(data)) == 0 {
		return nil, fmt.Errorf("no token for [%v]", key)
	}
	return tokens.ParseTokenInfo([]byte(data))
}

func (toks *TokenSet) SetTokenInfo(ctx context.Context, key string, tok *tokens.TokenInfo) error {
	if bytes, err := json.Marshal(tok); err != nil {
		return err
	} else {
		return toks.gsClient.SetString(ctx, tokens.FormatKey(key), string(bytes))
	}
}
