package main

import (
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"github.com/varunamachi/idx/core"
	"github.com/varunamachi/libx"
	"github.com/varunamachi/libx/errx"
	"github.com/varunamachi/libx/proc"
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
			stopAllCmd(),
			infoCmd(),
		)

	style := lipgloss.
		NewStyle().
		Foreground(lipgloss.Color("212")).
		Bold(true).
		Align(lipgloss.Left)

	beforeBefore := app.Before
	app.Before = func(ctx *cli.Context) error {
		if err := beforeBefore(ctx); err != nil {
			return err
		}
		log.Logger = log.Output(
			zerolog.ConsoleWriter{
				Out: proc.NewWriter("*****", os.Stderr, style, false)}).
			With().Logger()
		return nil
	}

	if err := app.RunContext(gtx, os.Args); err != nil {
		errx.PrintSomeStack(err)
		log.Fatal().Msg("Exiting due to errors")
	}
}
