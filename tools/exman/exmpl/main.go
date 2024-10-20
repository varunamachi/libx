package main

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/varunamachi/libx/rt"
)

func main() {
	log.Logger = log.Output(
		zerolog.ConsoleWriter{Out: os.Stderr}).
		With().Logger()

	gtx, cancel := rt.Gtx()
	defer cancel()

	if len(os.Args) < 2 {
		log.Fatal().Msg("Instance name is required")
	}

	tkr := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-gtx.Done():
			log.Warn().Msg("Termination initiated ;)")
			os.Exit(0)
		case <-tkr.C:
			log.Info().Str("proc", os.Args[1]).Msg("")
		}
	}

}
