package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "paymentapp",
	Short: "paymentapp - simple payment API gateway",

	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Help()
		if err != nil {

			log.Fatalf("Error executing command: %s", err.Error())
		}
	},
}
