package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/KevenMarioN/all-stars/server"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	debug := flag.Bool("debug", false, "sets log level to debug")
	flag.Parse()
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	output.FormatLevel = func(i any) string {
		var l string
		if ll, ok := i.(string); ok {
			l = ll
		}

		formattedLevel := strings.ToUpper(fmt.Sprintf("| %-6s|", l))
		switch l {
		case "debug":
			return "\x1b[33m" + formattedLevel + "\x1b[0m"
		case "info":
			return "\x1b[32m" + formattedLevel + "\x1b[0m"
		case "warn":
			return "\x1b[31m" + formattedLevel + "\x1b[0m"
		case "error", "fatal", "panic":
			return "\x1b[31m\x1b[1m" + formattedLevel + "\x1b[0m"
		default:
			return formattedLevel
		}
	}
	output.FormatMessage = func(i any) string {
		return fmt.Sprintf("%s", i)
	}
	output.FormatFieldName = func(i any) string {
		return fmt.Sprintf("\x1b[36m%s:\x1b[0m", i)
	}
	output.FormatFieldValue = func(i any) string {
		return fmt.Sprintf("\x1b[36m%s\x1b[0m", i)
	}

	multi := zerolog.MultiLevelWriter(output, os.Stdout)
	log.Logger = log.Output(multi)
	log.Debug().Msg("This logger is better!")
	log.Error().Msg("Okay any")
	srv := server.NewServer()
	srv.Get("health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusContinue)
	})
	v1 := srv.Group("v1")
	v1.Post("/auth", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	blocks := v1.Group("blocks")
	blocks.Get("", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
	})
	blocks.Get("/luck/{id}/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusContinue)
	})

	if err := srv.Run(7777); err != nil {
		log.Error().Err(err)
	}
}
