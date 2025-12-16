// Package setup is responsability by life of app
package setup

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"
)

// On is the function responsible for keeping the application running.
func On() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	log.Info().Msg("Service started successfully")
	log.Info().Msg("To finish, press <Ctrl + C>")

	<-c
	os.Exit(0)
}
