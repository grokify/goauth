package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/grokify/gotilla/fmt/fmtutil"
	"github.com/grokify/oauth2util-go/services/ringcentral"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

func handler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Hello World %s!\n", req.URL.Path[1:])
}

func printString(w http.ResponseWriter, s string) {
	fmt.Fprintln(w, s)
}

func NewClient(authCode string) *http.Client {
	return &http.Client{}
}

func oauth2Handler(w http.ResponseWriter, req *http.Request) {
	authCode := req.FormValue("code")
	log.WithFields(log.Fields{
		"oauth2": "authCodeReceived",
	}).Info(authCode)

	c := ringcentral.ApplicationCredentials{
		ClientID:     os.Getenv("RINGCENTRAL_CLIENT_ID"),
		ClientSecret: os.Getenv("RINGCENTRAL_CLIENT_SECRET"),
		ServerURL:    os.Getenv("RINGCENTRAL_SERVER_URL"),
		RedirectURL:  os.Getenv("RINGCENTRAL_REDIRECT_URL")}
	o2Config := c.Config()

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

func main() {
	err := godotenv.Load()
	if err != nil {
		log.WithFields(log.Fields{
			"config": "dotenvLoadingError",
		}).Fatal(err.Error())
	}

	fmt.Println(os.Getenv("RINGCENTRAL_REDIRECT_URL"))

	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)

	http.HandleFunc("/", handler)
	http.HandleFunc("/oauth2callback", oauth2Handler)
	http.ListenAndServe(":8080", nil)
}
