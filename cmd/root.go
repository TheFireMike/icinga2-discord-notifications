package cmd

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/thefiremike/icinga2-discord-notifications/notification"
	"os"
)

func init() {
	log.Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()

	rootCMD.Flags().String("webhook", "", "webhook URL")

	rootCMD.Flags().String("notification-type", "", "notification type")
	rootCMD.Flags().String("notification-author", "", "notification author")
	rootCMD.Flags().String("notification-comment", "", "notification comment")

	rootCMD.Flags().String("host-name", "", "host (display) name")
	rootCMD.Flags().String("host-state", "", "host state")
	rootCMD.Flags().String("host-output", "", "host output")

	rootCMD.Flags().String("service-name", "", "service (display) name")
	rootCMD.Flags().String("service-state", "", "service state")
	rootCMD.Flags().String("service-output", "", "service output")

	err := rootCMD.MarkFlagRequired("webhook")
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	err = rootCMD.MarkFlagRequired("notification-type")
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	err = rootCMD.MarkFlagRequired("host-name")
	if err != nil {
		log.Fatal().Err(err).Send()
	}
}

var rootCMD = &cobra.Command{
	Use:   "icinga2-discord-notifications",
	Short: "Support for Discord notifications for Icinga2.",
	Run: func(cmd *cobra.Command, args []string) {
		notification.SendNotification(notification.Event{
			NotificationType:    cmd.Flags().Lookup("notification-type").Value.String(),
			NotificationAuthor:  cmd.Flags().Lookup("notification-author").Value.String(),
			NotificationComment: cmd.Flags().Lookup("notification-comment").Value.String(),
			HostName:            cmd.Flags().Lookup("host-name").Value.String(),
			HostState:           cmd.Flags().Lookup("host-state").Value.String(),
			HostOutput:          cmd.Flags().Lookup("host-output").Value.String(),
			ServiceName:         cmd.Flags().Lookup("service-name").Value.String(),
			ServiceState:        cmd.Flags().Lookup("service-state").Value.String(),
			ServiceOutput:       cmd.Flags().Lookup("service-output").Value.String(),
		}, cmd.Flags().Lookup("webhook").Value.String())
	},
}

// Execute is the entrypoint for the CLI interface.
func Execute() {
	if err := rootCMD.Execute(); err != nil {
		os.Exit(1)
	}
}
