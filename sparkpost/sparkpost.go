package sparkpost

import (
	sp "github.com/SparkPost/gosparkpost"
)

const BaseUrl = "https://api.sparkpost.com"

func NewConfig(apiKey string) *sp.Config {
	return &sp.Config{
		BaseUrl:    BaseUrl,
		ApiKey:     apiKey,
		ApiVersion: 1}
}

func NewApiClient(apiKey string) (sp.Client, error) {
	var client sp.Client
	err := client.Init(NewConfig(apiKey))
	return client, err
}
