package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/bmizerany/pat"
	"github.com/caarlos0/env"
	"github.com/grokify/simplego/config"
	hum "github.com/grokify/simplego/net/httputilmore"
	log "github.com/sirupsen/logrus"

	"github.com/grokify/oauth2more/auth0"
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
	log.Debug("START_LOGIN_HANDLER")
	authUrlInfo := auth0.PKCEAuthorizationUrlInfo{
		Host:        cfg.Host,
		Scope:       cfg.Scope,
		ClientId:    cfg.ClientId,
		RedirectUri: cfg.RedirectUri}

	verifier, challenge, authUrl, err := authUrlInfo.Data()
	if err != nil {
		log.Fatal(err)
	}
	log.WithFields(log.Fields{
		"remoteAddr": r.RemoteAddr,
		"userAgent":  r.Header.Get(hum.HeaderUserAgent),
		"authUrl":    authUrl,
		"challenge":  challenge,
		"verifier":   verifier},
	).Debug("loginHandler")

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
	w.Header().Set(hum.HeaderContentType, hum.ContentTypeTextHtmlUtf8)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, tmpl, WebsiteTitle, WebsiteTitle, verifier, challenge, authUrl)
	log.Debug("END_LOGIN_HANDLER")
}

func (cfg *appConfig) Oauth2CallbackHandler(w http.ResponseWriter, r *http.Request) {
	log.Debug("START_OAUTH2CALLBACK_HANDLER")
	codeArr, ok := r.URL.Query()["code"]
	if !ok {
		log.Fatal("E_NO_CODE")
	}
	code := ""
	if len(codeArr) > 0 {
		code = strings.TrimSpace(codeArr[0])
	}

	cookie, err := r.Cookie(VerifierCookieName)
	if err != nil {
		log.Fatal(err)
	}
	tokenUrlInfo := auth0.PKCETokenUrlInfo{
		Host:         cfg.Host,
		GrantType:    "authorization_code",
		ClientId:     cfg.ClientId,
		CodeVerifier: cookie.Value,
		Code:         code,
		RedirectUri:  cfg.RedirectUri}

	resp, err := tokenUrlInfo.Exchange()
	if err != nil {
		log.Fatal(err)
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.WithFields(log.Fields{
		"verifier":           cookie.Value,
		"tokenResStatusCode": resp.StatusCode,
		"tokenResBody":       string(respBody)},
	).Info("oauth2CallbackHandler")

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

	w.Header().Set(hum.HeaderContentType, hum.ContentTypeTextHtmlUtf8)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, tmpl, WebsiteTitle, WebsiteTitle, string(respBody))
	log.Debug("END_LOGIN_HANDLER")
}

func main() {
	log.SetLevel(log.DebugLevel)
	if err := config.LoadDotEnvSkipEmpty(os.Getenv("ENV_PATH"), "./.env"); err != nil {
		panic(err)
	}

	cfg := appConfig{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err)
	}

	m := pat.New()
	m.Get("/", http.HandlerFunc(cfg.LoginHandler))
	m.Get("/oauth2callback", http.HandlerFunc(cfg.Oauth2CallbackHandler))
	http.Handle("/", m)

	log.Fatal(http.ListenAndServe(cfg.PortString(), nil))
}
