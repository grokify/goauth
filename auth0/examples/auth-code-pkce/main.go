package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	env "github.com/caarlos0/env/v11"
	"github.com/grokify/mogo/config"
	"github.com/grokify/mogo/net/http/httputilmore"
	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"

	"github.com/grokify/goauth/auth0"
)

const (
	DefaultPort        = 8080
	WebsiteTitle       = "Auth0 PKCE Demo in Go"
	VerifierCookieName = "auth0verifier"
)

type appConfig struct {
	Port        int    `env:"PORT"` // Set for use with Heroku
	Host        string `env:"AUTH0_HOST"`
	ClientId    string `env:"AUTH0_CLIENT_ID"`
	RedirectUri string `env:"AUTH0_REDIRECT_URI"`
	Scope       string `env:"AUTH0_SCOPE"`
}

func (cfg *appConfig) PortString() string {
	return fmt.Sprintf(":%v", cfg.Port)
}

func (cfg *appConfig) LoginHandler(w http.ResponseWriter, r *http.Request) {
	//log.Debug("START_LOGIN_HANDLER")
	zlog.Debug().Msg("START_LOGIN_HANDLER")
	authUrlInfo := auth0.PKCEAuthorizationURLInfo{
		Host:        cfg.Host,
		Scope:       cfg.Scope,
		ClientID:    cfg.ClientId,
		RedirectURI: cfg.RedirectUri}

	verifier, challenge, authUrl, err := authUrlInfo.Data()
	if err != nil {
		zlog.Fatal().Err(err)
	}
	zlog.Debug().
		Str("remoteAddr", r.RemoteAddr).
		Str("userAgent", r.Header.Get(httputilmore.HeaderUserAgent)).
		Str("authUrl", authUrl).
		Str("challenge", challenge).
		Str("verifier", verifier).
		Msg("loginHandler")

	tmpl := `<!DOCTYPE html>
<html>
  <head>
    <meta charset="UTF-8">
    <title>%s</title>
  </head>
  <body>
    <h1>%s</h1>
    <p>Verifier: %s</p>
    <p>Challenge: %s</p>
    <p><a href="%s">Login</a></p>
  </body>
</html>`

	// Cookie is used for demo purposes only. Use a server-side store
	// in production.
	cookie := http.Cookie{Name: VerifierCookieName, Value: verifier}
	http.SetCookie(w, &cookie)
	w.Header().Set(httputilmore.HeaderContentType, httputilmore.ContentTypeTextHTMLUtf8)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, tmpl, WebsiteTitle, WebsiteTitle, verifier, challenge, authUrl)
	zlog.Debug().Msg("END_LOGIN_HANDLER")
}

func (cfg *appConfig) Oauth2CallbackHandler(w http.ResponseWriter, r *http.Request) {
	zlog.Debug().Msg("START_OAUTH2CALLBACK_HANDLER")
	codeArr, ok := r.URL.Query()["code"]
	if !ok {
		zlog.Fatal().Msg("E_NO_CODE")
	}
	code := ""
	if len(codeArr) > 0 {
		code = strings.TrimSpace(codeArr[0])
	}

	cookie, err := r.Cookie(VerifierCookieName)
	if err != nil {
		zlog.Fatal().Err(err)
	}
	tokenUrlInfo := auth0.PKCETokenURLInfo{
		Host:         cfg.Host,
		GrantType:    "authorization_code",
		ClientID:     cfg.ClientId,
		CodeVerifier: cookie.Value,
		Code:         code,
		RedirectURI:  cfg.RedirectUri}

	resp, err := tokenUrlInfo.Exchange()
	if err != nil {
		zlog.Fatal().Err(err)
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		zlog.Fatal().Err(err)
	}
	zlog.Info().
		Int("tokenResStatusCode", resp.StatusCode).
		Str("tokenResBody", string(respBody)).
		Str("verifier(cookie value)", cookie.Value).
		Msg("oauth2CallbackHandler")

	tmpl := `<!DOCTYPE html>
<html>
  <head>
    <meta charset="UTF-8">
    <title>%s</title>
  </head>
  <body>
    <h1>%s</h1>
    <h2>Token</h2>
    <textarea style="width:100%%;height:50px">%s</textarea>
    <p><a href="/">Try Again</a></p>
  </body>
</html>`

	w.Header().Set(httputilmore.HeaderContentType, httputilmore.ContentTypeTextHTMLUtf8)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, tmpl, WebsiteTitle, WebsiteTitle, string(respBody))
	zlog.Debug().Msg("END_LOGIN_HANDLER")
}

func main() {
	zlog.Logger = zlog.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	if _, err := config.LoadDotEnv([]string{os.Getenv("ENV_PATH"), "./.env"}, -1); err != nil {
		panic(err)
	}

	cfg := appConfig{}
	if err := env.Parse(&cfg); err != nil {
		zlog.Fatal().Err(err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(cfg.LoginHandler))
	mux.Handle("/oauth2callback", http.HandlerFunc(cfg.Oauth2CallbackHandler))

	/*
		m := pat.New()
		m.Get("/", http.HandlerFunc(cfg.LoginHandler))
		m.Get("/oauth2callback", http.HandlerFunc(cfg.Oauth2CallbackHandler))
		http.Handle("/", m)
	*/

	log.Fatal(httputilmore.ListenAndServeTimeouts(cfg.PortString(), mux, 10*time.Second))
}
