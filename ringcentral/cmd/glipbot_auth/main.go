package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/grokify/goauth"
	"github.com/grokify/goauth/ringcentral"
	"github.com/grokify/mogo/config"
	"github.com/grokify/mogo/fmt/fmtutil"
	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
)

type RingCentralConfig struct {
	AppID        string `env:"RINGCENTRAL_APP_ID"`
	ClientID     string `env:"RINGCENTRAL_CLIENT_ID"`
	ClientSecret string `env:"RINGCENTRAL_CLIENT_SECRET"`
	ServerURL    string `env:"RINGCENTRAL_SERVER_URL"`
	RedirectURL  string `env:"RINGCENTRAL_REDIRECT_URL"`
	LandingURL   string `env:"RINGCENTRAL_LANDING_URL"`
	AppPort      int64  `env:"PORT"`
}

/*
func handleHelloWorld(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "RingCentral Glip OAuth Bootstrap Bot %s!\n", req.URL.Path[1:])
}
*/

type AppHandler struct {
	AppConfig RingCentralConfig
}

func (app *AppHandler) HandleBotButton(w http.ResponseWriter, req *http.Request) {
	bodyFormat := `<!DOCTYPE html>
	<html><body>
	<h1>Glip Bootstrap Bot</h1>
	<a href="https://apps1.ringcentral.com/app/%s/install?landing_url=%v" target="_blank" style="box-sizing:border-box;display: inline-block;border: 1px solid #0073ae;border-radius: 4px;text-decoration: none;height: 60px;line-height: 60px;width: 160px;padding-left: 20px;font-size: 14px;color:#0073ae;font-family:"Lato",Helvetica,Arial,sans-serif"><span>Add to </span><img style="width: 68px;vertical-align: middle;display: inline-block;margin-left: 10px;" src="http://netstorage.ringcentral.com/dpw/common/glip/logo_glip.png"></a>
	</body></html>`

	fmt.Fprintf(w,
		bodyFormat,
		app.AppConfig.AppID,
		url.QueryEscape(app.AppConfig.LandingURL))
}

func (app *AppHandler) HandleOauth2(w http.ResponseWriter, req *http.Request) {
	// Retrieve auth code from URL
	authCode := req.FormValue("code")
	zlog.Info().
		Str("code", authCode).
		Msg("OAuth2 code receiveds")

	// Exchange auth code for token
	o2Config := getOauth2Config(app.AppConfig)

	tok, err := o2Config.Exchange(context.Background(), authCode)
	if err != nil {
		zlog.Info().
			Err(err).
			Msg("oauth2 tokenExchangeError")

		printString(w, err.Error())
		return
	}

	bytes, err := json.Marshal(tok)
	if err != nil {
		printString(w, err.Error())
		return
	}

	// Log token
	zlog.Info().Str("auth2_token", string(bytes))

	printString(w, fmt.Sprintf("TOKEN:\n%v\n", string(bytes)))

	client := o2Config.Client(context.Background(), tok)

	// API Call
	cu := ringcentral.NewClientUtil(client)
	u, err := cu.GetSCIMUser()
	if err != nil {
		printString(w, err.Error())
		return
	}
	fmtutil.MustPrintJSON(u)

	bytes, err = json.Marshal(u)
	if err != nil {
		printString(w, err.Error())
		return
	}
	printString(w, string(bytes))
}

func getOauth2Config(appCfg RingCentralConfig) oauth2.Config {
	app := goauth.CredentialsOAuth2{
		ClientID:     appCfg.ClientID,
		ClientSecret: appCfg.ClientSecret,
		ServerURL:    appCfg.ServerURL,
		RedirectURL:  appCfg.RedirectURL}
	o2Config := app.Config()
	return o2Config
}

func printString(w http.ResponseWriter, s string) {
	fmt.Fprintln(w, s)
}

func main() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	_, err := config.LoadDotEnv([]string{os.Getenv("ENV_PATH"), "./.env"}, -1)
	if err != nil {
		zlog.Fatal().Err(err).
			Str("config", "dotenvLoadingError")
	}

	appCfg := RingCentralConfig{}
	err = env.Parse(&appCfg)
	if err != nil {
		zlog.Fatal().Err(err)
	}
	appHandler := AppHandler{AppConfig: appCfg}

	zlog.Info().
		Str("BotRedirectUrl", appCfg.RedirectURL).
		Int64("BotPort Local", appCfg.AppPort)

	mux := http.NewServeMux()
	mux.HandleFunc("/", appHandler.HandleBotButton)
	mux.HandleFunc("/oauth2callback", appHandler.HandleOauth2)

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%v", appCfg.AppPort),
		ReadTimeout:       2 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		WriteTimeout:      2 * time.Second,
		IdleTimeout:       30 * time.Second,
		Handler:           mux}

	log.Fatal(srv.ListenAndServe())
}
