package main

import (
	"github.com/hekonsek/paymentapp/cmd"
	log "github.com/sirupsen/logrus"
)

func main() {
	err := cmd.RootCmd.Execute()
	if err != nil {
		log.Fatalf("Error executing command: %s", err.Error())
	}
}
