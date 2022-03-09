package tokensetredis

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/grokify/goauth/multiservice/tokens"
	"github.com/grokify/gostor"
	rds "github.com/grokify/gostor/redis"
	"golang.org/x/oauth2"
)

type TokenSet struct {
	gsClient gostor.Client
}

func NewTokenSet(client *rds.Client) *TokenSet {
	return &TokenSet{gsClient: client}
}

func (toks *TokenSet) GetToken(key string) (*oauth2.Token, error) {
	if tokInfo, err := toks.GetTokenInfo(key); err != nil {
		return nil, err
	} else {
		return tokInfo.Token, nil
	}
}

func (toks *TokenSet) GetTokenInfo(key string) (*tokens.TokenInfo, error) {
	key = tokens.FormatKey(key)
	data := toks.gsClient.GetOrEmptyString(key)
	if len(strings.TrimSpace(data)) == 0 {
		return nil, fmt.Errorf("no token for [%v]", key)
	}
	return tokens.ParseTokenInfo([]byte(data))
}

func (toks *TokenSet) SetTokenInfo(key string, tok *tokens.TokenInfo) error {
	if bytes, err := json.Marshal(tok); err != nil {
		return err
	} else {
		return toks.gsClient.SetString(tokens.FormatKey(key), string(bytes))
	}
}
