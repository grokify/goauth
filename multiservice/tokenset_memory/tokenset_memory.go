package memory

import (
	"errors"
	"fmt"

	"github.com/grokify/oauth2more/multiservice/common"
	"golang.org/x/oauth2"
)

type TokenSet struct {
	tokenMap map[string]*common.TokenInfo
}

func NewTokenSet() *TokenSet {
	return &TokenSet{tokenMap: map[string]*common.TokenInfo{}}
}

func (toks *TokenSet) GetToken(key string) (*oauth2.Token, error) {
	key = common.FormatKey(key)
	if tok, ok := toks.tokenMap[key]; ok {
		return tok.Token, nil
	}
	return nil, fmt.Errorf("AppConfig not found for %v", key)
}

func (toks *TokenSet) GetTokenInfo(key string) (*common.TokenInfo, error) {
	key = common.FormatKey(key)
	if tok, ok := toks.tokenMap[key]; ok {
		return tok, nil
	}
	return nil, fmt.Errorf("AppConfig not found for %v", key)
}

func (toks *TokenSet) SetTokenInfo(key string, tok *common.TokenInfo) error {
	key = common.FormatKey(key)
	if len(key) == 0 {
		return errors.New("Set Token Requires Non-Empty Key")
	}
	toks.tokenMap[key] = tok
	return nil
}
