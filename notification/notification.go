package notification

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
)

//Event includes all information about the event which should be reported.
type Event struct {
	HostName         string
	ServiceName      string
	HostState        string
	ServiceState     string
	HostLastState    string
	ServiceLastState string
	ServiceOutput    string
}

type output struct {
	Content string `json:"content"`
}

//SendNotification sends the event to the specified discord webhook.
func SendNotification(event Event, webhook string) {
	var output output

	if event.ServiceState != "" {
		if event.ServiceState == event.ServiceLastState || event.ServiceLastState == "" {
			output.Content += "INFO: "
		} else if event.ServiceState != "OK" {
			output.Content += "PROBLEM: "
		} else if event.ServiceState == "OK" {
			output.Content += "RECOVER: "
		}

		output.Content += fmt.Sprintf("%s on %s is %s! Output:\n%s",
			event.ServiceName,
			event.HostName,
			event.ServiceState,
			event.ServiceOutput)
	} else if event.HostState != "" {
		if event.HostState == event.HostLastState || event.HostLastState == "" {
			output.Content += "INFO: "
		} else if event.HostState != "UP" {
			output.Content += "PROBLEM: "
		} else if event.HostState == "UP" {
			output.Content += "RECOVER: "
		}

		output.Content += fmt.Sprintf("Host %s is %s!",
			event.HostName,
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
