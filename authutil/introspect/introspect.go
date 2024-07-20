package introspect

import (
	"encoding/json"
	"net/http"
)

// IntrospectResponse is defined in RFC-7662: https://datatracker.ietf.org/doc/html/rfc7662
type IntrospectResponse struct {
	Active   bool   `json:"active"`
	Scope    string `json:"scope,omitempty"`
	ClientID string `json:"client_id,omitempty"`
	Username string `json:"username,omitempty"`
	Exp      int    `json:"exp,omitempty"`
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
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token := r.PostFormValue("token")
	if _, ok := ms.ActiveTokens[token]; ok {
		ms.Response.Active = true
	} else {
		ms.Response.Active = false
	}

	body, err := json.Marshal(ms.Response)
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func (ms MockServer) NewServeMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /introspect", ms.PostIntrospect)
	return mux
}

func (ms MockServer) ListenAndServe(addr string) {
	http.ListenAndServe(addr, ms.NewServeMux())
}
