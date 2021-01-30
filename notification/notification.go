package notification

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
	"strconv"
	"time"
)

type Event struct {
	HostName          string
	ServiceName       string
	HostState         string
	ServiceState      string
	HostLastChange    string
	ServiceLastChange string
	ServiceOutput     string
}

type output struct {
	Content string `json:"content"`
}

func SendNotification(event Event, webhook string) {
	var output output

	if event.ServiceState != "" {
		timestamp, err := strconv.Atoi(event.ServiceLastChange)
		if err != nil {
			timestamp = int(time.Now().Unix())
		}

		output.Content = fmt.Sprintf("%s: %s on %s is %s! Output:\n%s",
			time.Unix(int64(timestamp), 0).Format(time.RFC822),
			event.ServiceName,
			event.HostName,
			event.ServiceState,
			event.ServiceOutput)
	} else if event.HostState != "" {
		timestamp, err := strconv.Atoi(event.HostLastChange)
		if err != nil {
			timestamp = int(time.Now().Unix())
		}

		output.Content = fmt.Sprintf("%s: Host %s is %s!",
			time.Unix(int64(timestamp), 0).Format(time.RFC822),
			event.ServiceName,
			event.HostState)
	} else {
		log.Fatal().Msg("unknown event type")
	}

	outputJSON, err := json.Marshal(output)
	if err != nil {
		log.Fatal().Err(err).Msg("marshalling output failed")
	}

	client := resty.New()

	resp, err := client.R().
		SetHeader("Content-type", "application/json").
		SetBody(string(outputJSON)).
		Post(webhook)

	if err != nil {
		log.Fatal().Err(err).Msg("sending failed")
	}

	if resp.StatusCode() != 204 {
		log.Fatal().Msg(string(resp.Body()))
	}
}
