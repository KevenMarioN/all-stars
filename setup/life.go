// Package setup is responsability by life of app
package setup

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

// On is the function responsible for keeping the application running.
func On(logger *slog.Logger) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	logger.Info("Service started successfully")
	logger.Info("To finish, press <Ctrl + C>")

	<-c
	os.Exit(0)
}
