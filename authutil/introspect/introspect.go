package introspect

import (
	"encoding/json"
	"net/http"
)

// IntrospectResponse is defined in RFC-7662: https://datatracker.ietf.org/doc/html/rfc7662
type IntrospectResponse struct {
	Active    bool   `json:"active"`
	Aud       string `json:"aud,omitempty"`
	ClientID  string `json:"client_id,omitempty"`
	Exp       int    `json:"exp,omitempty"`
	Iat       int    `json:"iat,omitempty"`
	Iss       string `json:"iss,omitempty"`
	Jti       string `json:"jti,omitempty"`
	Nbf       int    `json:"nbf,omitempty"`
	Scope     string `json:"scope,omitempty"`
	Sub       string `json:"sub,omitempty"`
	TokenType string `json:"token_type,omitempty"`
	Username  string `json:"username,omitempty"`
}

func (ir IntrospectResponse) Clone() IntrospectResponse {
	return IntrospectResponse{
		Active:    ir.Active,
		Aud:       ir.Aud,
		ClientID:  ir.ClientID,
		Exp:       ir.Exp,
		Iat:       ir.Iat,
		Iss:       ir.Iss,
		Jti:       ir.Jti,
		Nbf:       ir.Nbf,
		Scope:     ir.Scope,
		Sub:       ir.Sub,
		TokenType: ir.TokenType,
		Username:  ir.Username}
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
		w.Write(body)
	}
}

func (ms MockServer) NewServeMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /introspect", ms.PostIntrospect)
	return mux
}

func (ms MockServer) ListenAndServe(addr string) {
	http.ListenAndServe(addr, ms.NewServeMux())
}
