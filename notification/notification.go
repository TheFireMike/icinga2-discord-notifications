package notification

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
	"strconv"
)

const (
	colorRed   = "15158332"
	colorGreen = "3066993"
	colorGrey  = "9807270"
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
	Content string  `json:"content"`
	Embeds  []embed `json:"embeds"`
}

type embed struct {
	Title       string `json:"title"`
	Color       string `json:"color"`
	Description string `json:"description"`
}

//SendNotification sends the event to the specified discord webhook.
func SendNotification(event Event, webhook string) {
	if event.HostName == "" {
		log.Fatal().Msg("host name is missing")
	}

	var output output

	if event.ServiceState != "" {
		if event.ServiceName == "" {
			log.Fatal().Msg("service name is missing")
		}

		embed := embed{
			Title: "Service Output",
		}

		if event.ServiceState == event.ServiceLastState || event.ServiceLastState == "" {
			output.Content += "**INFO**: "
			embed.Color = colorGrey
		} else if event.ServiceState != "OK" {
			output.Content += "**PROBLEM**: "
			embed.Color = colorRed
		} else if event.ServiceState == "OK" {
			output.Content += "**RECOVER**: "
			embed.Color = colorGreen
		}

		output.Content += fmt.Sprintf("%s on %s is %s!",
			event.ServiceName,
			event.HostName,
			event.ServiceState)

		if event.ServiceOutput != "" {
			// the backticks put the service output in a code block
			embed.Description = "```" + event.ServiceOutput + "```"

			output.Embeds = append(output.Embeds, embed)
		}
	} else if event.HostState != "" {
		if event.HostState == event.HostLastState || event.HostLastState == "" {
			output.Content += "**INFO**: "
		} else if event.HostState != "UP" {
			output.Content += "**PROBLEM**: "
		} else if event.HostState == "UP" {
			output.Content += "**RECOVER**: "
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

	resp, err := resty.New().R().
		SetHeader("Content-type", "application/json").
		SetBody(string(outputJSON)).
		Post(webhook)
	if err != nil {
		log.Fatal().Err(err).Msg("sending message failed")
	}

	if resp.StatusCode() != 204 {
		log.Fatal().Str("response", string(resp.Body())).Str("status_code", strconv.Itoa(resp.StatusCode())).Msg("sending message failed")
	}
}
