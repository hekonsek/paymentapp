package cmd

import (
	"github.com/hekonsek/paymentapp/api"
	"github.com/hekonsek/paymentapp/payments"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"strconv"
)

const UnixExitCodeGeneralError = 1

var startCommandPort string
var startCommandPersistence string

func init() {
	startCommand.Flags().StringVarP(&startCommandPort, "port", "p", "8080",
		"Specifies HTTP port for REST API. Can be overriden with PORT environment variable.")
	startCommand.Flags().StringVarP(&startCommandPersistence, "persistence", "", "mem",
		"Specifies persistence engine to be used (mem | awsdocdb). Can be overriden with PERSISTENCE environment variable.")
	RootCmd.AddCommand(startCommand)
}

var startCommand = &cobra.Command{
	Use:   "start",
	Short: "Starts gateway API application and blocks until Unix interruption signal is received.",
	Run: func(cmd *cobra.Command, args []string) {
		port := os.Getenv("PORT")
		if port == "" {
			port = startCommandPort
		}
		portInt, err := strconv.Atoi(startCommandPort)
		if err != nil {
			log.Fatalf("Error when starting HTTP server: %s", err.Error())
			os.Exit(UnixExitCodeGeneralError)
		}

		persistence := resolvePersistence()
		var persistenceStore payments.PaymentStore
		if persistence == "mem" {
			persistenceStore = payments.NewInMemoryPaymentStore()
		} else {
			persistenceStore, err = payments.NewDocdbPaymentStore("", -1)
			if err != nil {
				log.Fatalf("Error when starting HTTP server: %s", err.Error())
				os.Exit(UnixExitCodeGeneralError)
			}
		}

		a := api.ApiServer{
			Port:  portInt,
			Store: persistenceStore,
		}
		err = a.Start()
		if err != nil {
			log.Fatalf("Error when starting HTTP server: %s", err.Error())
			os.Exit(UnixExitCodeGeneralError)
		}
		a.WaitForInterruptSignal()
	},
}

func resolvePersistence() string {
	persistence := os.Getenv("PERSISTENCE")
	if persistence == "" {
		persistence = startCommandPersistence
	}
	if persistence != "mem" && persistence != "awsdocdb" {
		log.Fatalf("Error when starting HTTP server: Unknown persistence provider specified.")
		os.Exit(UnixExitCodeGeneralError)
	}
	return persistence
}
