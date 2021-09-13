# Glipbot Example

`glipbot_auth.go` is a simple HTTP server that handles Glip bot provisioning. Upon successfully adding this bot to Glip, this will log the permanent access token provided by the RingCentral bot provisioner. The token can then be used with the bot of your choice.

This is to bootstrap your bot during the development process. If you deploy a public bot for many users, it is still important to implement the bot provisioner flow directly into your bot.

## Instructions

### Create a bot app

Go to the RingCentral Developer Portal and create a bot:

1. Go to: https://developer.ringcentral.com
1. Create an app
1. Set **Platform Type** to `Server/Bot`
1. Set permissions to `Glip`, `Read Accounts`, `Webhook Subscriptions`
1. Add the redirect URL
  1. Note: while other apps can have more than one redirect URL, a bot can only have one. You will receive an error during provisioning if more than one is provided.
  1. Note: If you are developing locally and using ngrok tunneling, start ngrok for the port you will be running locally. By default, this app runs on port 8080. For example: https://abcd1234.ngrok.io/oauth2callback

### Start the app

1. Clone the app:

`$ go get github.com/grokify/oauth2more`

2. Change directory to the example file:

`$ cd $GOPATH/src/github.com/grokify/oauth2more/ringcentral/examples/glipbot_auth`

3. Create the `.env` file

`$ cp sample.env .env`
`$ vi .env`

4. Start the app

`$ go run glipbot_auth.go`

### Deploy the bot and get the token