package introspect

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/grokify/mogo/net/http/httputilmore"
)

// IntrospectResponse is defined in RFC-7662: https://datatracker.ietf.org/doc/html/rfc7662
type IntrospectResponse struct {
	Active     bool   `json:"active"`
	Audience   string `json:"aud,omitempty"`
	ClientID   string `json:"client_id,omitempty"`
	Expiration int    `json:"exp,omitempty"`
	IssuedAt   int    `json:"iat,omitempty"`
	Issuer     string `json:"iss,omitempty"`
	JWTID      string `json:"jti,omitempty"`
	NotBefore  int    `json:"nbf,omitempty"`
	Scope      string `json:"scope,omitempty"`
	Subject    string `json:"sub,omitempty"`
	TokenType  string `json:"token_type,omitempty"`
	Username   string `json:"username,omitempty"`
}

func (ir IntrospectResponse) Clone() IntrospectResponse {
	return IntrospectResponse{
		Active:     ir.Active,
		Audience:   ir.Audience,
		ClientID:   ir.ClientID,
		Expiration: ir.Expiration,
		IssuedAt:   ir.IssuedAt,
		Issuer:     ir.Issuer,
		JWTID:      ir.JWTID,
		NotBefore:  ir.NotBefore,
		Scope:      ir.Scope,
		Subject:    ir.Subject,
		TokenType:  ir.TokenType,
		Username:   ir.Username}
}

// MockServer is a mock server that implements RFC-7662 OAuth 2.0 Introspection API endpoint for testing purposes.
type MockServer struct {
	Response     IntrospectResponse
	ActiveTokens map[string]int
}

func NewMockServer(r IntrospectResponse, activeTokens []string) MockServer {
	m := map[string]int{}
	for _, t := range activeTokens {
		m[t]++
	}
	return MockServer{
		Response:     r,
		ActiveTokens: m}
}

func (ms MockServer) PostIntrospect(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resp := ms.Response.Clone()
	if _, ok := ms.ActiveTokens[r.PostFormValue("token")]; ok {
		resp.Active = true
	} else {
		resp.Active = false
	}

	if body, err := json.Marshal(resp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write(body)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

func (ms MockServer) NewServeMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /introspect", ms.PostIntrospect)
	return mux
}

func (ms MockServer) ListenAndServe(addr string) error {
	return httputilmore.NewServerTimeouts(addr, ms.NewServeMux(), time.Second).ListenAndServe()
}
