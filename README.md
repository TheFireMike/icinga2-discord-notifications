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
  -h, --help                               help for icinga2-discord-notifications
      --host-last-state-change string      Host last changed timestamp
      --host-name string                   Host name
      --host-state string                  Host state
      --service-last-state-change string   Service last changed timestamp
      --service-name string                Service name
      --service-output string              Service output
      --service-state string               Service state
      --webhook string                     The webhook URL
```
Use it as a notification plugin command in Icinga. Sample configuration:
```
object NotificationCommand "discord-webhook" {
    import "plugin-notification-command"
    command = [ "/usr/lib/nagios/icinga2-discord-notifications" ]
    arguments += {
        "--host-last-state-change" = "$host.last-state-change$"
        "--host-name" = "$host.name$"
        "--host-state" = "$host.state$"
        "--service-last-state-change" = "$service.last_state_change$"
        "--service-name" = "$service.name$"
        "--service-output" = "$service.output$"
        "--service-state" = "$service.state$"
        "--webhook" = <YOUR GITHUB WEBHOOK>
    }
}
```