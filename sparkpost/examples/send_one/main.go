package main

import (
	"fmt"
	"os"

	sp "github.com/SparkPost/gosparkpost"
	"github.com/grokify/gotilla/config"
	"github.com/grokify/oauth2more/sparkpost"
	log "github.com/sirupsen/logrus"
)

func sendTestEmail(client sp.Client) {
	// Create a Transmission using an inline Recipient List
	// and inline email Content.
	tx := &sp.Transmission{
		Recipients: []string{os.Getenv("SPARKPOST_EMAIL_RECIPIENT")},
		Content: sp.Content{
			HTML:    `<p>Hello World <b>Body</b>!</p>`,
			From:    os.Getenv("SPARKPOST_EMAIL_SENDER"),
			Subject: `Hello World Subject!`,
		},
	}
	id, _, err := client.Send(tx)
	if err != nil {
		log.Fatal(err)
	}
	log.WithFields(log.Fields{"email-id": id}).Info("email")
}

func main() {
	err := config.LoadDotEnvSkipEmpty(os.Getenv("ENV_PATH"), "./.env")
	if err != nil {
		log.Fatalf("Load env files failed: %s\n", err)
	}

	client, err := sparkpost.NewApiClient(os.Getenv("SPARKPOST_API_KEY"))
	if err != nil {
		log.Fatalf("SparkPost client init failed: %s\n", err)
	}

	sendTestEmail(client)

	fmt.Println("DONE")
}
