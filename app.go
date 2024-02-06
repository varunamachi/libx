package libx

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"github.com/varunamachi/libx/auth"
	"github.com/varunamachi/libx/httpx"
)

type BuildInfo struct {
	GitTag    string `json:"gitTag"`
	GitHash   string `json:"gitHash"`
	GitBranch string `json:"gitBranch"`
	BuildTime string `json:"buildTime"`
	BuildHost string `json:"buildHost"`
	BuildUser string `json:"buildUser"`
}

type App struct {
	*cli.App
	server    *httpx.Server
	buildInfo *BuildInfo
}

func NewApp(name, description, versionStr, author string) *App {
	app := &App{
		App: &cli.App{
			Name:        name,
			Description: description,
			Commands:    make([]*cli.Command, 0, 100),
			Version:     versionStr,
			Authors:     []*cli.Author{{Name: author}},
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "log-level",
					Value: "info",
					Usage: "Give log level, one of: 'trace', 'debug', " +
						"'info', 'warn', 'error'",
				},
			},
			Before: func(ctx *cli.Context) error {
				log.Logger = log.Output(
					zerolog.ConsoleWriter{Out: os.Stderr}).
					With().Logger()
				// With().Caller().Logger()
				logLevel := ctx.String("log-level")
				if logLevel != "" {
					level := zerolog.InfoLevel
					switch logLevel {
					case "trace":
						level = zerolog.TraceLevel
					case "debug":
						level = zerolog.DebugLevel
					case "info":
						level = zerolog.InfoLevel
					case "warn":
						level = zerolog.WarnLevel
					case "error":
						level = zerolog.ErrorLevel
					}
					zerolog.SetGlobalLevel(level)
				}
				return nil
			},
		},
	}
	app.Metadata = map[string]any{"app": app}
	return app
}

func NewCustomApp(cApp *cli.App) *App {
	app := &App{
		App: cApp,
	}
	app.Metadata = map[string]any{"app": app}
	return app
}

func (app *App) WithServer(port int, userGetter auth.UserRetrieverFunc) *App {
	app.server = httpx.NewServer(os.Stdout, userGetter)
	return app
}

func (app *App) WithEndpoints(ep ...*httpx.Endpoint) *App {
	app.server.WithAPIs(ep...)
	return app
}

func (app *App) WithCommands(cmds ...*cli.Command) *App {
	app.Commands = append(app.Commands, cmds...)
	return app
}

func (app *App) WithBuildInfo(bi *BuildInfo) *App {
	app.buildInfo = bi
	return app
}

func (app *App) BuildInfo() *BuildInfo {
	return app.buildInfo
}

func (app *App) Serve(port uint32) error {
	return app.server.Start(port)
}

func (app *App) StopServer() error {
	if app.server == nil {
		log.Trace().Msg("no running server found to stop")
		return nil
	}
	if err := app.server.Close(); err != nil {
		return err
	}
	app.server = nil
	return nil
}

func MustGetApp(ctx *cli.Context) *App {
	app, ok := ctx.App.Metadata["app"].(*App)
	if !ok {
		panic("invalid application type, please use libx.app.App")
	}
	return app
}
