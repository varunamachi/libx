package pg

import (
	"net/url"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"github.com/varunamachi/libx/errx"
)

func Wrap(cmd *cli.Command) *cli.Command {
	cmd.Flags = append(cmd.Flags,
		&cli.StringFlag{
			Name: "pg-url",
			Usage: "URL (NOT CONNECTION_STRING) of the host running postgres" +
				", this has the highest  precidence",
			EnvVars:  []string{"PG_URL"},
			Required: false,
		},
		&cli.StringFlag{
			Name:     "pg-host",
			Value:    "localhost",
			Usage:    "Address of the host running postgres",
			EnvVars:  []string{"PG_HOST"},
			Required: false,
		},
		&cli.IntFlag{
			Name:     "pg-port",
			Value:    5432,
			Usage:    "Port on which postgres is listening",
			EnvVars:  []string{"PG_PORT"},
			Required: false,
		},
		&cli.StringFlag{
			Name:     "pg-db",
			Value:    "",
			Usage:    "Database name",
			EnvVars:  []string{"PG_DB"},
			Required: false,
		},
		&cli.StringFlag{
			Name:     "pg-user",
			Usage:    "Postgres user name",
			EnvVars:  []string{"PG_USER"},
			Required: false,
		},
		&cli.StringFlag{
			Name:     "pg-pass",
			Value:    "",
			Usage:    "Postgres password for connection",
			EnvVars:  []string{"PG_PASS"},
			Required: false,
		},
		// TODO - use this value in opts
		&cli.StringFlag{
			Name:     "pg-timezone",
			Value:    "Asia/Kolkata",
			Usage:    "Postgres client time zone",
			EnvVars:  []string{"PG_TIMEZONE"},
			Required: false,
		},
	)

	if cmd.Before == nil {
		cmd.Before = RequirePostgres
	} else {
		otherBefore := cmd.Before
		cmd.Before = func(ctx *cli.Context) (err error) {
			err = RequirePostgres(ctx)
			if err == nil {
				err = otherBefore(ctx)
			}
			return errx.Wrap(err)
		}
	}

	return cmd
}

func RequirePostgres(ctx *cli.Context) error {
	urlStr := ctx.String("pg-url")
	if urlStr != "" {
		u, err := url.Parse(urlStr)
		if err != nil {
			return errx.Errf(err, "invalid URL '%s' given")
		}

		db, err := Connect(ctx.Context, u, ctx.String("pg-timezone"))
		if err != nil {
			log.Fatal().Err(err).Msg("failed to connect to database")
			return errx.Wrap(err)
		}
		SetDefaultConn(db)
	} else {
		db, err := ConnectWithOpts(ctx.Context, &ConnOpts{
			Host:     ctx.String("pg-host"),
			Port:     ctx.Int("pg-port"),
			User:     ctx.String("pg-user"),
			DBName:   ctx.String("pg-db"),
			Password: ctx.String("pg-pass"),
			TimeZone: ctx.String("pg-timezone"),
		})
		if err != nil {
			log.Fatal().Err(err).Msg("failed to connect to database")
			return errx.Wrap(err)
		}
		SetDefaultConn(db)
	}
	if err := defDB.Ping(); err != nil {
		log.Fatal().Err(err).Msg("failed to ping database")
	}
	return nil
}
