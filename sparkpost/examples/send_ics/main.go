// This package posts an ICS file to Gmail. You can update the sequence
// to update the times of the invite.
// See more here: https://stackoverflow.com/questions/50422635/
package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/SparkPost/gosparkpost"
	"github.com/caarlos0/env/v9"
	"github.com/grokify/goauth/sparkpost"
	"github.com/grokify/mogo/config"
	"github.com/grokify/mogo/net/http/httputilmore"
	"github.com/grokify/mogo/time/timeutil"
	"github.com/rs/zerolog/log"
)

const DefaultExampleEmail = "john@example.com"

type appConfig struct {
	SparkPostApiKey             string `env:"SPARKPOST_API_KEY"`
	SparkPostEmailSender        string `env:"SPARKPOST_EMAIL_SENDER"`
	SparkPostEmailRecipientDemo string `env:"SPARKPOST_EMAIL_RECIPIENT_DEMO"`
}

func sendTestEmail(cfg appConfig, client gosparkpost.Client) {
	day := 811 // Change this to set the day. This currently uses single digit month + 2 digit day
	seq := 0   // Sequence index. Start with 0 and increment. This demo doesn't support changing the day when changing sequence.

	if len(strings.TrimSpace(cfg.SparkPostEmailRecipientDemo)) == 0 {
		cfg.SparkPostEmailRecipientDemo = DefaultExampleEmail
	}

	format := `BEGIN:VCALENDAR
VERSION:2.0
PRODID:-//MY COMPANY//Calendar//EN
CALSCALE:GREGORIAN
METHOD:%v
BEGIN:VEVENT
%vUID:shift-%v-emp-128@mycompany.com
DTSTART:20180%vT010000
DTEND:20180%vT020000
DTSTAMP:%v
SEQUENCE:%v
SUMMARY:Morning shift
LOCATION:Morning Location
DESCRIPTION:Morning shift
END:VEVENT
BEGIN:VEVENT
%vUID:shift-%v-emp-128@mycompany.com
DTSTART:20180%vT130000
DTEND:20180%vT140000
DTSTAMP:%v
SEQUENCE:%v
SUMMARY:Night shift
LOCATION:Night Location
DESCRIPTION:Night
END:VEVENT
END:VCALENDAR`

	//attenddee needed for both events
	attendee := fmt.Sprintf(
		"ATTENDEE;ROLE=REQ-PARTICIPANT;CN=%s:MAILTO:%s\n",
		cfg.SparkPostEmailRecipientDemo, cfg.SparkPostEmailRecipientDemo)

	uid1 := day + 1000
	uid2 := uid1 + 1
	dtNow := time.Now().Format(timeutil.ISO8601CompactNoTZ)
	data := fmt.Sprintf(
		format, "REQUEST",
		attendee, uid1, day, day, dtNow, seq,
		attendee, uid2, day, day, dtNow, seq)

	fmt.Println(data)
	panic("z")
	attach := gosparkpost.Attachment{
		MIMEType: httputilmore.ContentTypeTextCalendarUtf8Request,
		B64Data:  base64.StdEncoding.EncodeToString([]byte(data))}

	tx := &gosparkpost.Transmission{
		Recipients: []string{cfg.SparkPostEmailRecipientDemo},
		Content: gosparkpost.Content{
			HTML:        `<p>Your Shifts!</p>`,
			From:        cfg.SparkPostEmailSender,
			Subject:     fmt.Sprintf(`MS%v Shift Schedule!`, day),
			Attachments: []gosparkpost.Attachment{attach},
		},
	}
	id, _, err := client.Send(tx)
	if err != nil {
		log.Fatal().
			Err(err).
			Str("method", "sparkpost client Send()")
	}
	log.Info().
		Err(err).
		Str("method", "sparkpost email client Send()").
		Str("email-id", id)
}

func buildExampleIcs(recipientSmtpAddr string, dtStart, dtEnd time.Time, seq uint) string {
	if len(strings.TrimSpace(recipientSmtpAddr)) == 0 {
		recipientSmtpAddr = DefaultExampleEmail
	}
	// seq = Sequence index. Start with 0 and increment. This demo doesn't support changing the day when changing sequence.

	day := 811 // Change this to set the day. This currently uses single digit month + 2 digit day

	format := `BEGIN:VCALENDAR
VERSION:2.0
PRODID:-//MY COMPANY//Calendar//EN
CALSCALE:GREGORIAN
METHOD:%v
BEGIN:VEVENT
%vUID:shift-%v-emp-128@mycompany.com
DTSTART:20180%vT010000
DTEND:20180%vT020000
DTSTAMP:%v
SEQUENCE:%v
SUMMARY:Morning shift
LOCATION:Morning Location
DESCRIPTION:Morning shift
END:VEVENT
BEGIN:VEVENT
%vUID:shift-%v-emp-128@mycompany.com
DTSTART:20180%vT130000
DTEND:20180%vT140000
DTSTAMP:%v
SEQUENCE:%v
SUMMARY:Night shift
LOCATION:Night Location
DESCRIPTION:Night
END:VEVENT
END:VCALENDAR`

	//attenddee needed for both events
	attendee := fmt.Sprintf(
		"ATTENDEE;ROLE=REQ-PARTICIPANT;CN=%s:MAILTO:%s\n",
		recipientSmtpAddr, recipientSmtpAddr)

	uid1 := day + 1000
	uid2 := uid1 + 1
	dtNow := time.Now().Format(timeutil.ISO8601CompactNoTZ)
	icsData := fmt.Sprintf(
		format, "REQUEST",
		attendee, uid1, day, day, dtNow, seq,
		attendee, uid2, day, day, dtNow, seq)
	return icsData
}

func main() {
	envFiles := []string{os.Getenv("ENV_PATH"), "./.env"}
	if _, err := config.LoadDotEnv(envFiles, -1); err != nil {
		log.Fatal().Err(err).
			Str("files", strings.Join(envFiles, ",")).
			Msg("Load env files failed")
	}

	cfg := appConfig{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatal().Err(err).
			Str("lib", "github.com/caarlos0/env/v9").
			Msgf("Cannot parse env to %s", "appConfig{}")
	}

	client, err := sparkpost.NewAPIClient(cfg.SparkPostApiKey)
	if err != nil {
		log.Fatal().Err(err).
			Str("lib", "github.com/goauth/sparkpost").
			Msg("SparkPost client init faile")
	}

	sendTestEmail(cfg, client)

	fmt.Println("DONE")
}
