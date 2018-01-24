package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/grokify/gotilla/fmt/fmtutil"
	"github.com/grokify/oauth2more/ringcentral"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

func handleHelloWorld(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "RingCentral Glip OAuth Bootstrap Bot %s!\n", req.URL.Path[1:])
}

func handleBotButton(w http.ResponseWriter, req *http.Request) {
	body := `<!DOCTYPE html>
	<html><body>
	<h1>Glip Bootstrap Bot</h1>
	<a href="https://apps1.ringcentral.com/app/Vh3KiNFGQ86JclwIIBcqIA~i4w85hFTTtaLpLwch0z_OA/install?landing_url=https%3A%2F%2F0a6f4754.ngrok.io" target="_blank" style="box-sizing:border-box;display: inline-block;border: 1px solid #0073ae;border-radius: 4px;text-decoration: none;height: 60px;line-height: 60px;width: 160px;padding-left: 20px;font-size: 14px;color:#0073ae;font-family:"Lato",Helvetica,Arial,sans-serif"><span>Add to </span><img style="width: 68px;vertical-align: middle;display: inline-block;margin-left: 10px;" src="http://netstorage.ringcentral.com/dpw/common/glip/logo_glip.png"></a>
	</body></html>`
	fmt.Fprintf(w, body, req.URL.Path[1:])
}

func handleOauth2(w http.ResponseWriter, req *http.Request) {
	// Retrieve auth code from URL
	authCode := req.FormValue("code")
	log.WithFields(log.Fields{
		"oauth2": "authCodeReceived",
	}).Info(authCode)

	// Exchange auth code for token
	o2Config := getOauth2Config()

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

func getOauth2Config() oauth2.Config {
	c := ringcentral.ApplicationCredentials{
		ClientID:     os.Getenv("RINGCENTRAL_CLIENT_ID"),
		ClientSecret: os.Getenv("RINGCENTRAL_CLIENT_SECRET"),
		ServerURL:    os.Getenv("RINGCENTRAL_SERVER_URL"),
		RedirectURL:  os.Getenv("RINGCENTRAL_REDIRECT_URL")}
	o2Config := c.Config()
	return o2Config
}

func printString(w http.ResponseWriter, s string) {
	fmt.Fprintln(w, s)
}

func loadEnv() error {
	envPaths := []string{}
	if len(os.Getenv("ENV_PATH")) > 0 {
		log.WithFields(log.Fields{
			"Note": "Found dotenv path",
		}).Info(os.Getenv("ENV_PATH"))

		envPaths = append(envPaths, os.Getenv("ENV_PATH"))
	}
	return godotenv.Load(envPaths...)
}

func main() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)

	err := loadEnv()
	if err != nil {
		log.WithFields(log.Fields{
			"config": "dotenvLoadingError",
		}).Fatal(err.Error())
	}

	log.WithFields(log.Fields{
		"BotRedirectUrl": "redirect URL",
	}).Info(os.Getenv("RINGCENTRAL_REDIRECT_URL"))

	http.HandleFunc("/", handleBotButton)
	http.HandleFunc("/oauth2callback", handleOauth2)
	http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("RINGCENTRAL_PORT")), nil)
}
