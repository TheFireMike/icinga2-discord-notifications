package notification

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
	"strconv"
)

const (
	colorRed    = "15158332"
	colorOrange = "15105570"
	colorGreen  = "3066993"
	colorPurple = "10181046"
	colorGrey   = "9807270"
)

//Event includes all information about the event which should be reported.
type Event struct {
	NotificationType    string
	NotificationAuthor  string
	NotificationComment string
	HostName            string
	HostState           string
	HostOutput          string
	ServiceName         string
	ServiceState        string
	ServiceOutput       string
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
	var output output

	if event.ServiceState != "" {
		if event.ServiceName == "" {
			log.Fatal().Msg("service name is missing")
		}
		output = getServiceOutput(event)
	} else if event.HostState != "" {
		output = getHostOutput(event)
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

func getServiceOutput(event Event) output {
	var output output

	output.Content = fmt.Sprintf("**%s**: %s on %s is %s!",
		event.NotificationType,
		event.ServiceName,
		event.HostName,
		event.ServiceState)

	if event.NotificationType == "CUSTOM" ||
		event.NotificationType == "ACKNOWLEDGEMENT" ||
		event.NotificationType == "DOWNTIMESTART" ||
		event.NotificationType == "DOWNTIMEEND" {
		if event.NotificationComment != "" {
			output.Embeds = append(output.Embeds, getComment(event))
		}
	} else if event.ServiceOutput != "" {
		output.Embeds = append(output.Embeds, embed{
			Title: "Service Output",
			Color: getColor(event.ServiceState),
			// the backticks put the service output in a code block
			Description: fmt.Sprintf("```%s```", event.ServiceOutput),
		})
	}

	return output
}

func getHostOutput(event Event) output {
	var output output

	output.Content = fmt.Sprintf("**%s**: Host %s is %s!",
		event.NotificationType,
		event.HostName,
		event.HostState)

	if event.NotificationType == "CUSTOM" ||
		event.NotificationType == "ACKNOWLEDGEMENT" ||
		event.NotificationType == "DOWNTIMESTART" ||
		event.NotificationType == "DOWNTIMEEND" {
		if event.NotificationComment != "" {
			output.Embeds = append(output.Embeds, getComment(event))
		}
	} else if event.HostOutput != "" {
		output.Embeds = append(output.Embeds, embed{
			Title: "Host Output",
			Color: getColor(event.HostState),
			// the backticks put the host output in a code block
			Description: fmt.Sprintf("```%s```", event.HostOutput),
		})
	}

	return output
}

func getComment(event Event) embed {
	title := "Comment"
	if event.NotificationAuthor != "" {
		title += fmt.Sprintf(" by %s", event.NotificationAuthor)
	}

	return embed{
		Title:       title,
		Color:       colorGrey,
		Description: event.NotificationComment,
	}
}

func getColor(state string) string {
	switch state {
	case "UP", "OK":
		return colorGreen
	case "WARNING":
		return colorOrange
	case "DOWN", "CRITICAL":
		return colorRed
	case "UNKNOWN":
		return colorPurple
	default:
		return colorGrey
	}
}
