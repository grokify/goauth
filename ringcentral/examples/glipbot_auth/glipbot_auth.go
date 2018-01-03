package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/grokify/gotilla/fmt/fmtutil"
	"github.com/grokify/oauth2more/ringcentral"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

func handler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Hello World %s!\n", req.URL.Path[1:])
}

func printString(w http.ResponseWriter, s string) {
	fmt.Fprintln(w, s)
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

func oauth2Handler(w http.ResponseWriter, req *http.Request) {
	authCode := req.FormValue("code")
	log.WithFields(log.Fields{
		"oauth2": "authCodeReceived",
	}).Info(authCode)

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
	log.WithFields(log.Fields{
		"oauth2": "token",
	}).Info(string(bytes))
	printString(w, fmt.Sprintf("TOKEN:\n%v\n", string(bytes)))

	client := o2Config.Client(oauth2.NoContext, tok)

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

	http.HandleFunc("/", handler)
	http.HandleFunc("/oauth2callback", oauth2Handler)
	http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("RINGCENTRAL_PORT")), nil)
}
