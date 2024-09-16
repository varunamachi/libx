package main

import (
	"os"

	"github.com/rs/zerolog/log"
	"github.com/varunamachi/idx/core"
	"github.com/varunamachi/idx/pg/schema"
	"github.com/varunamachi/libx"
	"github.com/varunamachi/libx/errx"
	"github.com/varunamachi/libx/rt"
)

func main() {
	gtx, cancel := rt.Gtx()
	defer cancel()

	app := libx.NewApp(
		"exman", "A process manager", "0.0.1", "varunamachi@gmail.com").
		WithBuildInfo(core.GetBuildInfo()).
		WithCommands(
			serveCmd(gtx),
			execCmd(),
			listCmd(),
			stopCmd(),
		)

	if err := schema.Init(gtx, "onServerStart"); err != nil {
		log.Fatal().Err(err).Msg("DB init failed")
	}

	if err := app.RunContext(gtx, os.Args); err != nil {
		errx.PrintSomeStack(err)
		log.Fatal().Msg("Exiting due to errors")
	}
}
