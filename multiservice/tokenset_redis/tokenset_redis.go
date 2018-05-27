package redis

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/grokify/gostor"
	rds "github.com/grokify/gostor/redis"
	"github.com/grokify/oauth2more/multiservice/common"
	"golang.org/x/oauth2"
)

type TokenSet struct {
	gsClient gostor.Client
}

func NewTokenSet(client rds.Client) *TokenSet {
	return &TokenSet{gsClient: client}
}

func (toks *TokenSet) GetToken(key string) (*oauth2.Token, error) {
	tokInfo, err := toks.GetTokenInfo(key)
	if err != nil {
		return nil, err
	}
	return tokInfo.Token, err
}

func (toks *TokenSet) GetTokenInfo(key string) (*common.TokenInfo, error) {
	key = common.FormatKey(key)
	data := toks.gsClient.GetOrEmptyString(key)
	if len(strings.TrimSpace(data)) == 0 {
		return nil, fmt.Errorf("No token for [%v]", key)
	}
	return common.ParseTokenInfo([]byte(data))
}

func (toks *TokenSet) SetTokenInfo(key string, tok *common.TokenInfo) error {
	key = common.FormatKey(key)
	if bytes, err := json.Marshal(tok); err != nil {
		return err
	} else {
		return toks.gsClient.SetString(key, string(bytes))
	}
}
