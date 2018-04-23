package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/caarlos0/env"
	"github.com/grokify/gotilla/config"
	"github.com/grokify/gotilla/fmt/fmtutil"
	"github.com/grokify/oauth2more/ringcentral"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

type RingCentralConfig struct {
	AppId        string `env:"RINGCENTRAL_APP_ID"`
	ClientId     string `env:"RINGCENTRAL_CLIENT_ID"`
	ClientSecret string `env:"RINGCENTRAL_CLIENT_SECRET"`
	ServerURL    string `env:"RINGCENTRAL_SERVER_URL"`
	RedirectURL  string `env:"RINGCENTRAL_REDIRECT_URL"`
	LandingURL   string `env:"RINGCENTRAL_LANDING_URL"`
	AppPort      int64  `env:"PORT"`
}

func handleHelloWorld(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "RingCentral Glip OAuth Bootstrap Bot %s!\n", req.URL.Path[1:])
}

type AppHandler struct {
	AppConfig RingCentralConfig
}

func (app *AppHandler) HandleBotButton(w http.ResponseWriter, req *http.Request) {
	bodyFormat := `<!DOCTYPE html>
	<html><body>
	<h1>Glip Bootstrap Bot</h1>
	<a href="https://apps1.ringcentral.com/app/%v/install?landing_url=%v" target="_blank" style="box-sizing:border-box;display: inline-block;border: 1px solid #0073ae;border-radius: 4px;text-decoration: none;height: 60px;line-height: 60px;width: 160px;padding-left: 20px;font-size: 14px;color:#0073ae;font-family:"Lato",Helvetica,Arial,sans-serif"><span>Add to </span><img style="width: 68px;vertical-align: middle;display: inline-block;margin-left: 10px;" src="http://netstorage.ringcentral.com/dpw/common/glip/logo_glip.png"></a>
	</body></html>`
	body := fmt.Sprintf(
		bodyFormat,
		app.AppConfig.AppId,
		url.QueryEscape(app.AppConfig.LandingURL),
	)

	fmt.Fprintf(w, bodyFormat, app.AppConfig.AppId, url.QueryEscape(app.AppConfig.LandingURL))
}

func (app *AppHandler) HandleOauth2(w http.ResponseWriter, req *http.Request) {
	// Retrieve auth code from URL
	authCode := req.FormValue("code")
	log.WithFields(log.Fields{
		"oauth2": "authCodeReceived",
	}).Info(authCode)

	// Exchange auth code for token
	o2Config := getOauth2Config(app.AppConfig)

	tok, err := o2Config.Exchange(oauth2.NoContext, authCode)
	if err != nil {
		log.WithFields(log.Fields{
			"oauth2": "tokenExchangeError",
		}).Info(err.Error())

		printString(w, err.Error())
		return
	}

	bytes, err := json.Marshal(tok)
	if err != nil {
		printString(w, err.Error())
		return
	}

	// Log token
	log.WithFields(log.Fields{
		"oauth2": "token",
	}).Info(string(bytes))
	printString(w, fmt.Sprintf("TOKEN:\n%v\n", string(bytes)))

	client := o2Config.Client(oauth2.NoContext, tok)

	// API Call
	cu := ringcentral.NewClientUtil(client)
	u, err := cu.GetSCIMUser()
	if err != nil {
		printString(w, err.Error())
		return
	}
	fmtutil.PrintJSON(u)

	bytes, err = json.Marshal(u)
	if err != nil {
		printString(w, err.Error())
		return
	}
	printString(w, string(bytes))
}

func getOauth2Config(appCfg RingCentralConfig) oauth2.Config {
	app := ringcentral.ApplicationCredentials{
		ClientID:     appCfg.ClientId,
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
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)

	err := config.LoadDotEnvSkipEmpty(os.Getenv("ENV_PATH"), "./.env")
	if err != nil {
		log.WithFields(log.Fields{
			"config": "dotenvLoadingError",
		}).Fatal(err.Error())
	}

	appCfg := RingCentralConfig{}
	err = env.Parse(&appCfg)
	if err != nil {
		log.Fatal(err)
	}
	appHandler := AppHandler{AppConfig: appCfg}

	log.WithFields(log.Fields{
		"BotRedirectUrl": "redirect URL",
	}).Info(appCfg.RedirectURL)
	log.WithFields(log.Fields{
		"BotPort": "Local Server Port URL",
	}).Info(appCfg.AppPort)

	http.HandleFunc("/", appHandler.HandleBotButton)
	http.HandleFunc("/oauth2callback", appHandler.HandleOauth2)
	log.Fatal(
		http.ListenAndServe(
			fmt.Sprintf(":%v", appCfg.AppPort), nil,
		),
	)
}
