# icinga2-discord-notifications
[![Go Report Card](https://goreportcard.com/badge/github.com/thefiremike/icinga2-discord-notifications)](https://goreportcard.com/report/github.com/thefiremike/icinga2-discord-notifications)
[![GitHub license](https://img.shields.io/badge/license-BSD-blue.svg)](https://github.com/thefiremike/icinga2-discord-notifications/blob/main/LICENSE)
[![GitHub code style](https://img.shields.io/badge/code%20style-uber--go-brightgreen)](https://github.com/uber-go/guide/blob/master/style.md)
[![GoDoc doc](https://img.shields.io/badge/godoc-reference-blue)](https://godoc.org/github.com/thefiremike/icinga2-discord-notifications)
```
Support for Discord notifications for Icinga2.

Usage:
  icinga2-discord-notifications [flags]

Flags:
  -h, --help                          help for icinga2-discord-notifications
      --host-name string              host (display) name
      --host-output string            host output
      --host-state string             host state
      --notification-author string    notification author
      --notification-comment string   notification comment
      --notification-type string      notification type
      --service-name string           service (display) name
      --service-output string         service output
      --service-state string          service state
      --webhook string                webhook URL
```

## Installation
You can download the latest release under the `Releases` tab or build it yourself with `go build`.

Use it as a notification plugin command in Icinga. Sample configuration:
```
object NotificationCommand "discord-webhook" {
    import "plugin-notification-command"
    command = [ ConfigDir + "/scripts/icinga2-discord-notifications" ]
    arguments += {
        "--host-name" = "$host.display_name$"
        "--host-output" = "$host.output$"
        "--host-state" = "$host.state$"
        "--notification-author" = "$notification.author$"
        "--notification-comment" = "$notification.comment$"
        "--notification-type" = "$notification.type$"
        "--service-name" = "$service.display_name$"
        "--service-output" = "$service.output$"
        "--service-state" = "$service.state$"
        "--webhook" = <YOUR DISCORD WEBHOOK>
    }
}
```
