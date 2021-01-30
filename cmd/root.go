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

	rootCMD.Flags().String("webhook", "", "The webhook URL")
	err := rootCMD.MarkFlagRequired("webhook")
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	rootCMD.Flags().String("host-name", "", "Host name")
	rootCMD.Flags().String("host-state", "", "Host state")
	rootCMD.Flags().String("host-last-state-change", "", "Host last changed timestamp")

	rootCMD.Flags().String("service-name", "", "Service name")
	rootCMD.Flags().String("service-state", "", "Service state")
	rootCMD.Flags().String("service-last-state-change", "", "Service last changed timestamp")
	rootCMD.Flags().String("service-output", "", "Service output")
}

var rootCMD = &cobra.Command{
	Use:   "icinga2-discord-notifications",
	Short: "Support for Discord notifications for Icinga2.",
	Run: func(cmd *cobra.Command, args []string) {
		event := notification.Event{
			HostName:          cmd.Flags().Lookup("host-name").Value.String(),
			HostState:         cmd.Flags().Lookup("host-state").Value.String(),
			HostLastChange:    cmd.Flags().Lookup("host-last-state-change").Value.String(),
			ServiceName:       cmd.Flags().Lookup("service-name").Value.String(),
			ServiceState:      cmd.Flags().Lookup("service-state").Value.String(),
			ServiceLastChange: cmd.Flags().Lookup("service-last-state-change").Value.String(),
			ServiceOutput:     cmd.Flags().Lookup("service-output").Value.String(),
		}

		notification.SendNotification(event, cmd.Flags().Lookup("webhook").Value.String())
	},
}

// Execute is the entrypoint for the CLI interface.
func Execute() {
	if err := rootCMD.Execute(); err != nil {
		os.Exit(1)
	}
}
