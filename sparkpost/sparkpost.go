package sparkpost

import (
	"github.com/SparkPost/gosparkpost"
)

const BaseUrl = "https://api.sparkpost.com"

func NewConfig(apiKey string) *gosparkpost.Config {
	return &gosparkpost.Config{
		BaseUrl:    BaseUrl, //nolint:gosec
		ApiKey:     apiKey,
		ApiVersion: 1}
}

func NewAPIClient(apiKey string) (gosparkpost.Client, error) {
	var client gosparkpost.Client
	err := client.Init(NewConfig(apiKey))
	return client, err
}
