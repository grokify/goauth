package main

import (
	"fmt"
	"os"

	sp "github.com/SparkPost/gosparkpost"
	"github.com/grokify/oauth2more/sparkpost"
	"github.com/grokify/simplego/config"

	"github.com/rs/zerolog/log"
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
		log.Fatal().Err(err)
	}
	log.Info().
		Str("email-id", id).
		Msg("sparkpost email")
}

func main() {
	err := config.LoadDotEnvSkipEmpty(os.Getenv("ENV_PATH"), "./.env")
	if err != nil {
		log.Fatal().Err(err).
			Msg("Load env files failed")
	}

	client, err := sparkpost.NewApiClient(os.Getenv("SPARKPOST_API_KEY"))
	if err != nil {
		log.Fatal().Err(err).
			Msg("SparkPost client init faile")
	}

	sendTestEmail(client)

	fmt.Println("DONE")
}
